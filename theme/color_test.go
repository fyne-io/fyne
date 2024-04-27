package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func Test_ButtonColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ButtonColor()
	assert.Equal(t, DarkTheme().Color(ColorNameButton, VariantDark), c, "wrong button color")
}

func Test_TextColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ForegroundColor()
	assert.Equal(t, DarkTheme().Color(ColorNameForeground, VariantDark), c, "wrong text color")
}

func Test_DisabledTextColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := DisabledColor()
	assert.Equal(t, DarkTheme().Color(ColorNameDisabled, VariantDark), c, "wrong disabled text color")
}

func Test_PlaceHolderColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := PlaceHolderColor()
	assert.Equal(t, DarkTheme().Color(ColorNamePlaceHolder, VariantDark), c, "wrong placeholder color")
}

func Test_PrimaryColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := PrimaryColor()
	assert.Equal(t, DarkTheme().Color(ColorNamePrimary, VariantDark), c, "wrong primary color")
}

func Test_HoverColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := HoverColor()
	assert.Equal(t, DarkTheme().Color(ColorNameHover, VariantDark), c, "wrong hover color")
}

func Test_FocusColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := FocusColor()
	assert.Equal(t, DarkTheme().Color(ColorNameFocus, VariantDark), c, "wrong focus color")
}

func Test_ScrollBarColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ScrollBarColor()
	assert.Equal(t, DarkTheme().Color(ColorNameScrollBar, VariantDark), c, "wrong scrollbar color")
}
