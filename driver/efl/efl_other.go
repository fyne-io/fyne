// +build !ci,efl

// +build !linux,!darwin,!windows,!freebsd,!openbsd,!netbsd

package efl

import "fyne.io/fyne"

func oSEngineName() string {
	return oSEngineOther
}

func oSWindowInit(w *window) {
}

func setCursor(w *window, object fyne.CanvasObject) {
}

func unsetCursor(w *window, object fyne.CanvasObject) {
}
