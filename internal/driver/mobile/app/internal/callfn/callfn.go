//go:build android
// +build android

// Package callfn provides an android entry point.
//
// It is a separate package from app because it contains Go assembly,
// which does not compile in a package using cgo.
package callfn

// CallFn calls a zero-argument function by its program counter.
// It is only intended for calling main.main. Using it for
// anything else will not end well.
func CallFn(fn uintptr)
