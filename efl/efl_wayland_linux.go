// +build wayland

package efl

// #cgo pkg-config: ecore-evas ecore-wl2
// #include <Ecore_Evas.h>
// #include <Ecore_Wl2.h>
import "C"

import "log"
import "strings"

var engine string

func oSEngineName() string {
	if engine != "" {
		return engine
	}

	env := C.getenv(C.CString("WAYLAND_DISPLAY"))

	if env == nil {
		log.Println("Unable to connect to Wayland - attempting X")
		engine = "software_x11"
	} else {
		engine = "wayland_shm"
	}

	return engine
}

func waylandWindowInit(w *window) {
	win := C.ecore_evas_wayland2_window_get(w.ee)
	C.ecore_wl2_window_type_set(win, C.ECORE_WL2_WINDOW_TYPE_TOPLEVEL)
}

func x11WindowInit(w *window) {
}

func oSWindowInit(w *window) {
	if strings.HasPrefix(oSEngineName(), "wayland_") {
		waylandWindowInit(w)
	} else {
		x11WindowInit(w)
	}
}
