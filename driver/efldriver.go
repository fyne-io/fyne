package driver

// #cgo pkg-config: ecore
// #include <Ecore.h>
import "C"

type EFLDriver struct {
	running bool
}

func (d *EFLDriver) Run() {
	d.running = true
	C.ecore_main_loop_begin()
}

func (d *EFLDriver) Quit() {
	C.ecore_main_loop_quit()
	d.running = false
}
