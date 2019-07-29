package driver

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// WalkObjectTree will walk an object tree executing the passed functions following the following
// rules:
// - beforeChildren is called for the start obj before traversing its children
// - the obj's children are traversed by calling walkObjects on each of them
// - afterChildren is called for the obj after traversing the obj's children
// The walk can be aborted by returning true in one of the functions:
// - if beforeChildren returns true, further traversing is stopped immediately, the after function
//   will not be called for the obj where the walk stopped, however, it will be called for all its
//   parents
func WalkObjectTree(
	obj fyne.CanvasObject,
	beforeChildren func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	afterChildren func(fyne.CanvasObject, fyne.CanvasObject),
) bool {
	clipSize := fyne.NewSize(math.MaxInt32, math.MaxInt32)
	return walkObjectTree(obj, nil, fyne.NewPos(0, 0), fyne.NewPos(0, 0), clipSize, beforeChildren, afterChildren)
}

func walkObjectTree(
	obj fyne.CanvasObject,
	parent fyne.CanvasObject,
	offset, clipPos fyne.Position,
	clipSize fyne.Size,
	beforeChildren func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	afterChildren func(fyne.CanvasObject, fyne.CanvasObject),
) bool {
	pos := obj.Position().Add(offset)

	var children []fyne.CanvasObject
	switch co := obj.(type) {
	case *fyne.Container:
		children = co.Objects
	case fyne.Widget:
		children = widget.Renderer(co).Objects()

		if scroll, ok := obj.(*widget.ScrollContainer); ok {
			clipPos = pos
			clipSize = scroll.Size()
		}
	}

	if beforeChildren != nil {
		if beforeChildren(obj, pos, clipPos, clipSize) {
			return true
		}
	}

	cancelled := false
	for _, child := range children {
		if walkObjectTree(child, obj, pos, clipPos, clipSize, beforeChildren, afterChildren) {
			cancelled = true
			break
		}
	}

	if afterChildren != nil {
		afterChildren(obj, parent)
	}
	return cancelled
}

// ScaleInt converts a fyne coordinate into a device specific output value using the canvas scale
func ScaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(math.Round(float64(v) * float64(c.Scale())))
	}
}

// UnscaleInt converts a device specific output value into a fyne coordinate using the canvas scale
func UnscaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) / c.Scale())
	}
}

// FindObjectAtPositionMatching is used to find an object in a canvas at the specified position.
// The matches function determines of the type of object that is found at this position is of a suitable type.
// The various canvas roots and overlays that can be searched are also passed in.
func FindObjectAtPositionMatching(mouse fyne.Position, matches func(object fyne.CanvasObject) bool,
	overlay fyne.CanvasObject, roots ...fyne.CanvasObject) (fyne.CanvasObject, int, int) {
	var found fyne.CanvasObject
	foundX, foundY := 0, 0

	findFunc := func(walked fyne.CanvasObject, pos fyne.Position, clipPos fyne.Position, clipSize fyne.Size) bool {
		if !walked.Visible() {
			return false
		}

		if mouse.X < clipPos.X || mouse.Y < clipPos.Y {
			return false
		}

		if mouse.X >= clipPos.X+clipSize.Width || mouse.Y >= clipPos.Y+clipSize.Height {
			return false
		}

		if mouse.X < pos.X || mouse.Y < pos.Y {
			return false
		}

		if mouse.X >= pos.X+walked.Size().Width || mouse.Y >= pos.Y+walked.Size().Height {
			return false
		}

		if matches(walked) {
			found = walked
			foundX, foundY = mouse.X-pos.X, mouse.Y-pos.Y
		}
		return false
	}

	if overlay != nil {
		WalkObjectTree(overlay, findFunc, nil)
	} else {
		for _, root := range roots {
			WalkObjectTree(root, findFunc, nil)
			if found != nil {
				break
			}
		}
	}

	return found, foundX, foundY
}
