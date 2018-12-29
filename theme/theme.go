// Package theme defines how a Fyne app should look when rendered
package theme

import (
	"image/color"

	"github.com/fyne-io/fyne"
)

type builtinTheme struct {
	background color.Color

	button, text, placeholder, primary color.Color
}

var lightBackground = color.RGBA{0xf5, 0xf5, 0xf5, 0xff}

// LightTheme defines the built in light theme colours and sizes
func LightTheme() fyne.Theme {
	return &builtinTheme{
		background:  lightBackground,
		button:      color.RGBA{0xd9, 0xd9, 0xd9, 0xff},
		text:        color.RGBA{0x0, 0x0, 0x0, 0xdd},
		placeholder: color.RGBA{0x88, 0x88, 0x88, 0xff},
		primary:     color.RGBA{0x9f, 0xa8, 0xda, 0xff},
	}
}

// DarkTheme defines the built in dark theme colours and sizes
func DarkTheme() fyne.Theme {
	return &builtinTheme{
		background:  color.RGBA{0x42, 0x42, 0x42, 0xff},
		button:      color.RGBA{0x21, 0x21, 0x21, 0xff},
		text:        color.RGBA{0xff, 0xff, 0xff, 0xff},
		placeholder: color.RGBA{0x88, 0x88, 0x88, 0xff},
		primary:     color.RGBA{0x1a, 0x23, 0x7e, 0xff},
	}
}

func (t *builtinTheme) BackgroundColor() color.Color {
	return t.background
}

// ButtonColor returns the theme's standard button colour
func (t *builtinTheme) ButtonColor() color.Color {
	return t.button
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

// TextSize returns the standard text size
func (t *builtinTheme) TextSize() int {
	return 14
}

// TextFont returns the font path for the regular font style
func (t *builtinTheme) TextFont() fyne.Resource {
	return regular
}

// TextBoldFont retutns the font path for the bold font style
func (t *builtinTheme) TextBoldFont() fyne.Resource {
	return bold
}

// TextItalicFont returns the font path for the italic font style
func (t *builtinTheme) TextItalicFont() fyne.Resource {
	return italic
}

// TextBoldItalicFont returns the font path for the bold and italic font style
func (t *builtinTheme) TextBoldItalicFont() fyne.Resource {
	return bolditalic
}

// TextMonospaceFont retutns the font path for the monospace font face
func (t *builtinTheme) TextMonospaceFont() fyne.Resource {
	return monospace
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
