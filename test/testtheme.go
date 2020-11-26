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
		background:         red,
		bold:               theme.DefaultTextItalicFont(),
		boldItalic:         theme.DefaultTextMonospaceFont(),
		button:             color.Black,
		disabledButton:     color.White,
		disabledText:       color.Black,
		focus:              green,
		hover:              green,
		iconInlineSize:     24,
		italic:             theme.DefaultTextBoldItalicFont(),
		monospace:          theme.DefaultTextFont(),
		padding:            10,
		placeholder:        blue,
		primary:            green,
		regular:            theme.DefaultTextBoldFont(),
		scrollBar:          blue,
		scrollBarSize:      10,
		scrollBarSmallSize: 2,
		shadow:             blue,
		text:               color.White,
		textSize:           18,
	}
}
