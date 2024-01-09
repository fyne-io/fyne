package fyne

// TextAlign represents the horizontal alignment of text within a widget or
// canvas object.
type TextAlign int

const (
	// TextAlignLeading specifies a left alignment for left-to-right languages.
	TextAlignLeading TextAlign = iota
	// TextAlignCenter places the text centrally within the available space.
	TextAlignCenter
	// TextAlignTrailing will align the text right for a left-to-right language.
	TextAlignTrailing
)

// TextTruncation controls how a `Label` or `Entry` will truncate its text.
// The default value is `TextTruncateOff` which will not truncate.
//
// Since: 2.4
type TextTruncation int

const (
	// TextTruncateOff means no truncation will be applied, it is the default.
	// This means that the minimum size of a text block will be the space required to display it fully.
	TextTruncateOff TextTruncation = iota
	// TextTruncateClip will truncate text when it reaches the end of the available space.
	TextTruncateClip
	// TextTruncateEllipsis is like regular truncation except that an ellipses (…) will be inserted
	// wherever text has been shortened to fit.
	//
	// Since: 2.4
	TextTruncateEllipsis
)

// TextWrap represents how text longer than the widget's width will be wrapped.
type TextWrap int

const (
	// TextWrapOff extends the widget's width to fit the text, no wrapping is applied.
	TextWrapOff TextWrap = iota
	// TextTruncate trims the text to the widget's width, no wrapping is applied.
	// If an entry is asked to truncate it will provide scrolling capabilities.
	//
	// Deprecated: Use `TextTruncateClip` value of the widget `Truncation` field instead
	TextTruncate
	// TextWrapBreak trims the line of characters to the widget's width adding the excess as new line.
	// An Entry with text wrapping will scroll vertically if there is not enough space for all the text.
	TextWrapBreak
	// TextWrapWord trims the line of words to the widget's width adding the excess as new line.
	// An Entry with text wrapping will scroll vertically if there is not enough space for all the text.
	TextWrapWord
)

// TextStyle represents the styles that can be applied to a text canvas object
// or text based widget.
type TextStyle struct {
	Bold      bool // Should text be bold
	Italic    bool // Should text be italic
	Monospace bool // Use the system monospace font instead of regular
	// Since: 2.2
	Symbol bool // Use the system symbol font.
	// Since: 2.1
	TabWidth int // Width of tabs in spaces
}

// MeasureText uses the current driver to calculate the size of text when rendered.
func MeasureText(text string, size float32, style TextStyle) Size {
	s, _ := CurrentApp().Driver().RenderedTextSize(text, size, style)
	return s
}
