package binding

import (
	"fyne.io/fyne/widget"
)

func BindSliderChanged(slider *widget.Slider, data *Float64Binding) {
	existing := slider.OnChanged
	slider.OnChanged = func(f float64) {
		if data.GetFloat64() != f {
			data.SetFloat64(f)
		}
		if existing != nil {
			existing(f)
		}
	}
	data.AddFloat64Listener(func(f float64) {
		if slider.Value != f {
			slider.SetValue(f)
			if existing != nil {
				existing(f)
			}
		}
	})
}
