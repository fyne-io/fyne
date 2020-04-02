package binding

import (
	"fyne.io/fyne/widget"
)

func BindRadioChanged(radio *widget.Radio, data *StringBinding) {
	existing := radio.OnChanged
	radio.OnChanged = func(s string) {
		if data.GetString() != s {
			data.SetString(s)
		}
		if existing != nil {
			existing(s)
		}
	}
	data.AddStringListener(func(s string) {
		if radio.Selected != s {
			radio.SetSelected(s)
			if existing != nil {
				existing(s)
			}
		}
	})
}
