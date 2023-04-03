// Package container provides containers that are used to lay out and organise applications.
package container

import (
	"fyne.io/fyne/v2"
)

// New returns a new Container instance holding the specified CanvasObjects which will be laid out according to the specified Layout.
//
// Since: 2.0
func New(layout fyne.Layout, objects ...fyne.CanvasObject) *fyne.Container {
	return &fyne.Container{Layout: layout, Objects: objects}
}

// NewWithoutLayout returns a new Container instance holding the specified CanvasObjects that are manually arranged.
//
// Since: 2.0
func NewWithoutLayout(objects ...fyne.CanvasObject) *fyne.Container {
	return &fyne.Container{Objects: objects}
}
