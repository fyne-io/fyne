package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Try to keep these in sync with the existing color names at theme/colors.go.
var knownColorNames = []fyne.ThemeColorName{
	theme.ColorNameBackground,
	theme.ColorNameButton,
	theme.ColorNameDisabledButton,
	theme.ColorNameDisabled,
	theme.ColorNameError,
	theme.ColorNameFocus,
	theme.ColorNameForeground,
	theme.ColorNameHeaderBackground,
	theme.ColorNameHover,
	theme.ColorNameHyperlink,
	theme.ColorNameInputBackground,
	theme.ColorNameInputBorder,
	theme.ColorNameMenuBackground,
	theme.ColorNameOnPrimary,
	theme.ColorNameOverlayBackground,
	theme.ColorNamePlaceHolder,
	theme.ColorNamePressed,
	theme.ColorNamePrimary,
	theme.ColorNameScrollBar,
	theme.ColorNameSelection,
	theme.ColorNameSeparator,
	theme.ColorNameShadow,
	theme.ColorNameSuccess,
	theme.ColorNameWarning,
}

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

func (s *configurableThemeTestSuite) TestAllColorsDefined() {
	t := s.T()
	th := s.constructor()
	for _, cn := range knownColorNames {
		assert.NotNil(t, th.Color(cn, theme.VariantDark), "undefined color %s in theme %s", cn, th.name)
	}
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
