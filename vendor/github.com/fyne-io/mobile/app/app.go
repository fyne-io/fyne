// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin windows

package app

import (
	"github.com/fyne-io/mobile/event/lifecycle"
	"github.com/fyne-io/mobile/event/size"
	"github.com/fyne-io/mobile/gl"
	_ "github.com/fyne-io/mobile/internal/mobileinit"
)

// Main is called by the main.main function to run the mobile application.
//
// It calls f on the App, in a separate goroutine, as some OS-specific
// libraries require being on 'the main thread'.
func Main(f func(App)) {
	main(f)
}

// App is how a GUI mobile application interacts with the OS.
type App interface {
	// Events returns the events channel. It carries events from the system to
	// the app. The type of such events include:
	//  - lifecycle.Event
	//  - mouse.Event
	//  - paint.Event
	//  - size.Event
	//  - touch.Event
	// from the golang.org/x/mobile/event/etc packages. Other packages may
	// define other event types that are carried on this channel.
	Events() <-chan interface{}

	// Send sends an event on the events channel. It does not block.
	Send(event interface{})

	// Publish flushes any pending drawing commands, such as OpenGL calls, and
	// swaps the back buffer to the screen.
	Publish() PublishResult

	// TODO: replace filters (and the Events channel) with a NextEvent method?

	// Filter calls each registered event filter function in sequence.
	Filter(event interface{}) interface{}

	// RegisterFilter registers a event filter function to be called by Filter. The
	// function can return a different event, or return nil to consume the event,
	// but the function can also return its argument unchanged, where its purpose
	// is to trigger a side effect rather than modify the event.
	RegisterFilter(f func(interface{}) interface{})

	ShowVirtualKeyboard(KeyboardType)
	HideVirtualKeyboard()
	ShowFileOpenPicker(func(string, func()), *FileFilter)
}

type FileFilter struct {
	Extensions []string
	MimeTypes []string
}

// PublishResult is the result of an App.Publish call.
type PublishResult struct {
	// BackBufferPreserved is whether the contents of the back buffer was
	// preserved. If false, the contents are undefined.
	BackBufferPreserved bool
}

var theApp = &app{
	eventsOut:      make(chan interface{}),
	lifecycleStage: lifecycle.StageDead,
	publish:        make(chan struct{}),
	publishResult:  make(chan PublishResult),
}

func init() {
	theApp.eventsIn = pump(theApp.eventsOut)
	theApp.glctx, theApp.worker = gl.NewContext()
}

func (a *app) sendLifecycle(to lifecycle.Stage) {
	if a.lifecycleStage == to {
		return
	}
	a.eventsIn <- lifecycle.Event{
		From:        a.lifecycleStage,
		To:          to,
		DrawContext: a.glctx,
	}
	a.lifecycleStage = to
}

type app struct {
	filters []func(interface{}) interface{}

	eventsOut      chan interface{}
	eventsIn       chan interface{}
	lifecycleStage lifecycle.Stage
	publish        chan struct{}
	publishResult  chan PublishResult

	glctx  gl.Context
	worker gl.Worker
}

func (a *app) Events() <-chan interface{} {
	return a.eventsOut
}

func (a *app) Send(event interface{}) {
	a.eventsIn <- event
}

func (a *app) Publish() PublishResult {
	// gl.Flush is a lightweight (on modern GL drivers) blocking call
	// that ensures all GL functions pending in the gl package have
	// been passed onto the GL driver before the app package attempts
	// to swap the screen buffer.
	//
	// This enforces that the final receive (for this paint cycle) on
	// gl.WorkAvailable happens before the send on endPaint.
	a.glctx.Flush()
	a.publish <- struct{}{}
	return <-a.publishResult
}

func (a *app) Filter(event interface{}) interface{} {
	for _, f := range a.filters {
		event = f(event)
	}
	return event
}

func (a *app) RegisterFilter(f func(interface{}) interface{}) {
	a.filters = append(a.filters, f)
}

func (a *app) ShowVirtualKeyboard(keyboard KeyboardType) {
	driverShowVirtualKeyboard(keyboard)
}

func (a *app) HideVirtualKeyboard() {
	driverHideVirtualKeyboard()
}

func (a *app) ShowFileOpenPicker(callback func(string, func()), filter *FileFilter) {
	driverShowFileOpenPicker(callback, filter)
}

type stopPumping struct{}

// pump returns a channel src such that sending on src will eventually send on
// dst, in order, but that src will always be ready to send/receive soon, even
// if dst currently isn't. It is effectively an infinitely buffered channel.
//
// In particular, goroutine A sending on src will not deadlock even if goroutine
// B that's responsible for receiving on dst is currently blocked trying to
// send to A on a separate channel.
//
// Send a stopPumping on the src channel to close the dst channel after all queued
// events are sent on dst. After that, other goroutines can still send to src,
// so that such sends won't block forever, but such events will be ignored.
func pump(dst chan interface{}) (src chan interface{}) {
	src = make(chan interface{})
	go func() {
		// initialSize is the initial size of the circular buffer. It must be a
		// power of 2.
		const initialSize = 16
		i, j, buf, mask := 0, 0, make([]interface{}, initialSize), initialSize-1

		srcActive := true
		for {
			maybeDst := dst
			if i == j {
				maybeDst = nil
			}
			if maybeDst == nil && !srcActive {
				// Pump is stopped and empty.
				break
			}

			select {
			case maybeDst <- buf[i&mask]:
				buf[i&mask] = nil
				i++

			case e := <-src:
				if _, ok := e.(stopPumping); ok {
					srcActive = false
					continue
				}

				if !srcActive {
					continue
				}

				// Allocate a bigger buffer if necessary.
				if i+len(buf) == j {
					b := make([]interface{}, 2*len(buf))
					n := copy(b, buf[j&mask:])
					copy(b[n:], buf[:j&mask])
					i, j = 0, len(buf)
					buf, mask = b, len(b)-1
				}

				buf[j&mask] = e
				j++
			}
		}

		close(dst)
		// Block forever.
		for range src {
		}
	}()
	return src
}

// TODO: do this for all build targets, not just linux (x11 and Android)? If
// so, should package gl instead of this package call RegisterFilter??
//
// TODO: does Android need this?? It seems to work without it (Nexus 7,
// KitKat). If only x11 needs this, should we move this to x11.go??
func (a *app) registerGLViewportFilter() {
	a.RegisterFilter(func(e interface{}) interface{} {
		if e, ok := e.(size.Event); ok {
			a.glctx.Viewport(0, 0, e.WidthPx, e.HeightPx)
		}
		return e
	})
}

func screenOrientation(width, height int) size.Orientation {
	if width > height {
		return size.OrientationLandscape
	}

	return size.OrientationPortrait
}
