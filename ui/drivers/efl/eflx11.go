package efl

// #cgo pkg-config: ecore-evas ecore-x
// #include <Ecore_Evas.h>
// #include <Ecore_X.h>
import "C"

func X11EngineName() string {
	return "software_x11"
}

func X11WindowInit(w *window) {
}
