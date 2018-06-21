// +build !ci

// +build freebsd openbsd netbsd

package desktop

// #cgo pkg-config: ecore-evas
// #include <Ecore_Evas.h>
import "C"

var engine string

func oSEngineName() string {
	return "software_x11"
}

func oSWindowInit(w *window) {
}
