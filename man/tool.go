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
	"path/filepath"

	"github.com/spf13/cobra"
)

// DocGenTool is an opaque type created by CreateDocGenCmdLineTool
type DocGenTool struct {
	installDirectory string
	docCmd           *cobra.Command
	appCmd           *cobra.Command
}

// CreateDocGenCmdLineTool creates a command line parser that can be used
// in a utility tool to generate documentation for a companion application.
func CreateDocGenCmdLineTool(appCmd *cobra.Command) *DocGenTool {
	dg := &DocGenTool{
		appCmd: appCmd,
	}

	dg.docCmd = &cobra.Command{
		Use:   "doc",
		Args:  cobra.NoArgs,
		Short: "Generate documentation, etc.",
	}
	dg.docCmd.PersistentFlags().StringVar(&dg.installDirectory, "directory", ".", "Directory to install generated files")

	return dg
}

// AddBashCompletionGenerator will create a subcommand for the utility tool
// that will generate a Bash Completion file for the companion app.  It will
// support a --directory flag and use the fileName passed into this function.
func (dg *DocGenTool) AddBashCompletionGenerator(fileName string) *DocGenTool {

	completeCmd := &cobra.Command{
		Use:   "generate-auto-complete",
		Args:  cobra.NoArgs,
		Short: "Generate bash auto complete script",
		RunE: func(myCmd *cobra.Command, args []string) error {
			path := filepath.Join(dg.installDirectory, fileName)
			return dg.appCmd.GenBashCompletionFile(path)
		},
	}

	dg.docCmd.AddCommand(completeCmd)

	return dg
}

// AddDocGenerator will create a subcommand for the utility tool that will
// generate documentation with the passed in CobraManOptions and templateName.
// It supports a --directory flag for where to place the generated files.  The
// subcommand will be named generate-<templateName> where templateName is the
// same as the template used to generate the documentation.
func (dg *DocGenTool) AddDocGenerator(opts *CobraManOptions, templateName string) *DocGenTool {
	// Make sure template exists or we will later get runtime panic
	_, ok := templateMap[templateName]
	if !ok {
		panic("the given template has not been registered: " + templateName)
	}

	genCmd := &cobra.Command{
		Use:   "generate-" + templateName,
		Args:  cobra.NoArgs,
		Short: "Generate docs with the " + templateName + " template",
		RunE: func(myCmd *cobra.Command, args []string) error {
			return GenerateDocs(dg.appCmd, opts, dg.installDirectory, templateName)
		},
	}

	dg.docCmd.AddCommand(genCmd)

	return dg
}

// Execute will parse args and execute the command line
func (dg *DocGenTool) Execute() error {
	return dg.docCmd.Execute()
}
