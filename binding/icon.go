package binding

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func BindIconResource(icon *widget.Icon, data *ResourceBinding) {
	data.AddResourceListener(func(r fyne.Resource) {
		if icon.Resource != r {
			icon.SetResource(r)
		}
	})
}
