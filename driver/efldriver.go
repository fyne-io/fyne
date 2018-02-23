package driver

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

import "log"
import "runtime"

type EFLDriver struct {
	running bool
}

func findEngineName() string {
	if runtime.GOOS == "darwin" {
		return CocoaEngineName()
	}

	env := C.getenv(C.CString("WAYLAND_DISPLAY"))

	if env != nil {
		log.Println("Wayland support is currently disabled - attempting XWayland")
	}

	return X11EngineName()
}

func (d *EFLDriver) Run() {
	d.running = true
	C.ecore_main_loop_begin()
}

func (d *EFLDriver) Quit() {
	C.ecore_main_loop_quit()
	d.running = false
}
