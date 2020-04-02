package binding

import (
	"fyne.io/fyne/widget"
)

func BindCheckChanged(check *widget.Check, data *BoolBinding) {
	existing := check.OnChanged
	check.OnChanged = func(b bool) {
		if data.GetBool() != b {
			data.SetBool(b)
		}
		if existing != nil {
			existing(b)
		}
	}
	data.AddBoolListener(func(b bool) {
		if check.Checked != b {
			check.SetChecked(b)
			if existing != nil {
				existing(b)
			}
		}
	})
}

func BindCheckText(check *widget.Check, data *StringBinding) {
	data.AddStringListener(func(s string) {
		if check.Text != s {
			check.SetText(s)
		}
	})
}
