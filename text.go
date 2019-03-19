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

// TextStyle represents the styles that can be applied to a text canvas object
// or text based widget.
type TextStyle struct {
	Bold      bool // Should text be bold
	Italic    bool // Should text be italic
	Monospace bool // Use the system monospace font instead of regular
}
