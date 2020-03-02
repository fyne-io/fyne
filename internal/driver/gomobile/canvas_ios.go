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

func (*device) SystemScale() float32 {
	if currentDPI >= 450 {
		return 3
	} else if currentDPI >= 340 {
		return 2.5
	}
	return 2
}
