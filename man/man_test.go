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
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func checkForFile(t *testing.T, str string) {
	if _, err := os.Stat(str); err == nil {
		os.Remove(str)
		return
	}
	assert.Fail(t, "Expected file does not exist: "+str)
}

func TestGenerateManPages(t *testing.T) {
	var err error

	opts := CobraManOptions{}
	cmd := &cobra.Command{}
	err = GenerateDocs(cmd, &opts, "", "troff")
	assert.Equal(t, "you need a command name to have a man page", err.Error())

	cmd = &cobra.Command{Use: "foo"}
	assert.Nil(t, GenerateDocs(cmd, &opts, "", "troff"))
	checkForFile(t, "foo.1")

	opts = CobraManOptions{Section: "8"}
	assert.Nil(t, GenerateDocs(cmd, &opts, "", "troff"))
	checkForFile(t, "foo.8")

	cmd = &cobra.Command{Use: "foo"}
	cmd2 := &cobra.Command{Use: "bar", Run: func(cmd *cobra.Command, args []string) {}}
	cmd.AddCommand(cmd2)
	opts = CobraManOptions{}
	assert.Nil(t, GenerateDocs(cmd, &opts, "", "troff"))
	checkForFile(t, "foo.1")
	checkForFile(t, "foo-bar.1")

	cmd = &cobra.Command{Use: "foo"}
	cmd2 = &cobra.Command{Use: "bar", Run: func(cmd *cobra.Command, args []string) {}}
	cmd.AddCommand(cmd2)
	assert.Nil(t, GenerateDocs(cmd, &opts, "", "markdown"))
	checkForFile(t, "foo.md")
	checkForFile(t, "foo_bar.md")

}

func TestSetCobraManOptDefaults(t *testing.T) {
	opts := CobraManOptions{}

	validate(&opts, "troff")
	assert.Equal(t, opts.Section, "1")
	assert.Equal(t, opts.fileCmdSeparator, "-")
	assert.Equal(t, opts.fileSuffix, "1")

	delta := time.Now().Sub(*opts.Date)
	if delta.Seconds() >= 1 {
		assert.Fail(t, "time difference too large")
	}

	opts = CobraManOptions{}
	validate(&opts, "markdown")
	assert.Equal(t, opts.Section, "1")
	assert.Equal(t, opts.fileCmdSeparator, "_")
	assert.Equal(t, opts.fileSuffix, "md")

	opts = CobraManOptions{}
	assert.Panics(t, func() { validate(&opts, "no exist") }, "should have paniced")
}

func TestGenerateManPageRequired(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{Use: "foo"}
	opts := CobraManOptions{}

	// Test header options
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".TH \"FOO\" \"1\" \".*\" \"\" \"\"", buf.String())

	buf.Reset()
	opts = CobraManOptions{LeftFooter: "kitty kat", CenterHeader: "Hello", CenterFooter: "meow", Section: "3"}
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".TH \"FOO\" \"3\" \"meow\" \"kitty kat\" \"Hello\"", buf.String())

	buf.Reset()
	date, _ := time.Parse(time.RFC3339, "1968-06-21T15:04:05Z")
	opts = CobraManOptions{Date: &date}
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".TH \"FOO\" \"1\" \"Jun 1968\" \"\" \"\"", buf.String())

	// Test name
	cmd = &cobra.Command{Use: "bar"}
	opts = CobraManOptions{}
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH NAME\nbar\n", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "bar", Short: "going to"}
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH NAME\nbar - going to", buf.String())

	// Test Synopsis
	assert.Regexp(t, ".SH SYNOPSIS\n.sp\n.+bar", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "foo"}
	cmd2 := &cobra.Command{Use: "cat", Run: func(cmd *cobra.Command, args []string) {}}
	cmd3 := &cobra.Command{Use: "dog", Run: func(cmd *cobra.Command, args []string) {}}
	cmd.AddCommand(cmd2, cmd3)
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH SYNOPSIS\n.sp\n.+foo cat.+flags.+\n.br\n.+foo dog", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "foo"}
	cmd.Flags().String("thing", "", "string with no default")
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, "SH SYNOPSIS\n.sp\n.+foo.+\\\\-\\\\-thing.+<args>]", buf.String())

	// Test DESCRIPTION
	buf.Reset()
	cmd = &cobra.Command{Use: "bar", Short: "a short one"}
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, "SH DESCRIPTION\n.PP\na short one", buf.String())

	cmd.Long = `Long desc

This is long & stuff.`
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH DESCRIPTION\n.PP\nLong desc\n.PP\nThis is long \\\\& stuff.", buf.String())

}

func TestCobraManOptions(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{Use: "foo"}
	opts := CobraManOptions{}

	cmd = &cobra.Command{Use: "foo"}
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.NotRegexp(t, ".SH OPTIONS\n", buf.String()) // No OPTIONS section if no flags

	cmd.Flags().String("flag", "", "string with no default")
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH OPTIONS\n.TP\n.+flag.+\nstring with no default", buf.String()) // No OPTIONS section if no flags

	cmd.Flags().String("hello", "world", "default is world")
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".TP\n.+flag.+\nstring with no default", buf.String()) // No OPTIONS section if no flags
}

func TestGenerateManPageAltSections(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{Use: "foo"}
	opts := CobraManOptions{}

	// ENVIRONMENT
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.NotRegexp(t, ".SH ENVIRONMENT\n", buf.String()) // No ENVIRONMENT section if not in opts

	opts = CobraManOptions{Environment: "This uses ENV"}
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH ENVIRONMENT\n.PP\nThis uses ENV\n", buf.String())

	annotations := make(map[string]string)
	annotations["man-environment-section"] = "Override at cmd level"
	cmd.Annotations = annotations
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH ENVIRONMENT\n.PP\nOverride at cmd", buf.String())

	// FILES
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.NotRegexp(t, ".SH FILES\n", buf.String()) // No FILES section if not in opts

	opts = CobraManOptions{Files: "This uses files"}
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH FILES\n.PP\nThis uses files\n", buf.String())

	annotations = make(map[string]string)
	annotations["man-files-section"] = "Override at cmd level"
	cmd.Annotations = annotations
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH FILES\n.PP\nOverride at cmd", buf.String())

	// BUGS
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.NotRegexp(t, ".SH BUGS\n", buf.String()) // No BUGS section if not in opts

	opts = CobraManOptions{Bugs: "This has bugs"}
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH BUGS\n.PP\nThis has bugs\n", buf.String())

	annotations = make(map[string]string)
	annotations["man-bugs-section"] = "Override at cmd level"
	cmd.Annotations = annotations
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH BUGS\n.PP\nOverride at cmd", buf.String())

	// EXAMPLES
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.NotRegexp(t, ".SH EXAMPLES\n", buf.String()) // No EXAMPLES section if not in opts

	cmd.Example = "Here is example"
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH EXAMPLES\n.PP\nHere is example\n", buf.String())

	annotations = make(map[string]string)
	annotations["man-examples-section"] = "Override at cmd level"
	cmd.Annotations = annotations
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH EXAMPLES\n.PP\nOverride at cmd", buf.String())

	// AUTHOR
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH AUTHOR\n.PP\n.SM Page auto-generated by rayjohnson/cobra-man", buf.String()) // Always have AUTHOR SECTION

	opts = CobraManOptions{Author: "Written by Ray Johnson"}
	buf.Reset()
	assert.NoError(t, GenerateOnePage(cmd, &opts, "troff", buf))
	assert.Regexp(t, ".SH AUTHOR\nWritten by Ray Johnson\n.PP\n.SM Page auto-generated", buf.String()) // No OPTIONS section if not in opts
}
