package ui

import "image/color"

type RectangleObject struct {
	Color color.RGBA
}

func (r RectangleObject) MinSize() (int, int) {
	return 1, 1
}

func NewRectangle(color color.RGBA) *RectangleObject {
	return &RectangleObject{
		Color: color,
	}
}
