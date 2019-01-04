// +build !ci,efl

// Package efl provides a full Fyne render implementation using an EFL installation.
// This supports Windows, Mac OS X and Linux using the EFL (Evas) render pipeline.
package efl

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

import (
	"log"

	"fyne.io/fyne"
)

const (
	oSEngineOther = "unknown"
)

type eFLDriver struct {
}

func (d *eFLDriver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	for _, win := range windows {
		canv := win.canvas.(*eflCanvas)
		if canv.native[obj] != nil {
			return canv
		}
	}

	return nil
}

func (d *eFLDriver) Run() {
	runEFL()
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
