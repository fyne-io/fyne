// Package canvas contains all of the primitive CanvasObjects that make up a Fyne GUI
package canvas

import "image/color"

// Rectangle describes a coloured rectangle primitive in a Fyne canvas
type Rectangle struct {
	baseObject

	FillColor   color.RGBA // The rectangle fill colour
	StrokeColor color.RGBA // The rectangle stroke colour
	StrokeWidth float32    // The stroke width of the rectangle
}

// NewRectangle returns a new Rectangle instance
func NewRectangle(color color.RGBA) *Rectangle {
	return &Rectangle{
		FillColor: color,
	}
}
