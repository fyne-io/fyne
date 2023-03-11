package container // import "fyne.io/fyne/v2/container"

import (
	"fyne.io/fyne/v2"
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
// The top, bottom, left and right parameters specify the items that should be placed around edges,
// the remaining elements will be in the center. Nil can be used to an edge if it should not be filled.
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
// The objects will be placed in the container from left to right.
//
// Since: 1.4
func NewHBox(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewHBoxLayout(), objects...)
}

// NewMax creates a new container with the specified objects filling the available space.
//
// Since: 1.4
func NewMax(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewMaxLayout(), objects...)
}

// NewPadded creates a new container with the specified objects inset by standard padding size.
//
// Since: 1.4
func NewPadded(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewPaddedLayout(), objects...)
}

// NewVBox creates a new container with the specified objects and using the VBox layout.
// The objects will be stacked in the container from top to bottom.
//
// Since: 1.4
func NewVBox(objects ...fyne.CanvasObject) *fyne.Container {
	return New(layout.NewVBoxLayout(), objects...)
}
