// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux,!android

package app

/*
Simple on-screen app debugging for X11. Not an officially supported
development target for apps, as screens with mice are very different
than screens with touch panels.
*/

/*
#cgo LDFLAGS: -lEGL -lGLESv2 -lX11

void createWindow(void);
void processEvents(void);
void swapBuffers(void);
*/
import "C"
import (
	"runtime"
	"time"

	"github.com/fyne-io/mobile/event/lifecycle"
	"github.com/fyne-io/mobile/event/paint"
	"github.com/fyne-io/mobile/event/size"
	"github.com/fyne-io/mobile/event/touch"
	"github.com/fyne-io/mobile/geom"
)

func init() {
	theApp.registerGLViewportFilter()
}

func main(f func(App)) {
	runtime.LockOSThread()

	workAvailable := theApp.worker.WorkAvailable()
	heartbeat := time.NewTicker(time.Second / 60)

	C.createWindow()

	// TODO: send lifecycle events when e.g. the X11 window is iconified or moved off-screen.
	theApp.sendLifecycle(lifecycle.StageFocused)

	// TODO: translate X11 expose events to shiny paint events, instead of
	// sending this synthetic paint event as a hack.
	theApp.eventsIn <- paint.Event{}

	donec := make(chan struct{})
	go func() {
		f(theApp)
		close(donec)
	}()

	// TODO: can we get the actual vsync signal?
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	var tc <-chan time.Time

	for {
		select {
		case <-donec:
			return
		case <-heartbeat.C:
			C.processEvents()
		case <-workAvailable:
			theApp.worker.DoWork()
		case <-theApp.publish:
			C.swapBuffers()
			tc = ticker.C
		case <-tc:
			tc = nil
			theApp.publishResult <- PublishResult{}
		}
	}
}

//export onResize
func onResize(w, h int) {
	// TODO(nigeltao): don't assume 72 DPI. DisplayWidth and DisplayWidthMM
	// is probably the best place to start looking.
	pixelsPerPt := float32(1)
	theApp.eventsIn <- size.Event{
		WidthPx:     w,
		HeightPx:    h,
		WidthPt:     geom.Pt(w),
		HeightPt:    geom.Pt(h),
		PixelsPerPt: pixelsPerPt,
		Orientation: screenOrientation(w, h),
	}
}

func sendTouch(t touch.Type, x, y float32) {
	theApp.eventsIn <- touch.Event{
		X:        x,
		Y:        y,
		Sequence: 0, // TODO: button??
		Type:     t,
	}
}

//export onTouchBegin
func onTouchBegin(x, y float32) { sendTouch(touch.TypeBegin, x, y) }

//export onTouchMove
func onTouchMove(x, y float32) { sendTouch(touch.TypeMove, x, y) }

//export onTouchEnd
func onTouchEnd(x, y float32) { sendTouch(touch.TypeEnd, x, y) }

var stopped bool

//export onStop
func onStop() {
	if stopped {
		return
	}
	stopped = true
	theApp.sendLifecycle(lifecycle.StageDead)
	theApp.eventsIn <- stopPumping{}
}

// driverShowVirtualKeyboard does nothing on desktop
func driverShowVirtualKeyboard(KeyboardType) {
}

// driverHideVirtualKeyboard does nothing on desktop
func driverHideVirtualKeyboard() {
}

// driverShowFileOpenPicker does nothing on desktop
func driverShowFileOpenPicker(func(string, func()), *FileFilter) {
}

// driverShowFileSavePicker does nothing on desktop
func driverShowFileSavePicker(func(string, func()), *FileFilter, string) {
}
