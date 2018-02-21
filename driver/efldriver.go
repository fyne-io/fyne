package driver

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

import "log"

type EFLDriver struct {
}

func findEngineName() string {
	env := C.getenv(C.CString("WAYLAND_DISPLAY"))

	if env == nil {
		log.Println("Unable to connect to Wayland - attempting X")
		return X11EngineName()
	}

	return WaylandEngineName()
}

func (d *EFLDriver) Run() {
	C.ecore_main_loop_begin()
}
