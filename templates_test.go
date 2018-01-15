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
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRegisterTemplate(t *testing.T) {
	assert.Panics(t, func() { RegisterTemplate("bad", "-", "txt", "what {{ ") }, "The code did not panic")
	assert.NotPanics(t, func() { RegisterTemplate("good", "-", "txt", "Hello {{ \"world\" }} ") }, "The code should not panic")
}

func TestCustomerTemplate(t *testing.T) {
	buf := new(bytes.Buffer)

	RegisterTemplate("good", "-", "txt", "Hello {{ \"world\" }} ")
	cmd := &cobra.Command{Use: "foo"}
	opts := CobraManOptions{}
	assert.NoError(t, GenerateOnePage(cmd, &opts, "good", buf))
	assert.Regexp(t, "Hello world", buf.String()) // No OPTIONS section if not in opts

}
