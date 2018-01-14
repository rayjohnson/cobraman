# Generating Man Pages For Your cobra.Command
![Image of Cobra Man](https://raw.githubusercontent.com/rayjohnson/cobra-man/master/cobra-man.jpeg)

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
	manOpts := &man.CobraManOptions{
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
```

## Annotations

This library uses the Annotations fields cobra.Cmd and pFlag to give some hints for the
generation of the documentation.

The following annotations on the cobra.Command object provides a way to provide content
for additional sections in the man page.  The first three override the global Options in 
case you want some of these sections only on some command man pages.
* man-files-section
* man-bugs-section
* man-environment-section
* man-examples-section

The **man-examples-section** is a way to override the content of the cmd.Examples field.
This is paticularly useful if you want to provide raw Troff code to make it look a bit 
better.

Here is an example of how you can set the annotations on the command:
```go
	annotations := make(map[string]string)
	annotations["man-files-section"] = "We use lots of files!"
	cmd.Annotations = annotations
```

In addition, there is an annotation you can put on individual flags:
* man-arg-hints

This provides a way to give a short description to the value expected by an flag.  This
is used by the built-in template in the OPTIONS section.  For example, setting the
annotation like this:
```go
	annotation := []string{"path"}
	flags.SetAnnotation("file", "man-arg-hints", annotation)
```

Will generate a option description like this:
```
-f, --file = <path>
```

## Templates

Cobra Man uses Go templates to generate the documentation it generates.  You can replace the
template used by setting the **UseTemplate** variable in man options.  A couple of templates
are defined that can be used out of the box.  They include:

* man.TroffManTemplate - which generates a man page with basic troff macros
* man.MdocManTemplate - which generates a man page using the mdoc macro package
* man.MarkdownTemplate - which generates a page using Markdown

But, of course, you can provide your own template if you like for maximum power!


