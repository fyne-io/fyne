// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin
// +build ios

package app

/*
#cgo CFLAGS: -x objective-c -DGL_SILENCE_DEPRECATION
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework GLKit -framework OpenGLES -framework QuartzCore
#include <sys/utsname.h>
#include <stdint.h>
#include <pthread.h>
#include <UIKit/UIDevice.h>
#import <GLKit/GLKit.h>

extern struct utsname sysInfo;

void runApp(void);
void makeCurrentContext(GLintptr ctx);
void swapBuffers(GLintptr ctx);
uint64_t threadID();

void showKeyboard();
void hideKeyboard();
*/
import "C"
import (
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/geom"
)

var initThreadID uint64

func init() {
	// Lock the goroutine responsible for initialization to an OS thread.
	// This means the goroutine running main (and calling the run function
	// below) is locked to the OS thread that started the program. This is
	// necessary for the correct delivery of UIKit events to the process.
	//
	// A discussion on this topic:
	// https://groups.google.com/forum/#!msg/golang-nuts/IiWZ2hUuLDA/SNKYYZBelsYJ
	runtime.LockOSThread()
	initThreadID = uint64(C.threadID())
}

func main(f func(App)) {
	//if tid := uint64(C.threadID()); tid != initThreadID {
	//	log.Fatalf("app.Run called on thread %d, but app.init ran on %d", tid, initThreadID)
	//}

	go func() {
		f(theApp)
		// TODO(crawshaw): trigger runApp to return
	}()
	C.runApp()
	panic("unexpected return from app.runApp")
}

var pixelsPerPt float32
var screenScale int // [UIScreen mainScreen].scale, either 1, 2, or 3.

var DisplayMetrics struct{
	WidthPx int
	HeightPx int
}

//export setDisplayMetrics
func setDisplayMetrics(width, height int, scale int) {
	DisplayMetrics.WidthPx = width * scale
	DisplayMetrics.HeightPx = height * scale
}

//export setScreen
func setScreen(scale int) {
	C.uname(&C.sysInfo)
	name := C.GoString(&C.sysInfo.machine[0])

	var v float32

	switch {
	case strings.HasPrefix(name, "iPhone"):
		v = 163
	case strings.HasPrefix(name, "iPad"):
		// TODO: is there a better way to distinguish the iPad Mini?
		switch name {
		case "iPad2,5", "iPad2,6", "iPad2,7", "iPad4,4", "iPad4,5", "iPad4,6", "iPad4,7":
			v = 163 // iPad Mini
		default:
			v = 132
		}
	default:
		v = 163 // names like i386 and x86_64 are the simulator
	}

	if v == 0 {
		log.Printf("unknown machine: %s", name)
		v = 163 // emergency fallback
	}

	pixelsPerPt = v * float32(scale) / 72
	screenScale = scale
}

//export updateConfig
func updateConfig(width, height, orientation int32) {
	o := size.OrientationUnknown
	switch orientation {
	case C.UIDeviceOrientationPortrait, C.UIDeviceOrientationPortraitUpsideDown:
		o = size.OrientationPortrait
	case C.UIDeviceOrientationLandscapeLeft, C.UIDeviceOrientationLandscapeRight:
		o = size.OrientationLandscape
	}
	widthPx := screenScale * int(width)
	heightPx := screenScale * int(height)
	theApp.eventsIn <- size.Event{
		WidthPx:     widthPx,
		HeightPx:    heightPx,
		WidthPt:     geom.Pt(float32(widthPx) / pixelsPerPt),
		HeightPt:    geom.Pt(float32(heightPx) / pixelsPerPt),
		PixelsPerPt: pixelsPerPt,
		Orientation: o,
	}
	theApp.eventsIn <- paint.Event{External: true}
}

// touchIDs is the current active touches. The position in the array
// is the ID, the value is the UITouch* pointer value.
//
// It is widely reported that the iPhone can handle up to 5 simultaneous
// touch events, while the iPad can handle 11.
var touchIDs [11]uintptr

var touchEvents struct {
	sync.Mutex
	pending []touch.Event
}

//export sendTouch
func sendTouch(cTouch, cTouchType uintptr, x, y float32) {
	id := -1
	for i, val := range touchIDs {
		if val == cTouch {
			id = i
			break
		}
	}
	if id == -1 {
		for i, val := range touchIDs {
			if val == 0 {
				touchIDs[i] = cTouch
				id = i
				break
			}
		}
		if id == -1 {
			panic("out of touchIDs")
		}
	}

	t := touch.Type(cTouchType)
	if t == touch.TypeEnd {
		touchIDs[id] = 0
	}

	theApp.eventsIn <- touch.Event{
		X:        x,
		Y:        y,
		Sequence: touch.Sequence(id),
		Type:     t,
	}
}

//export lifecycleDead
func lifecycleDead() { theApp.sendLifecycle(lifecycle.StageDead) }

//export lifecycleAlive
func lifecycleAlive() { theApp.sendLifecycle(lifecycle.StageAlive) }

//export lifecycleVisible
func lifecycleVisible() { theApp.sendLifecycle(lifecycle.StageVisible) }

//export lifecycleFocused
func lifecycleFocused() { theApp.sendLifecycle(lifecycle.StageFocused) }

//export drawloop
func drawloop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	for workAvailable := theApp.worker.WorkAvailable();;{
		select {
		case <-workAvailable:
			theApp.worker.DoWork()
		case <-theApp.publish:
			theApp.publishResult <- PublishResult{}
			return
		case <-time.After(100 * time.Millisecond): // incase the method blocked!!
			return
		}
	}
}

//export startloop
func startloop(ctx C.GLintptr) {
	go theApp.loop(ctx)
}

// loop is the primary drawing loop.
//
// After UIKit has captured the initial OS thread for processing UIKit
// events in runApp, it starts loop on another goroutine. It is locked
// to an OS thread for its OpenGL context.
func (a *app) loop(ctx C.GLintptr) {
	runtime.LockOSThread()
	C.makeCurrentContext(ctx)

	workAvailable := a.worker.WorkAvailable()

	for {
		select {
		case <-workAvailable:
			a.worker.DoWork()
		case <-theApp.publish:
		loop1:
			for {
				select {
				case <-workAvailable:
					a.worker.DoWork()
				default:
					break loop1
				}
			}
			C.swapBuffers(ctx)
			theApp.publishResult <- PublishResult{}
		}
	}
}

// driverShowVirtualKeyboard requests the driver to show a virtual keyboard for text input
func driverShowVirtualKeyboard() {
	C.showKeyboard()
}

// driverHideVirtualKeyboard requests the driver to hide any visible virtual keyboard
func driverHideVirtualKeyboard() {
	C.hideKeyboard()
}

//export keyboardTyped
func keyboardTyped(str *C.char) {
	for _, r := range C.GoString(str) {
		k := key.Event{
			Rune: r,
			Code: getCodeFromRune(r),
			Direction: key.DirPress,
		}
		theApp.eventsIn <- k

		k.Direction = key.DirRelease
		theApp.eventsIn <- k
	}
}

//export keyboardDelete
func keyboardDelete() {
	theApp.eventsIn <- key.Event{
		Code: key.CodeDeleteBackspace,
		Direction: key.DirPress,
	}
	theApp.eventsIn <- key.Event{
		Code: key.CodeDeleteBackspace,
		Direction: key.DirRelease,
	}
}

func getCodeFromRune(r rune) key.Code {
	switch r {
	case '0':
		return key.Code0
	case '1':
		return key.Code1
	case '2':
		return key.Code2
	case '3':
		return key.Code3
	case '4':
		return key.Code4
	case '5':
		return key.Code5
	case '6':
		return key.Code6
	case '7':
		return key.Code7
	case '8':
		return key.Code8
	case '9':
		return key.Code9
	case 'a', 'A':
		return key.CodeA
	case 'b', 'B':
		return key.CodeB
	case 'c', 'C':
		return key.CodeC
	case 'd', 'D':
		return key.CodeD
	case 'e', 'E':
		return key.CodeE
	case 'f', 'F':
		return key.CodeF
	case 'g', 'G':
		return key.CodeG
	case 'h', 'H':
		return key.CodeH
	case 'i', 'I':
		return key.CodeI
	case 'j', 'J':
		return key.CodeJ
	case 'k', 'K':
		return key.CodeK
	case 'l', 'L':
		return key.CodeL
	case 'm', 'M':
		return key.CodeM
	case 'n', 'N':
		return key.CodeN
	case 'o', 'O':
		return key.CodeO
	case 'p', 'P':
		return key.CodeP
	case 'q', 'Q':
		return key.CodeQ
	case 'r', 'R':
		return key.CodeR
	case 's', 'S':
		return key.CodeS
	case 't', 'T':
		return key.CodeT
	case 'u', 'U':
		return key.CodeU
	case 'v', 'V':
		return key.CodeV
	case 'w', 'W':
		return key.CodeW
	case 'x', 'X':
		return key.CodeX
	case 'y', 'Y':
		return key.CodeY
	case 'z', 'Z':
		return key.CodeZ
	case ',':
		return key.CodeComma
	case '.':
		return key.CodeFullStop
	case ' ':
		return key.CodeSpacebar
	case '\n':
		return key.CodeReturnEnter
	case '`':
		return key.CodeGraveAccent
	case '-':
		return key.CodeHyphenMinus
	case '=':
		return key.CodeEqualSign
	case '[':
		return key.CodeLeftSquareBracket
	case ']':
		return key.CodeRightSquareBracket
	case '\\':
		return key.CodeBackslash
	case ';':
		return key.CodeSemicolon
	case '\'':
		return key.CodeApostrophe
	case '/':
		return key.CodeSlash
	}
	return key.CodeUnknown
}