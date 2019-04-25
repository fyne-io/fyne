package theme

import (
	"testing"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

// TODO figure how to test the default theme if we are messing with it...
//func TestThemeDefault(t *testing.T) {
//	loaded := fyne.GlobalSettings().Theme()
//	assert.NotEqual(t, lightBackground, loaded.BackgroundColor())
//}

func TestThemeChange_Basic(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	bg := BackgroundColor()

	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	assert.NotEqual(t, bg, BackgroundColor())
}

func TestTheme_Dark_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	bg := BackgroundColor()
	hex := ColorToHexString(bg)
	assert.Equal(t, "#424242", hex, "wrong dark theme background color")
}

func TestTheme_Light_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	bg := BackgroundColor()
	hex := ColorToHexString(bg)
	assert.Equal(t, "#f5f5f5", hex, "wrong light theme background color")
}

func TestTheme_SolarizedLight_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(SolarizedLight())
	bg := BackgroundColor()
	hex := ColorToHexString(bg)
	assert.Equal(t, "#fdf6e3", hex, "wrong solarized dark background color")
}

func TestTheme_SolarizedDark_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(SolarizedDark())
	bg := BackgroundColor()
	hex := ColorToHexString(bg)
	assert.Equal(t, "#002b36", hex, "wrong solarized light background color")
}

func Test_ButtonColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ButtonColor()
	hex := ColorToHexString(c)
	assert.Equal(t, "#212121", hex, "wrong button color")
}

func Test_HyperlinkColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := HyperlinkColor()
	hex := ColorToHexString(c)
	assert.Equal(t, "#9999ff", hex, "wrong hyperlink color")
}

func Test_TextColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := TextColor()
	hex := ColorToHexString(c)
	assert.Equal(t, "#ffffff", hex, "wrong hyperlink color")
}

func Test_PlaceHolderColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := PlaceHolderColor()
	hex := ColorToHexString(c)
	assert.Equal(t, "#b2b2b2", hex, "wrong placeholder color")
}

func Test_PrimaryColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := PrimaryColor()
	hex := ColorToHexString(c)
	assert.Equal(t, "#1a237e", hex, "wrong primary color")
}

func Test_HoverColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := HoverColor()
	hex := ColorToHexString(c)
	assert.Equal(t, "#313131", hex, "wrong hover color")
}

func Test_FocusColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := FocusColor()
	hex := ColorToHexString(c)
	assert.Equal(t, ColorToHexString(PrimaryColor()), hex, "wrong focus color")
}

func Test_ScrollBarColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ScrollBarColor()
	hex := ColorToHexString(c)
	assert.Equal(t, "#000000", hex, "wrong scrollbar color")
}

func Test_IconColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := IconColor()
	hex := ColorToHexString(c)
	assert.Equal(t, "#ffffff", hex, "wrong icon color")
}

func Test_TextSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, 14, TextSize(), "wrong text size")
}

func Test_TextFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Regular.ttf"
	result := TextFont().Name()
	assert.Equal(t, expect, result, "wrong regular text font")
}

func Test_TextBoldFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Bold.ttf"
	result := TextBoldFont().Name()
	assert.Equal(t, expect, result, "wrong bold text font")
}

func Test_TextItalicFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Italic.ttf"
	result := TextItalicFont().Name()
	assert.Equal(t, expect, result, "wrong italic text font")
}

func Test_TextBoldItalicFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-BoldItalic.ttf"
	result := TextBoldItalicFont().Name()
	assert.Equal(t, expect, result, "wrong bold italic text font")
}

func Test_TextMonospaceFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoMono-Regular.ttf"
	result := TextMonospaceFont().Name()
	assert.Equal(t, expect, result, "wrong monospace font")
}

func Test_Padding(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, 4, Padding(), "wrong padding")
}

func Test_IconInlineSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, 20, IconInlineSize(), "wrong inline icon size")
}

func Test_ScrollBarSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, 16, ScrollBarSize(), "wrong inline icon size")
}

func Test_DefaultTextFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Regular.ttf"
	result := DefaultTextFont().Name()
	assert.Equal(t, expect, result, "wrong default text font")
}

func Test_DefaultTextBoldFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Bold.ttf"
	result := DefaultTextBoldFont().Name()
	assert.Equal(t, expect, result, "wrong default text bold font")
}

func Test_DefaultTextBoldItalicFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-BoldItalic.ttf"
	result := DefaultTextBoldItalicFont().Name()
	assert.Equal(t, expect, result, "wrong default text bold italic font")
}

func Test_DefaultTextItalicFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoSans-Italic.ttf"
	result := DefaultTextItalicFont().Name()
	assert.Equal(t, expect, result, "wrong default text italic font")
}

func Test_DefaultTextMonospaceFont(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	expect := "NotoMono-Regular.ttf"
	result := DefaultTextMonospaceFont().Name()
	assert.Equal(t, expect, result, "wrong default monospace font")
}
