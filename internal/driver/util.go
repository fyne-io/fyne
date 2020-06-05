package driver

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"
)

// WalkVisibleObjectTree will walk an object tree for all visible objects executing the passed functions following
// the following rules:
// - beforeChildren is called for the start obj before traversing its children
// - the obj's children are traversed by calling walkObjects on each of the visible items
// - afterChildren is called for the obj after traversing the obj's children
// The walk can be aborted by returning true in one of the functions:
// - if beforeChildren returns true, further traversing is stopped immediately, the after function
//   will not be called for the obj where the walk stopped, however, it will be called for all its
//   parents
func WalkVisibleObjectTree(
	obj fyne.CanvasObject,
	beforeChildren func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	afterChildren func(fyne.CanvasObject, fyne.CanvasObject),
) bool {
	clipSize := fyne.NewSize(math.MaxInt32, math.MaxInt32)
	return walkObjectTree(obj, nil, fyne.NewPos(0, 0), fyne.NewPos(0, 0), clipSize, beforeChildren, afterChildren, true)
}

// WalkCompleteObjectTree will walk an object tree for all objects (ignoring visible state) executing the passed
// functions following the following rules:
// - beforeChildren is called for the start obj before traversing its children
// - the obj's children are traversed by calling walkObjects on each of the items
// - afterChildren is called for the obj after traversing the obj's children
// The walk can be aborted by returning true in one of the functions:
// - if beforeChildren returns true, further traversing is stopped immediately, the after function
//   will not be called for the obj where the walk stopped, however, it will be called for all its
//   parents
func WalkCompleteObjectTree(
	obj fyne.CanvasObject,
	beforeChildren func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	afterChildren func(fyne.CanvasObject, fyne.CanvasObject),
) bool {
	clipSize := fyne.NewSize(math.MaxInt32, math.MaxInt32)
	return walkObjectTree(obj, nil, fyne.NewPos(0, 0), fyne.NewPos(0, 0), clipSize, beforeChildren, afterChildren, false)
}

func walkObjectTree(
	obj fyne.CanvasObject,
	parent fyne.CanvasObject,
	offset, clipPos fyne.Position,
	clipSize fyne.Size,
	beforeChildren func(fyne.CanvasObject, fyne.Position, fyne.Position, fyne.Size) bool,
	afterChildren func(fyne.CanvasObject, fyne.CanvasObject),
	requireVisible bool,
) bool {
	if requireVisible && !obj.Visible() {
		return false
	}
	pos := obj.Position().Add(offset)

	var children []fyne.CanvasObject
	switch co := obj.(type) {
	case *fyne.Container:
		children = co.Objects
	case fyne.Widget:
		children = cache.Renderer(co).Objects()

		if _, ok := obj.(fyne.Scrollable); ok {
			clipPos = pos
			clipSize = obj.Size()
		}
	}

	if beforeChildren != nil {
		if beforeChildren(obj, pos, clipPos, clipSize) {
			return true
		}
	}

	cancelled := false
	for _, child := range children {
		if walkObjectTree(child, obj, pos, clipPos, clipSize, beforeChildren, afterChildren, requireVisible) {
			cancelled = true
			break
		}
	}

	if afterChildren != nil {
		afterChildren(obj, parent)
	}
	return cancelled
}

// FindObjectAtPositionMatching is used to find an object in a canvas at the specified position.
// The matches function determines of the type of object that is found at this position is of a suitable type.
// The various canvas roots and overlays that can be searched are also passed in.
func FindObjectAtPositionMatching(mouse fyne.Position, matches func(object fyne.CanvasObject) bool, overlay fyne.CanvasObject, roots ...fyne.CanvasObject) (fyne.CanvasObject, fyne.Position, int) {
	var found fyne.CanvasObject
	var foundPos fyne.Position

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
			foundPos = fyne.NewPos(mouse.X-pos.X, mouse.Y-pos.Y)
		}
		return false
	}

	layer := 0
	if overlay != nil {
		WalkVisibleObjectTree(overlay, findFunc, nil)
	} else {
		for _, root := range roots {
			layer++
			if root == nil {
				continue
			}
			WalkVisibleObjectTree(root, findFunc, nil)
			if found != nil {
				break
			}
		}
	}

	return found, foundPos, layer
}

// AbsolutePositionForObject returns the absolute position of an object in a set of object trees.
// If the object is not part of any of the trees, the position (0,0) is returned.
func AbsolutePositionForObject(object fyne.CanvasObject, trees []fyne.CanvasObject) fyne.Position {
	var pos fyne.Position
	findPos := func(o fyne.CanvasObject, p fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if o == object {
			pos = p
			return true
		}
		return false
	}
	for _, tree := range trees {
		if WalkVisibleObjectTree(tree, findPos, nil) {
			break
		}
	}
	return pos
}
