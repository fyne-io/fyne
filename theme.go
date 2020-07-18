package fyne

import "image/color"

// Theme defines the requirements of any Fyne theme.
type Theme interface {
	BackgroundColor() color.Color
	ButtonColor() color.Color
	DisabledButtonColor() color.Color
	// Deprecated: Hyperlinks now use the primary color for consistency.
	HyperlinkColor() color.Color
	TextColor() color.Color
	DisabledTextColor() color.Color
	// Deprecated: Icons now use the text colour for consistency.
	IconColor() color.Color
	// Deprecated: Disabled icons match disabled text color for consistency.
	DisabledIconColor() color.Color
	PlaceHolderColor() color.Color
	PrimaryColor() color.Color
	HoverColor() color.Color
	FocusColor() color.Color
	ScrollBarColor() color.Color
	ShadowColor() color.Color

	TextSize() int
	TextFont() Resource
	TextBoldFont() Resource
	TextItalicFont() Resource
	TextBoldItalicFont() Resource
	TextMonospaceFont() Resource

	Padding() int
	IconInlineSize() int
	ScrollBarSize() int
	ScrollBarSmallSize() int
}
