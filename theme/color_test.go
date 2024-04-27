package theme_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func Test_ButtonColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	c := theme.ButtonColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNameButton, theme.VariantDark), c, "wrong button color")
}

func Test_DisabledTextColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	c := theme.DisabledColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNameDisabled, theme.VariantDark), c, "wrong disabled text color")
}

func Test_FocusColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	c := theme.FocusColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNameFocus, theme.VariantDark), c, "wrong focus color")
}

func Test_HoverColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	c := theme.HoverColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNameHover, theme.VariantDark), c, "wrong hover color")
}

func Test_PlaceHolderColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	c := theme.PlaceHolderColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNamePlaceHolder, theme.VariantDark), c, "wrong placeholder color")
}

func Test_PrimaryColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	c := theme.PrimaryColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNamePrimary, theme.VariantDark), c, "wrong primary color")
}

func Test_ScrollBarColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	c := theme.ScrollBarColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNameScrollBar, theme.VariantDark), c, "wrong scrollbar color")
}

func Test_TextColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
	c := theme.ForegroundColor()
	assert.Equal(t, theme.DarkTheme().Color(theme.ColorNameForeground, theme.VariantDark), c, "wrong text color")
}
