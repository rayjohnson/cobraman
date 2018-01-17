package main

import (
	"github.com/rayjohnson/cobraman"
	"github.com/rayjohnson/cobraman/example/cmd"
	"os"
)

func main() {
	// Get the root cobra command for the zap application
	appCmds := cmd.GetRootCmd()

	docGenerator := cobraman.CreateDocGenCmdLineTool(appCmds)
	docGenerator.AddBashCompletionGenerator("zap.sh")

	manOpts := &cobraman.CobraManOptions{
		LeftFooter:   "Example",
		CenterHeader: "Example Manual",
		Author:       "Ray Johnson <ray.johnson@gmail.com>",
		Bugs:         `Bugs related to cobraman can be filed at https://github.com/rayjohnson/cobraman`,
	}
	docGenerator.AddDocGenerator(manOpts, "mdoc")
	docGenerator.AddDocGenerator(manOpts, "troff")
	docGenerator.AddDocGenerator(manOpts, "markdown")

	if err := docGenerator.Execute(); err != nil {
		os.Exit(1)
	}
}

