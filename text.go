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
	// If an entry is asked to truncate it will provide scrolling capabilities.
	TextTruncate
	// TextWrapBreak trims the line of characters to the widget's width adding the excess as new line.
	// An Entry with text wrapping will scroll vertically if there is not enough space for all the text.
	TextWrapBreak
	// TextWrapWord trims the line of words to the widget's width adding the excess as new line.
	// An Entry with text wrapping will scroll vertically if there is not enough space for all the text.
	TextWrapWord
)

const (
	// DefaultTabWidth is the default width in spaces
	DefaultTabWidth = 4
)

// TextStyle represents the styles that can be applied to a text canvas object
// or text based widget.
type TextStyle struct {
	Bold      bool // Should text be bold
	Italic    bool // Should text be italic
	Monospace bool // Use the system monospace font instead of regular
	tabWidth  int  // Width of tabs in spaces
}

// TabWidth either returns the set tab width or if not set the returns the DefaultTabWidth
func (ts *TextStyle) TabWidth() int {
	if ts.tabWidth == 0 {
		return DefaultTabWidth
	}
	return ts.tabWidth
}

// SetTabWidth sets the tab width if the supplied value is greater than 0
func (ts *TextStyle) SetTabWidth(tabWidth int) {
	if tabWidth > 0 {
		ts.tabWidth = tabWidth
	}
}

// MeasureText uses the current driver to calculate the size of text when rendered.
func MeasureText(text string, size float32, style TextStyle) Size {
	return CurrentApp().Driver().RenderedTextSize(text, size, style)
}
