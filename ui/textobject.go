package ui

import "image/color"

import "github.com/fyne-io/fyne/ui/theme"

type TextObject struct {
	Color    color.RGBA
	Size     Size
	Position Position

	Text     string
	FontSize int
	Bold     bool
	Italic   bool

	minSize Size
}

func (t *TextObject) CurrentSize() Size {
	return t.Size
}

func (t *TextObject) Resize(size Size) {
	t.Size = size
}

func (t *TextObject) CurrentPosition() Position {
	return t.Position
}

func (t *TextObject) Move(pos Position) {
	t.Position = pos
}

func (t *TextObject) SetMinSize(size Size) {
	t.minSize = size
}

func (t *TextObject) MinSize() Size {
	return t.minSize
}

func NewText(text string) *TextObject {
	return &TextObject{
		Color:    theme.TextColor(),
		Text:     text,
		FontSize: theme.TextSize(),
	}
}
