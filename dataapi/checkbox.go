package dataapi

import (
	"strings"

	"fyne.io/fyne/widget"
)

// NewCheck returns a checkbox widget that is 2-way bound to the dataItem
func NewCheck(data DataItem, label string, callback func(checked bool)) *widget.Check {

	// create the base widget
	w := widget.NewCheck(label, callback)
	if bb, ok := data.(*Bool); ok {
		w.SetChecked(bb.value)
	}

	maskID := data.AddListener(func(b DataItem) {
		if bb, ok := b.(*Bool); ok {
			w.SetChecked(bb.value)
			callback(bb.value)
			return
		}
		// Bound to a number, so set the field based on 0 = false, else true
		if bb, ok := b.(*Int); ok {
			value := bb.value > 0
			w.SetChecked(value)
			callback(value)
			return
		}
		// dataItem is not a bool, get the value and test
		// vs the string value
		switch strings.ToLower(b.String()) {
		case "true", "1", "ok", "yes", "on":
			w.SetChecked(true)
			callback(true)
		default:
			w.SetChecked(false)
			callback(false)
		}
	})

	if s, ok := data.(SettableBool); ok {
		// The DataItem is settable, so add an onchange handler
		// to the widget to trigger the set function on the dataItem
		w.OnChanged = func(b bool) {
			s.SetBool(b, maskID)
		}
		return w
	}

	if s, ok := data.(SettableInt); ok {
		// The DataItem is settable as a number, so add an onchange handler
		// to the widget to trigger the set function on the dataItem
		w.OnChanged = func(b bool) {
			value := 0
			if b {
				value = 1
			}
			s.SetInt(value, maskID)
		}
		return w
	}

	if s, ok := data.(Settable); ok {
		// The DataItem is settable as a string, so add an onchange handler
		// to the widget to trigger the set function on the dataItem
		// setting the string to true/false
		w.OnChanged = func(b bool) {
			value := "false"
			if b {
				value = "true"
			}
			s.Set(value, maskID)
		}
	}

	return w
}
