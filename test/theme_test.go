package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Try to keep these in sync with the existing color names at theme/color.go.
var knownColorNames = []fyne.ThemeColorName{
	theme.ColorNameBackground,
	theme.ColorNameButton,
	theme.ColorNameDisabled,
	theme.ColorNameDisabledButton,
	theme.ColorNameError,
	theme.ColorNameFocus,
	theme.ColorNameForeground,
	theme.ColorNameForegroundOnError,
	theme.ColorNameForegroundOnPrimary,
	theme.ColorNameForegroundOnSuccess,
	theme.ColorNameForegroundOnWarning,
	theme.ColorNameHeaderBackground,
	theme.ColorNameHover,
	theme.ColorNameHyperlink,
	theme.ColorNameInputBackground,
	theme.ColorNameInputBorder,
	theme.ColorNameMenuBackground,
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

// Try to keep this in sync with the existing variants at theme/theme.go
var knownVariants = []fyne.ThemeVariant{
	theme.VariantDark,
	theme.VariantLight,
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
	for _, variant := range knownVariants {
		for _, cn := range knownColorNames {
			assert.NotNil(t, th.Color(cn, variant), "undefined color %s variant %d in theme %s", cn, variant, s.name)
		}
	}
}

func (s *configurableThemeTestSuite) TestUniqueColorValues() {
	t := s.T()
	th := s.constructor()
	seenByVariant := map[fyne.ThemeVariant]map[string]fyne.ThemeColorName{}
	for _, variant := range knownVariants {
		seen := seenByVariant[variant]
		if seen == nil {
			seen = map[string]fyne.ThemeColorName{}
			seenByVariant[variant] = seen
		}
		for _, cn := range knownColorNames {
			c := th.Color(cn, theme.VariantDark)
			r, g, b, a := c.RGBA()
			key := fmt.Sprintf("%d %d %d %d", r, g, b, a)
			assert.True(t, seen[key] == "", "color value %#v for color %s variant %d already used for color %s in theme %s", c, cn, variant, seen[key], s.name)
			seen[key] = cn
		}
	}
}
