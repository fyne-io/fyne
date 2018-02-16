package ui

import "image/color"

type RectangleObject struct {
	Size     Size
	Position Position

	Color color.RGBA
}

func (r *RectangleObject) CurrentSize() Size {
	return r.Size
}

func (r *RectangleObject) Resize(size Size) {
	r.Size = size
}

func (r *RectangleObject) CurrentPosition() Position {
	return r.Position
}

func (r *RectangleObject) Move(pos Position) {
	r.Position = pos
}

func (r *RectangleObject) MinSize() Size {
	return NewSize(1, 1)
}

func NewRectangle(color color.RGBA) *RectangleObject {
	return &RectangleObject{
		Color: color,
	}
}
