// +build !ci,efl

package efl

// #cgo pkg-config: ecore-evas ecore-win32
// #include <Ecore_Evas.h>
// #include <Ecore_Win32.h>
import "C"

var (
	defaultCursor, entryCursor, hyperlinkCursor *C.Ecore_Win32_Cursor
)

func init() {
	defaultCursor = ecore_win32_cursor_shaped_new(ECORE_WIN32_CURSOR_SHAPE_ARROW)
	entryCursor = ecore_win32_cursor_shaped_new(ECORE_WIN32_CURSOR_SHAPE_I_BEAM)
	hyperlinkCursor = ecore_win32_cursor_shaped_new(ECORE_WIN32_CURSOR_SHAPE_HAND)
}

func oSEngineName() string {
	return "software_gdi"
}

func oSWindowInit(w *window) {
}

func setCursor(w *window, object fyne.CanvasObject) {
	win := C.ecore_evas_win32_window_get(w.ee)

	cursor := defaultCursor
	switch wid := object.(type) {
	case *widget.Entry:
		if !wid.ReadOnly {
			cursor = entryCursor
		}
	case *widget.Hyperlink:
		cursor = hyperlinkCursor
	}

	C.ecore_win32_window_cursor_set(win, C.uint(cursor))
}

func unsetCursor(w *window, object fyne.CanvasObject) {
	win := C.ecore_evas_win32_window_get(w.ee)
	C.ecore_win32_window_cursor_set(win, C.ECORE_X_CURSOR_X)
}
