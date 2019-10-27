package dataapi

import "fyne.io/fyne/widget"

// NewCheck returns a checkbox widget that is 2-way bound to the dataItem
func NewCheck(data DataItem, label string, callback func(checked bool)) *widget.Check {

	// create the base widget
	w := widget.NewCheck(label, callback)

	i := data.AddListener(func(b DataItem) {
		if bb, ok := b.(*Bool); ok {
			w.SetChecked(bb.value)
			callback(bb.value)
		}
	})

	if s, ok := data.(SettableBool); ok {
		// The DataItem is settable, so add an onchange handler
		// to the widget to trigger the set function on the dataItem
		w.OnChanged = func(b bool) {
			s.SetBool(b, i)
		}
	}
	return w
}
