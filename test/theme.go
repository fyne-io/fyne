package test

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
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
				theme.ColorNameBackground:        color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff},
				theme.ColorNameButton:            color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff},
				theme.ColorNameDisabled:          color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
				theme.ColorNameDisabledButton:    color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
				theme.ColorNameError:             color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff},
				theme.ColorNameFocus:             color.NRGBA{R: 0x78, G: 0x3a, B: 0x3a, A: 0xff},
				theme.ColorNameForeground:        color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
				theme.ColorNameHover:             color.NRGBA{R: 0x88, G: 0xff, B: 0xff, A: 0x22},
				theme.ColorNameHeaderBackground:  color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
				theme.ColorNameHyperlink:         color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
				theme.ColorNameInputBackground:   color.NRGBA{R: 0x66, G: 0x66, B: 0x66, A: 0xff},
				theme.ColorNameInputBorder:       color.NRGBA{R: 0x86, G: 0x86, B: 0x86, A: 0xff},
				theme.ColorNameMenuBackground:    color.NRGBA{R: 0x56, G: 0x56, B: 0x56, A: 0xff},
				theme.ColorNameOverlayBackground: color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
				theme.ColorNamePlaceHolder:       color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff},
				theme.ColorNamePressed:           color.NRGBA{A: 0x33},
				theme.ColorNamePrimary:           color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
				theme.ColorNameScrollBar:         color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa},
				theme.ColorNameSeparator:         color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
				theme.ColorNameSelection:         color.NRGBA{R: 0x78, G: 0x3a, B: 0x3a, A: 0x99},
				theme.ColorNameShadow:            color.NRGBA{A: 0x88},
			},
			fonts: map[fyne.TextStyle]fyne.Resource{
				{}:                         theme.DefaultTextFont(),
				{Bold: true}:               theme.DefaultTextBoldFont(),
				{Bold: true, Italic: true}: theme.DefaultTextBoldItalicFont(),
				{Italic: true}:             theme.DefaultTextItalicFont(),
				{Monospace: true}:          theme.DefaultTextMonospaceFont(),
			},
			sizes: map[fyne.ThemeSizeName]float32{
				theme.SizeNameInlineIcon:         float32(20),
				theme.SizeNameInnerPadding:       float32(8),
				theme.SizeNameLineSpacing:        float32(4),
				theme.SizeNamePadding:            float32(4),
				theme.SizeNameScrollBar:          float32(16),
				theme.SizeNameScrollBarSmall:     float32(3),
				theme.SizeNameSeparatorThickness: float32(1),
				theme.SizeNameText:               float32(14),
				theme.SizeNameHeadingText:        float32(23.8),
				theme.SizeNameSubHeadingText:     float32(18),
				theme.SizeNameCaptionText:        float32(11),
				theme.SizeNameInputBorder:        float32(2),
				theme.SizeNameInputRadius:        float32(4),
				theme.SizeNameSelectionRadius:    float32(4),
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
