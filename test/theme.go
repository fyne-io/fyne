package test

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var defaultTheme fyne.Theme

var _ fyne.Theme = (*configurableTheme)(nil)

type configurableTheme struct {
	background         color.Color
	bold               fyne.Resource
	boldItalic         fyne.Resource
	button             color.Color
	disabledButton     color.Color
	disabledText       color.Color
	focus              color.Color
	hover              color.Color
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
	if defaultTheme == nil {
		defaultTheme = &configurableTheme{
			background:         color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff},
			bold:               theme.DefaultTextBoldFont(),
			boldItalic:         theme.DefaultTextBoldItalicFont(),
			button:             color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
			disabledButton:     color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
			disabledText:       color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
			focus:              color.NRGBA{R: 0x81, G: 0xc7, B: 0x84, A: 0xff},
			hover:              color.NRGBA{R: 0x88, G: 0xff, B: 0xff, A: 0x22},
			iconInlineSize:     20,
			italic:             theme.DefaultTextItalicFont(),
			monospace:          theme.DefaultTextMonospaceFont(),
			padding:            4,
			placeholder:        color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff},
			primary:            color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
			regular:            theme.DefaultTextFont(),
			scrollBar:          color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa},
			scrollBarSize:      16,
			scrollBarSmallSize: 3,
			shadow:             color.NRGBA{A: 0x88},
			text:               color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			textSize:           14,
		}
	}
	return defaultTheme
}

func (t *configurableTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch n {
	case theme.Colors.Background:
		return t.background
	case theme.Colors.Text:
		return t.text
	case theme.Colors.Button:
		return t.button
	case theme.Colors.DisabledButton:
		return t.disabledButton
	case theme.Colors.DisabledText:
		return t.disabledText
	case theme.Colors.Focus:
		return t.focus
	case theme.Colors.Hover:
		return t.hover
	case theme.Colors.PlaceHolder:
		return t.placeholder
	case theme.Colors.Primary:
		return t.primary
	case theme.Colors.ScrollBar:
		return t.scrollBar
	case theme.Colors.Shadow:
		return t.shadow
	default:
		return color.Transparent
	}
}

func (t *configurableTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return t.monospace
	}
	if style.Bold {
		if style.Italic {
			return t.boldItalic
		}
		return t.bold
	}
	if style.Italic {
		return t.italic
	}
	return t.regular
}

func (t *configurableTheme) Size(s fyne.ThemeSizeName) int {
	switch s {
	case theme.Sizes.InlineIcon:
		return t.iconInlineSize
	case theme.Sizes.Padding:
		return t.padding
	case theme.Sizes.ScrollBar:
		return t.scrollBarSize
	case theme.Sizes.ScrollBarSmall:
		return t.scrollBarSmallSize
	case theme.Sizes.Text:
		return t.textSize
	default:
		return 0
	}
}
