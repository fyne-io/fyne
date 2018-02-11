package efl

// #cgo pkg-config: ecore-evas ecore-wl2
// #cgo CFLAGS: -DEFL_BETA_API_SUPPORT=1
// #include <Ecore_Evas.h>
// #include <Ecore_Wl2.h>
import "C"

func WaylandEngineName() string {
	return "wayland_shm"
}

func WaylandWindowInit(w *window) {
	win := C.ecore_evas_wayland2_window_get(w.ee)
	C.ecore_wl2_window_type_set(win, C.ECORE_WL2_WINDOW_TYPE_TOPLEVEL)
}
