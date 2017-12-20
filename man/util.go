package man

import (
	"strings"

	"github.com/spf13/cobra"
)

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }


var sanatizeReplacer *strings.Replacer
func simpleToTroff(str string) string {
	// Guessing this is already troff - so let it pass through
	if str[0] == '.' {
		return str
	}

	// TODO: this could certainly be more sophisticated.  Pull requests welcome!
	// Right now it is good enough for the most simple cases.
	if sanatizeReplacer == nil {
		sanatizeReplacer = strings.NewReplacer("\n\n", "\n.PP\n")
	}
	return backslashify(sanatizeReplacer.Replace(str))
}

var backslashReplacer *strings.Replacer
func backslashify(str string) string {
	if backslashReplacer == nil {
		backslashReplacer = strings.NewReplacer("-", "\\-", "_", "\\_", "&", "\\&", "\\", "\\\\", "~", "\\~")
	}
	return backslashReplacer.Replace(str)
}

