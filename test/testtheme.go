package test

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var (
	red   = &color.RGBA{R: 200, G: 0, B: 0, A: 255}
	green = &color.RGBA{R: 0, G: 255, B: 0, A: 255}
	blue  = &color.RGBA{R: 0, G: 0, B: 255, A: 255}
)

// NewTheme returns a new testTheme.
func NewTheme() fyne.Theme {
	return &configurableTheme{
		colors: map[fyne.ThemeColorName]color.Color{
			theme.ColorNameBackground:     red,
			theme.ColorNameButton:         color.Black,
			theme.ColorNameDisabled:       color.Black,
			theme.ColorNameDisabledButton: color.White,
			theme.ColorNameError:          blue,
			theme.ColorNameFocus:          green,
			theme.ColorNameForeground:     color.White,
			theme.ColorNameHover:          green,
			theme.ColorNamePlaceHolder:    blue,
			theme.ColorNamePressed:        blue,
			theme.ColorNamePrimary:        green,
			theme.ColorNameScrollBar:      blue,
			theme.ColorNameShadow:         blue,
		},
		fonts: map[fyne.TextStyle]fyne.Resource{
			fyne.TextStyle{}:                         theme.DefaultTextBoldFont(),
			fyne.TextStyle{Bold: true}:               theme.DefaultTextItalicFont(),
			fyne.TextStyle{Bold: true, Italic: true}: theme.DefaultTextMonospaceFont(),
			fyne.TextStyle{Italic: true}:             theme.DefaultTextBoldItalicFont(),
			fyne.TextStyle{Monospace: true}:          theme.DefaultTextFont(),
		},
		sizes: map[fyne.ThemeSizeName]float32{
			theme.SizeNameInlineIcon:         float32(24),
			theme.SizeNamePadding:            float32(10),
			theme.SizeNameScrollBar:          float32(10),
			theme.SizeNameScrollBarSmall:     float32(2),
			theme.SizeNameSeparatorThickness: float32(1),
			theme.SizeNameText:               float32(18),
		},
	}
}
