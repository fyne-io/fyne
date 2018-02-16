package ui

import "image/color"

type RectangleObject struct {
	Color color.RGBA
	Size  Size
}

func (t *RectangleObject) CurrentSize() Size {
	return t.Size
}

func (t *RectangleObject) Resize(size Size) {
	t.Size = size
}

func (r *RectangleObject) MinSize() Size {
	return NewSize(1, 1)
}

func NewRectangle(color color.RGBA) *RectangleObject {
	return &RectangleObject{
		Color: color,
	}
}
