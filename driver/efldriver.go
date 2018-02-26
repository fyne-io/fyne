package driver

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

import "log"

const (
	oSEngineOther = "unknown"
)

type EFLDriver struct {
	running bool
}

func NewEFLDriver() *EFLDriver {
	driver := new(EFLDriver)

	if oSEngineName() == oSEngineOther {
		log.Fatalln("Unsupported operating system")
	}

	return driver
}

func (d *EFLDriver) Run() {
	d.running = true
	C.ecore_main_loop_begin()
}

func (d *EFLDriver) Quit() {
	C.ecore_main_loop_quit()
	d.running = false
}
