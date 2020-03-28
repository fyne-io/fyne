package gomobile

import "fyne.io/fyne"

type hasPicker interface {
	ShowFileOpenPicker(callback func(string))
}

// ShowFileOpenPicker loads the native file open dialog and returns the chosen file path via the callback func.
func ShowFileOpenPicker(callback func(string)) {
	if a, ok := fyne.CurrentApp().Driver().(*mobileDriver).app.(hasPicker); ok {
		a.ShowFileOpenPicker(callback)
	}
}
