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

	// SizeNameScrollPadding is the name of theme lookup for scroll padding. Can be used to make the scrollbar non-overlapping.
	SizeNameScrollPadding fyne.ThemeSizeName = "scrollPadding"

	// SizeNameScrollBarPaddingAlong is the name of theme lookup for scrollbar padding along. Can be used to make the scrollbar-thumb shorter than.
	SizeNameScrollBarPaddingAlong fyne.ThemeSizeName = "scrollBarPaddingAlong"

	// SizeNameScrollBarPaddingAcross is the name of theme lookup for scrollbar padding across. Can be used to make the scrollbar-thumb thinner than the scrollbar.
	SizeNameScrollBarPaddingAcross fyne.ThemeSizeName = "scrollBarPaddingAcross"

	// SizeNameScrollBarCornerRadius is the name of theme lookup for scrollbar corner radius..
	SizeNameScrollBarCornerRadius fyne.ThemeSizeName = "scrollBarCornerRadius"

	// SizeNameScrollBar is the name of theme lookup for the scrollbar size.
	//
	// Since: 2.0
	SizeNameScrollBar fyne.ThemeSizeName = "scrollBar"

	// SizeNameScrollBarLimit is the name of theme lookup for the scroll-thumb size limit. Scrollbar-thumbs will not get shorter than this.
	SizeNameScrollBarLimit fyne.ThemeSizeName = "scrollBarLimit"

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
)

// CaptionTextSize returns the size for caption text.
func CaptionTextSize() float32 {
	return current().Size(SizeNameCaptionText)
}

// IconInlineSize is the standard size of icons which appear within buttons, labels etc.
func IconInlineSize() float32 {
	return current().Size(SizeNameInlineIcon)
}

// InnerPadding is the standard gap between element content and the outside edge of a widget.
//
// Since: 2.3
func InnerPadding() float32 {
	return current().Size(SizeNameInnerPadding)
}

// InputBorderSize returns the input border size (or underline size for an entry).
//
// Since: 2.0
func InputBorderSize() float32 {
	return current().Size(SizeNameInputBorder)
}

// InputRadiusSize returns the input radius size.
//
// Since: 2.4
func InputRadiusSize() float32 {
	return current().Size(SizeNameInputRadius)
}

// LineSpacing is the default gap between multiple lines of text.
//
// Since: 2.3
func LineSpacing() float32 {
	return current().Size(SizeNameLineSpacing)
}

// Padding is the standard gap between elements and the border around interface elements.
func Padding() float32 {
	return current().Size(SizeNamePadding)
}

// ScrollPadding is the gap under scrollbars, can be used to make scrollbars non-overlapping.
func ScrollPadding() float32 {
	return current().Size(SizeNameScrollPadding)
}

// ScrollBarPaddingAlong is the gap at the start and end of the scrollbar.
func ScrollBarPaddingAlong() float32 {
	return current().Size(SizeNameScrollBarPaddingAlong)
}

// ScrollBarPaddingAcross is the gap at the sides of the scrollbar, making the scrollbar-thumb thinner.
func ScrollBarPaddingAcross() float32 {
	return current().Size(SizeNameScrollBarPaddingAcross)
}

// ScrollBarCornerRadius is the amount of rounding of scrollbar-thumbs.
func ScrollBarCornerRadius() float32 {
	return current().Size(SizeNameScrollBarCornerRadius)
}

// ScrollBarSize is the width (or height) of the bars on a ScrollContainer.
func ScrollBarSize() float32 {
	return current().Size(SizeNameScrollBar)
}

// ScrollBarLimitSize is the minimal length of scrollbar-thumbs.
func ScrollBarLimitSize() float32 {
	return current().Size(SizeNameScrollBarLimit)
}

// ScrollBarSmallSize is the width (or height) of the minimized bars on a ScrollContainer.
func ScrollBarSmallSize() float32 {
	return current().Size(SizeNameScrollBarSmall)
}

// SelectionRadiusSize returns the selection highlight radius size.
//
// Since: 2.4
func SelectionRadiusSize() float32 {
	return current().Size(SizeNameSelectionRadius)
}

// SeparatorThicknessSize is the standard thickness of the separator widget.
//
// Since: 2.0
func SeparatorThicknessSize() float32 {
	return current().Size(SizeNameSeparatorThickness)
}

// TextHeadingSize returns the text size for header text.
//
// Since: 2.1
func TextHeadingSize() float32 {
	return current().Size(SizeNameHeadingText)
}

// TextSize returns the standard text size.
func TextSize() float32 {
	return current().Size(SizeNameText)
}

// TextSubHeadingSize returns the text size for sub-header text.
//
// Since: 2.1
func TextSubHeadingSize() float32 {
	return current().Size(SizeNameSubHeadingText)
}
