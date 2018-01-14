package man

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackslashify(t *testing.T) {
	cases := [][]string{
		{`foo-bar`, `foo\-bar`},
		{`foo&bar`, `foo\&bar`},
		{`foo_bar`, `foo\_bar`},
		{`foo\bar`, `foo\\bar`},
		{`foo~bar`, `foo\~bar`},
		{`-_&\~`, `\-\_\&\\\~`},
	}

	for i := 0; i < len(cases); i++ {
		str := backslashify(cases[i][0])
		expected := cases[i][1]
		assert.Equal(t, expected, str)
	}
}

func TestDashify(t *testing.T) {
	cases := [][]string{
		{`foo bar`, `foo-bar`},
		{`foo bar cat`, `foo-bar-cat`},
		{` foo bar `, `-foo-bar-`},
	}

	for i := 0; i < len(cases); i++ {
		str := dashify(cases[i][0])
		expected := cases[i][1]
		assert.Equal(t, expected, str)
	}
}

func TestUnderscoreify(t *testing.T) {
	cases := [][]string{
		{`foo bar`, `foo_bar`},
		{`foo bar cat`, `foo_bar_cat`},
		{` foo bar `, `_foo_bar_`},
	}

	for i := 0; i < len(cases); i++ {
		str := underscoreify(cases[i][0])
		expected := cases[i][1]
		assert.Equal(t, expected, str)
	}
}

func TestSimpleToTroff(t *testing.T) {
	cases := [][]string{
		{"Some test\none a line", "Some test\none a line"},
		{"Some test\n\nwith empty line", "Some test\n.PP\nwith empty line"},
		{".ignore me\n\none a line", ".ignore me\n\none a line"},
		{"Some test\n\n\nwith empty line", "Some test\n.PP\nwith empty line"},
		{"Some test\n\n\n\nwith empty line", "Some test\n.PP\nwith empty line"},
	}

	for i := 0; i < len(cases); i++ {
		str := simpleToTroff(cases[i][0])
		expected := cases[i][1]
		assert.Equal(t, expected, str)
	}
}

func TestSimpleToMdoc(t *testing.T) {
	cases := [][]string{
		{"Some test\none a line", "Some test\none a line"},
		{"Some test\n\nwith empty line", "Some test\n.Pp\nwith empty line"},
		{".ignore me\n\none a line", ".ignore me\n\none a line"},
		{"Some test\n\n\nwith empty line", "Some test\n.Pp\nwith empty line"},
		{"Some test\n\n\n\nwith empty line", "Some test\n.Pp\nwith empty line"},
	}

	for i := 0; i < len(cases); i++ {
		str := simpleToMdoc(cases[i][0])
		expected := cases[i][1]
		assert.Equal(t, expected, str)
	}
}
