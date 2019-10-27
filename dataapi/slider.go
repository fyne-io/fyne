package dataapi

import "fyne.io/fyne/widget"

// NewSlider returns a slider control that is 2-way bound to the dataItem
func NewSlider(data DataItem, min, max float64) *widget.Slider {

	// create the base widget
	w := widget.NewSlider(min, max)
	if ii, ok := data.(*Float); ok {
		w.Set(ii.value)
	}

	maskID := data.AddListener(func(i DataItem) {
		if ii, ok := i.(*Float); ok {
			w.Set(ii.value)
		}
	})

	if s, ok := data.(SettableFloat); ok {
		// The DataItem is settable, so add an onchange handler
		// to the widget to trigger the set function on the dataItem
		w.OnChanged = func(v float64) {
			s.SetFloat(v, maskID)
		}
	}
	return w
}
