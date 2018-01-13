package man

import (
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

var multiNewlineRegex *regexp.Regexp

func simpleToMdoc(str string) string {
	// Guessing this is already troff - so let it pass through
	if len(str) > 1 && str[0] == '.' {
		return str
	}

	// TODO: this could certainly be more sophisticated.  Pull requests welcome!
	// Right now it is good enough for the most simple cases.
	if multiNewlineRegex == nil {
		multiNewlineRegex = regexp.MustCompile(`\n+\n`)
	}

	return backslashify(multiNewlineRegex.ReplaceAllString(str, "\n.Pp\n"))
}

func simpleToTroff(str string) string {
	// Guessing this is already troff - so let it pass through
	if len(str) > 1 && str[0] == '.' {
		return str
	}

	// TODO: this could certainly be more sophisticated.  Pull requests welcome!
	// Right now it is good enough for the most simple cases.
	if multiNewlineRegex == nil {
		multiNewlineRegex = regexp.MustCompile(`\n+\n`)
	}

	return backslashify(multiNewlineRegex.ReplaceAllString(str, "\n.PP\n"))
}

var backslashReplacer *strings.Replacer

func backslashify(str string) string {
	if backslashReplacer == nil {
		backslashReplacer = strings.NewReplacer("-", "\\-", "_", "\\_", "&", "\\&", "\\", "\\\\", "~", "\\~")
	}
	return backslashReplacer.Replace(str)
}

func dashify(str string) string {
	return strings.Replace(str, " ", "-", -1)
}
