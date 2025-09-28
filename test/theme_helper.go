//go:build !tamago && !noos

package test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

// AssertAllColorNamesDefined asserts that all known color names are defined for the given theme.
func AssertAllColorNamesDefined(t *testing.T, th fyne.Theme, themeName string) {
	oldApp := fyne.CurrentApp()
	defer fyne.SetCurrentApp(oldApp)

	for _, primaryName := range theme.PrimaryColorNames() {
		testApp := NewTempApp(t)
		testApp.Settings().(*testSettings).primaryColor = primaryName
		for variantName, variant := range KnownThemeVariants() {
			for _, cn := range knownColorNames {
				assert.NotNil(t, th.Color(cn, variant), "undefined color %s variant %s in theme %s", cn, variantName, themeName)
				// Transparent is used by the default theme as fallback for unknown color names.
				// Built-in color names should have well-defined non-transparent values.
				assert.NotEqual(t, color.Transparent, th.Color(cn, variant), "undefined color %s variant %s in theme %s", cn, variantName, themeName)
			}
		}
	}
}
