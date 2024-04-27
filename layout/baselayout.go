package layout

import "fyne.io/fyne/v2/theme"

type BaseLayout struct {
	paddingFn func() (float32, float32, float32, float32)
}

func (b *BaseLayout) SetPaddingFn(fn func() (float32, float32, float32, float32)) {
	b.paddingFn = fn
}

func (b *BaseLayout) GetPaddings() (float32, float32, float32, float32) {
	if b.paddingFn == nil {
		padding := theme.Padding()
		return padding, padding, padding, padding
	}
	return b.paddingFn()
}
