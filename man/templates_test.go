package man

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
	opts := CobraManOptions{TemplateName: "good"}
	assert.NoError(t, GenerateOnePage(cmd, &opts, buf))
	assert.Regexp(t, "Hello world", buf.String()) // No OPTIONS section if not in opts

}
