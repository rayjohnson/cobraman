# Generating Documentation For Your cobra.Command
![Image of Cobra Man](https://raw.githubusercontent.com/rayjohnson/cobraman/master/cobra-man.jpeg)

This is a replacement for the man generator used by spf13/cobra.  The code in spf13/cobra/doc has different generators that hard-code what gets output for man pages, markdown, etc.  It
also calls a lot of other 3rd party libraries.  This package uses the Go template facility
to generate documentation.  It is much more powerful and flexible thus giving a lot more
control over the format and style of documentation you would like to generate.

Here is a simple example to get you started:

```go
package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/rjohnson/cobraman"
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
		Bugs:        `Bugs related to cobra-man can be filed at https://github.com/rjohnson/cobraman`,
	}
	err := man.GenerateManPages(cmd.Root(), manOpts)
	if err != nil {
		log.Fatal(err)
	}
}
```

That will get you a man page `/tmp/dofoo.1`

GoDoc has the full API documentation [here](https://godoc.org/github.com/rayjohnson/cobraman).  Be sure to checkout the documentation for CobraManOptions as it provides many options to control the output.

There is also an example directory with a simple dummy application that shows some of the features of this package.  See the [README](example/README.md).

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

Cobra Man uses Go templates to generate the documentation.  You can replace the template used by setting the **TemplateName** variable in CobraManOptions.  A couple of templates are defined that can be used out of the box.  They include:

* "troff" - which generates a man page with basic troff macros
* "mdoc" - which generates a man page using the mdoc macro package
* "markdown" - which generates a page using Markdown

But, of course, you can provide your own template if you like for maximum power!

See [Writing your own template](WRITING_A_TEMPLATE.md) for more information.


