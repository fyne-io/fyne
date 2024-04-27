package theme_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultTheme_ShadowColor(t *testing.T) {
	t.Run("light", func(t *testing.T) {
		_, _, _, a := theme.DefaultTheme().Color(theme.ColorNameShadow, theme.VariantLight).RGBA()
		assert.NotEqual(t, 255, a, "should not be transparent")
	})

	t.Run("dark", func(t *testing.T) {
		_, _, _, a := theme.DefaultTheme().Color(theme.ColorNameShadow, theme.VariantDark).RGBA()
		assert.NotEqual(t, 255, a, "should not be transparent")
	})
}

func TestEmptyTheme(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(&emptyTheme{})
	assert.NotNil(t, theme.ForegroundColor())
	assert.NotNil(t, theme.TextFont())
	assert.NotNil(t, theme.HelpIcon())
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
}

func TestThemeChange(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	bg := theme.BackgroundColor()

	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	assert.NotEqual(t, bg, theme.BackgroundColor())
}

func TestTheme_Bootstrapping(t *testing.T) {
	current := fyne.CurrentApp().Settings().Theme()
	fyne.CurrentApp().Settings().SetTheme(nil)

	// this should not crash
	theme.BackgroundColor()

	fyne.CurrentApp().Settings().SetTheme(current)
}

func TestTheme_Dark_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	bg := theme.BackgroundColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNameBackground, theme.VariantDark), bg, "wrong dark theme background color")
}

func TestTheme_Light_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	bg := theme.BackgroundColor()
	assert.Equal(t, theme.LightTheme().Color(theme.ColorNameBackground, theme.VariantLight), bg, "wrong light theme background color")
}

type emptyTheme struct {
}

func (e *emptyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return nil
}

func (e *emptyTheme) Font(s fyne.TextStyle) fyne.Resource {
	return nil
}

func (e *emptyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return nil
}

func (e *emptyTheme) Size(n fyne.ThemeSizeName) float32 {
	return 0
}
