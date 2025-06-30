package test

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

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
	AssertAllColorNamesDefined(s.T(), s.constructor(), s.name)
}

func (s *configurableThemeTestSuite) TestUniqueColorValues() {
	t := s.T()
	th := s.constructor()
	seenByVariant := map[fyne.ThemeVariant]map[string]fyne.ThemeColorName{}
	for variantName, variant := range KnownThemeVariants() {
		seen := seenByVariant[variant]
		if seen == nil {
			seen = map[string]fyne.ThemeColorName{}
			seenByVariant[variant] = seen
		}
		for _, cn := range knownColorNames {
			c := th.Color(cn, variant)
			r, g, b, a := c.RGBA()
			key := fmt.Sprintf("%d %d %d %d", r, g, b, a)
			assert.True(t, seen[key] == "", "color value %#v for color %s variant %s already used for color %s in theme %s", c, cn, variantName, seen[key], s.name)
			seen[key] = cn
		}
	}
}
