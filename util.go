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
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

var multiNewlineRegex = regexp.MustCompile(`\n+\n`)

func simpleToMdoc(str string) string {
	// Guessing this is already troff - so let it pass through
	if len(str) > 1 && str[0] == '.' {
		return str
	}

	// TODO: this could certainly be more sophisticated.  Pull requests welcome!
	// Right now it is good enough for the most simple cases.
	return backslashify(multiNewlineRegex.ReplaceAllString(str, "\n.Pp\n"))
}

func simpleToTroff(str string) string {
	// Guessing this is already troff - so let it pass through
	if len(str) > 1 && str[0] == '.' {
		return str
	}

	// TODO: this could certainly be more sophisticated.  Pull requests welcome!
	// Right now it is good enough for the most simple cases.
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

func underscoreify(str string) string {
	return strings.Replace(str, " ", "_", -1)
}
