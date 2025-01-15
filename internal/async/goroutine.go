package async

import (
	"log"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/build"
)

// mainGoroutineID stores the main goroutine ID.
// This ID must be initialized in main.init because
// a main goroutine may not equal to 1 due to the
// influence of a garbage collector.
var mainGoroutineID uint64

func init() {
	mainGoroutineID = goroutineID()
}

func IsMainGoroutine() bool {
	return goroutineID() == mainGoroutineID
}

func EnsureMain(fn func()) {
	if build.DisableThreadChecks || IsMainGoroutine() {
		fn()
		return
	}

	log.Println("*** Error in Fyne call thread, this should have been called in fyne.Do ***")

	pc := make([]uintptr, 2)
	count := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc)
	frame, more := frames.Next()
	if more && count > 1 {
		nextFrame, _ := frames.Next()                     // skip an occasional driver call to itself
		if !strings.Contains(nextFrame.File, "runtime") { // don't descend into Go
			frame = nextFrame
		}
	}
	log.Printf("  From: %s:%d", frame.File, frame.Line)

	fyne.Do(fn)
}
