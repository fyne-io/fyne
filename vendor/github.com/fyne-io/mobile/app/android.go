// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build android

/*
Android Apps are built with -buildmode=c-shared. They are loaded by a
running Java process.

Before any entry point is reached, a global constructor initializes the
Go runtime, calling all Go init functions. All cgo calls will block
until this is complete. Next JNI_OnLoad is called. When that is
complete, one of two entry points is called.

All-Go apps built using NativeActivity enter at ANativeActivity_onCreate.

Go libraries (for example, those built with gomobile bind) do not use
the app package initialization.
*/

package app

/*
#cgo LDFLAGS: -landroid -llog -lEGL -lGLESv2

#include <android/configuration.h>
#include <android/input.h>
#include <android/keycodes.h>
#include <android/looper.h>
#include <android/native_activity.h>
#include <android/native_window.h>
#include <EGL/egl.h>
#include <jni.h>
#include <pthread.h>
#include <stdlib.h>

EGLDisplay display;
EGLSurface surface;

char* createEGLSurface(ANativeWindow* window);
char* destroyEGLSurface();
int32_t getKeyRune(JNIEnv* env, AInputEvent* e);

void showKeyboard(JNIEnv* env, int keyboardType);
void hideKeyboard(JNIEnv* env);
void showFileOpen(JNIEnv* env, char* mimes);

void Java_org_golang_app_GoNativeActivity_filePickerReturned(JNIEnv *env, jclass clazz, jstring str);
*/
import "C"
import (
	"fmt"
	"log"
	"mime"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/fyne-io/mobile/app/internal/callfn"
	"github.com/fyne-io/mobile/event/key"
	"github.com/fyne-io/mobile/event/lifecycle"
	"github.com/fyne-io/mobile/event/paint"
	"github.com/fyne-io/mobile/event/size"
	"github.com/fyne-io/mobile/event/touch"
	"github.com/fyne-io/mobile/geom"
	"github.com/fyne-io/mobile/internal/mobileinit"
)

// mimeMap contains standard mime entries that are missing on Android
var mimeMap = map[string]string{
	".txt": "text/plain",
}

// RunOnJVM runs fn on a new goroutine locked to an OS thread with a JNIEnv.
//
// RunOnJVM blocks until the call to fn is complete. Any Java
// exception or failure to attach to the JVM is returned as an error.
//
// The function fn takes vm, the current JavaVM*,
// env, the current JNIEnv*, and
// ctx, a jobject representing the global android.context.Context.
func RunOnJVM(fn func(vm, jniEnv, ctx uintptr) error) error {
	return mobileinit.RunOnJVM(fn)
}

//export setCurrentContext
func setCurrentContext(vm *C.JavaVM, ctx C.jobject) {
	mobileinit.SetCurrentContext(unsafe.Pointer(vm), uintptr(ctx))
}

//export callMain
func callMain(mainPC uintptr) {
	for _, name := range []string{"FILESDIR", "TMPDIR", "PATH", "LD_LIBRARY_PATH"} {
		n := C.CString(name)
		os.Setenv(name, C.GoString(C.getenv(n)))
		C.free(unsafe.Pointer(n))
	}

	// Set timezone.
	//
	// Note that Android zoneinfo is stored in /system/usr/share/zoneinfo,
	// but it is in some kind of packed TZiff file that we do not support
	// yet. As a stopgap, we build a fixed zone using the tm_zone name.
	var curtime C.time_t
	var curtm C.struct_tm
	C.time(&curtime)
	C.localtime_r(&curtime, &curtm)
	tzOffset := int(curtm.tm_gmtoff)
	tz := C.GoString(curtm.tm_zone)
	time.Local = time.FixedZone(tz, tzOffset)

	go callfn.CallFn(mainPC)
}

//export onStart
func onStart(activity *C.ANativeActivity) {
}

//export onResume
func onResume(activity *C.ANativeActivity) {
}

//export onSaveInstanceState
func onSaveInstanceState(activity *C.ANativeActivity, outSize *C.size_t) unsafe.Pointer {
	return nil
}

//export onPause
func onPause(activity *C.ANativeActivity) {
}

//export onStop
func onStop(activity *C.ANativeActivity) {
}

//export onCreate
func onCreate(activity *C.ANativeActivity) {
	// Set the initial configuration.
	//
	// Note we use unbuffered channels to talk to the activity loop, and
	// NativeActivity calls these callbacks sequentially, so configuration
	// will be set before <-windowRedrawNeeded is processed.
	windowConfigChange <- windowConfigRead(activity)
}

//export onDestroy
func onDestroy(activity *C.ANativeActivity) {
	activityDestroyed <- struct{}{}
}

//export onWindowFocusChanged
func onWindowFocusChanged(activity *C.ANativeActivity, hasFocus C.int) {
}

//export onNativeWindowCreated
func onNativeWindowCreated(activity *C.ANativeActivity, window *C.ANativeWindow) {
}

//export onNativeWindowRedrawNeeded
func onNativeWindowRedrawNeeded(activity *C.ANativeActivity, window *C.ANativeWindow) {
	// Called on orientation change and window resize.
	// Send a request for redraw, and block this function
	// until a complete draw and buffer swap is completed.
	// This is required by the redraw documentation to
	// avoid bad draws.
	windowRedrawNeeded <- window
	<-windowRedrawDone
}

//export onNativeWindowDestroyed
func onNativeWindowDestroyed(activity *C.ANativeActivity, window *C.ANativeWindow) {
	windowDestroyed <- window
}

//export onInputQueueCreated
func onInputQueueCreated(activity *C.ANativeActivity, q *C.AInputQueue) {
	inputQueue <- q
	<-inputQueueDone
}

//export onInputQueueDestroyed
func onInputQueueDestroyed(activity *C.ANativeActivity, q *C.AInputQueue) {
	inputQueue <- nil
	<-inputQueueDone
}

//export onContentRectChanged
func onContentRectChanged(activity *C.ANativeActivity, rect *C.ARect) {
}

type windowConfig struct {
	orientation size.Orientation
	pixelsPerPt float32
}

func windowConfigRead(activity *C.ANativeActivity) windowConfig {
	aconfig := C.AConfiguration_new()
	C.AConfiguration_fromAssetManager(aconfig, activity.assetManager)
	orient := C.AConfiguration_getOrientation(aconfig)
	density := C.AConfiguration_getDensity(aconfig)
	C.AConfiguration_delete(aconfig)

	// Calculate the screen resolution. This value is approximate. For example,
	// a physical resolution of 200 DPI may be quantized to one of the
	// ACONFIGURATION_DENSITY_XXX values such as 160 or 240.
	//
	// A more accurate DPI could possibly be calculated from
	// https://developer.android.com/reference/android/util/DisplayMetrics.html#xdpi
	// but this does not appear to be accessible via the NDK. In any case, the
	// hardware might not even provide a more accurate number, as the system
	// does not apparently use the reported value. See golang.org/issue/13366
	// for a discussion.
	var dpi int
	switch density {
	case C.ACONFIGURATION_DENSITY_DEFAULT:
		dpi = 160
	case C.ACONFIGURATION_DENSITY_LOW,
		C.ACONFIGURATION_DENSITY_MEDIUM,
		213, // C.ACONFIGURATION_DENSITY_TV
		C.ACONFIGURATION_DENSITY_HIGH,
		320, // ACONFIGURATION_DENSITY_XHIGH
		480, // ACONFIGURATION_DENSITY_XXHIGH
		640: // ACONFIGURATION_DENSITY_XXXHIGH
		dpi = int(density)
	case C.ACONFIGURATION_DENSITY_NONE:
		log.Print("android device reports no screen density")
		dpi = 72
	default:
		log.Printf("android device reports unknown density: %d", density)
		// All we can do is guess.
		if density > 0 {
			dpi = int(density)
		} else {
			dpi = 72
		}
	}

	o := size.OrientationUnknown
	switch orient {
	case C.ACONFIGURATION_ORIENTATION_PORT:
		o = size.OrientationPortrait
	case C.ACONFIGURATION_ORIENTATION_LAND:
		o = size.OrientationLandscape
	}

	return windowConfig{
		orientation: o,
		pixelsPerPt: float32(dpi) / 72,
	}
}

//export onConfigurationChanged
func onConfigurationChanged(activity *C.ANativeActivity) {
	// A rotation event first triggers onConfigurationChanged, then
	// calls onNativeWindowRedrawNeeded. We extract the orientation
	// here and save it for the redraw event.
	windowConfigChange <- windowConfigRead(activity)
}

//export onLowMemory
func onLowMemory(activity *C.ANativeActivity) {
}

var (
	inputQueue         = make(chan *C.AInputQueue)
	inputQueueDone     = make(chan struct{})
	windowDestroyed    = make(chan *C.ANativeWindow)
	windowRedrawNeeded = make(chan *C.ANativeWindow)
	windowRedrawDone   = make(chan struct{})
	windowConfigChange = make(chan windowConfig)
	activityDestroyed  = make(chan struct{})

	screenInsetTop, screenInsetBottom, screenInsetLeft, screenInsetRight int
)

func init() {
	theApp.registerGLViewportFilter()
}

func main(f func(App)) {
	mainUserFn = f
	// TODO: merge the runInputQueue and mainUI functions?
	go func() {
		if err := mobileinit.RunOnJVM(runInputQueue); err != nil {
			log.Fatalf("app: %v", err)
		}
	}()
	// Preserve this OS thread for:
	//	1. the attached JNI thread
	//	2. the GL context
	if err := mobileinit.RunOnJVM(mainUI); err != nil {
		log.Fatalf("app: %v", err)
	}
}

// driverShowVirtualKeyboard requests the driver to show a virtual keyboard for text input
func driverShowVirtualKeyboard(keyboard KeyboardType) {
	err := mobileinit.RunOnJVM(func(vm, jniEnv, ctx uintptr) error {
		env := (*C.JNIEnv)(unsafe.Pointer(jniEnv)) // not a Go heap pointer
		C.showKeyboard(env, C.int(int32(keyboard)))
		return nil
	})
	if err != nil {
		log.Fatalf("app: %v", err)
	}
}

// driverHideVirtualKeyboard requests the driver to hide any visible virtual keyboard
func driverHideVirtualKeyboard() {
	if err := mobileinit.RunOnJVM(hideSoftInput); err != nil {
		log.Fatalf("app: %v", err)
	}
}

func hideSoftInput(vm, jniEnv, ctx uintptr) error {
	env := (*C.JNIEnv)(unsafe.Pointer(jniEnv)) // not a Go heap pointer
	C.hideKeyboard(env)
	return nil
}

var fileCallback func(string, func())

//export filePickerReturned
func filePickerReturned(str *C.char) {
	if fileCallback == nil {
		return
	}

	fileCallback(C.GoString(str), nil)
	fileCallback = nil
}

//export insetsChanged
func insetsChanged(top, bottom, left, right int) {
	screenInsetTop, screenInsetBottom, screenInsetLeft, screenInsetRight = top, bottom, left, right
}

func driverShowFileOpenPicker(callback func(string, func()), filter *FileFilter) {
	fileCallback = callback

	mimes := "*/*"
	if filter.MimeTypes != nil {
		mimes = strings.Join(filter.MimeTypes, "|")
	} else if filter.Extensions != nil {
		var mimeTypes []string
		for _, ext := range filter.Extensions {
			if mimeEntry, ok := mimeMap[ext]; ok {
				mimeTypes = append(mimeTypes, mimeEntry)

				continue
			}

			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				continue
			}

			mimeTypes = append(mimeTypes, mimeType)
		}
		mimes = strings.Join(mimeTypes, "|")
	}
	mimeStr := C.CString(mimes)
	defer C.free(unsafe.Pointer(mimeStr))

	open := func(vm, jniEnv, ctx uintptr) error {
		// TODO pass in filter...
		env := (*C.JNIEnv)(unsafe.Pointer(jniEnv)) // not a Go heap pointer
		C.showFileOpen(env, mimeStr)
		return nil
	}

	if err := mobileinit.RunOnJVM(open); err != nil {
		log.Fatalf("app: %v", err)
	}
}

var mainUserFn func(App)

var DisplayMetrics struct {
	WidthPx  int
	HeightPx int
}

func mainUI(vm, jniEnv, ctx uintptr) error {
	workAvailable := theApp.worker.WorkAvailable()

	donec := make(chan struct{})
	go func() {
		mainUserFn(theApp)
		close(donec)
	}()

	var pixelsPerPt float32

	for {
		select {
		case <-donec:
			return nil
		case cfg := <-windowConfigChange:
			pixelsPerPt = cfg.pixelsPerPt
		case w := <-windowRedrawNeeded:
			if C.surface == nil {
				if errStr := C.createEGLSurface(w); errStr != nil {
					return fmt.Errorf("%s (%s)", C.GoString(errStr), eglGetError())
				}
				DisplayMetrics.WidthPx = int(C.ANativeWindow_getWidth(w))
				DisplayMetrics.HeightPx = int(C.ANativeWindow_getHeight(w))
			}
			theApp.sendLifecycle(lifecycle.StageFocused)
			widthPx := int(C.ANativeWindow_getWidth(w))
			heightPx := int(C.ANativeWindow_getHeight(w))
			theApp.eventsIn <- size.Event{
				WidthPx:       widthPx,
				HeightPx:      heightPx,
				WidthPt:       geom.Pt(float32(widthPx) / pixelsPerPt),
				HeightPt:      geom.Pt(float32(heightPx) / pixelsPerPt),
				InsetTopPx:    screenInsetTop,
				InsetBottomPx: screenInsetBottom,
				InsetLeftPx:   screenInsetLeft,
				InsetRightPx:  screenInsetRight,
				PixelsPerPt:   pixelsPerPt,
				Orientation:   screenOrientation(widthPx, heightPx), // we are guessing orientation here as it was not always working
			}
			theApp.eventsIn <- paint.Event{External: true}
		case <-windowDestroyed:
			if C.surface != nil {
				if errStr := C.destroyEGLSurface(); errStr != nil {
					return fmt.Errorf("%s (%s)", C.GoString(errStr), eglGetError())
				}
			}
			C.surface = nil
			theApp.sendLifecycle(lifecycle.StageAlive)
		case <-activityDestroyed:
			theApp.sendLifecycle(lifecycle.StageDead)
		case <-workAvailable:
			theApp.worker.DoWork()
		case <-theApp.publish:
			// TODO: compare a generation number to redrawGen for stale paints?
			if C.surface != nil {
				// eglSwapBuffers blocks until vsync.
				if C.eglSwapBuffers(C.display, C.surface) == C.EGL_FALSE {
					log.Printf("app: failed to swap buffers (%s)", eglGetError())
				}
			}
			select {
			case windowRedrawDone <- struct{}{}:
			default:
			}
			theApp.publishResult <- PublishResult{}
		}
	}
}

func runInputQueue(vm, jniEnv, ctx uintptr) error {
	env := (*C.JNIEnv)(unsafe.Pointer(jniEnv)) // not a Go heap pointer

	// Android loopers select on OS file descriptors, not Go channels, so we
	// translate the inputQueue channel to an ALooper_wake call.
	l := C.ALooper_prepare(C.ALOOPER_PREPARE_ALLOW_NON_CALLBACKS)
	pending := make(chan *C.AInputQueue, 1)
	go func() {
		for q := range inputQueue {
			pending <- q
			C.ALooper_wake(l)
		}
	}()

	var q *C.AInputQueue
	for {
		if C.ALooper_pollAll(-1, nil, nil, nil) == C.ALOOPER_POLL_WAKE {
			select {
			default:
			case p := <-pending:
				if q != nil {
					processEvents(env, q)
					C.AInputQueue_detachLooper(q)
				}
				q = p
				if q != nil {
					C.AInputQueue_attachLooper(q, l, 0, nil, nil)
				}
				inputQueueDone <- struct{}{}
			}
		}
		if q != nil {
			processEvents(env, q)
		}
	}
}

func processEvents(env *C.JNIEnv, q *C.AInputQueue) {
	var e *C.AInputEvent
	for C.AInputQueue_getEvent(q, &e) >= 0 {
		if C.AInputQueue_preDispatchEvent(q, e) != 0 {
			continue
		}
		processEvent(env, e)
		C.AInputQueue_finishEvent(q, e, 0)
	}
}

func processEvent(env *C.JNIEnv, e *C.AInputEvent) {
	switch C.AInputEvent_getType(e) {
	case C.AINPUT_EVENT_TYPE_KEY:
		processKey(env, e)
	case C.AINPUT_EVENT_TYPE_MOTION:
		// At most one of the events in this batch is an up or down event; get its index and change.
		upDownIndex := C.size_t(C.AMotionEvent_getAction(e)&C.AMOTION_EVENT_ACTION_POINTER_INDEX_MASK) >> C.AMOTION_EVENT_ACTION_POINTER_INDEX_SHIFT
		upDownType := touch.TypeMove
		switch C.AMotionEvent_getAction(e) & C.AMOTION_EVENT_ACTION_MASK {
		case C.AMOTION_EVENT_ACTION_DOWN, C.AMOTION_EVENT_ACTION_POINTER_DOWN:
			upDownType = touch.TypeBegin
		case C.AMOTION_EVENT_ACTION_UP, C.AMOTION_EVENT_ACTION_POINTER_UP:
			upDownType = touch.TypeEnd
		}

		for i, n := C.size_t(0), C.AMotionEvent_getPointerCount(e); i < n; i++ {
			t := touch.TypeMove
			if i == upDownIndex {
				t = upDownType
			}
			theApp.eventsIn <- touch.Event{
				X:        float32(C.AMotionEvent_getX(e, i)),
				Y:        float32(C.AMotionEvent_getY(e, i)),
				Sequence: touch.Sequence(C.AMotionEvent_getPointerId(e, i)),
				Type:     t,
			}
		}
	default:
		log.Printf("unknown input event, type=%d", C.AInputEvent_getType(e))
	}
}

func processKey(env *C.JNIEnv, e *C.AInputEvent) {
	deviceID := C.AInputEvent_getDeviceId(e)
	if deviceID == 0 {
		// Software keyboard input, leaving for scribe/IME.
		return
	}

	k := key.Event{
		Rune: rune(C.getKeyRune(env, e)),
		Code: convAndroidKeyCode(int32(C.AKeyEvent_getKeyCode(e))),
	}
	switch C.AKeyEvent_getAction(e) {
	case C.AKEY_STATE_DOWN:
		k.Direction = key.DirPress
	case C.AKEY_STATE_UP:
		k.Direction = key.DirRelease
	default:
		k.Direction = key.DirNone
	}
	// TODO(crawshaw): set Modifiers.
	theApp.eventsIn <- k
}

func eglGetError() string {
	switch errNum := C.eglGetError(); errNum {
	case C.EGL_SUCCESS:
		return "EGL_SUCCESS"
	case C.EGL_NOT_INITIALIZED:
		return "EGL_NOT_INITIALIZED"
	case C.EGL_BAD_ACCESS:
		return "EGL_BAD_ACCESS"
	case C.EGL_BAD_ALLOC:
		return "EGL_BAD_ALLOC"
	case C.EGL_BAD_ATTRIBUTE:
		return "EGL_BAD_ATTRIBUTE"
	case C.EGL_BAD_CONTEXT:
		return "EGL_BAD_CONTEXT"
	case C.EGL_BAD_CONFIG:
		return "EGL_BAD_CONFIG"
	case C.EGL_BAD_CURRENT_SURFACE:
		return "EGL_BAD_CURRENT_SURFACE"
	case C.EGL_BAD_DISPLAY:
		return "EGL_BAD_DISPLAY"
	case C.EGL_BAD_SURFACE:
		return "EGL_BAD_SURFACE"
	case C.EGL_BAD_MATCH:
		return "EGL_BAD_MATCH"
	case C.EGL_BAD_PARAMETER:
		return "EGL_BAD_PARAMETER"
	case C.EGL_BAD_NATIVE_PIXMAP:
		return "EGL_BAD_NATIVE_PIXMAP"
	case C.EGL_BAD_NATIVE_WINDOW:
		return "EGL_BAD_NATIVE_WINDOW"
	case C.EGL_CONTEXT_LOST:
		return "EGL_CONTEXT_LOST"
	default:
		return fmt.Sprintf("Unknown EGL err: %d", errNum)
	}
}

func convAndroidKeyCode(aKeyCode int32) key.Code {
	// Many Android key codes do not map into USB HID codes.
	// For those, key.CodeUnknown is returned. This switch has all
	// cases, even the unknown ones, to serve as a documentation
	// and search aid.
	switch aKeyCode {
	case C.AKEYCODE_UNKNOWN:
	case C.AKEYCODE_SOFT_LEFT:
	case C.AKEYCODE_SOFT_RIGHT:
	case C.AKEYCODE_HOME:
		return key.CodeHome
	case C.AKEYCODE_BACK:
	case C.AKEYCODE_CALL:
	case C.AKEYCODE_ENDCALL:
	case C.AKEYCODE_0:
		return key.Code0
	case C.AKEYCODE_1:
		return key.Code1
	case C.AKEYCODE_2:
		return key.Code2
	case C.AKEYCODE_3:
		return key.Code3
	case C.AKEYCODE_4:
		return key.Code4
	case C.AKEYCODE_5:
		return key.Code5
	case C.AKEYCODE_6:
		return key.Code6
	case C.AKEYCODE_7:
		return key.Code7
	case C.AKEYCODE_8:
		return key.Code8
	case C.AKEYCODE_9:
		return key.Code9
	case C.AKEYCODE_STAR:
	case C.AKEYCODE_POUND:
	case C.AKEYCODE_DPAD_UP:
	case C.AKEYCODE_DPAD_DOWN:
	case C.AKEYCODE_DPAD_LEFT:
	case C.AKEYCODE_DPAD_RIGHT:
	case C.AKEYCODE_DPAD_CENTER:
	case C.AKEYCODE_VOLUME_UP:
		return key.CodeVolumeUp
	case C.AKEYCODE_VOLUME_DOWN:
		return key.CodeVolumeDown
	case C.AKEYCODE_POWER:
	case C.AKEYCODE_CAMERA:
	case C.AKEYCODE_CLEAR:
	case C.AKEYCODE_A:
		return key.CodeA
	case C.AKEYCODE_B:
		return key.CodeB
	case C.AKEYCODE_C:
		return key.CodeC
	case C.AKEYCODE_D:
		return key.CodeD
	case C.AKEYCODE_E:
		return key.CodeE
	case C.AKEYCODE_F:
		return key.CodeF
	case C.AKEYCODE_G:
		return key.CodeG
	case C.AKEYCODE_H:
		return key.CodeH
	case C.AKEYCODE_I:
		return key.CodeI
	case C.AKEYCODE_J:
		return key.CodeJ
	case C.AKEYCODE_K:
		return key.CodeK
	case C.AKEYCODE_L:
		return key.CodeL
	case C.AKEYCODE_M:
		return key.CodeM
	case C.AKEYCODE_N:
		return key.CodeN
	case C.AKEYCODE_O:
		return key.CodeO
	case C.AKEYCODE_P:
		return key.CodeP
	case C.AKEYCODE_Q:
		return key.CodeQ
	case C.AKEYCODE_R:
		return key.CodeR
	case C.AKEYCODE_S:
		return key.CodeS
	case C.AKEYCODE_T:
		return key.CodeT
	case C.AKEYCODE_U:
		return key.CodeU
	case C.AKEYCODE_V:
		return key.CodeV
	case C.AKEYCODE_W:
		return key.CodeW
	case C.AKEYCODE_X:
		return key.CodeX
	case C.AKEYCODE_Y:
		return key.CodeY
	case C.AKEYCODE_Z:
		return key.CodeZ
	case C.AKEYCODE_COMMA:
		return key.CodeComma
	case C.AKEYCODE_PERIOD:
		return key.CodeFullStop
	case C.AKEYCODE_ALT_LEFT:
		return key.CodeLeftAlt
	case C.AKEYCODE_ALT_RIGHT:
		return key.CodeRightAlt
	case C.AKEYCODE_SHIFT_LEFT:
		return key.CodeLeftShift
	case C.AKEYCODE_SHIFT_RIGHT:
		return key.CodeRightShift
	case C.AKEYCODE_TAB:
		return key.CodeTab
	case C.AKEYCODE_SPACE:
		return key.CodeSpacebar
	case C.AKEYCODE_SYM:
	case C.AKEYCODE_EXPLORER:
	case C.AKEYCODE_ENVELOPE:
	case C.AKEYCODE_ENTER:
		return key.CodeReturnEnter
	case C.AKEYCODE_DEL:
		return key.CodeDeleteBackspace
	case C.AKEYCODE_GRAVE:
		return key.CodeGraveAccent
	case C.AKEYCODE_MINUS:
		return key.CodeHyphenMinus
	case C.AKEYCODE_EQUALS:
		return key.CodeEqualSign
	case C.AKEYCODE_LEFT_BRACKET:
		return key.CodeLeftSquareBracket
	case C.AKEYCODE_RIGHT_BRACKET:
		return key.CodeRightSquareBracket
	case C.AKEYCODE_BACKSLASH:
		return key.CodeBackslash
	case C.AKEYCODE_SEMICOLON:
		return key.CodeSemicolon
	case C.AKEYCODE_APOSTROPHE:
		return key.CodeApostrophe
	case C.AKEYCODE_SLASH:
		return key.CodeSlash
	case C.AKEYCODE_AT:
	case C.AKEYCODE_NUM:
	case C.AKEYCODE_HEADSETHOOK:
	case C.AKEYCODE_FOCUS:
	case C.AKEYCODE_PLUS:
	case C.AKEYCODE_MENU:
	case C.AKEYCODE_NOTIFICATION:
	case C.AKEYCODE_SEARCH:
	case C.AKEYCODE_MEDIA_PLAY_PAUSE:
	case C.AKEYCODE_MEDIA_STOP:
	case C.AKEYCODE_MEDIA_NEXT:
	case C.AKEYCODE_MEDIA_PREVIOUS:
	case C.AKEYCODE_MEDIA_REWIND:
	case C.AKEYCODE_MEDIA_FAST_FORWARD:
	case C.AKEYCODE_MUTE:
	case C.AKEYCODE_PAGE_UP:
		return key.CodePageUp
	case C.AKEYCODE_PAGE_DOWN:
		return key.CodePageDown
	case C.AKEYCODE_PICTSYMBOLS:
	case C.AKEYCODE_SWITCH_CHARSET:
	case C.AKEYCODE_BUTTON_A:
	case C.AKEYCODE_BUTTON_B:
	case C.AKEYCODE_BUTTON_C:
	case C.AKEYCODE_BUTTON_X:
	case C.AKEYCODE_BUTTON_Y:
	case C.AKEYCODE_BUTTON_Z:
	case C.AKEYCODE_BUTTON_L1:
	case C.AKEYCODE_BUTTON_R1:
	case C.AKEYCODE_BUTTON_L2:
	case C.AKEYCODE_BUTTON_R2:
	case C.AKEYCODE_BUTTON_THUMBL:
	case C.AKEYCODE_BUTTON_THUMBR:
	case C.AKEYCODE_BUTTON_START:
	case C.AKEYCODE_BUTTON_SELECT:
	case C.AKEYCODE_BUTTON_MODE:
	case C.AKEYCODE_ESCAPE:
		return key.CodeEscape
	case C.AKEYCODE_FORWARD_DEL:
		return key.CodeDeleteForward
	case C.AKEYCODE_CTRL_LEFT:
		return key.CodeLeftControl
	case C.AKEYCODE_CTRL_RIGHT:
		return key.CodeRightControl
	case C.AKEYCODE_CAPS_LOCK:
		return key.CodeCapsLock
	case C.AKEYCODE_SCROLL_LOCK:
	case C.AKEYCODE_META_LEFT:
		return key.CodeLeftGUI
	case C.AKEYCODE_META_RIGHT:
		return key.CodeRightGUI
	case C.AKEYCODE_FUNCTION:
	case C.AKEYCODE_SYSRQ:
	case C.AKEYCODE_BREAK:
	case C.AKEYCODE_MOVE_HOME:
	case C.AKEYCODE_MOVE_END:
	case C.AKEYCODE_INSERT:
		return key.CodeInsert
	case C.AKEYCODE_FORWARD:
	case C.AKEYCODE_MEDIA_PLAY:
	case C.AKEYCODE_MEDIA_PAUSE:
	case C.AKEYCODE_MEDIA_CLOSE:
	case C.AKEYCODE_MEDIA_EJECT:
	case C.AKEYCODE_MEDIA_RECORD:
	case C.AKEYCODE_F1:
		return key.CodeF1
	case C.AKEYCODE_F2:
		return key.CodeF2
	case C.AKEYCODE_F3:
		return key.CodeF3
	case C.AKEYCODE_F4:
		return key.CodeF4
	case C.AKEYCODE_F5:
		return key.CodeF5
	case C.AKEYCODE_F6:
		return key.CodeF6
	case C.AKEYCODE_F7:
		return key.CodeF7
	case C.AKEYCODE_F8:
		return key.CodeF8
	case C.AKEYCODE_F9:
		return key.CodeF9
	case C.AKEYCODE_F10:
		return key.CodeF10
	case C.AKEYCODE_F11:
		return key.CodeF11
	case C.AKEYCODE_F12:
		return key.CodeF12
	case C.AKEYCODE_NUM_LOCK:
		return key.CodeKeypadNumLock
	case C.AKEYCODE_NUMPAD_0:
		return key.CodeKeypad0
	case C.AKEYCODE_NUMPAD_1:
		return key.CodeKeypad1
	case C.AKEYCODE_NUMPAD_2:
		return key.CodeKeypad2
	case C.AKEYCODE_NUMPAD_3:
		return key.CodeKeypad3
	case C.AKEYCODE_NUMPAD_4:
		return key.CodeKeypad4
	case C.AKEYCODE_NUMPAD_5:
		return key.CodeKeypad5
	case C.AKEYCODE_NUMPAD_6:
		return key.CodeKeypad6
	case C.AKEYCODE_NUMPAD_7:
		return key.CodeKeypad7
	case C.AKEYCODE_NUMPAD_8:
		return key.CodeKeypad8
	case C.AKEYCODE_NUMPAD_9:
		return key.CodeKeypad9
	case C.AKEYCODE_NUMPAD_DIVIDE:
		return key.CodeKeypadSlash
	case C.AKEYCODE_NUMPAD_MULTIPLY:
		return key.CodeKeypadAsterisk
	case C.AKEYCODE_NUMPAD_SUBTRACT:
		return key.CodeKeypadHyphenMinus
	case C.AKEYCODE_NUMPAD_ADD:
		return key.CodeKeypadPlusSign
	case C.AKEYCODE_NUMPAD_DOT:
		return key.CodeKeypadFullStop
	case C.AKEYCODE_NUMPAD_COMMA:
	case C.AKEYCODE_NUMPAD_ENTER:
		return key.CodeKeypadEnter
	case C.AKEYCODE_NUMPAD_EQUALS:
		return key.CodeKeypadEqualSign
	case C.AKEYCODE_NUMPAD_LEFT_PAREN:
	case C.AKEYCODE_NUMPAD_RIGHT_PAREN:
	case C.AKEYCODE_VOLUME_MUTE:
		return key.CodeMute
	case C.AKEYCODE_INFO:
	case C.AKEYCODE_CHANNEL_UP:
	case C.AKEYCODE_CHANNEL_DOWN:
	case C.AKEYCODE_ZOOM_IN:
	case C.AKEYCODE_ZOOM_OUT:
	case C.AKEYCODE_TV:
	case C.AKEYCODE_WINDOW:
	case C.AKEYCODE_GUIDE:
	case C.AKEYCODE_DVR:
	case C.AKEYCODE_BOOKMARK:
	case C.AKEYCODE_CAPTIONS:
	case C.AKEYCODE_SETTINGS:
	case C.AKEYCODE_TV_POWER:
	case C.AKEYCODE_TV_INPUT:
	case C.AKEYCODE_STB_POWER:
	case C.AKEYCODE_STB_INPUT:
	case C.AKEYCODE_AVR_POWER:
	case C.AKEYCODE_AVR_INPUT:
	case C.AKEYCODE_PROG_RED:
	case C.AKEYCODE_PROG_GREEN:
	case C.AKEYCODE_PROG_YELLOW:
	case C.AKEYCODE_PROG_BLUE:
	case C.AKEYCODE_APP_SWITCH:
	case C.AKEYCODE_BUTTON_1:
	case C.AKEYCODE_BUTTON_2:
	case C.AKEYCODE_BUTTON_3:
	case C.AKEYCODE_BUTTON_4:
	case C.AKEYCODE_BUTTON_5:
	case C.AKEYCODE_BUTTON_6:
	case C.AKEYCODE_BUTTON_7:
	case C.AKEYCODE_BUTTON_8:
	case C.AKEYCODE_BUTTON_9:
	case C.AKEYCODE_BUTTON_10:
	case C.AKEYCODE_BUTTON_11:
	case C.AKEYCODE_BUTTON_12:
	case C.AKEYCODE_BUTTON_13:
	case C.AKEYCODE_BUTTON_14:
	case C.AKEYCODE_BUTTON_15:
	case C.AKEYCODE_BUTTON_16:
	case C.AKEYCODE_LANGUAGE_SWITCH:
	case C.AKEYCODE_MANNER_MODE:
	case C.AKEYCODE_3D_MODE:
	case C.AKEYCODE_CONTACTS:
	case C.AKEYCODE_CALENDAR:
	case C.AKEYCODE_MUSIC:
	case C.AKEYCODE_CALCULATOR:
	}
	/* Defined in an NDK API version beyond what we use today:
	C.AKEYCODE_ASSIST
	C.AKEYCODE_BRIGHTNESS_DOWN
	C.AKEYCODE_BRIGHTNESS_UP
	C.AKEYCODE_EISU
	C.AKEYCODE_HENKAN
	C.AKEYCODE_KANA
	C.AKEYCODE_KATAKANA_HIRAGANA
	C.AKEYCODE_MEDIA_AUDIO_TRACK
	C.AKEYCODE_MUHENKAN
	C.AKEYCODE_RO
	C.AKEYCODE_YEN
	C.AKEYCODE_ZENKAKU_HANKAKU
	*/
	return key.CodeUnknown
}
