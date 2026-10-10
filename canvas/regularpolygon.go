package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*RegularPolygon)(nil)

// RegularPolygon describes a colored regular polygon primitive in a Fyne canvas.
// The rendered portion will be in the center of the available space.
//
// Since: 2.8
type RegularPolygon struct {
	baseObject

	FillColor    color.Color // The polygon fill color
	StrokeColor  color.Color // The polygon stroke color
	StrokeWidth  float32     // The stroke width of the polygon
	CornerRadius float32     // The radius of the polygon corners
	Angle        float32     // Angle of polygon, in degrees (positive means clockwise, negative means counter-clockwise direction).
	Sides        uint        //	Amount of sides of polygon.
}

// Hide will set this polygon to not be visible
func (r *RegularPolygon) Hide() {
	r.baseObject.Hide()

	repaint(r)
}

// Move the polygon to a new position, relative to its parent / canvas
func (r *RegularPolygon) Move(pos fyne.Position) {
	if r.Position() == pos {
		return
	}

	r.baseObject.Move(pos)

	repaint(r)
}

// Refresh causes this polygon to be redrawn with its configured state.
func (r *RegularPolygon) Refresh() {
	Refresh(r)
}

// Resize on a polygon updates the new size of this object.
// If it has a stroke width this will cause it to Refresh.
func (r *RegularPolygon) Resize(s fyne.Size) {
	if s == r.Size() {
		return
	}

	r.baseObject.Resize(s)

	Refresh(r)
}

// NewRegularPolygon returns a new RegularPolygon instance
//
// Since: 2.8
func NewRegularPolygon(sides uint, c color.Color) *RegularPolygon {
	return &RegularPolygon{
		Sides:     sides,
		FillColor: c,
	}
}
