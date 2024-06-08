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
	suite.Run(t, &configurableThemeTestSuite{
		constructor: NewTheme,
		name:        "Ugly Test Theme",
	})
}

func Test_Theme(t *testing.T) {
	suite.Run(t, &configurableThemeTestSuite{
		constructor: Theme,
		name:        "Default Test Theme",
	})
}

type configurableThemeTestSuite struct {
	suite.Suite
	constructor func() fyne.Theme
	name        string
}

func (s *configurableThemeTestSuite) TestAllColorsDefined() {
	t := s.T()
	th := s.constructor()
	for _, cn := range knownColorNames {
		assert.NotNil(t, th.Color(cn, theme.VariantDark), "undefined color %s in theme %s", cn, s.name)
	}
}

func (s *configurableThemeTestSuite) TestUniqueColorValues() {
	t := s.T()
	th := s.constructor()
	seen := map[string]fyne.ThemeColorName{}
	for _, cn := range knownColorNames {
		c := th.Color(cn, theme.VariantDark)
		r, g, b, a := c.RGBA()
		key := fmt.Sprintf("%d %d %d %d", r, g, b, a)
		assert.True(t, seen[key] == "", "color value %#v for color %s already used for color %s in theme %s", c, cn, seen[key], s.name)
		seen[key] = cn
	}
}
