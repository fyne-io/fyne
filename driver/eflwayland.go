package driver

// #cgo pkg-config: ecore-evas
// #include <Ecore_Evas.h>
import "C"

func WaylandEngineName() string {
	return "wayland_shm"
}

func WaylandWindowInit(w *window) {
}
