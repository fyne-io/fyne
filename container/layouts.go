// Package container provides container widgets that are used to lay out and organise applications
package container // import "fyne.io/fyne/container"

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
)

// NewAdaptiveGrid creates a new container with the specified objects and using the grid layout.
// When in a horizontal arrangement the rowcols parameter will specify the column count, when in vertical
// it will specify the rows. On mobile this will dynamically refresh when device is rotated.
func NewAdaptiveGrid(rowcols int, objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(rowcols), objects...)
}

// NewBorder creates a new container with the specified objects and using the border layout.
// The top, bottom, left and right parameters specify the items that should be placed around edges,
// the remaining elements will be in the center. Nil can be used to an edge if it should not be filled.
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
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(top, bottom, left, right), all...)
}

// NewCenter creates a new container with the specified objects centered in the available space.
func NewCenter(objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewCenterLayout(), objects...)
}

// NewGridWithColumns creates a new container with the specified objects and using the grid layout with
// a specified number of columns. The number of rows will depend on how many children are in the container.
func NewGridWithColumns(cols int, objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(cols), objects...)
}

// NewGridWithRows creates a new container with the specified objects and using the grid layout with
// a specified number of columns. The number of columns will depend on how many children are in the container.
func NewGridWithRows(rows int, objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewGridLayoutWithRows(rows), objects...)
}

// NewGridWrap creates a new container with the specified objects and using the gridwrap layout.
// Every element will be resized to the size parameter and the content will arrange along a row and flow to a
// new row if the elements don't fit.
func NewGridWrap(size fyne.Size, objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewGridWrapLayout(size), objects...)
}

// NewHBox creates a new container with the specified objects and using the HBox layout.
// The objects will be placed in the container from left to right.
func NewHBox(objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), objects...)
}

// NewMax creates a new container with the specified objects filling the available space.
func NewMax(objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewMaxLayout(), objects...)
}

// NewPadded creates a new container with the specified objects inset by standard padding size.
func NewPadded(objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewPaddedLayout(), objects...)
}

// NewVBox creates a new container with the specified objects and using the VBox layout.
// The objects will be stacked in the container from top to bottom.
func NewVBox(objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), objects...)
}
