//go:build ios || android

package app

import "C"

import (
	"runtime"
	"runtime/debug"
)

// clean caches - called when the OS wants some memory back
func cleanCaches() {
	runtime.GC()
	debug.FreeOSMemory()
}
