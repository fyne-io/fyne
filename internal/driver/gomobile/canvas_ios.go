// +build ios

package gomobile

import "fyne.io/fyne"

func devicePadding() (topLeft, bottomRight fyne.Size) {
	landscape := fyne.IsHorizontal(fyne.CurrentDevice().Orientation())

	if landscape {
		return fyne.NewSize(44, 0), fyne.NewSize(44, 21)
	} else {
		return fyne.NewSize(0, 44), fyne.NewSize(0, 34)
	}
}

func deviceScale() float32 {
	return 3 // TODO return 2 for < iPhone X, Xs but 3 for Plus/Max size
}
