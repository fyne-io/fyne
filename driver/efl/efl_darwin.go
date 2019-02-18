// +build !ci,efl

package efl

// #cgo pkg-config: ecore-evas ecore-cocoa
// #include <Ecore_Evas.h>
// #include <Ecore_Cocoa.h>
import "C"

import "fyne.io/fyne"
import "fyne.io/fyne/widget"

func oSEngineName() string {
	return "opengl_cocoa"
}

func oSWindowInit(w *window) {
}

func setCursor(w *window, object fyne.CanvasObject) {
	win := C.ecore_evas_cocoa_window_get(w.ee)

	cursor := C.ECORE_COCOA_CURSOR_ARROW
	switch wid := object.(type) {
	case *widget.Entry:
		if !wid.ReadOnly {
			cursor = C.ECORE_COCOA_CURSOR_IBEAM
		}
	case *widget.Hyperlink:
		cursor = C.ECORE_COCOA_CURSOR_POINTING_HAND
	}

	C.ecore_cocoa_window_cursor_set(win, C.Ecore_Cocoa_Cursor(cursor))
}

func unsetCursor(w *window, object fyne.CanvasObject) {
	win := C.ecore_evas_cocoa_window_get(w.ee)
	C.ecore_cocoa_window_cursor_set(win, C.ECORE_COCOA_CURSOR_ARROW)
}
