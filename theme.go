package fyne

import "image/color"

// ThemeVariant indicates a variation of a theme, such as light or dark.
//
// Since 2.0.0
type ThemeVariant uint

// ThemeColorName is used to look up a colour based on it's name.
//
// Since 2.0.0
type ThemeColorName string

// ThemeSizeName is used to look up a size based on it's name.
//
// Since 2.0.0
type ThemeSizeName string

// Theme defines the method to look up colors, sizes and fonts that make up a Fyne theme.
//
// Since 2.0.0
type Theme interface {
	Color(ThemeColorName, ThemeVariant) color.Color
	Size(ThemeSizeName) int
	Font(TextStyle) Resource
}

// LegacyTheme defines the requirements of any Fyne theme.
// This was previously called Theme and is kept for simpler transition of applications built before 2.0.0.
//
// Since 2.0.0
type LegacyTheme interface {
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
