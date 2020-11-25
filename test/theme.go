package test

import (
	"image/color"

	"fyne.io/fyne"
)

var defaulTheme fyne.Theme

var _ fyne.Theme = (*configurableTheme)(nil)

type configurableTheme struct {
	background         color.Color
	bold               fyne.Resource
	boldItalic         fyne.Resource
	button             color.Color
	disabledButton     color.Color
	disabledIcon       color.Color
	disabledText       color.Color
	focus              color.Color
	hover              color.Color
	hyperlink          color.Color
	icon               color.Color
	iconInlineSize     int
	italic             fyne.Resource
	monospace          fyne.Resource
	padding            int
	placeholder        color.Color
	primary            color.Color
	regular            fyne.Resource
	scrollBar          color.Color
	scrollBarSize      int
	scrollBarSmallSize int
	shadow             color.Color
	text               color.Color
	textSize           int
}

// Theme returns a theme useful for image based tests.
func Theme() fyne.Theme {
	if defaulTheme == nil {
		defaulTheme = &configurableTheme{
			background:         color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff},
			bold:               bold,
			boldItalic:         bolditalic,
			button:             color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
			disabledButton:     color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
			disabledIcon:       color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}, // deprecated: bright red variant to make it visible
			disabledText:       color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
			focus:              color.NRGBA{R: 0x81, G: 0xc7, B: 0x84, A: 0xff},
			hover:              color.NRGBA{R: 0x88, G: 0xff, B: 0xff, A: 0x22},
			hyperlink:          color.NRGBA{R: 0xee, G: 0x00, B: 0x00, A: 0xff}, // deprecated: bright red variant to make it visible
			icon:               color.NRGBA{R: 0xdd, G: 0x00, B: 0x00, A: 0xff}, // deprecated: bright red variant to make it visible
			iconInlineSize:     20,
			italic:             italic,
			monospace:          monospace,
			padding:            4,
			placeholder:        color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff},
			primary:            color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
			regular:            regular,
			scrollBar:          color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa},
			scrollBarSize:      16,
			scrollBarSmallSize: 3,
			shadow:             color.NRGBA{A: 0x88},
			text:               color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			textSize:           14,
		}
	}
	return defaulTheme
}

func (t *configurableTheme) BackgroundColor() color.Color {
	return t.background
}

func (t *configurableTheme) ButtonColor() color.Color {
	return t.button
}

func (t *configurableTheme) DisabledButtonColor() color.Color {
	return t.disabledButton
}

func (t *configurableTheme) DisabledIconColor() color.Color {
	return t.disabledIcon
}

func (t *configurableTheme) DisabledTextColor() color.Color {
	return t.disabledText
}

func (t *configurableTheme) FocusColor() color.Color {
	return t.focus
}

func (t *configurableTheme) HoverColor() color.Color {
	return t.hover
}

func (t *configurableTheme) HyperlinkColor() color.Color {
	return t.hyperlink
}

func (t *configurableTheme) IconColor() color.Color {
	return t.icon
}

func (t *configurableTheme) IconInlineSize() int {
	return t.iconInlineSize
}

func (t *configurableTheme) Padding() int {
	return t.padding
}

func (t *configurableTheme) PlaceHolderColor() color.Color {
	return t.placeholder
}

func (t *configurableTheme) PrimaryColor() color.Color {
	return t.primary
}

func (t *configurableTheme) ScrollBarColor() color.Color {
	return t.scrollBar
}

func (t *configurableTheme) ScrollBarSize() int {
	return t.scrollBarSize
}

func (t *configurableTheme) ScrollBarSmallSize() int {
	return t.scrollBarSmallSize
}

func (t *configurableTheme) ShadowColor() color.Color {
	return t.shadow
}

func (t *configurableTheme) TextColor() color.Color {
	return t.text
}

func (t *configurableTheme) TextSize() int {
	return t.textSize
}

func (t *configurableTheme) TextBoldFont() fyne.Resource {
	return t.bold
}

func (t *configurableTheme) TextBoldItalicFont() fyne.Resource {
	return t.boldItalic
}

func (t *configurableTheme) TextFont() fyne.Resource {
	return t.regular
}

func (t *configurableTheme) TextItalicFont() fyne.Resource {
	return t.italic
}

func (t *configurableTheme) TextMonospaceFont() fyne.Resource {
	return t.monospace
}
