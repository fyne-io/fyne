package container // import "fyne.io/fyne/v2/container"

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/layout"
)

// NewAdaptiveGrid creates a new container with the specified objects and using the grid layout.
// When in a horizontal arrangement the rowcols parameter will specify the column count, when in vertical
// it will specify the rows. On mobile this will dynamically refresh when device is rotated.
//
// Since: 1.4
func NewAdaptiveGrid(rowcols int, objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewAdaptiveGridLayout(rowcols), objects...)
}

// NewBorder creates a new container with the specified objects and using the border layout.
// The top, bottom, left and right parameters specify the items that should be placed around edges.
// Nil can be used to an edge if it should not be filled.
// Passed objects not assigned to any edge (parameters 5 onwards) will be used to fill the space
// remaining in the middle.
// Parameters 6 onwards will be stacked over the middle content in the specified order as a Stack container.
//
// Since: 1.4
func NewBorder(top, bottom, left, right fyne.CanvasObject, objects ...fyne.CanvasObject) *fyne.Container {
	all := objects
	if top != nil {
		all = append(all, top)
	}
	if bottom != nil {
		all = append(all, bottom)
	}
	if left != nil {
		all = append(all, left)
	}
	if right != nil {
		all = append(all, right)
	}

	if len(objects) == 1 && objects[0] == nil {
		internal.LogHint("Border layout requires only 4 parameters, optional items cannot be nil")
		all = all[1:]
	}
	return New(layout.NewBorderLayout(top, bottom, left, right), all...)
}

// NewCenter creates a new container with the specified objects centered in the available space.
//
// Since: 1.4
func NewCenter(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewCenterLayout(), objects...)
}

// NewGridWithColumns creates a new container with the specified objects and using the grid layout with
// a specified number of columns. The number of rows will depend on how many children are in the container.
//
// Since: 1.4
func NewGridWithColumns(cols int, objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewGridLayoutWithColumns(cols), objects...)
}

// NewGridWithRows creates a new container with the specified objects and using the grid layout with
// a specified number of rows. The number of columns will depend on how many children are in the container.
//
// Since: 1.4
func NewGridWithRows(rows int, objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewGridLayoutWithRows(rows), objects...)
}

// NewGridWrap creates a new container with the specified objects and using the gridwrap layout.
// Every element will be resized to the size parameter and the content will arrange along a row and flow to a
// new row if the elements don't fit.
//
// Since: 1.4
func NewGridWrap(size fyne.Size, objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewGridWrapLayout(size), objects...)
}

// NewHBox creates a new container with the specified objects and using the HBox layout.
// The objects will be placed in the container from left to right and always displayed
// at their horizontal MinSize. Use a different layout if the objects are intended
// to be larger then their horizontal MinSize.
//
// Since: 1.4
func NewHBox(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewHBoxLayout(), objects...)
}

// NewMax creates a new container with the specified objects filling the available space.
//
// Since: 1.4
//
// Deprecated: Use container.NewStack() instead.
func NewMax(objects ...fyne.CanvasObject) *fyne.Container {
	return NewStack(objects...)
}

// NewPadded creates a new container with the specified objects inset by standard padding size.
//
// Since: 1.4
func NewPadded(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewPaddedLayout(), objects...)
}

// NewStack returns a new container that stacks objects on top of each other.
// Objects at the end of the container will be stacked on top of objects before.
// Having only a single object has no impact as CanvasObjects will
// fill the available space even without a Stack.
//
// Since: 2.4
func NewStack(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewStackLayout(), objects...)
}

// NewVBox creates a new container with the specified objects and using the VBox layout.
// The objects will be stacked in the container from top to bottom and always displayed
// at their vertical MinSize. Use a different layout if the objects are intended
// to be larger then their vertical MinSize.
//
// Since: 1.4
func NewVBox(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewVBoxLayout(), objects...)
}
