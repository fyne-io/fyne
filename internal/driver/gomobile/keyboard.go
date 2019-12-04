package gomobile

import (
	"fyne.io/fyne"
)

type hasKeyboard interface {
	ShowVirtualKeyboard()
	HideVirtualKeyboard()
}

func showVirtualKeyboard() {
	if adv, ok := fyne.CurrentApp().Driver().(*mobileDriver).app.(hasKeyboard); ok {
		adv.ShowVirtualKeyboard()
	}
}

func hideVirtualKeyboard() {
	if adv, ok := fyne.CurrentApp().Driver().(*mobileDriver).app.(hasKeyboard); ok {
		adv.HideVirtualKeyboard()
	}
}
