// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package app

import (
	"fmt"
)

func main(f func(a App)) {
	fmt.Errorf("Running mobile simulation mode does not currently work on Windows.")
}

// ShowVirtualKeyboard does nothing on desktop
func ShowVirtualKeyboard() {
}

// HideVirtualKeyboard does nothing on desktop
func HideVirtualKeyboard() {
}
