// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package app

import (
	"log"

	"golang.org/x/exp/shiny/driver/gldriver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func main(f func(a App)) {
	gldriver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(nil)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		theApp.glctx = nil
		theApp.worker = nil // handled by shiny

		go func() {
			for range theApp.publish {
				res := w.Publish()
				theApp.publishResult <- PublishResult{
					BackBufferPreserved: res.BackBufferPreserved,
				}
			}
		}()

		go f(theApp)

		for {
			theApp.Send(convertEvent(w.NextEvent()))
		}
	})
}

func convertEvent(e interface{}) interface{} {
	switch e := e.(type) {
	case lifecycle.Event:
		if theApp.glctx == nil {
			theApp.glctx = e.DrawContext.(gl.Context)
		}
	case mouse.Event:
		te := touch.Event{
			X: e.X,
			Y: e.Y,
		}
		switch e.Direction {
		case mouse.DirNone:
			te.Type = touch.TypeMove
		case mouse.DirPress:
			te.Type = touch.TypeBegin
		case mouse.DirRelease:
			te.Type = touch.TypeEnd
		}
		return te
	}
	return e
}
