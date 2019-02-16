// +build !ci,efl

// +build linux openbsd freebsd netbsd

// +build !wayland

package efl

// #cgo pkg-config: ecore-evas ecore-x
// #include <Ecore_Evas.h>
// #include <Ecore_X.h>
// #include <Ecore_X_Cursor.h>
import "C"
import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func oSEngineName() string {
	return "software_x11"
}

func oSWindowInit(w *window) {
}

func setCursor(w *window, object fyne.CanvasObject) {
	win := C.ecore_evas_software_x11_window_get(w.ee)

	cursor := C.ECORE_X_CURSOR_ARROW
	switch wid := object.(type) {
	case *widget.Entry:
		if !wid.ReadOnly {
			cursor = C.ECORE_X_CURSOR_XTERM
		}
	case *widget.Hyperlink:
		cursor = C.ECORE_X_CURSOR_HAND2
	}

	C.ecore_x_window_cursor_set(win, C.ecore_x_cursor_shape_get(C.int(cursor)))
}

func unsetCursor(w *window, object fyne.CanvasObject) {
	win := C.ecore_evas_software_x11_window_get(w.ee)
	C.ecore_x_window_cursor_set(win, C.ECORE_X_CURSOR_ARROW)
}
