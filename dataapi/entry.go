package dataapi

import "fyne.io/fyne/widget"

// NewEntry returns an entry widget that is 2-way bound to the dataItem
func NewEntry(data DataItem) *widget.Entry {

	// create the base widget
	w := widget.NewEntry()

	i := data.AddListener(func(d DataItem) {
		w.SetText(d.String())
	})

	if s, ok := data.(Settable); ok {
		// The DataItem is settable, so add an onchange handler
		// to the widget to trigger the set function on the dataItem
		w.OnChanged = func(txt string) {
			s.Set(txt, i)
		}
	}
	return w
}
