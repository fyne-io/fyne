//go:build wasm

package glfw

import (
	"fyne.io/fyne/v2"
)

func hideVirtualKeyboard() {
	if d, ok := fyne.CurrentDevice().(*glDevice); ok {
		d.hideVirtualKeyboard()
	}
}

func handleKeyboard(obj fyne.Focusable) {
	isDisabled := false
	if disWid, ok := obj.(fyne.Disableable); ok {
		isDisabled = disWid.Disabled()
	}
	if obj != nil && !isDisabled {
		showVirtualKeyboard()
	} else {
		hideVirtualKeyboard()
	}
}

func showVirtualKeyboard() {
	if d, ok := fyne.CurrentDevice().(*glDevice); ok {
		d.showVirtualKeyboard()
	}
}
