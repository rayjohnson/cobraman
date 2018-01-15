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
	"strings"
	"text/template"
)

type manTemplate struct {
	separator string
	extension string
	template  *template.Template
}

var templateMap = make(map[string]manTemplate)

// RegisterTemplate takes a template string creates a template for use with CobraMan.  It
// also takes a separator and file extension to be used when generating the file names for
// the generated files.
func RegisterTemplate(name string, separator string, extension string, templateString string) {
	// Build the template
	funcMap := template.FuncMap{
		"upper":         strings.ToUpper,
		"backslashify":  backslashify,
		"dashify":       dashify,
		"underscoreify": underscoreify,
		"simpleToTroff": simpleToTroff,
		"simpleToMdoc":  simpleToMdoc,
	}
	parsedTemplate := template.Must(template.New(name).Funcs(funcMap).Parse(templateString))

	t := manTemplate{
		separator: separator,
		extension: extension,
		template:  parsedTemplate,
	}
	templateMap[name] = t
}

func getTemplate(name string) (string, string, *template.Template) {
	t := templateMap[name]
	return t.separator, t.extension, t.template
}
