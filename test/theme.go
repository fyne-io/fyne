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
	sizes  map[fyne.ThemeSizeName]float32
}

// Theme returns a theme useful for image based tests.
func Theme() fyne.Theme {
	if defaultTheme == nil {
		defaultTheme = &configurableTheme{
			colors: map[fyne.ThemeColorName]color.Color{
				theme.ColorNameBackground:     color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff},
				theme.ColorNameButton:         color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
				theme.ColorNameDisabled:       color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
				theme.ColorNameDisabledButton: color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
				theme.ColorNameError:          color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff},
				theme.ColorNameFocus:          color.NRGBA{R: 0x81, G: 0xc7, B: 0x84, A: 0xff},
				theme.ColorNameForeground:     color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
				theme.ColorNameHover:          color.NRGBA{R: 0x88, G: 0xff, B: 0xff, A: 0x22},
				theme.ColorNamePlaceHolder:    color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff},
				theme.ColorNamePrimary:        color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
				theme.ColorNameScrollBar:      color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa},
				theme.ColorNameShadow:         color.NRGBA{A: 0x88},
			},
			fonts: map[fyne.TextStyle]fyne.Resource{
				fyne.TextStyle{}:                         theme.DefaultTextFont(),
				fyne.TextStyle{Bold: true}:               theme.DefaultTextBoldFont(),
				fyne.TextStyle{Bold: true, Italic: true}: theme.DefaultTextBoldItalicFont(),
				fyne.TextStyle{Italic: true}:             theme.DefaultTextItalicFont(),
				fyne.TextStyle{Monospace: true}:          theme.DefaultTextMonospaceFont(),
			},
			sizes: map[fyne.ThemeSizeName]float32{
				theme.SizeNameInlineIcon:     float32(20),
				theme.SizeNamePadding:        float32(4),
				theme.SizeNameScrollBar:      float32(16),
				theme.SizeNameScrollBarSmall: float32(3),
				theme.SizeNameText:           float32(14),
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

func (t *configurableTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *configurableTheme) Size(s fyne.ThemeSizeName) float32 {
	return t.sizes[s]
}
