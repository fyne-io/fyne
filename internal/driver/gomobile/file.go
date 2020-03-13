package gomobile

import "fyne.io/fyne"

type hasPicker interface {
	ShowFileOpenPicker(callback func(string))
}

func ShowFileOpenPicker(callback func(string)) {
	if a, ok := fyne.CurrentApp().Driver().(*mobileDriver).app.(hasPicker); ok {
		a.ShowFileOpenPicker(callback)
	}
}
