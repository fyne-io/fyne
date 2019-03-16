// Package theme defines how a Fyne app should look when rendered
package theme // import "fyne.io/fyne/theme"

import (
	"image/color"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"fyne.io/fyne"
)

type builtinTheme struct {
	background color.Color

	button, text, hyperlink, placeholder, primary, scrollBar color.Color
	regular, bold, italic, bolditalic, monospace             fyne.Resource
}

var lightBackground = color.RGBA{0xf5, 0xf5, 0xf5, 0xff}

// LightTheme defines the built in light theme colours and sizes
func LightTheme() fyne.Theme {
	theme := &builtinTheme{
		background:  lightBackground,
		button:      color.RGBA{0xd9, 0xd9, 0xd9, 0xff},
		text:        color.RGBA{0x21, 0x21, 0x21, 0xff},
		hyperlink:   color.RGBA{0x0, 0x0, 0xd9, 0xff},
		placeholder: color.RGBA{0x88, 0x88, 0x88, 0xff},
		primary:     color.RGBA{0x9f, 0xa8, 0xda, 0xff},
		scrollBar:   color.RGBA{0x0, 0x0, 0x0, 0x99},
	}

	theme.initFonts()
	return theme
}

// DarkTheme defines the built in dark theme colours and sizes
func DarkTheme() fyne.Theme {
	theme := &builtinTheme{
		background:  color.RGBA{0x42, 0x42, 0x42, 0xff},
		button:      color.RGBA{0x21, 0x21, 0x21, 0xff},
		text:        color.RGBA{0xff, 0xff, 0xff, 0xff},
		hyperlink:   color.RGBA{0x99, 0x99, 0xff, 0xff},
		placeholder: color.RGBA{0xb2, 0xb2, 0xb2, 0xff},
		primary:     color.RGBA{0x1a, 0x23, 0x7e, 0xff},
		scrollBar:   color.RGBA{0x0, 0x0, 0x0, 0x99},
	}

	theme.initFonts()
	return theme
}

func (t *builtinTheme) BackgroundColor() color.Color {
	return t.background
}

// ButtonColor returns the theme's standard button colour
func (t *builtinTheme) ButtonColor() color.Color {
	return t.button
}

// HyperlinkColor returns the theme's standard hyperlink colour
func (t *builtinTheme) HyperlinkColor() color.Color {
	return t.hyperlink
}

// TextColor returns the theme's standard text colour
func (t *builtinTheme) TextColor() color.Color {
	return t.text
}

// PlaceHolderColor returns the theme's placeholder text colour
func (t *builtinTheme) PlaceHolderColor() color.Color {
	return t.placeholder
}

// PrimaryColor returns the colour used to highlight primary features
func (t *builtinTheme) PrimaryColor() color.Color {
	return t.primary
}

// FocusColor returns the colour used to highlight a focussed widget
func (t *builtinTheme) FocusColor() color.Color {
	return t.primary
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func (t *builtinTheme) ScrollBarColor() color.Color {
	return t.scrollBar
}

// TextSize returns the standard text size
func (t *builtinTheme) TextSize() int {
	return 14
}

func loadCustomFont(env, variant string, fallback fyne.Resource) fyne.Resource {
	variantPath := strings.Replace(env, "Regular", variant, 0)

	file, err := os.Open(variantPath)
	if err != nil {
		fyne.LogError("Error loading specified font", err)
		return fallback
	}
	ret, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		fyne.LogError("Error loading specified font", err2)
		return fallback
	}

	name := path.Base(variantPath)
	return &fyne.StaticResource{StaticName: name, StaticContent: ret}

}

func (t *builtinTheme) initFonts() {
	t.regular = regular
	t.bold = bold
	t.italic = italic
	t.bolditalic = bolditalic
	t.monospace = monospace

	font := os.Getenv("FYNE_FONT")
	if font != "" {
		t.regular = loadCustomFont(font, "Regular", regular)
		t.bold = loadCustomFont(font, "Bold", bold)
		t.italic = loadCustomFont(font, "Italic", italic)
		t.bolditalic = loadCustomFont(font, "BoldItalic", bolditalic)
	}
	font = os.Getenv("FYNE_FONT_MONOSPACE")
	if font != "" {
		t.monospace = loadCustomFont(font, "Regular", monospace)
	}
}

// TextFont returns the font path for the regular font style
func (t *builtinTheme) TextFont() fyne.Resource {
	return t.regular
}

// TextBoldFont retutns the font path for the bold font style
func (t *builtinTheme) TextBoldFont() fyne.Resource {
	return t.bold
}

// TextItalicFont returns the font path for the italic font style
func (t *builtinTheme) TextItalicFont() fyne.Resource {
	return t.italic
}

// TextBoldItalicFont returns the font path for the bold and italic font style
func (t *builtinTheme) TextBoldItalicFont() fyne.Resource {
	return t.bolditalic
}

// TextMonospaceFont retutns the font path for the monospace font face
func (t *builtinTheme) TextMonospaceFont() fyne.Resource {
	return t.monospace
}

// Padding is the standard gap between elements and the border around interface
// elements
func (t *builtinTheme) Padding() int {
	return 4
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func (t *builtinTheme) IconInlineSize() int {
	return 20
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer
func (t *builtinTheme) ScrollBarSize() int {
	return 16
}

func current() fyne.Theme {
	//	if fyne.CurrentApp().Theme() != nil
	return fyne.CurrentApp().Settings().Theme()
}

// BackgroundColor returns the theme's background colour
func BackgroundColor() color.Color {
	return current().BackgroundColor()
}

// ButtonColor returns the theme's standard button colour
func ButtonColor() color.Color {
	return current().ButtonColor()
}

// HyperlinkColor returns the theme's standard hyperlink colour
func HyperlinkColor() color.Color {
	return current().HyperlinkColor()
}

// TextColor returns the theme's standard text colour
func TextColor() color.Color {
	return current().TextColor()
}

// PlaceHolderColor returns the theme's standard text colour
func PlaceHolderColor() color.Color {
	return current().PlaceHolderColor()
}

// PrimaryColor returns the colour used to highlight primary features
func PrimaryColor() color.Color {
	return current().PrimaryColor()
}

// FocusColor returns the colour used to highlight a focussed widget
func FocusColor() color.Color {
	return current().FocusColor()
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func ScrollBarColor() color.Color {
	return current().ScrollBarColor()
}

// TextSize returns the standard text size
func TextSize() int {
	return current().TextSize()
}

// TextFont returns the font path for the regular font style
func TextFont() fyne.Resource {
	return current().TextFont()
}

// TextBoldFont retutns the font path for the bold font style
func TextBoldFont() fyne.Resource {
	return current().TextBoldFont()
}

// TextItalicFont returns the font path for the italic font style
func TextItalicFont() fyne.Resource {
	return current().TextItalicFont()
}

// TextBoldItalicFont returns the font path for the bold and italic font style
func TextBoldItalicFont() fyne.Resource {
	return current().TextBoldItalicFont()
}

// TextMonospaceFont retutns the font path for the monospace font face
func TextMonospaceFont() fyne.Resource {
	return current().TextMonospaceFont()
}

// Padding is the standard gap between elements and the border around interface
// elements
func Padding() int {
	return current().Padding()
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func IconInlineSize() int {
	return current().IconInlineSize()
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer
func ScrollBarSize() int {
	return current().ScrollBarSize()
}

// DefaultTextFont returns the font path for the built-in regular font style
func DefaultTextFont() fyne.Resource {
	return regular
}

// DefaultTextBoldFont retutns the font path for the built-in bold font style
func DefaultTextBoldFont() fyne.Resource {
	return bold
}

// DefaultTextItalicFont returns the font path for the built-in italic font style
func DefaultTextItalicFont() fyne.Resource {
	return italic
}

// DefaultTextBoldItalicFont returns the font path for the built-in bold and italic font style
func DefaultTextBoldItalicFont() fyne.Resource {
	return bolditalic
}

// DefaultTextMonospaceFont retutns the font path for the built-in monospace font face
func DefaultTextMonospaceFont() fyne.Resource {
	return monospace
}
