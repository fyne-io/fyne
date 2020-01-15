// Package widget defines the UI widgets within the Fyne toolkit
package widget // import "fyne.io/fyne/widget"

import "fyne.io/fyne/dataapi"

// BaseWidget provides a helper that handles basic widget behaviours.
type DataWidget struct {
	BaseWidget

	listenerID int
}

// Bind will bind this widget to the given DataItem
func (w *DataWidget) Bind(data dataapi.DataItem) {
	w.listenerID = data.AddListener(w.SetFromData)
}

// SetFromData is called on the widget whenever the DataItem changes
func (w *DataWidget) SetFromData(data dataapi.DataItem) {
	// Do not update if hidden
	if !w.Hidden {
		println("widget update from data")
	}
}
