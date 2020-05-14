package test

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

var (
	red   = &color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green = &color.RGBA{R: 0, G: 255, B: 0, A: 255}
	blue  = &color.RGBA{R: 0, G: 0, B: 255, A: 255}
)

const testTextSize = 18

var _ fyne.Theme = testTheme{}

// testTheme is a simple theme variation used for testing that widgets adapt correctly
type testTheme struct {
}

// NewTheme returns a new testTheme.
func NewTheme() fyne.Theme {
	return &testTheme{}
}

func (testTheme) BackgroundColor() color.Color {
	return red
}

func (testTheme) ButtonColor() color.Color {
	return color.Black
}

func (testTheme) DisabledButtonColor() color.Color {
	return color.White
}

func (testTheme) DisabledIconColor() color.Color {
	return color.Black
}

func (testTheme) DisabledTextColor() color.Color {
	return color.Black
}

func (testTheme) FocusColor() color.Color {
	return green
}

func (testTheme) HoverColor() color.Color {
	return green
}

func (testTheme) HyperlinkColor() color.Color {
	return green
}

func (testTheme) IconColor() color.Color {
	return color.White
}

func (testTheme) IconInlineSize() int {
	return 24
}

func (testTheme) Padding() int {
	return 10
}

func (testTheme) PlaceHolderColor() color.Color {
	return blue
}

func (testTheme) PrimaryColor() color.Color {
	return green
}

func (testTheme) ScrollBarColor() color.Color {
	return blue
}

func (testTheme) ScrollBarSize() int {
	return 10
}

func (testTheme) ScrollBarSmallSize() int {
	return 2
}

func (testTheme) ShadowColor() color.Color {
	return blue
}

func (testTheme) TextBoldFont() fyne.Resource {
	return theme.DefaultTextItalicFont()
}

func (testTheme) TextBoldItalicFont() fyne.Resource {
	return theme.DefaultTextMonospaceFont()
}

func (testTheme) TextColor() color.Color {
	return color.White
}

func (testTheme) TextFont() fyne.Resource {
	return theme.DefaultTextBoldFont()
}

func (testTheme) TextItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

func (testTheme) TextMonospaceFont() fyne.Resource {
	return theme.DefaultTextFont()
}

func (testTheme) TextSize() int {
	return testTextSize
}
