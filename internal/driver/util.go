package driver

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// WalkObjectTree will walk an object tree executing the passed functions following the following
// rules:
// - beforeChildren is called for the start obj before traversing its children
// - the obj's children are traversed by calling walkObjects on each of them
// - afterChildren is called for the obj after traversing the obj's children
// The walk can be aborted by returning true in one of the functions:
// - if beforeChildren returns true, further traversing is stopped immediatly, the after function
//   will not be called for the obj where the walk stopped, however, it will be called for all its
//   parents
// - if a walk has been stopped, the after function is called with the third argument set to true
func WalkObjectTree(
	obj fyne.CanvasObject,
	pos fyne.Position,
	beforeChildren func(fyne.CanvasObject, fyne.Position) bool,
	afterChildren func(fyne.CanvasObject, fyne.Position, bool),
) bool {
	var children []fyne.CanvasObject
	var offset fyne.Position
	switch co := obj.(type) {
	case *fyne.Container:
		offset = co.Position().Add(pos)
		children = co.Objects
	case fyne.Widget:
		offset = co.Position().Add(pos)
		children = widget.Renderer(co).Objects()
	default:
		offset = pos
	}

	if beforeChildren != nil {
		if beforeChildren(obj, offset) {
			return true
		}
	}

	cancelled := false
	for _, child := range children {
		if WalkObjectTree(child, offset, beforeChildren, afterChildren) {
			cancelled = true
			break
		}
	}

	if afterChildren != nil {
		afterChildren(obj, offset, cancelled)
	}
	return cancelled
}
