package fyne

// TextAlign represents the horizontal alignment of text within a widget or
// canvas object
type TextAlign int

const (
	// TextAlignLeading specifies a left alignment for left-to-right languages
	TextAlignLeading = iota
	// TextAlignCenter places the text centrally within the available space
	TextAlignCenter
	// TextAlignTrailing will align the text right for a left-to-right language
	TextAlignTrailing
)
