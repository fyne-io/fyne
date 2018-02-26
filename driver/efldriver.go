package driver

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

import "log"

const (
	oSEngineOther = "unknown"
)

type eFLDriver struct {
	running bool
}

func NewEFLDriver() Driver {
	driver := new(eFLDriver)

	if oSEngineName() == oSEngineOther {
		log.Fatalln("Unsupported operating system")
	}

	return driver
}

func (d *eFLDriver) Run() {
	d.running = true
	C.ecore_main_loop_begin()
}

func (d *eFLDriver) Quit() {
	C.ecore_main_loop_quit()
	d.running = false
}
