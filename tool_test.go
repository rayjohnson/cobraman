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

package cobraman

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateDocGenCmdLineTool(t *testing.T) {
	appCmd := &cobra.Command{}
	dg := CreateDocGenCmdLineTool(appCmd)
	assert.NotNil(t, dg, "CreateDocGenCmdLineTool should not return nil")
}

func TestAddDocGenerator(t *testing.T) {
	appCmd := &cobra.Command{}
	dg := CreateDocGenCmdLineTool(appCmd)

	opts := &CobraManOptions{}

	// Template does not exist
	assert.Panics(t, func() { dg.AddDocGenerator(opts, "foo") })
	dg.AddDocGenerator(opts, "mdoc") // This one does exist

	// bad sub-command
	args := []string{"generate-bar"}
	dg.docCmd.SetArgs(args)
	buf := new(bytes.Buffer)
	dg.docCmd.SetOutput(buf)

	// No error is thrown instead usage string is shown
	assert.NoError(t, dg.Execute())
	assert.Regexp(t, "Available Commands.+\n.+generate-mdoc", buf)

	buf.Reset()
	args = []string{"generate-mdoc"}
	dg.docCmd.SetArgs(args)
	assert.Error(t, dg.Execute()) // Will generate an error because our command is dumb
}

func TestAddBashCompletionGenerator(t *testing.T) {
	appCmd := &cobra.Command{}
	dg := CreateDocGenCmdLineTool(appCmd)
	dg.AddBashCompletionGenerator("foo.txt")
}

func TestExecute(t *testing.T) {
	appCmd := &cobra.Command{Use: "foo"}
	cmd2 := &cobra.Command{Use: "child1", Run: func(cmd *cobra.Command, args []string) {}}
	cmd3 := &cobra.Command{Use: "child2", Run: func(cmd *cobra.Command, args []string) {}}
	appCmd.AddCommand(cmd2, cmd3)

	dg := CreateDocGenCmdLineTool(appCmd)
	opts := &CobraManOptions{}
	dg.AddDocGenerator(opts, "mdoc")
	dg.AddDocGenerator(opts, "markdown")
	dg.AddBashCompletionGenerator("foo.txt")

	args := []string{"generate-markdown"}
	dg.docCmd.SetArgs(args)
	assert.NoError(t, dg.Execute())
	checkForFile(t, "foo.md")
	checkForFile(t, "foo_child1.md")
	checkForFile(t, "foo_child2.md")

	args = []string{"generate-auto-complete"}
	dg.docCmd.SetArgs(args)
	assert.NoError(t, dg.Execute())
	checkForFile(t, "foo.txt")
}
