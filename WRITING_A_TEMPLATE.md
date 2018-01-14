# Writing your own template

Cobra-man uses Go's template system to generate documents.  If you are new to them you should
first check out the core go documentation about templates:
https://golang.org/pkg/text/template/

## Register a template

If you provide your own template you need to register it before calling CobraMan.
The **RegisterTemplate** function allows you to register the template string and
you also set some hints on how to generate the file names.

Here is an example call:
```
	RegisterTemplate("markdown", "_", "md", MarkdownTemplate)
```

The first argument is the name of the template.  You will pass that into cobraManOptions.TemplateName.  The second is a separator to use when generating a file name.  it used between the base name and the name of sub-commands.  The third argument is the extension to give the file name.  Finally, the last argument is a string that defines yiour template.

*Note: the extension argument can also take the special string "use_section" and the extension used will be the value set in cobraManOptions.Section.*

## Variables

The following variables are available for generating documentation.

* .Date - The date passed in to CobraManOptions (or Now() if it was not set)
* .Section - The section number set in CobraManOptions (defaults to "1")
* .CenterFooter - Text to put in the center part of a footer.
* .LeftFooter - Text to use in the left part of a footer
* .CenterHeader - Text to use in the center part of a header
* .UseLine - Cobra UseLine text
* .CommandPath - the space separated path for current command (e.g. "git commit")
* .ShortDescription - The ShortDescription set on a Cobra command
* .Description - The Description set on a Cobra command
* .NoArgs - A boolean set to true if the cobra.NoArgs is used for the command
* .AllFlags - an array of Flag objects defining all flags available for this command
* .InheritedFlags - an array of Flag objects defining flags inherited from parent commands
* .NonInheritedFlags - an array of Flag objects defining flags NOT inherited from parent commands
* .SeeAlsos - an array of the SeeAlso struct containing info about related commands
* .SubCommands - an array of child command names
* .Author - Text of Author variable set by CobraManOptions
* .Environment - Text of Environment variable set by CobraManOptions
* .Files - Text of Files variable set by CobraManOptions
* .Bugs - Text of Bugs variable set by CobraManOptions
* .Examples - Text of Example variable set on the cobra command

#### Flag struct (found in the various Flags arrays)

* .Shorthand - The "short" name for a flag (e.g. "h")
* .Name - The "long" name for a flag (e.g. "help")
* .Usage - The usage string set on the pflag.Flag
* .NoOptDefVal - (TODO - how best to describe)
* .DefValue - The default value set on the pflag
* .ArgHint - The value of an annotation on the pflag named "man-arg-hints"

#### SeeAlso struct (used in the SeeAlsos array)

* .CmdPath - the space separated path of a related path
* .Section - the man Section which will usually be the same as .Section above
* .IsParent - a boolean denoting this entry is the parent
* .IsChild - a boolean denoting this entry is a child sub-command
* .IsSibling - a boolean denoting this entry is a sibling sub-command

## Functions

The following functions are also available within templates to trasform the text
to meet your needs.

* upper - Transforms the text to upper case
* dashify - Converts any spaces in the text to dashes "-"
* underscoreify - Converts any spaces in the text to underscores "_"
* backslahify - Puts a backslash "\\" in front of any of the following characters:
	-, _, \&, \\, ~
* simpleToTroff - Inserts .PP where one or more blank newlines appear
* simpleToMdoc - Inserts .Pp where one or more blank newlines appear

## Example

Here is an abridged version of the MarkdownTemplate to see how to use the above 
variables and functions in a template.  

```
## {{.CommandPath | upper }}

{{ .ShortDescription }}

### Synopsis

{{ .Description }}

{{- if .AllFlags }}

### Options

The following options are supported:

{{ range .AllFlags -}}
* {{ if .Shorthand }}{{ print "-" .Shorthand }}, {{ end -}}{{ print "--" .Name }}
{{- if not .NoOptDefVal }}{{if .ArgHint }}=<{{ .ArgHint }}>{{ else }}=<{{ .DefValue }}>{{ end }}{{ end }}
{{- print " - " .Usage }}
{{ end }}
{{- end }}

### See Also

{{- range $index, $element := .SeeAlsos}}
* [{{ $element.CmdPath }}]({{ $element.CmdPath | underscoreify }}.md)
{{- end }}
{{- end }}
```
