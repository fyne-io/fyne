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

// TextWrap represents how text longer than the widget's width will be wrapped.
type TextWrap int

const (
	// TextWrapOff extends the widget's width to fit the text, no wrapping is applied.
	TextWrapOff TextWrap = iota
	// TextTruncate trims the text to the widget's width, no wrapping is applied.
	TextTruncate
	// TextWrapBreak trims the line of characters to the widget's width adding the excess as new line.
	TextWrapBreak
	// TextWrapWord trims the line of words to the widget's width adding the excess as new line.
	TextWrapWord
)

// TextStyle represents the styles that can be applied to a text canvas object
// or text based widget.
type TextStyle struct {
	Bold      bool // Should text be bold
	Italic    bool // Should text be italic
	Monospace bool // Use the system monospace font instead of regular
}

// MeasureText uses the current driver to calculate the size of text when rendered.
func MeasureText(text string, size int, style TextStyle) Size {
	return CurrentApp().Driver().RenderedTextSize(text, size, style)
}
