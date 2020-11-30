package test

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var defaultTheme fyne.Theme

var _ fyne.Theme = (*configurableTheme)(nil)

type configurableTheme struct {
	colors map[fyne.ThemeColorName]color.Color
	fonts  map[fyne.TextStyle]fyne.Resource
	sizes  map[fyne.ThemeSizeName]int
}

// Theme returns a theme useful for image based tests.
func Theme() fyne.Theme {
	if defaultTheme == nil {
		defaultTheme = &configurableTheme{
			colors: map[fyne.ThemeColorName]color.Color{
				theme.Colors.Background:     color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff},
				theme.Colors.Button:         color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
				theme.Colors.Disabled:       color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
				theme.Colors.DisabledButton: color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
				theme.Colors.Focus:          color.NRGBA{R: 0x81, G: 0xc7, B: 0x84, A: 0xff},
				theme.Colors.Foreground:     color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
				theme.Colors.Hover:          color.NRGBA{R: 0x88, G: 0xff, B: 0xff, A: 0x22},
				theme.Colors.PlaceHolder:    color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff},
				theme.Colors.Primary:        color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
				theme.Colors.ScrollBar:      color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa},
				theme.Colors.Shadow:         color.NRGBA{A: 0x88},
			},
			fonts: map[fyne.TextStyle]fyne.Resource{
				fyne.TextStyle{}:                         theme.DefaultTextFont(),
				fyne.TextStyle{Bold: true}:               theme.DefaultTextBoldFont(),
				fyne.TextStyle{Bold: true, Italic: true}: theme.DefaultTextBoldItalicFont(),
				fyne.TextStyle{Italic: true}:             theme.DefaultTextItalicFont(),
				fyne.TextStyle{Monospace: true}:          theme.DefaultTextMonospaceFont(),
			},
			sizes: map[fyne.ThemeSizeName]int{
				theme.Sizes.InlineIcon:     20,
				theme.Sizes.Padding:        4,
				theme.Sizes.ScrollBar:      16,
				theme.Sizes.ScrollBarSmall: 3,
				theme.Sizes.Text:           14,
			},
		}
	}
	return defaultTheme
}

func (t *configurableTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.colors[n]
}

func (t *configurableTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.fonts[style]
}

func (t *configurableTheme) Size(s fyne.ThemeSizeName) int {
	return t.sizes[s]
}
