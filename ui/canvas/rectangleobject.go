package canvas

import "image/color"

import "github.com/fyne-io/fyne/ui"

type RectangleObject struct {
	Size     ui.Size
	Position ui.Position

	Color color.RGBA
}

func (r *RectangleObject) CurrentSize() ui.Size {
	return r.Size
}

func (r *RectangleObject) Resize(size ui.Size) {
	r.Size = size
}

func (r *RectangleObject) CurrentPosition() ui.Position {
	return r.Position
}

func (r *RectangleObject) Move(pos ui.Position) {
	r.Position = pos
}

func (r *RectangleObject) MinSize() ui.Size {
	return ui.NewSize(1, 1)
}

func NewRectangle(color color.RGBA) *RectangleObject {
	return &RectangleObject{
		Color: color,
	}
}
