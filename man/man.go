// Copyright © 2018 Ray Johnson <ray.johnson@gmail.com>.
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

	// TemplateName allows you to set the template used to generate
	// documentation.  Templates need to be registered via RegisterTemplate.
	// The default template is "troff"
	TemplateName string

	// Private fields

	// fileCmdSeparator defines what character to use to separate the
	// sub commands in the man page file name.  The '-' char is the default.
	fileCmdSeparator string

	// fileSuffix is the file extension to use for file name.  Defaults to the section
	// for man templates and .md for the MarkdownTemplate template.
	fileSuffix string
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
	basename := strings.Replace(cmd.CommandPath(), " ", opts.fileCmdSeparator, -1)
	if basename == "" {
		return fmt.Errorf("you need a command name to have a man page")
	}
	filename := filepath.Join(opts.Directory, basename+"."+opts.fileSuffix)
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

	if opts.TemplateName == "" {
		opts.TemplateName = "troff"
	}
	sep, ext, t := getTemplate(opts.TemplateName)
	if t == nil {
		panic("template could not be found: " + opts.TemplateName)
	}
	opts.fileCmdSeparator = sep
	opts.fileSuffix = ext
	if ext == "use_section" {
		opts.fileSuffix = opts.Section
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

	AllFlags          []manFlag
	InheritedFlags    []manFlag
	NonInheritedFlags []manFlag
	SeeAlsos          []seeAlso
	SubCommands       []string

	Author      string
	Environment string
	Files       string
	Bugs        string
	Examples    string
}

type manFlag struct {
	Shorthand   string
	Name        string
	NoOptDefVal string
	DefValue    string
	Usage       string
	ArgHint     string
}

type seeAlso struct {
	CmdPath   string
	Section   string
	IsParent  bool
	IsChild   bool
	IsSibling bool
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

	// Get template and generate the documentation page
	_, _, t := getTemplate(opts.TemplateName)
	err := t.Execute(w, values)
	if err != nil {
		return err
	}
	return nil
}

func genFlagArray(flags *pflag.FlagSet) []manFlag {
	flagArray := make([]manFlag, 0, 15)
	flags.VisitAll(func(flag *pflag.Flag) {
		if len(flag.Deprecated) > 0 || flag.Hidden {
			return
		}
		thisFlag := manFlag{
			Name:        flag.Name,
			NoOptDefVal: flag.NoOptDefVal,
			DefValue:    flag.DefValue,
			Usage:       flag.Usage,
		}
		if len(flag.ShorthandDeprecated) == 0 {
			thisFlag.Shorthand = flag.Shorthand
		}
		hintArr, exists := flag.Annotations["man-arg-hints"]
		if exists && len(hintArr) > 0 {
			thisFlag.ArgHint = hintArr[0]
		}
		flagArray = append(flagArray, thisFlag)
	})

	return flagArray
}

func generateSeeAlsos(cmd *cobra.Command, section string) []seeAlso {
	seealsos := make([]seeAlso, 0)
	if cmd.HasParent() {
		see := seeAlso{
			CmdPath:  cmd.Parent().CommandPath(),
			Section:  section,
			IsParent: true,
		}
		seealsos = append(seealsos, see)
		siblings := cmd.Parent().Commands()
		sort.Sort(byName(siblings))
		for _, c := range siblings {
			if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() || c.Name() == cmd.Name() {
				continue
			}
			see := seeAlso{
				CmdPath:   c.CommandPath(),
				Section:   section,
				IsSibling: true,
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
		see := seeAlso{
			CmdPath: c.CommandPath(),
			Section: section,
			IsChild: true,
		}
		seealsos = append(seealsos, see)
	}

	return seealsos
}
