// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows

package app

import "log"

func GoBack() {
	// When simulating mobile there are no other activities open (and we can't just force background)
}

// Main is called by the main.main function to run the mobile application.
//
// It calls f on the App, in a separate goroutine, as some OS-specific
// libraries require being on 'the main thread'.
func Main(f func(App)) {
	log.Fatalln("Running mobile simulation mode does not currently work on Windows.") //revive:disable-line:deep-exit
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
