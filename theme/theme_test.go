package theme

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"

	"github.com/stretchr/testify/assert"
)

func TestBuiltinTheme_ShadowColor(t *testing.T) {
	shadow := ShadowColor()

	_, _, _, a := shadow.RGBA()
	assert.NotEqual(t, 255, a)
}

func TestEmptyTheme(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(&emptyTheme{})
	assert.NotNil(t, ForegroundColor())
	assert.NotNil(t, TextFont())
	assert.NotNil(t, HelpIcon())
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
}

func TestThemeChange(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	bg := BackgroundColor()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, bg, BackgroundColor())
}

func TestTheme_Bootstrapping(t *testing.T) {
	current := fyne.CurrentApp().Settings().Theme()
	fyne.CurrentApp().Settings().SetTheme(nil)

	// this should not crash
	BackgroundColor()

	fyne.CurrentApp().Settings().SetTheme(current)
}

func TestTheme_Dark_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	bg := BackgroundColor()
	assert.Equal(t, DarkTheme().Color(ColorNameBackground, VariantDark), bg, "wrong dark theme background color")
}

func TestTheme_Light_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	bg := BackgroundColor()
	assert.Equal(t, LightTheme().Color(ColorNameBackground, VariantLight), bg, "wrong light theme background color")
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
