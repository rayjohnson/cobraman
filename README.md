# Generating Man Pages For Your cobra.Command

This is a replacement for the man generator used by spf13/cobra.  The one in spf13/cobra/doc first
generates markdown and then calls another package to convert mark-down to roff.  This one generates
the roff directly and a little more cleanly.  It also has a few more options.

Here is a simple example to get you started:

```go
package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/rjohnson/cobra-man/man"
)

func main() {
	cmd := &cobra.Command{
		Use:   "dofoo",
		Short: "my dofoo program",
	}
	manOpts := &man.GenerateManOptions{
		LeftFooter:  "Dofoo " + version,
		Author:      "Foo Bar <foo@bar.com>",
		Directory:   "/tmp",
		Bugs:        `Bugs related to cobra-man can be filed at https://github.com/rjohnson/cobra-man`,
	}
	err := man.GenerateManPages(cmd.Root(), manOpts)
	if err != nil {
		log.Fatal(err)
	}
}
```

That will get you a man page `/tmp/dofoo.1`

Here is the full set of options you may use:
```
	// ProgramName is used in the man page header across all pages
	// The default is to generate an all CAPS path like CMD-SUBCMD
	// for each page.  Because this would instead make them the same
	// for all pages it is probably best not to override.
	ProgramName string

	// What section to generate the pages 4 (1 is the default if not set)
	Section string

	// CenterFooter used across all pages (defaults to current month and year)
	// If you just want to set the date used in the center footer use Date
	CenterFooter string

	// If you just want to set the date used in the center footer use Date
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
	// it as an annotation: cmd.Annotations["man-files-section"]
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

	// CommandSperator defines what character to use to separate the
	// sub commands in the man page file name.  The '-' char is the default.
	CommandSeparator string

	// GenSeprateInheiratedFlags will generate a separate section for
	// inherited flags.  By default they will all be in the same OPZTIONS
	// section.
	GenSeprateInheritedFlags bool

	// UseTemplate allows you to override the default go template used to
	// generate the man pages with your own version.
	UseTemplate string
```