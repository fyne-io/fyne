// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android && (arm || 386 || amd64 || arm64)
// +build android
// +build arm 386 amd64 arm64

// Package callfn provides an android entry point.
//
// It is a separate package from app because it contains Go assembly,
// which does not compile in a package using cgo.
package callfn

// CallFn calls a zero-argument function by its program counter.
// It is only intended for calling main.main. Using it for
// anything else will not end well.
func CallFn(fn uintptr)
