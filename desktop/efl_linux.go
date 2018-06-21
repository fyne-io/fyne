// +build !ci

// +build !wayland

package desktop

// #cgo pkg-config: ecore-evas
// #include <Ecore_Evas.h>
import "C"

func oSEngineName() string {
	return "software_x11"
}

func oSWindowInit(w *window) {
}
