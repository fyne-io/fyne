// Package efl provides an EFL render provider for the Fyne implementation
package efl

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

// NewEFLDriver returns a new Driver instance implemented using the
// Enlightenment Foundation Libraries (EFL)
func NewEFLDriver() *eFLDriver {
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
