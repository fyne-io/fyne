package async

import (
	"log"
	"runtime"
	"strings"
	"sync/atomic"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/build"
)

var (
	// mainGoroutineID stores the main goroutine ID.
	// This ID must be initialized during setup by calling `SetMainGoroutine` because
	// a main goroutine may not equal to 1 due to the influence of a garbage collector.
	mainGoroutineID atomic.Uint64

	warnedNotMigrated bool
)

func SetMainGoroutine() {
	mainGoroutineID.Store(goroutineID())
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
	PrintFyneDoWarning()

	log.Println("*** Error in Fyne call thread, fyne.Do[AndWait] called from main goroutine ***")
	log.Println("*** This error needs to be addressed as it will become a deadlock soon ***")
	log.Println("*** The next major Fyne release will disable this safety by default! ***")

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
	PrintFyneDoWarning()

	log.Println("*** Error in Fyne call thread, this should have been called in fyne.Do[AndWait] ***")
	log.Println("*** This error needs to be addressed as it could cause concurrent errors soon ***")
	log.Println("*** The next major Fyne release will disable this safety by default! ***")

	logStackTop(1)
	fyne.DoAndWait(fn)
}

func logStackTop(skip int) {
	pc := make([]uintptr, 16) //revive:disable-line:add-constant -- TODO: What is this 16 about?
	_ = runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc)
	frame, more := frames.Next()

	var nextFrame runtime.Frame
	for more {
		nextFrame, more = frames.Next()
		if nextFrame.File == "" || strings.Contains(nextFrame.File, "runtime") { // don't descend into Go
			break
		}

		frame = nextFrame
		if !strings.Contains(nextFrame.File, "/fyne/") { // skip library lines
			break
		}
	}
	log.Printf("  From: %s:%d", frame.File, frame.Line)
}

func goroutineID() (id uint64) {
	//revive:disable:add-constant
	var buf [30]byte
	runtime.Stack(buf[:], false)
	for i := 10; buf[i] != ' '; i++ {
		id = id*10 + uint64(buf[i]&15)
	}
	//revive:enable:add-constant

	return id
}

func PrintFyneDoWarning() {
	if warnedNotMigrated { // shown before
		return
	}

	warnedNotMigrated = true
	log.Println("*** This application has not been migrated to the fyne.Do threading model ***")
	log.Println("*** The next major Fyne release will remove this safety! ***")
	log.Println("*** Read more at https://docs.fyne.io/started/goroutines ***")
}
