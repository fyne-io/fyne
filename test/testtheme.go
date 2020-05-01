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

// BackgroundColor satisfies the fyne.Theme interface.
func (testTheme) BackgroundColor() color.Color {
	return red
}

// ButtonColor satisfies the fyne.Theme interface.
func (testTheme) ButtonColor() color.Color {
	return color.Black
}

// DisabledButtonColor satisfies the fyne.Theme interface.
func (testTheme) DisabledButtonColor() color.Color {
	return color.White
}

// DisabledIconColor satisfies the fyne.Theme interface.
func (testTheme) DisabledIconColor() color.Color {
	return color.Black
}

// DisabledTextColor satisfies the fyne.Theme interface.
func (testTheme) DisabledTextColor() color.Color {
	return color.Black
}

// FocusColor satisfies the fyne.Theme interface.
func (testTheme) FocusColor() color.Color {
	return green
}

// HoverColor satisfies the fyne.Theme interface.
func (testTheme) HoverColor() color.Color {
	return green
}

// HyperlinkColor satisfies the fyne.Theme interface.
func (testTheme) HyperlinkColor() color.Color {
	return green
}

// IconColor satisfies the fyne.Theme interface.
func (testTheme) IconColor() color.Color {
	return color.White
}

// IconInlineSize satisfies the fyne.Theme interface.
func (testTheme) IconInlineSize() int {
	return 24
}

// Padding satisfies the fyne.Theme interface.
func (testTheme) Padding() int {
	return 10
}

// PlaceHolderColor satisfies the fyne.Theme interface.
func (testTheme) PlaceHolderColor() color.Color {
	return blue
}

// PrimaryColor satisfies the fyne.Theme interface.
func (testTheme) PrimaryColor() color.Color {
	return green
}

// ScrollBarColor satisfies the fyne.Theme interface.
func (testTheme) ScrollBarColor() color.Color {
	return blue
}

// ScrollBarSize satisfies the fyne.Theme interface.
func (testTheme) ScrollBarSize() int {
	return 10
}

// ScrollBarSmallSize satisfies the fyne.Theme interface.
func (testTheme) ScrollBarSmallSize() int {
	return 2
}

// ShadowColor satisfies the fyne.Theme interface.
func (testTheme) ShadowColor() color.Color {
	return blue
}

// TextBoldFont satisfies the fyne.Theme interface.
func (testTheme) TextBoldFont() fyne.Resource {
	return theme.DefaultTextItalicFont()
}

// TextBoldItalicFont satisfies the fyne.Theme interface.
func (testTheme) TextBoldItalicFont() fyne.Resource {
	return theme.DefaultTextMonospaceFont()
}

// TextColor satisfies the fyne.Theme interface.
func (testTheme) TextColor() color.Color {
	return color.White
}

// TextFont satisfies the fyne.Theme interface.
func (testTheme) TextFont() fyne.Resource {
	return theme.DefaultTextBoldFont()
}

// TextItalicFont satisfies the fyne.Theme interface.
func (testTheme) TextItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

// TextMonospaceFont satisfies the fyne.Theme interface.
func (testTheme) TextMonospaceFont() fyne.Resource {
	return theme.DefaultTextFont()
}

// TextSize satisfies the fyne.Theme interface.
func (testTheme) TextSize() int {
	return testTextSize
}
