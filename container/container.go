// Package container provides container widgets that are used to lay out and organise applications
package container

import (
	"fyne.io/fyne"
)

// New returns a new Container instance holding the specified CanvasObjects which will be laid out according to the specified Layout.
func New(layout fyne.Layout, objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout, objects...)
}

// NewManual returns a new Container instance holding the specified CanvasObjects that are manually arranged.
func NewManual(objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithoutLayout(objects...)
}
