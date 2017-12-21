package man

import (
	"os"
	"bytes"
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

func TestGenerateManPage(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{Use: "foo"}
	opts := GenerateManOptions{}

	// Test header options
	generateManPage(cmd, &opts, buf)
	assert.Regexp(t, ".TH \"FOO\" \"1\" \".*\" \"\" \"\"", buf.String())

	buf.Reset()
	opts = GenerateManOptions{ LeftFooter: "kitty kat", CenterHeader: "Hello", CenterFooter: "meow", ProgramName: "Bobby", Section: "3" }
	generateManPage(cmd, &opts, buf)
	assert.Regexp(t, ".TH \"Bobby\" \"3\" \"meow\" \"kitty kat\" \"Hello\"", buf.String())

	buf.Reset()
	date, _ := time.Parse(time.RFC3339, "1968-06-21T15:04:05Z")
	opts = GenerateManOptions{ Date: &date }
	generateManPage(cmd, &opts, buf)
	assert.Regexp(t, ".TH \"FOO\" \"1\" \"Jun 1968\" \"\" \"\"", buf.String())

	// Test name
	cmd = &cobra.Command{Use: "bar"}
	opts = GenerateManOptions{}
	generateManPage(cmd, &opts, buf)
	assert.Regexp(t, ".SH NAME\nbar\n", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "bar", Short: "going to"}
	generateManPage(cmd, &opts, buf)
	assert.Regexp(t, ".SH NAME\nbar .. going to", buf.String())

	// Test Synopsis
	assert.Regexp(t, ".SH SYNOPSIS\n.sp\n.SY bar", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "foo"}
	cmd2 := &cobra.Command{Use: "cat", Run: func(cmd *cobra.Command, args []string) {}}
	cmd3 := &cobra.Command{Use: "dog", Run: func(cmd *cobra.Command, args []string) {}}
	cmd.AddCommand(cmd2, cmd3)
	generateManPage(cmd, &opts, buf)
	assert.Regexp(t, ".SH SYNOPSIS\n.sp\n.SY foo cat\n.RI . flags .\n.YS\n.SY foo dog", buf.String())

	buf.Reset()
	cmd = &cobra.Command{Use: "foo"}
	cmd.Flags().String("thing", "", "string with no default")
	generateManPage(cmd, &opts, buf)
	assert.Regexp(t, ".SH SYNOPSIS\n.sp\n.SY foo\n.OP thing", buf.String())
}
