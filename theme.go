package fyne

import "image/color"

// Theme defines the requirements of any Fyne theme
type Theme interface {
	BackgroundColor() color.Color
	ButtonColor() color.Color
	TextColor() color.Color
	PrimaryColor() color.Color
	FocusColor() color.Color

	TextSize() int
	TextFont() Resource
	TextBoldFont() Resource
	TextItalicFont() Resource
	TextBoldItalicFont() Resource
	TextMonospaceFont() Resource

	Padding() int
	IconInlineSize() int
}
