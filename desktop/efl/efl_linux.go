// +build !ci,!gl

// +build !wayland

package efl

// #cgo pkg-config: ecore-evas
// #include <Ecore_Evas.h>
import "C"

func oSEngineName() string {
	return "software_x11"
}

func oSWindowInit(w *window) {
}
