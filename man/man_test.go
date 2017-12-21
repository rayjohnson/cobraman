package man

import (
	"os"
	"testing"

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
