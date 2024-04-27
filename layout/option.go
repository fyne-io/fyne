package layout

import "fyne.io/fyne/v2"

type LayoutOption func(fyne.Layout)

// WithCustomPaddings is a LayoutOption that allows setting custom paddings for the layout.
// It takes a function that returns the top, bottom, left, and right paddings.
func WithCustomPaddings(fn func() (float32, float32, float32, float32)) LayoutOption {
	return func(layout fyne.Layout) {
		layout.SetPaddingFn(fn)
	}
}
