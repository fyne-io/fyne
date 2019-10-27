package dataapi

import "fyne.io/fyne/widget"

// NewRadio returns a radio button widget that is 2-way bound to the dataItem
func NewRadio(data DataItem, labels []string, callback func(string)) *widget.Radio {

	// safely get the string value based on a number
	getValue := func(i int) string {
		if i < 0 || i >= len(labels) {
			return ""
		}
		return labels[i]
	}

	getIndex := func(lbl string) int {
		for k, v := range labels {
			if v == lbl {
				return k
			}
		}
		return 0
	}

	// create the base widget
	w := widget.NewRadio(labels, callback)
	if ii, ok := data.(*Int); ok {
		w.SetSelected(getValue(ii.value))
	}

	maskID := data.AddListener(func(i DataItem) {
		if ii, ok := i.(*Int); ok {
			lbl := getValue(ii.value)
			w.SetSelected(lbl)
			callback(lbl)
		}
	})

	if s, ok := data.(SettableInt); ok {
		// The DataItem is settable, so add an onchange handler
		// to the widget to trigger the set function on the dataItem
		w.OnChanged = func(value string) {
			ii := getIndex(value)
			s.SetInt(ii, maskID)
		}
	}
	return w
}
