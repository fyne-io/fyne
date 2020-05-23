package test

import (
	"image/color"

	"fyne.io/fyne/theme"

	"fyne.io/fyne"
)

// MonoTheme is a monochrome theme
// This should be used for running test cases so that any future changes
// to the colored themes do not require regenerating all the test images
func MonoTheme() fyne.Theme {
	theme := &monochrom{
		background:     color.NRGBA{0x44, 0x44, 0x44, 0xff},
		button:         color.NRGBA{0x33, 0x33, 0x33, 0xff},
		hover:          color.NRGBA{0x55, 0x55, 0x55, 0xff},
		disabledButton: color.NRGBA{0x22, 0x22, 0x22, 0xff},
		text:           color.NRGBA{0xff, 0xff, 0xff, 0xff},
		disabledText:   color.NRGBA{0x88, 0x88, 0x88, 0xff},
		icon:           color.NRGBA{0xee, 0xee, 0xee, 0xff},
		disabledIcon:   color.NRGBA{0xaa, 0xaa, 0xaa, 0xff},
		hyperlink:      color.NRGBA{0x99, 0x99, 0x99, 0xff},
		placeholder:    color.NRGBA{0xb2, 0xb2, 0xb2, 0xff},
		primary:        color.NRGBA{0x66, 0x66, 0x66, 0xff},
		scrollBar:      color.NRGBA{0x11, 0x11, 0x11, 0xff},
		shadow:         color.NRGBA{0x0, 0x0, 0x0, 0xff},
	}
	return theme
}

type monochrom struct {
	background color.Color

	button, primary, text, icon, hyperlink, placeholder, hover, scrollBar, shadow color.Color
	regular, bold, italic, bolditalic, monospace                                  fyne.Resource
	disabledButton, disabledIcon, disabledText                                    color.Color
}

func (t *monochrom) BackgroundColor() color.Color {
	return t.background
}

// ButtonColor returns the theme's standard button colour
func (t *monochrom) ButtonColor() color.Color {
	return t.button
}

// DisabledButtonColor returns the theme's disabled button colour
func (t *monochrom) DisabledButtonColor() color.Color {
	return t.disabledButton
}

// HyperlinkColor returns the theme's standard hyperlink colour
func (t *monochrom) HyperlinkColor() color.Color {
	return t.hyperlink
}

// TextColor returns the theme's standard text colour
func (t *monochrom) TextColor() color.Color {
	return t.text
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *monochrom) DisabledTextColor() color.Color {
	return t.disabledText
}

// IconColor returns the theme's standard text colour
func (t *monochrom) IconColor() color.Color {
	return t.icon
}

// DisabledIconColor returns the color for a disabledIcon UI element
func (t *monochrom) DisabledIconColor() color.Color {
	return t.disabledIcon
}

// PlaceHolderColor returns the theme's placeholder text colour
func (t *monochrom) PlaceHolderColor() color.Color {
	return t.placeholder
}

// PrimaryColor returns the colour used to highlight primary features
func (t *monochrom) PrimaryColor() color.Color {
	return t.primary
}

// HoverColor returns the colour used to highlight interactive elements currently under a cursor
func (t *monochrom) HoverColor() color.Color {
	return t.hover
}

// FocusColor returns the colour used to highlight a focused widget
func (t *monochrom) FocusColor() color.Color {
	return t.primary
}

// ScrollBarColor returns the color (and translucency) for a scrollBar
func (t *monochrom) ScrollBarColor() color.Color {
	return t.scrollBar
}

// ShadowColor returns the color (and translucency) for shadows used for indicating elevation
func (t *monochrom) ShadowColor() color.Color {
	return t.shadow
}

// TextSize returns the standard text size
func (t *monochrom) TextSize() int {
	return 14
}

// TextFont returns the standard text font
func (t *monochrom) TextFont() fyne.Resource {
	return theme.DefaultTextFont()
}

// TextBoldFont returns the standard bold font
func (t *monochrom) TextBoldFont() fyne.Resource {
	return theme.DefaultTextBoldFont()
}

// TextItalicFont returns the standard italic font
func (t *monochrom) TextItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

// TextBoldItalicFont returns the standard italic font
func (t *monochrom) TextBoldItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

// TextMonospaceFont returns the mono font
func (t *monochrom) TextMonospaceFont() fyne.Resource {
	return theme.DefaultTextMonospaceFont()
}

// Padding is the standard gap between elements and the border around interface
// elements
func (t *monochrom) Padding() int {
	return 4
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func (t *monochrom) IconInlineSize() int {
	return 20
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer
func (t *monochrom) ScrollBarSize() int {
	return 16
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer
func (t *monochrom) ScrollBarSmallSize() int {
	return 3
}
