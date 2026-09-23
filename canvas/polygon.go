package canvas

import "image/color"

// Polygon describes a colored regular polygon primitive in a Fyne canvas.
// The rendered portion will be in the center of the available space.
// Deprecated: Use [RegularPolygon] instead
//
// Since: 2.7
type Polygon = RegularPolygon

// NewPolygon returns a new Polygon instance
// Deprecated: Use [NewRegularPolygon] instead
//
// Since: 2.7
func NewPolygon(sides uint, c color.Color) *Polygon {
	return &RegularPolygon{
		Sides:     sides,
		FillColor: c,
	}
}
