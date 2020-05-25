// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !cgo !darwin,!linux,!openbsd,!freebsd
// +build !windows

package gl

// This file contains stub implementations of what the other work*.go files
// provide. These stubs don't do anything, other than compile (e.g. when cgo is
// disabled).

type context struct{}

func (*context) enqueue(c call) uintptr {
	panic("unimplemented; GOOS/CGO combination not supported")
}

func (*context) cString(str string) (uintptr, func()) {
	panic("unimplemented; GOOS/CGO combination not supported")
}

func (*context) cStringPtr(str string) (uintptr, func()) {
	panic("unimplemented; GOOS/CGO combination not supported")
}

type context3 = context

func NewContext() (Context, Worker) {
	panic("unimplemented; GOOS/CGO combination not supported")
}

func Version() string {
	panic("unimplemented; GOOS/CGO combination not supported")
}
