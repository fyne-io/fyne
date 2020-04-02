package binding

import (
	"fyne.io/fyne/widget"
)

func BindSelectChanged(selec *widget.Select, data *StringBinding) {
	existing := selec.OnChanged
	selec.OnChanged = func(s string) {
		if data.GetString() != s {
			data.SetString(s)
		}
		if existing != nil {
			existing(s)
		}
	}
	data.AddStringListener(func(s string) {
		if selec.Selected != s {
			selec.SetSelected(s)
			if existing != nil {
				existing(s)
			}
		}
	})
}
