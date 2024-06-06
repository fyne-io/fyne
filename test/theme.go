package test

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var defaultTheme fyne.Theme

// NewTheme returns a new test theme using quiet ugly colors.
func NewTheme() fyne.Theme {
	blue := func(alpha uint8) color.Color {
		return &color.RGBA{R: 0, G: 0, B: 255, A: alpha}
	}
	gray := func(level uint8) color.Color {
		return &color.Gray{Y: level}
	}
	green := func(alpha uint8) color.Color {
		return &color.RGBA{R: 0, G: 255, B: 0, A: alpha}
	}
	red := func(alpha uint8) color.Color {
		return &color.RGBA{R: 200, G: 0, B: 0, A: alpha}
	}

	return &configurableTheme{
		colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:        red(255),
			theme.ColorNameButton:            gray(0),
			theme.ColorNameDisabled:          gray(0),
			theme.ColorNameDisabledButton:    gray(255),
			theme.ColorNameError:             blue(255),
			theme.ColorNameFocus:             red(66),
			theme.ColorNameForeground:        gray(255),
			theme.ColorNameHover:             green(255),
			theme.ColorNameHeaderBackground:  red(22),
			theme.ColorNameInputBackground:   red(30),
			theme.ColorNameInputBorder:       gray(0),
			theme.ColorNameMenuBackground:    red(30),
			theme.ColorNameOnPrimary:         red(200),
			theme.ColorNameOverlayBackground: red(44),
			theme.ColorNamePlaceHolder:       blue(255),
			theme.ColorNamePressed:           blue(255),
			theme.ColorNamePrimary:           green(255),
			theme.ColorNameScrollBar:         blue(255),
			theme.ColorNameSeparator:         gray(0),
			theme.ColorNameSelection:         red(44),
			theme.ColorNameShadow:            blue(255),
		},
		name: "Ugly Test Theme",
		fonts: map[fyne.TextStyle]fyne.Resource{
			{}:                         theme.DefaultTextBoldFont(),
			{Bold: true}:               theme.DefaultTextItalicFont(),
			{Bold: true, Italic: true}: theme.DefaultTextMonospaceFont(),
			{Italic: true}:             theme.DefaultTextBoldItalicFont(),
			{Monospace: true}:          theme.DefaultTextFont(),
		},
		sizes: map[fyne.ThemeSizeName]float32{
			theme.SizeNameInlineIcon:         float32(24),
			theme.SizeNameInnerPadding:       float32(20),
			theme.SizeNameLineSpacing:        float32(6),
			theme.SizeNamePadding:            float32(10),
			theme.SizeNameScrollBar:          float32(10),
			theme.SizeNameScrollBarSmall:     float32(2),
			theme.SizeNameSeparatorThickness: float32(1),
			theme.SizeNameText:               float32(18),
			theme.SizeNameHeadingText:        float32(30.6),
			theme.SizeNameSubHeadingText:     float32(24),
			theme.SizeNameCaptionText:        float32(15),
			theme.SizeNameInputBorder:        float32(5),
			theme.SizeNameInputRadius:        float32(2),
			theme.SizeNameSelectionRadius:    float32(6),
		},
	}
}

// Theme returns a test theme useful for image based tests.
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
				theme.ColorNameOnPrimary:         color.NRGBA{R: 0x08, G: 0x0c, B: 0x0f, A: 0xff},
				theme.ColorNameOverlayBackground: color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff},
				theme.ColorNamePlaceHolder:       color.NRGBA{R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff},
				theme.ColorNamePressed:           color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33},
				theme.ColorNamePrimary:           color.NRGBA{R: 0xff, G: 0xcc, B: 0x80, A: 0xff},
				theme.ColorNameScrollBar:         color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xaa},
				theme.ColorNameSeparator:         color.NRGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff},
				theme.ColorNameSelection:         color.NRGBA{R: 0x78, G: 0x3a, B: 0x3a, A: 0x99},
				theme.ColorNameShadow:            color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x88},
			},
			fonts: map[fyne.TextStyle]fyne.Resource{
				{}:                         theme.DefaultTextFont(),
				{Bold: true}:               theme.DefaultTextBoldFont(),
				{Bold: true, Italic: true}: theme.DefaultTextBoldItalicFont(),
				{Italic: true}:             theme.DefaultTextItalicFont(),
				{Monospace: true}:          theme.DefaultTextMonospaceFont(),
			},
			name: "Default Test Theme",
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

type configurableTheme struct {
	colors map[fyne.ThemeColorName]color.Color
	fonts  map[fyne.TextStyle]fyne.Resource
	name   string
	sizes  map[fyne.ThemeSizeName]float32
}

var _ fyne.Theme = (*configurableTheme)(nil)

func (t *configurableTheme) Color(n fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	if t.colors[n] == nil {
		fyne.LogError(fmt.Sprintf("color %s not defined in theme %s", n, t.name), nil)
	}

	return t.colors[n]
}

func (t *configurableTheme) Font(style fyne.TextStyle) fyne.Resource {
	if t.fonts[style] == nil {
		fyne.LogError(fmt.Sprintf("font for style %#v not defined in theme %s", style, t.name), nil)
	}

	return t.fonts[style]
}

func (t *configurableTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *configurableTheme) Size(s fyne.ThemeSizeName) float32 {
	if t.sizes[s] == 0 {
		fyne.LogError(fmt.Sprintf("size %s not defined in theme %s", s, t.name), nil)
	}

	return t.sizes[s]
}
