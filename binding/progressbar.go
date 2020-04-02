package binding

import (
	"fyne.io/fyne/widget"
)

func BindProgressBarValue(progressBar *widget.ProgressBar, data *Float64Binding) {
	data.AddFloat64Listener(func(f float64) {
		if progressBar.Value != f {
			progressBar.SetValue(f)
		}
	})
}
