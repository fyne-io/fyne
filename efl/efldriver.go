// Package efl provides an EFL render provider for the Fyne implementation
package efl

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

import "log"

import "github.com/fyne-io/fyne/ui"

const (
	oSEngineOther = "unknown"
)

type eFLDriver struct {
}

// NewEFLDriver returns a new Driver instance implemented using the
// Enlightenment Foundation Libraries (EFL)
func NewEFLDriver() *eFLDriver {
	driver := new(eFLDriver)

	if oSEngineName() == oSEngineOther {
		log.Fatalln("Unsupported operating system")
	}

	ui.SetDriver(driver)
	return driver
}
