package binding

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func BindButtonTapped(button *widget.Button, data *BoolBinding) {
	existing := button.OnTapped
	button.OnTapped = func() {
		if data.GetBool() != true {
			data.SetBool(true)
		}
		if existing != nil {
			existing()
		}
	}
	data.AddBoolListener(func(b bool) {
		if b && existing != nil {
			existing()
		}
	})
}

func BindButtonText(button *widget.Button, data *StringBinding) {
	data.AddStringListener(func(s string) {
		if button.Text != s {
			button.SetText(s)
		}
	})
}

func BindButtonIcon(button *widget.Button, data *ResourceBinding) {
	data.AddResourceListener(func(r fyne.Resource) {
		if button.Icon != r {
			button.SetIcon(r)
		}
	})
}
