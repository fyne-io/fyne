package canvas

import "image/color"

// NewSquare returns a new Rectangle instance that has a square aspect ratio.
//
// Since: 2.7
func NewSquare(c color.Color) *Rectangle {
	return &Rectangle{
		Aspect:    1,
		FillColor: c,
	}
}
