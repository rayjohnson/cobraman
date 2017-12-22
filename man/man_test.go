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

	opts := GenerateManOptions{}
	cmd := &cobra.Command{}
	err = GenerateManPages(cmd, &opts)
	assert.Equal(t, "you need a command name to have a man page", err.Error())

	cmd = &cobra.Command{Use: "foo"}
	assert.Nil(t, GenerateManPages(cmd, &opts))
	checkForFile(t, "foo.1")

	opts = GenerateManOptions{Section: "8"}
	assert.Nil(t, GenerateManPages(cmd, &opts))
	checkForFile(t, "foo.8")

	cmd = &cobra.Command{Use: "foo"}
	cmd2 := &cobra.Command{Use: "bar", Run: func(cmd *cobra.Command, args []string) {}}
	cmd.AddCommand(cmd2)
	opts = GenerateManOptions{}
	assert.Nil(t, GenerateManPages(cmd, &opts))
	checkForFile(t, "foo.1")
	checkForFile(t, "foo-bar.1")

	cmd = &cobra.Command{Use: "foo"}
	cmd2 = &cobra.Command{Use: "bar", Run: func(cmd *cobra.Command, args []string) {}}
	cmd.AddCommand(cmd2)
	opts = GenerateManOptions{CommandSeparator: "_"}
	assert.Nil(t, GenerateManPages(cmd, &opts))
	checkForFile(t, "foo.1")
	checkForFile(t, "foo_bar.1")

}

func TestGenerateManPageRequired(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{Use: "foo"}
	opts := GenerateManOptions{}

	// Test header options
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".TH \"FOO\" \"1\" \".*\" \"\" \"\"", buf.String())

	buf.Reset()
	opts = GenerateManOptions{LeftFooter: "kitty kat", CenterHeader: "Hello", CenterFooter: "meow", ProgramName: "Bobby", Section: "3"}
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".TH \"Bobby\" \"3\" \"meow\" \"kitty kat\" \"Hello\"", buf.String())

	buf.Reset()
	date, _ := time.Parse(time.RFC3339, "1968-06-21T15:04:05Z")
	opts = GenerateManOptions{Date: &date}
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".TH \"FOO\" \"1\" \"Jun 1968\" \"\" \"\"", buf.String())

	// Test name
	cmd = &cobra.Command{Use: "bar"}
	opts = GenerateManOptions{}
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH NAME\nbar\n", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "bar", Short: "going to"}
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH NAME\nbar .. going to", buf.String())

	// Test Synopsis
	assert.Regexp(t, ".SH SYNOPSIS\n.sp\n.+bar", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "foo"}
	cmd2 := &cobra.Command{Use: "cat", Run: func(cmd *cobra.Command, args []string) {}}
	cmd3 := &cobra.Command{Use: "dog", Run: func(cmd *cobra.Command, args []string) {}}
	cmd.AddCommand(cmd2, cmd3)
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH SYNOPSIS\n.sp\n.+foo cat.+flags.+\n.br\n.+foo dog", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "foo"}
	cmd.Flags().String("thing", "", "string with no default")
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, "SH SYNOPSIS\n.sp\n.+foo.+\\\\-\\\\-thing.+<args>]", buf.String())

	// Test DESCRIPTION
	buf.Reset()
	cmd = &cobra.Command{Use: "bar", Short: "a short one"}
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, "SH DESCRIPTION\n.PP\na short one", buf.String())

	cmd.Long = `Long desc

This is long & stuff.`
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH DESCRIPTION\n.PP\nLong desc\n.PP\nThis is long \\\\& stuff.", buf.String())

}

func TestGenerateManPageOptions(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{Use: "foo"}
	opts := GenerateManOptions{}

	cmd = &cobra.Command{Use: "foo"}
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.NotRegexp(t, ".SH OPTIONS\n", buf.String()) // No OPTIONS section if no flags

	cmd.Flags().String("flag", "", "string with no default")
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH OPTIONS\n.TP\n.+flag.+\nstring with no default", buf.String()) // No OPTIONS section if no flags

	cmd.Flags().String("hello", "world", "default is world")
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".TP\n.+flag.+\nstring with no default", buf.String()) // No OPTIONS section if no flags

	// TODO: I's like to revisit the format of OPTIONs section
}

func TestGenerateManPageAltSections(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{Use: "foo"}
	opts := GenerateManOptions{}

	// ENVIRONMENT
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.NotRegexp(t, ".SH ENVIRONMENT\n", buf.String()) // No OPTIONS section if not in opts

	opts = GenerateManOptions{Environment: "This uses ENV"}
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH ENVIRONMENT\n.PP\nThis uses ENV\n", buf.String()) // No OPTIONS section if not in opts

	annotations := make(map[string]string)
	annotations["man-environment-section"] = "Override at cmd level"
	cmd.Annotations = annotations
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH ENVIRONMENT\n.PP\nOverride at cmd", buf.String()) // No OPTIONS section if not in opts

	// FILES
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.NotRegexp(t, ".SH FILES\n", buf.String()) // No OPTIONS section if not in opts

	opts = GenerateManOptions{Files: "This uses files"}
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH FILES\n.PP\nThis uses files\n", buf.String()) // No OPTIONS section if not in opts

	annotations = make(map[string]string)
	annotations["man-files-section"] = "Override at cmd level"
	cmd.Annotations = annotations
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH FILES\n.PP\nOverride at cmd", buf.String()) // No OPTIONS section if not in opts

	// BUGS
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.NotRegexp(t, ".SH BUGS\n", buf.String()) // No OPTIONS section if not in opts

	opts = GenerateManOptions{Bugs: "This has bugs"}
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH BUGS\n.PP\nThis has bugs\n", buf.String()) // No OPTIONS section if not in opts

	annotations = make(map[string]string)
	annotations["man-bugs-section"] = "Override at cmd level"
	cmd.Annotations = annotations
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH BUGS\n.PP\nOverride at cmd", buf.String()) // No OPTIONS section if not in opts

	// EXAMPLES
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.NotRegexp(t, ".SH EXAMPLES\n", buf.String()) // No OPTIONS section if not in opts

	cmd.Example = "Here is example"
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH EXAMPLES\n.PP\nHere is example\n", buf.String()) // No OPTIONS section if not in opts

	annotations = make(map[string]string)
	annotations["man-examples-section"] = "Override at cmd level"
	cmd.Annotations = annotations
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH EXAMPLES\n.PP\nOverride at cmd", buf.String()) // No OPTIONS section if not in opts

	// AUTHOR
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.NotRegexp(t, ".SH AUTHOR\n", buf.String()) // No OPTIONS section if not in opts

	opts = GenerateManOptions{Author: "Written by Ray Johnson"}
	buf.Reset()
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, ".SH AUTHOR\n.PP\nWritten by Ray Johnson\n.PP\n.SM Page auto-generated", buf.String()) // No OPTIONS section if not in opts
}

func TestGenerateManPageTemplate(t *testing.T) {
	buf := new(bytes.Buffer)

	// bad user template
	cmd := &cobra.Command{Use: "foo"}
	opts := GenerateManOptions{UseTemplate: "what {{ "}
	assert.Error(t, generateManPage(cmd, &opts, buf))

	buf.Reset()
	opts = GenerateManOptions{UseTemplate: "Hello {{ \"world\" }} "}
	assert.NoError(t, generateManPage(cmd, &opts, buf))
	assert.Regexp(t, "Hello world", buf.String()) // No OPTIONS section if not in opts

}
