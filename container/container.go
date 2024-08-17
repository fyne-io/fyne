// Package container provides containers that are used to lay out and organise applications.
package container

import (
	"fyne.io/fyne/v2"
)

// New returns a new Container instance holding the specified CanvasObjects which will be laid out according to the specified Layout.
//
// Since: 2.0
func New(layout fyne.Layout, objects ...fyne.CanvasObject) *fyne.Container {
	var cntr *fyne.Container = new(fyne.Container)
	cntr.Layout = layout
	cntr.InitObjectMap()
	for _, obj := range objects {
		cntr.Add(obj)
	}
	return cntr
}

// NewWithoutLayout returns a new Container instance holding the specified CanvasObjects that are manually arranged.
//
// Since: 2.0
func NewWithoutLayout(objects ...fyne.CanvasObject) *fyne.Container {
	return New(nil, objects...)
}
