// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package app lets you write portable all-Go apps for Android and iOS.

There are typically two ways to use Go on Android and iOS. The first
is to write a Go library and use `gomobile bind` to generate language
bindings for Java and Objective-C. Building a library does not
require the app package. The `gomobile bind` command produces output
that you can include in an Android Studio or Xcode project. For more
on language bindings, see https://golang.org/x/mobile/cmd/gobind.

The second way is to write an app entirely in Go. The APIs are limited
to those that are portable between both Android and iOS, in particular
OpenGL, audio, and other Android NDK-like APIs. An all-Go app should
use this app package to initialize the app, manage its lifecycle, and
receive events.

Building apps

Apps written entirely in Go have a main function, and can be built
with `gomobile build`, which directly produces runnable output for
Android and iOS.

The gomobile tool can get installed with go get. For reference, see
https://golang.org/x/mobile/cmd/gomobile.

For detailed instructions and documentation, see
https://golang.org/wiki/Mobile.

Event processing in Native Apps

The Go runtime is initialized on Android when NativeActivity onCreate is
called, and on iOS when the process starts. In both cases, Go init functions
run before the app lifecycle has started.

An app is expected to call the Main function in main.main. When the function
exits, the app exits. Inside the func passed to Main, call Filter on every
event received, and then switch on its type. Registered filters run when the
event is received, not when it is sent, so that filters run in the same
goroutine as other code that calls OpenGL.

	package main

	import (
		"log"

		"golang.org/x/mobile/app"
		"golang.org/x/mobile/event/lifecycle"
		"golang.org/x/mobile/event/paint"
	)

	func main() {
		app.Main(func(a app.App) {
			for e := range a.Events() {
				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					// ...
				case paint.Event:
					log.Print("Call OpenGL here.")
					a.Publish()
				}
			}
		})
	}

An event is represented by the empty interface type interface{}. Any value can
be an event. Commonly used types include Event types defined by the following
packages:
	- golang.org/x/mobile/event/lifecycle
	- golang.org/x/mobile/event/mouse
	- golang.org/x/mobile/event/paint
	- golang.org/x/mobile/event/size
	- golang.org/x/mobile/event/touch
For example, touch.Event is the type that represents touch events. Other
packages may define their own events, and send them on an app's event channel.

Other packages can also register event filters, e.g. to manage resources in
response to lifecycle events. Such packages should call:
	app.RegisterFilter(etc)
in an init function inside that package.
*/
package app // import "golang.org/x/mobile/app"
