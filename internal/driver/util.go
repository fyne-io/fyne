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
// - if beforeChildren returns true, further traversing is stopped immediately, the after function
//   will not be called for the obj where the walk stopped, however, it will be called for all its
//   parents
func WalkObjectTree(
	obj fyne.CanvasObject,
	beforeChildren func(fyne.CanvasObject, fyne.Position) bool,
	afterChildren func(fyne.CanvasObject),
) bool {
	return walkObjectTree(obj, fyne.NewPos(0, 0), beforeChildren, afterChildren)
}

func walkObjectTree(
	obj fyne.CanvasObject,
	offset fyne.Position,
	beforeChildren func(fyne.CanvasObject, fyne.Position) bool,
	afterChildren func(fyne.CanvasObject),
) bool {
	var children []fyne.CanvasObject
	switch co := obj.(type) {
	case *fyne.Container:
		children = co.Objects
	case fyne.Widget:
		children = widget.Renderer(co).Objects()
	}

	pos := obj.Position().Add(offset)
	if beforeChildren != nil {
		if beforeChildren(obj, pos) {
			return true
		}
	}

	cancelled := false
	for _, child := range children {
		if walkObjectTree(child, pos, beforeChildren, afterChildren) {
			cancelled = true
			break
		}
	}

	if afterChildren != nil {
		afterChildren(obj)
	}
	return cancelled
}
