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
			theme.Colors.Background:     red,
			theme.Colors.Button:         color.Black,
			theme.Colors.Disabled:       color.Black,
			theme.Colors.DisabledButton: color.White,
			theme.Colors.Focus:          green,
			theme.Colors.Foreground:     color.White,
			theme.Colors.Hover:          green,
			theme.Colors.PlaceHolder:    blue,
			theme.Colors.Primary:        green,
			theme.Colors.ScrollBar:      blue,
			theme.Colors.Shadow:         blue,
		},
		fonts: map[fyne.TextStyle]fyne.Resource{
			fyne.TextStyle{}:                         theme.DefaultTextBoldFont(),
			fyne.TextStyle{Bold: true}:               theme.DefaultTextItalicFont(),
			fyne.TextStyle{Bold: true, Italic: true}: theme.DefaultTextMonospaceFont(),
			fyne.TextStyle{Italic: true}:             theme.DefaultTextBoldItalicFont(),
			fyne.TextStyle{Monospace: true}:          theme.DefaultTextFont(),
		},
		sizes: map[fyne.ThemeSizeName]int{
			theme.Sizes.InlineIcon:     24,
			theme.Sizes.Padding:        10,
			theme.Sizes.ScrollBar:      10,
			theme.Sizes.ScrollBarSmall: 2,
			theme.Sizes.Text:           18,
		},
	}
}
