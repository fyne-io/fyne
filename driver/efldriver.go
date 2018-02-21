package driver

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

import "log"

type EFLDriver struct {
	running bool
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
	d.running = true
	C.ecore_main_loop_begin()
}

func (d *EFLDriver) Quit() {
	C.ecore_main_loop_quit()
	d.running = false
}
