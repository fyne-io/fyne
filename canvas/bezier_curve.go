package canvas

import (
	"image/color"

	"fyne.io/fyne/v2"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*BezierCurve)(nil)

// BezierCurve describes a colored bezier curve primitive in a Fyne canvas.
type BezierCurve struct {
	baseObject

	StartPoint    fyne.Position   // Starting point of the bezier curve
	EndPoint      fyne.Position   // Ending point of the bezier curve
	ControlPoints []fyne.Position // Max 2 control points (1 for quadratic, 2 for cubic)
	StrokeColor   color.Color     // The bezier curve stroke color
	StrokeWidth   float32         // The stroke width of the bezier curve
}

// Hide will set this bezier curve to not be visible
func (r *BezierCurve) Hide() {
	r.baseObject.Hide()

	repaint(r)
}

// Move the bezier curve to a new position, relative to its parent / canvas
func (r *BezierCurve) Move(pos fyne.Position) {
	if r.Position() == pos {
		return
	}

	r.baseObject.Move(pos)

	repaint(r)
}

// Refresh causes this bezier curve to be redrawn with its configured state.
func (r *BezierCurve) Refresh() {
	Refresh(r)
}

// Resize on a bezier curve updates the new size of this object.
// If it has a stroke width this will cause it to Refresh.
func (r *BezierCurve) Resize(s fyne.Size) {
	if s == r.Size() {
		return
	}

	r.baseObject.Resize(s)
	if r.StrokeWidth == 0 {
		return
	}

	Refresh(r)
}

// NewQuadraticBezierCurve creates a quadratic bezier curve with 1 control points
func NewQuadraticBezierCurve(startPoint, controlPoint, endPoint fyne.Position, color color.Color) *BezierCurve {
	return &BezierCurve{
		StartPoint:    startPoint,
		EndPoint:      endPoint,
		ControlPoints: []fyne.Position{controlPoint},
		StrokeColor:   color,
		StrokeWidth:   1,
	}
}

// NewCubicBezierCurve creates a cubic bezier curve with 2 control points
func NewCubicBezierCurve(startPoint, controlPoint1, controlPoint2, endPoint fyne.Position, color color.Color) *BezierCurve {
	return &BezierCurve{
		StartPoint:    startPoint,
		EndPoint:      endPoint,
		ControlPoints: []fyne.Position{controlPoint1, controlPoint2},
		StrokeColor:   color,
		StrokeWidth:   1,
	}
}
