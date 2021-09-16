package test

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
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
			theme.ColorNameBackground:      red,
			theme.ColorNameButton:          color.Black,
			theme.ColorNameDisabled:        color.Black,
			theme.ColorNameDisabledButton:  color.White,
			theme.ColorNameError:           blue,
			theme.ColorNameFocus:           color.RGBA{red.R, red.G, red.B, 66},
			theme.ColorNameForeground:      color.White,
			theme.ColorNameHover:           green,
			theme.ColorNameInputBackground: color.RGBA{red.R, red.G, red.B, 30},
			theme.ColorNamePlaceHolder:     blue,
			theme.ColorNamePressed:         blue,
			theme.ColorNamePrimary:         green,
			theme.ColorNameScrollBar:       blue,
			theme.ColorNameSelection:       color.RGBA{red.R, red.G, red.B, 44},
			theme.ColorNameShadow:          blue,
		},
		fonts: map[fyne.TextStyle]fyne.Resource{
			{}:                         theme.DefaultTextBoldFont(),
			{Bold: true}:               theme.DefaultTextItalicFont(),
			{Bold: true, Italic: true}: theme.DefaultTextMonospaceFont(),
			{Italic: true}:             theme.DefaultTextBoldItalicFont(),
			{Monospace: true}:          theme.DefaultTextFont(),
		},
		sizes: map[fyne.ThemeSizeName]float32{
			theme.SizeNameInlineIcon:         float32(24),
			theme.SizeNamePadding:            float32(10),
			theme.SizeNameScrollBar:          float32(10),
			theme.SizeNameScrollBarSmall:     float32(2),
			theme.SizeNameSeparatorThickness: float32(1),
			theme.SizeNameText:               float32(18),
			theme.SizeNameHeadingText:        float32(30.6),
			theme.SizeNameSubHeadingText:     float32(24),
			theme.SizeNameCaptionText:        float32(15),
			theme.SizeNameInputBorder:        float32(5),
		},
	}
}
