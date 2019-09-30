// +build !ios,!android

package gomobile

import "fyne.io/fyne"

func devicePadding() (topLeft, bottomRight fyne.Size) {
	return fyne.NewSize(0, 0), fyne.NewSize(0, 0)
}

func deviceScale() float32 {
	return 1
}
