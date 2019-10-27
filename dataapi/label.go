package dataapi

import (
	"fyne.io/fyne/widget"
)

// NewLabel returns a new *widget.Label that reacts to changes in that item
func NewLabel(data DataItem) *widget.Label {

	// create the base widget
	w := widget.NewLabel(data.String())
	// we are only interested in input changes
	data.AddListener(func(d DataItem) {
		w.SetText(d.String())
	})
	return w
}
