// Package theme defines how a Fyne app should look when rendered
package theme

import "image/color"

import "github.com/fyne-io/fyne"

var loadedColors *themeColors
var loadedTheme string

type themeColors struct {
	Background color.RGBA

	Button, Text, Primary color.RGBA
}

// Basic definition of light theme colours
func loadLightColors() *themeColors {
	return &themeColors{
		Background: color.RGBA{0xff, 0xff, 0xff, 0xff},
		Button:     color.RGBA{0xee, 0xee, 0xee, 0xff},
		Text:       color.RGBA{0x0, 0x0, 0x0, 0xdd},
		Primary:    color.RGBA{0x9f, 0xa8, 0xda, 0xff},
	}
}

// Basic definition of dark theme colours
func loadDarkColors() *themeColors {
	return &themeColors{
		Background: color.RGBA{0x42, 0x42, 0x42, 0xff},
		Button:     color.RGBA{0x21, 0x21, 0x21, 0xff},
		Text:       color.RGBA{0xff, 0xff, 0xff, 0xff},
		Primary:    color.RGBA{0x1a, 0x23, 0x7e, 0xff},
	}
}

// Load the right theme colours based on environment / settings
func colors() *themeColors {
	if loadedTheme != fyne.GetSettings().Theme() {
		loadedColors = nil
	}

	c := loadedColors
	if loadedColors == nil {

		if fyne.GetSettings().Theme() == "light" {
			c = loadLightColors()
		} else {
			c = loadDarkColors()
		}
	}

	loadedColors = c
	return c
}

// BackgroundColor returns the theme's background colour
func BackgroundColor() color.RGBA {
	return colors().Background
}

// ButtonColor returns the theme's standard button colour
func ButtonColor() color.RGBA {
	return colors().Button
}

// TextColor returns the theme's standard text colour
func TextColor() color.RGBA {
	return colors().Text
}

// PrimaryColor returns the colour used to highlight primary features
func PrimaryColor() color.RGBA {
	return colors().Primary
}

// FocusColor returns the colour used to highlight a focussed widget
func FocusColor() color.RGBA {
	return colors().Primary
}

// TextSize returns the standard text size
func TextSize() int {
	return 14
}

// TextFont returns the font path for the regular font style
func TextFont() *fyne.Resource {
	return regular
}

// TextBoldFont retutns the font path for the bold font style
func TextBoldFont() *fyne.Resource {
	return bold
}

// TextItalicFont returns the font path for the italic font style
func TextItalicFont() *fyne.Resource {
	return italic
}

// TextBoldItalicFont returns the font path for the bold and italic font style
func TextBoldItalicFont() *fyne.Resource {
	return bolditalic
}

// Padding is the standard gap between elements and the border around interface
// elements
func Padding() int {
	return 4
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func IconInlineSize() int {
	return 16
}
