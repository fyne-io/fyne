package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Polygon)(nil)

// Polygon describes a colored regular polygon primitive in a Fyne canvas.
// The rendered portion will be in the center of the available space.
//
// Since: 2.7
type Polygon struct {
	baseObject

	FillColor    color.Color // The polygon fill color
	StrokeColor  color.Color // The polygon stroke color
	StrokeWidth  float32     // The stroke width of the polygon
	CornerRadius float32     // The radius of the polygon corners
	Angle        float32     // Angle of polygon, in degrees (positive means clockwise, negative means counter-clockwise direction).
	Sides        uint        //	Amount of sides of polygon.
}

// Hide will set this polygon to not be visible
func (r *Polygon) Hide() {
	r.baseObject.Hide()

	repaint(r)
}

// Move the polygon to a new position, relative to its parent / canvas
func (r *Polygon) Move(pos fyne.Position) {
	if r.Position() == pos {
		return
	}

	r.baseObject.Move(pos)

	repaint(r)
}

// Refresh causes this polygon to be redrawn with its configured state.
func (r *Polygon) Refresh() {
	Refresh(r)
}

// Resize on a polygon updates the new size of this object.
// If it has a stroke width this will cause it to Refresh.
func (r *Polygon) Resize(s fyne.Size) {
	if s == r.Size() {
		return
	}

	r.baseObject.Resize(s)
	if r.StrokeWidth == 0 {
		return
	}

	Refresh(r)
}

// NewPolygon returns a new Polygon instance
func NewPolygon(sides uint, color color.Color) *Polygon {
	return &Polygon{
		Sides:     sides,
		FillColor: color,
	}
}
