// +build !ci,!gl

// +build freebsd openbsd netbsd

package efl

// #cgo pkg-config: ecore-evas
// #include <Ecore_Evas.h>
import "C"

func oSEngineName() string {
	return "software_x11"
}

func oSWindowInit(w *window) {
}
