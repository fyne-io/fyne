package ui

type CanvasTextObject interface {
	CanvasObject
	FontSize() int
	SetFontSize(int)
	Bold() bool
	SetBold(bool)
	Italic() bool
	SetItalic(bool)
}
