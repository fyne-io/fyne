package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"fyne.io/fyne/v2"
)

func Test_NewTheme(t *testing.T) {
	suite.Run(t, &configurableThemeTestSuite{constructor: func() *configurableTheme {
		return NewTheme().(*configurableTheme)
	}})
}

func Test_Theme(t *testing.T) {
	suite.Run(t, &configurableThemeTestSuite{constructor: func() *configurableTheme {
		return Theme().(*configurableTheme)
	}})
}

type configurableThemeTestSuite struct {
	suite.Suite
	constructor func() *configurableTheme
}

func (s *configurableThemeTestSuite) TestUniqueColorValues() {
	t := s.T()
	th := s.constructor()
	seen := map[string]fyne.ThemeColorName{}
	for cn, c := range th.colors {
		r, g, b, a := c.RGBA()
		key := fmt.Sprintf("%d %d %d %d", r, g, b, a)
		assert.True(t, seen[key] == "", "color value %#v for color %s already used for color %s in theme %s", c, cn, seen[key], th.name)
		seen[key] = cn
	}
}
