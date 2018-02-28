package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/theme"

type TextObject struct {
	Color    color.RGBA
	Size     ui.Size
	Position ui.Position

	Text     string
	FontSize int
	Bold     bool
	Italic   bool

	minSize ui.Size
}

func (t *TextObject) CurrentSize() ui.Size {
	return t.Size
}

func (t *TextObject) Resize(size ui.Size) {
	t.Size = size
}

func (t *TextObject) CurrentPosition() ui.Position {
	return t.Position
}

func (t *TextObject) Move(pos ui.Position) {
	t.Position = pos
}

func (t *TextObject) SetMinSize(size ui.Size) {
	t.minSize = size
}

func (t *TextObject) MinSize() ui.Size {
	return t.minSize
}

func NewText(text string) *TextObject {
	return &TextObject{
		Color:    theme.TextColor(),
		Text:     text,
		FontSize: theme.TextSize(),
	}
}
