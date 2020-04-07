// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package paint defines an event for the app being ready to paint.
//
// See the golang.org/x/mobile/app package for details on the event model.
package paint // import "github.com/fyne-io/mobile/event/paint"

// Event indicates that the app is ready to paint the next frame of the GUI.
//
//A frame is completed by calling the App's Publish method.
type Event struct {
	// External is true for paint events sent by the screen driver.
	//
	// An external event may be sent at any time in response to an
	// operating system event, for example the window opened, was
	// resized, or the screen memory was lost.
	//
	// Programs actively drawing to the screen as fast as vsync allows
	// should ignore external paint events to avoid a backlog of paint
	// events building up.
	External bool
}
