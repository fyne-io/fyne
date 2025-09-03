package theme

import "fyne.io/fyne/v2"

const (
	// SizeNameCaptionText is the name of theme lookup for helper text size, normally smaller than regular text size.
	//
	// Since: 2.0
	SizeNameCaptionText fyne.ThemeSizeName = "helperText"

	// SizeNameInlineIcon is the name of theme lookup for inline icons size.
	//
	// Since: 2.0
	SizeNameInlineIcon fyne.ThemeSizeName = "iconInline"

	// SizeNameInnerPadding is the name of theme lookup for internal widget padding size.
	//
	// Since: 2.3
	SizeNameInnerPadding fyne.ThemeSizeName = "innerPadding"

	// SizeNameLineSpacing is the name of theme lookup for between text line spacing.
	//
	// Since: 2.3
	SizeNameLineSpacing fyne.ThemeSizeName = "lineSpacing"

	// SizeNamePadding is the name of theme lookup for padding size.
	//
	// Since: 2.0
	SizeNamePadding fyne.ThemeSizeName = "padding"

	// SizeNameScrollBar is the name of theme lookup for the scrollbar size.
	//
	// Since: 2.0
	SizeNameScrollBar fyne.ThemeSizeName = "scrollBar"

	// SizeNameScrollBarSmall is the name of theme lookup for the shrunk scrollbar size.
	//
	// Since: 2.0
	SizeNameScrollBarSmall fyne.ThemeSizeName = "scrollBarSmall"

	// SizeNameSeparatorThickness is the name of theme lookup for the thickness of a separator.
	//
	// Since: 2.0
	SizeNameSeparatorThickness fyne.ThemeSizeName = "separator"

	// SizeNameText is the name of theme lookup for text size.
	//
	// Since: 2.0
	SizeNameText fyne.ThemeSizeName = "text"

	// SizeNameHeadingText is the name of theme lookup for text size of a heading.
	//
	// Since: 2.1
	SizeNameHeadingText fyne.ThemeSizeName = "headingText"

	// SizeNameSubHeadingText is the name of theme lookup for text size of a sub-heading.
	//
	// Since: 2.1
	SizeNameSubHeadingText fyne.ThemeSizeName = "subHeadingText"

	// SizeNameInputBorder is the name of theme lookup for input border size.
	//
	// Since: 2.0
	SizeNameInputBorder fyne.ThemeSizeName = "inputBorder"

	// SizeNameInputRadius is the name of theme lookup for input corner radius.
	//
	// Since: 2.4
	SizeNameInputRadius fyne.ThemeSizeName = "inputRadius"

	// SizeNameSelectionRadius is the name of theme lookup for selection corner radius.
	//
	// Since: 2.4
	SizeNameSelectionRadius fyne.ThemeSizeName = "selectionRadius"

	// SizeNameScrollBarRadius is the name of theme lookup for the scroll bar corner radius.
	//
	// Since: 2.5
	SizeNameScrollBarRadius fyne.ThemeSizeName = "scrollBarRadius"

	// SizeNameWindowButtonHeight is the name of the height for an inner window titleBar button.
	//
	// Since: 2.6
	SizeNameWindowButtonHeight fyne.ThemeSizeName = "windowButtonHeight"

	// SizeNameWindowButtonRadius is the name of the radius for an inner window titleBar button.
	//
	// Since: 2.6
	SizeNameWindowButtonRadius fyne.ThemeSizeName = "windowButtonRadius"

	// SizeNameWindowButtonIcon is the name of the width of an inner window titleBar button.
	//
	// Since: 2.6
	SizeNameWindowButtonIcon fyne.ThemeSizeName = "windowButtonIcon"

	// SizeNameWindowTitleBarHeight is the height for inner window titleBars.
	//
	// Since: 2.6
	SizeNameWindowTitleBarHeight fyne.ThemeSizeName = "windowTitleBarHeight"
)

// CaptionTextSize returns the size for caption text.
func CaptionTextSize() float32 {
	return Current().Size(SizeNameCaptionText)
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func IconInlineSize() float32 {
	return Current().Size(SizeNameInlineIcon)
}

// InnerPadding is the standard gap between element content and the outside edge of a widget.
//
// Since: 2.3
func InnerPadding() float32 {
	return Current().Size(SizeNameInnerPadding)
}

// InputBorderSize returns the input border size (or underline size for an entry).
//
// Since: 2.0
func InputBorderSize() float32 {
	return Current().Size(SizeNameInputBorder)
}

// InputRadiusSize returns the input radius size.
//
// Since: 2.4
func InputRadiusSize() float32 {
	return Current().Size(SizeNameInputRadius)
}

// LineSpacing is the default gap between multiple lines of text.
//
// Since: 2.3
func LineSpacing() float32 {
	return Current().Size(SizeNameLineSpacing)
}

// Padding is the standard gap between elements and the border around interface elements.
func Padding() float32 {
	return Current().Size(SizeNamePadding)
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer.
func ScrollBarSize() float32 {
	return Current().Size(SizeNameScrollBar)
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer.
func ScrollBarSmallSize() float32 {
	return Current().Size(SizeNameScrollBarSmall)
}

// SelectionRadiusSize returns the selection highlight radius size.
//
// Since: 2.4
func SelectionRadiusSize() float32 {
	return Current().Size(SizeNameSelectionRadius)
}

// SeparatorThicknessSize is the standard thickness of the separator widget.
//
// Since: 2.0
func SeparatorThicknessSize() float32 {
	return Current().Size(SizeNameSeparatorThickness)
}

// Size looks up the specified size for current theme.
//
// Since: 2.5
func Size(name fyne.ThemeSizeName) float32 {
	return Current().Size(name)
}

// SizeForWidget looks up the specified size for the requested widget using the current theme.
// If the widget theme has been overridden that theme will be used.
//
// Since: 2.5
func SizeForWidget(name fyne.ThemeSizeName, w fyne.Widget) float32 {
	return CurrentForWidget(w).Size(name)
}

// TextHeadingSize returns the text size for header text.
//
// Since: 2.1
func TextHeadingSize() float32 {
	return Current().Size(SizeNameHeadingText)
}

// TextSize returns the standard text size.
func TextSize() float32 {
	return Current().Size(SizeNameText)
}

// TextSubHeadingSize returns the text size for sub-header text.
//
// Since: 2.1
func TextSubHeadingSize() float32 {
	return Current().Size(SizeNameSubHeadingText)
}

func (t *builtinTheme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case SizeNameSeparatorThickness:
		return 1
	case SizeNameInlineIcon:
		return 20
	case SizeNameInnerPadding:
		return 8
	case SizeNameLineSpacing:
		return 4
	case SizeNamePadding:
		return 4
	case SizeNameScrollBar:
		return 12
	case SizeNameScrollBarSmall:
		return 3
	case SizeNameText:
		return 14
	case SizeNameHeadingText:
		return 24
	case SizeNameSubHeadingText:
		return 18
	case SizeNameCaptionText:
		return 11
	case SizeNameInputBorder:
		return 1
	case SizeNameInputRadius:
		return 5
	case SizeNameSelectionRadius:
		return 3
	case SizeNameScrollBarRadius:
		return 3
	case SizeNameWindowButtonHeight:
		return 16
	case SizeNameWindowButtonRadius:
		return 8
	case SizeNameWindowButtonIcon:
		return 14
	case SizeNameWindowTitleBarHeight:
		return 26

	default:
		return 0
	}
}
