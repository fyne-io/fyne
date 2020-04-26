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

func TestBuiltinTheme_ShadowColor(t *testing.T) {
	shadow := ShadowColor()

	_, _, _, a := shadow.RGBA()
	assert.NotEqual(t, 255, a)
}

func TestTheme_Dark_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	bg := BackgroundColor()
	assert.Equal(t, DarkTheme().BackgroundColor(), bg, "wrong dark theme background color")
}

func TestTheme_Light_ReturnsCorrectBackground(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(LightTheme())
	bg := BackgroundColor()
	assert.Equal(t, LightTheme().BackgroundColor(), bg, "wrong light theme background color")
}

func Test_ButtonColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ButtonColor()
	assert.Equal(t, DarkTheme().ButtonColor(), c, "wrong button color")
}

func Test_HyperlinkColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := HyperlinkColor()
	assert.Equal(t, DarkTheme().HyperlinkColor(), c, "wrong hyperlink color")
}

func Test_TextColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := TextColor()
	assert.Equal(t, DarkTheme().TextColor(), c, "wrong text color")
}

func Test_DisabledTextColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := DisabledTextColor()
	assert.Equal(t, DarkTheme().DisabledTextColor(), c, "wrong disabled text color")
}

func Test_PlaceHolderColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := PlaceHolderColor()
	assert.Equal(t, DarkTheme().PlaceHolderColor(), c, "wrong placeholder color")
}

func Test_PrimaryColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := PrimaryColor()
	assert.Equal(t, DarkTheme().PrimaryColor(), c, "wrong primary color")
}

func Test_HoverColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := HoverColor()
	assert.Equal(t, DarkTheme().HoverColor(), c, "wrong hover color")
}

func Test_FocusColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := FocusColor()
	assert.Equal(t, PrimaryColor(), c, "wrong focus color")
}

func Test_ScrollBarColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := ScrollBarColor()
	assert.Equal(t, DarkTheme().ScrollBarColor(), c, "wrong scrollbar color")
}

func Test_IconColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := IconColor()
	assert.Equal(t, DarkTheme().IconColor(), c, "wrong icon color")
}

func Test_DisabledIconColor(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	c := DisabledIconColor()
	assert.Equal(t, DarkTheme().DisabledIconColor(), c, "wrong disabled icon color")
}

func Test_TextSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().TextSize(), TextSize(), "wrong text size")
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
	expect := "DejaVuSansMono-Powerline.ttf"
	result := TextMonospaceFont().Name()
	assert.Equal(t, expect, result, "wrong monospace font")
}

func Test_Padding(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().Padding(), Padding(), "wrong padding")
}

func Test_IconInlineSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().IconInlineSize(), IconInlineSize(), "wrong inline icon size")
}

func Test_ScrollBarSize(t *testing.T) {
	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
	assert.Equal(t, DarkTheme().ScrollBarSize(), ScrollBarSize(), "wrong inline icon size")
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
	expect := "DejaVuSansMono-Powerline.ttf"
	result := DefaultTextMonospaceFont().Name()
	assert.Equal(t, expect, result, "wrong default monospace font")
}
