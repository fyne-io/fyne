// +build !ci

package desktop

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

import "log"

import "github.com/fyne-io/fyne"

const (
	oSEngineOther = "unknown"
)

type eFLDriver struct {
}

// NewEFLDriver sets up a new Driver instance implemented using the
// Enlightenment Foundation Libraries (EFL).
// It checks that there is a renderer available for the current operating system
// and causes a fatal error if not.
func NewEFLDriver() fyne.Driver {
	driver := new(eFLDriver)

	if oSEngineName() == oSEngineOther {
		log.Fatalln("Unsupported operating system")
	}

	return driver
}
