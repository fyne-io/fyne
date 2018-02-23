package driver

// #cgo pkg-config: ecore-evas
// #include <Ecore_Evas.h>
import "C"

func WaylandEngineName() string {
	return "wayland_shm"
}

func WaylandWindowInit(w *window) {
	// TODO find a way to init te ecore_wl2 window without custom code...
}
