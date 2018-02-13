package ui

import "image/color"

import "github.com/fyne-io/fyne/ui/theme"

type TextObject struct {
	Color color.RGBA

	Text     string
	FontSize int
	Bold     bool
	Italic   bool
}

func (t TextObject) MinSize() (int, int) {
	return 1, 1 // TODO calcualte from font/string
}

func NewText(text string) *TextObject {
	return &TextObject{
		Color:    theme.TextColor(),
		Text:     text,
		FontSize: theme.TextSize(),
	}
}
