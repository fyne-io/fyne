package container

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// Split defines a container whose size is split between two children.
type Split = widget.SplitContainer

// NewHSplit creates a horizontally arranged container with the specified leading and trailing elements.
// A vertical split bar that can be dragged will be added between the elements.
func NewHSplit(leading, trailing fyne.CanvasObject) *Split {
	return widget.NewHSplitContainer(leading, trailing)
}

// NewVSplit creates a vertically arranged container with the specified top and bottom elements.
// A horizontal split bar that can be dragged will be added between the elements.
func NewVSplit(top, bottom fyne.CanvasObject) *Split {
	return widget.NewVSplitContainer(top, bottom)
}

// TODO move the implementation into here in 2.0 when we delete the old API.
// we cannot do that right now due to Scroll dependency order.
