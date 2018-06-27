// +build !ci

package desktop

// #cgo pkg-config: eina evas ecore-evas ecore-input
// #cgo CFLAGS: -DEFL_BETA_API_SUPPORT=master-compatibility-hack
// #include <Eina.h>
// #include <Evas.h>
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <Ecore_Input.h>
//
// void onKeyDown_cgo(Ecore_Window, void *);
// void onExit_cgo(Ecore_Event_Signal_Exit *);
import "C"

import "runtime"
import "time"
import "unsafe"

import "github.com/fyne-io/fyne"

// channel to signal quitting
var quit = make(chan bool, 1)

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

// initEFL runs our mainthread loop to execute UI functions for EFL
func initEFL() {
	C.ecore_event_handler_add(C.ECORE_EVENT_SIGNAL_EXIT, (C.Ecore_Event_Handler_Cb)(unsafe.Pointer(C.onExit_cgo)), nil)
	C.ecore_event_handler_add(C.ECORE_EVENT_KEY_DOWN, (C.Ecore_Event_Handler_Cb)(unsafe.Pointer(C.onKeyDown_cgo)), nil)

	tick := time.NewTicker(time.Second / 120)
	for {
		select {
		case <-quit:
			close(quit)
			tick.Stop()
			return
		case <-tick.C:
			drawDirty()
			C.ecore_main_loop_iterate()
		}
	}
}

// DoQuit will cause the driver's Quit method to be called to terminate the app
//export DoQuit
func DoQuit() {
	fyne.GetDriver().Quit()
}

// Quit will cause the render loop to end and the application to exit
func (d *eFLDriver) Quit() {
	quit <- true
}

func drawDirty() {
	for _, canvas := range canvases {
		canvas.fitContent()
		for obj := range canvas.dirty {
			delete(canvas.dirty, obj)

			canvas.doRefresh(obj)
		}
	}
}

// do runs f on the main thread.
func queueRender(c *eflCanvas, co fyne.CanvasObject) {
	c.dirty[co] = true
}
