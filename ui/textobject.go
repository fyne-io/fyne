package ui

import "image/color"

import "github.com/fyne-io/fyne/ui/theme"

type TextObject struct {
	Color color.RGBA
	Size  Size

	Text     string
	FontSize int
	Bold     bool
	Italic   bool
}

func (t *TextObject) CurrentSize() Size {
	return t.Size
}

func (t *TextObject) Resize(size Size) {
	t.Size = size
}

func (t *TextObject) MinSize() Size {
	return NewSize(1, 1) // TODO calcualte from font/string
}

func NewText(text string) *TextObject {
	return &TextObject{
		Color:    theme.TextColor(),
		Text:     text,
		FontSize: theme.TextSize(),
	}
}
