// Copyright 2015 Red Hat Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package man

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// GenerateManOptions is used configure how GenerateManPages will
// do its job.
type CobraManOptions struct {
	// What section to generate the pages 4 (1 is the default if not set)
	Section string

	// CenterFooter used across all pages (defaults to current month and year)
	// If you just want to set the date used in the center footer use Date
	CenterFooter string

	// If you just want to set the date used in the center footer use Date
	// Will default to Now
	Date *time.Time

	// LeftFooter used across all pages
	LeftFooter string

	// CenterHeader used across all pages
	CenterHeader string

	// Files if set with content will create a FILES section for all
	// pages.  If you want this section only for a single command add
	// it as an annotation: cmd.Annotations["man-files-section"]
	// The field will be sanitized for troff output. However, if
	// it starts with a '.' we assume it is valid troff and pass it through.
	Files string

	// Bugs if set with content will create a BUGS section for all
	// pages.  If you want this section only for a single command add
	// it as an annotation: cmd.Annotations["man-bugs-section"]
	// The field will be sanitized for troff output. However, if
	// it starts with a '.' we assume it is valid troff and pass it through.
	Bugs string

	// Environment if set with content will create a ENVIRONMENT section for all
	// pages.  If you want this section only for a single command add
	// it as an annotation: cmd.Annotations["man-environment-section"]
	// The field will be sanitized for troff output. However, if
	// it starts with a '.' we assume it is valid troff and pass it through.
	Environment string

	// Author if set will create a Author section with this content.
	Author string

	// Directory location for where to generate the man pages
	Directory string

	// FileCmdSeparator defines what character to use to separate the
	// sub commands in the man page file name.  The '-' char is the default.
	FileCmdSeparator string

	// FileSuffix is the file extension to use for file name.  Defaults to the section
	// for man templates and .md for the MarkdownTemplate template.
	FileSuffix string

	// UseTemplate allows you to set the template used to generate
	// documentation.  The default is defined in man.TroffManTemplate
	UseTemplate string
}

// GenerateManPages - build man pages for the passed in cobra.Command
// and all of its children
func GenerateManPages(cmd *cobra.Command, opts *CobraManOptions) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := GenerateManPages(c, opts); err != nil {
			return err
		}
	}

	// Set defaults
	setCobraManOptDefaults(opts)

	// Generate file name and open the file
	basename := strings.Replace(cmd.CommandPath(), " ", opts.FileCmdSeparator, -1)
	if basename == "" {
		return fmt.Errorf("you need a command name to have a man page")
	}
	filename := filepath.Join(opts.Directory, basename+"."+opts.FileSuffix)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Generate the documentation
	return GenerateOnePage(cmd, opts, f)
}

func setCobraManOptDefaults(opts *CobraManOptions) {
	if opts.Section == "" {
		opts.Section = "1"
	}
	if opts.Date == nil {
		now := time.Now()
		opts.Date = &now
	}
	if opts.FileCmdSeparator == "" {
		if opts.UseTemplate == MarkdownTemplate {
			opts.FileCmdSeparator = "_"
		} else {
			opts.FileCmdSeparator = "-"
		}
	}
	if opts.UseTemplate == "" {
		opts.UseTemplate = TroffManTemplate
	}
	if opts.FileSuffix == "" {
		if opts.UseTemplate == MarkdownTemplate {
			opts.FileSuffix = "md"
		} else {
			opts.FileSuffix = opts.Section
		}
	}
}

type manStruct struct {
	Date             *time.Time
	Section          string
	CenterFooter     string
	LeftFooter       string
	CenterHeader     string
	UseLine          string
	CommandPath      string
	ShortDescription string
	Description      string
	NoArgs           bool

	AllFlags          []Flag
	InheritedFlags    []Flag
	NonInheritedFlags []Flag
	SeeAlsos          []SeeAlso
	SubCommands       []string

	Author      string
	Environment string
	Files       string
	Bugs        string
	Examples    string
}

type Flag struct {
	Shorthand   string
	Name        string
	NoOptDefVal string
	DefValue    string
	Usage       string
	ArgHint     string
}

type SeeAlso struct {
	CmdPath string
	Section string
}

// GenerateOnePage will generate one documentation page and output the result to w
// TODO: document use of this function in README
func GenerateOnePage(cmd *cobra.Command, opts *CobraManOptions, w io.Writer) error {
	// Set defaults - these would already be set unless GenerateOnePage called directly
	setCobraManOptDefaults(opts)

	values := manStruct{}

	// Header fields
	values.LeftFooter = opts.LeftFooter
	values.CenterHeader = opts.CenterHeader
	values.Section = opts.Section
	values.Date = opts.Date
	values.CenterFooter = opts.CenterFooter
	if opts.CenterFooter == "" {
		// TODO: should this be part of template instead?
		values.CenterFooter = values.Date.Format("Jan 2006")
	}

	values.ShortDescription = cmd.Short
	values.UseLine = cmd.UseLine()
	values.CommandPath = cmd.CommandPath()

	// Use reflection to see if cobra.NoArgs was set
	argFuncName := runtime.FuncForPC(reflect.ValueOf(cmd.Args).Pointer()).Name()
	values.NoArgs = strings.HasSuffix(argFuncName, "cobra.NoArgs")

	if cmd.HasSubCommands() {
		subCmdArr := make([]string, 0, 10)
		for _, c := range cmd.Commands() {
			if c.IsAdditionalHelpTopicCommand() {
				continue
			}
			subCmdArr = append(subCmdArr, c.CommandPath())
		}
		values.SubCommands = subCmdArr
	}

	// DESCRIPTION
	description := cmd.Long
	if len(description) == 0 {
		description = cmd.Short
	}
	values.Description = description

	// Flag arrays
	values.AllFlags = genFlagArray(cmd.Flags())
	values.InheritedFlags = genFlagArray(cmd.InheritedFlags())
	values.NonInheritedFlags = genFlagArray(cmd.NonInheritedFlags())

	// ENVIRONMENT section
	altEnvironmentSection, _ := cmd.Annotations["man-environment-section"]
	if opts.Environment != "" || altEnvironmentSection != "" {
		if altEnvironmentSection != "" {
			values.Environment = altEnvironmentSection
		} else {
			values.Environment = opts.Environment
		}
	}

	// FILES section
	altFilesSection, _ := cmd.Annotations["man-files-section"]
	if opts.Files != "" || altFilesSection != "" {
		if altFilesSection != "" {
			values.Files = altFilesSection
		} else {
			values.Files = opts.Files
		}
	}

	// BUGS section
	altBugsSection, _ := cmd.Annotations["man-bugs-section"]
	if opts.Bugs != "" || altBugsSection != "" {
		if altBugsSection != "" {
			values.Bugs = altBugsSection
		} else {
			values.Bugs = opts.Bugs
		}
	}

	// EXAMPLES section
	altExampleSection, _ := cmd.Annotations["man-examples-section"]
	if cmd.Example != "" || altExampleSection != "" {
		if altExampleSection != "" {
			values.Examples = altExampleSection
		} else {
			values.Examples = cmd.Example
		}
	}

	// AUTHOR section
	values.Author = opts.Author

	// SEE ALSO section
	values.SeeAlsos = generateSeeAlsos(cmd, values.Section)

	// Build the template and generate the man page
	funcMap := template.FuncMap{
		"upper":         strings.ToUpper,
		"backslashify":  backslashify,
		"dashify":       dashify,
		"underscoreify": underscoreify,
		"simpleToTroff": simpleToTroff,
		"simpleToMdoc":  simpleToMdoc,
	}
	parsedTemplate, err := template.New("man").Funcs(funcMap).Parse(opts.UseTemplate)
	if err != nil {
		return err
	}
	err = parsedTemplate.Execute(w, values)
	if err != nil {
		return err
	}
	return nil
}

func genFlagArray(flags *pflag.FlagSet) []Flag {
	flagArray := make([]Flag, 0, 15)
	flags.VisitAll(func(flag *pflag.Flag) {
		if len(flag.Deprecated) > 0 || flag.Hidden {
			return
		}
		manFlag := Flag{
			Name:        flag.Name,
			NoOptDefVal: flag.NoOptDefVal,
			DefValue:    flag.DefValue,
			Usage:       flag.Usage,
		}
		if len(flag.ShorthandDeprecated) == 0 {
			manFlag.Shorthand = flag.Shorthand
		}
		hintArr, exists := flag.Annotations["man-arg-hints"]
		if exists && len(hintArr) > 0 {
			manFlag.ArgHint = hintArr[0]
		}
		flagArray = append(flagArray, manFlag)
	})

	return flagArray
}

func generateSeeAlsos(cmd *cobra.Command, section string) []SeeAlso {
	seealsos := make([]SeeAlso, 0)
	if cmd.HasParent() {
		see := SeeAlso{
			CmdPath: cmd.Parent().CommandPath(),
			Section: section,
		}
		seealsos = append(seealsos, see)
		// TODO: may want to control if siblings are shown or not
		siblings := cmd.Parent().Commands()
		sort.Sort(byName(siblings))
		for _, c := range siblings {
			if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() || c.Name() == cmd.Name() {
				continue
			}
			see := SeeAlso{
				CmdPath: c.CommandPath(),
				Section: section,
			}
			seealsos = append(seealsos, see)
		}
	}
	children := cmd.Commands()
	sort.Sort(byName(children))
	for _, c := range children {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		see := SeeAlso{
			CmdPath: c.CommandPath(),
			Section: section,
		}
		seealsos = append(seealsos, see)
	}

	return seealsos
}
