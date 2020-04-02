package binding

import (
	"fyne.io/fyne/widget"
)

func BindLabelText(label *widget.Label, data *StringBinding) {
	data.AddStringListener(func(s string) {
		if label.Text != s {
			label.SetText(s)
		}
	})
}
