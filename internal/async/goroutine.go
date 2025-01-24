package async

import (
	"log"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/build"
)

// mainGoroutineID stores the main goroutine ID.
// This ID must be initialized during setup by calling `SetMainGoroutine` because
// a main goroutine may not equal to 1 due to the influence of a garbage collector.
var mainGoroutineID uint64

func SetMainGoroutine() {
	mainGoroutineID = goroutineID()
}

// IsMainGoroutine returns true if it is called from the main goroutine, false otherwise.
func IsMainGoroutine() bool {
	return goroutineID() == mainGoroutineID
}

// EnsureNotMain is part of our thread transition and makes sure that the passed function runs off main.
// If the context is running on a goroutine or the transition has been disabled this will blindly run.
// Otherwise, an error will be logged and the function will be called on a new goroutine.
//
// This will be removed later and should never be public
func EnsureNotMain(fn func()) {
	if build.MigratedToFyneDo() || !IsMainGoroutine() {
		fn()
		return
	}

	log.Println("*** Error in Fyne call thread, fyne.Do[AndWait] called from main goroutine ***")

	logStackTop(2)
	go fn()
}

// EnsureMain is part of our thread transition and makes sure that the passed function runs on main.
// If the context is main or the transition has been disabled this will blindly run.
// Otherwise, an error will be logged and the function will be called on the main goroutine.
//
// This will be removed later and should never be public
func EnsureMain(fn func()) {
	if build.MigratedToFyneDo() || IsMainGoroutine() {
		fn()
		return
	}

	log.Println("*** Error in Fyne call thread, this should have been called in fyne.Do[AndWait] ***")

	logStackTop(1)
	fyne.DoAndWait(fn)
}

func logStackTop(skip int) {
	pc := make([]uintptr, 2)
	count := runtime.Callers(2+skip, pc)
	frames := runtime.CallersFrames(pc)
	frame, more := frames.Next()
	if more && count > 1 {
		nextFrame, _ := frames.Next()                     // skip an occasional driver call to itself
		if !strings.Contains(nextFrame.File, "runtime") { // don't descend into Go
			frame = nextFrame
		}
	}
	log.Printf("  From: %s:%d", frame.File, frame.Line)
}
