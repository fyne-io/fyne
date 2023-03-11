// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android
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
#include <stdbool.h>

extern EGLDisplay display;
extern EGLSurface surface;

char* createEGLSurface(ANativeWindow* window);
char* destroyEGLSurface();
int32_t getKeyRune(JNIEnv* env, AInputEvent* e);

void showKeyboard(JNIEnv* env, int keyboardType);
void hideKeyboard(JNIEnv* env);
void showFileOpen(JNIEnv* env, char* mimes);
void showFileSave(JNIEnv* env, char* mimes, char* filename);

void Java_org_golang_app_GoNativeActivity_filePickerReturned(JNIEnv *env, jclass clazz, jstring str);
*/
import "C"
import (
	"fmt"
	"log"
	"mime"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
	"unsafe"

	"fyne.io/fyne/v2/internal/driver/mobile/app/callfn"
	"fyne.io/fyne/v2/internal/driver/mobile/event/key"
	"fyne.io/fyne/v2/internal/driver/mobile/event/lifecycle"
	"fyne.io/fyne/v2/internal/driver/mobile/event/paint"
	"fyne.io/fyne/v2/internal/driver/mobile/event/size"
	"fyne.io/fyne/v2/internal/driver/mobile/event/touch"
	"fyne.io/fyne/v2/internal/driver/mobile/mobileinit"
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

//export setDarkMode
func setDarkMode(dark C.bool) {
	darkMode = bool(dark)
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
	runtime.GC()
	debug.FreeOSMemory()
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
	darkMode                                                             bool
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

func mimeStringFromFilter(filter *FileFilter) string {
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
				log.Println("Could not find mime for extension " + ext + ", allowing all")
				return "*/*" // could not find one, so allow all
			}

			mimeTypes = append(mimeTypes, mimeType)
		}
		mimes = strings.Join(mimeTypes, "|")
	}
	return mimes
}

func driverShowFileOpenPicker(callback func(string, func()), filter *FileFilter) {
	fileCallback = callback

	mimes := mimeStringFromFilter(filter)
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

func driverShowFileSavePicker(callback func(string, func()), filter *FileFilter, filename string) {
	fileCallback = callback

	mimes := mimeStringFromFilter(filter)
	mimeStr := C.CString(mimes)
	defer C.free(unsafe.Pointer(mimeStr))
	filenameStr := C.CString(filename)
	defer C.free(unsafe.Pointer(filenameStr))

	save := func(vm, jniEnv, ctx uintptr) error {
		env := (*C.JNIEnv)(unsafe.Pointer(jniEnv)) // not a Go heap pointer
		C.showFileSave(env, mimeStr, filenameStr)
		return nil
	}

	if err := mobileinit.RunOnJVM(save); err != nil {
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
			theApp.events.In() <- size.Event{
				WidthPx:       widthPx,
				HeightPx:      heightPx,
				WidthPt:       float32(widthPx) / pixelsPerPt,
				HeightPt:      float32(heightPx) / pixelsPerPt,
				InsetTopPx:    screenInsetTop,
				InsetBottomPx: screenInsetBottom,
				InsetLeftPx:   screenInsetLeft,
				InsetRightPx:  screenInsetRight,
				PixelsPerPt:   pixelsPerPt,
				Orientation:   screenOrientation(widthPx, heightPx), // we are guessing orientation here as it was not always working
				DarkMode:      darkMode,
			}
			theApp.events.In() <- paint.Event{External: true}
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
			theApp.events.In() <- touch.Event{
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
	if k.Rune >= '0' && k.Rune <= '9' { // GBoard generates key events for numbers, but we see them in textChanged
		return
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
	theApp.events.In() <- k
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

var androidKeycoe = map[int32]key.Code{
	C.AKEYCODE_HOME:            key.CodeHome,
	C.AKEYCODE_0:               key.Code0,
	C.AKEYCODE_1:               key.Code1,
	C.AKEYCODE_2:               key.Code2,
	C.AKEYCODE_3:               key.Code3,
	C.AKEYCODE_4:               key.Code4,
	C.AKEYCODE_5:               key.Code5,
	C.AKEYCODE_6:               key.Code6,
	C.AKEYCODE_7:               key.Code7,
	C.AKEYCODE_8:               key.Code8,
	C.AKEYCODE_9:               key.Code9,
	C.AKEYCODE_VOLUME_UP:       key.CodeVolumeUp,
	C.AKEYCODE_VOLUME_DOWN:     key.CodeVolumeDown,
	C.AKEYCODE_A:               key.CodeA,
	C.AKEYCODE_B:               key.CodeB,
	C.AKEYCODE_C:               key.CodeC,
	C.AKEYCODE_D:               key.CodeD,
	C.AKEYCODE_E:               key.CodeE,
	C.AKEYCODE_F:               key.CodeF,
	C.AKEYCODE_G:               key.CodeG,
	C.AKEYCODE_H:               key.CodeH,
	C.AKEYCODE_I:               key.CodeI,
	C.AKEYCODE_J:               key.CodeJ,
	C.AKEYCODE_K:               key.CodeK,
	C.AKEYCODE_L:               key.CodeL,
	C.AKEYCODE_M:               key.CodeM,
	C.AKEYCODE_N:               key.CodeN,
	C.AKEYCODE_O:               key.CodeO,
	C.AKEYCODE_P:               key.CodeP,
	C.AKEYCODE_Q:               key.CodeQ,
	C.AKEYCODE_R:               key.CodeR,
	C.AKEYCODE_S:               key.CodeS,
	C.AKEYCODE_T:               key.CodeT,
	C.AKEYCODE_U:               key.CodeU,
	C.AKEYCODE_V:               key.CodeV,
	C.AKEYCODE_W:               key.CodeW,
	C.AKEYCODE_X:               key.CodeX,
	C.AKEYCODE_Y:               key.CodeY,
	C.AKEYCODE_Z:               key.CodeZ,
	C.AKEYCODE_COMMA:           key.CodeComma,
	C.AKEYCODE_PERIOD:          key.CodeFullStop,
	C.AKEYCODE_ALT_LEFT:        key.CodeLeftAlt,
	C.AKEYCODE_ALT_RIGHT:       key.CodeRightAlt,
	C.AKEYCODE_SHIFT_LEFT:      key.CodeLeftShift,
	C.AKEYCODE_SHIFT_RIGHT:     key.CodeRightShift,
	C.AKEYCODE_TAB:             key.CodeTab,
	C.AKEYCODE_SPACE:           key.CodeSpacebar,
	C.AKEYCODE_ENTER:           key.CodeReturnEnter,
	C.AKEYCODE_DEL:             key.CodeDeleteBackspace,
	C.AKEYCODE_GRAVE:           key.CodeGraveAccent,
	C.AKEYCODE_MINUS:           key.CodeHyphenMinus,
	C.AKEYCODE_EQUALS:          key.CodeEqualSign,
	C.AKEYCODE_LEFT_BRACKET:    key.CodeLeftSquareBracket,
	C.AKEYCODE_RIGHT_BRACKET:   key.CodeRightSquareBracket,
	C.AKEYCODE_BACKSLASH:       key.CodeBackslash,
	C.AKEYCODE_SEMICOLON:       key.CodeSemicolon,
	C.AKEYCODE_APOSTROPHE:      key.CodeApostrophe,
	C.AKEYCODE_SLASH:           key.CodeSlash,
	C.AKEYCODE_PAGE_UP:         key.CodePageUp,
	C.AKEYCODE_PAGE_DOWN:       key.CodePageDown,
	C.AKEYCODE_ESCAPE:          key.CodeEscape,
	C.AKEYCODE_FORWARD_DEL:     key.CodeDeleteForward,
	C.AKEYCODE_CTRL_LEFT:       key.CodeLeftControl,
	C.AKEYCODE_CTRL_RIGHT:      key.CodeRightControl,
	C.AKEYCODE_CAPS_LOCK:       key.CodeCapsLock,
	C.AKEYCODE_META_LEFT:       key.CodeLeftGUI,
	C.AKEYCODE_META_RIGHT:      key.CodeRightGUI,
	C.AKEYCODE_INSERT:          key.CodeInsert,
	C.AKEYCODE_F1:              key.CodeF1,
	C.AKEYCODE_F2:              key.CodeF2,
	C.AKEYCODE_F3:              key.CodeF3,
	C.AKEYCODE_F4:              key.CodeF4,
	C.AKEYCODE_F5:              key.CodeF5,
	C.AKEYCODE_F6:              key.CodeF6,
	C.AKEYCODE_F7:              key.CodeF7,
	C.AKEYCODE_F8:              key.CodeF8,
	C.AKEYCODE_F9:              key.CodeF9,
	C.AKEYCODE_F10:             key.CodeF10,
	C.AKEYCODE_F11:             key.CodeF11,
	C.AKEYCODE_F12:             key.CodeF12,
	C.AKEYCODE_NUM_LOCK:        key.CodeKeypadNumLock,
	C.AKEYCODE_NUMPAD_0:        key.CodeKeypad0,
	C.AKEYCODE_NUMPAD_1:        key.CodeKeypad1,
	C.AKEYCODE_NUMPAD_2:        key.CodeKeypad2,
	C.AKEYCODE_NUMPAD_3:        key.CodeKeypad3,
	C.AKEYCODE_NUMPAD_4:        key.CodeKeypad4,
	C.AKEYCODE_NUMPAD_5:        key.CodeKeypad5,
	C.AKEYCODE_NUMPAD_6:        key.CodeKeypad6,
	C.AKEYCODE_NUMPAD_7:        key.CodeKeypad7,
	C.AKEYCODE_NUMPAD_8:        key.CodeKeypad8,
	C.AKEYCODE_NUMPAD_9:        key.CodeKeypad9,
	C.AKEYCODE_NUMPAD_DIVIDE:   key.CodeKeypadSlash,
	C.AKEYCODE_NUMPAD_MULTIPLY: key.CodeKeypadAsterisk,
	C.AKEYCODE_NUMPAD_SUBTRACT: key.CodeKeypadHyphenMinus,
	C.AKEYCODE_NUMPAD_ADD:      key.CodeKeypadPlusSign,
	C.AKEYCODE_NUMPAD_DOT:      key.CodeKeypadFullStop,
	C.AKEYCODE_NUMPAD_ENTER:    key.CodeKeypadEnter,
	C.AKEYCODE_NUMPAD_EQUALS:   key.CodeKeypadEqualSign,
	C.AKEYCODE_VOLUME_MUTE:     key.CodeMute,
}

func convAndroidKeyCode(aKeyCode int32) key.Code {
	if code, ok := androidKeycoe[aKeyCode]; ok {
		return code
	}
	return key.CodeUnknown
}

/*
	Many Android key codes do not map into USB HID codes.
	For those, key.CodeUnknown is returned. This switch has all
	cases, even the unknown ones, to serve as a documentation
	and search aid.
	C.AKEYCODE_UNKNOWN
	C.AKEYCODE_SOFT_LEFT
	C.AKEYCODE_SOFT_RIGHT
	C.AKEYCODE_BACK
	C.AKEYCODE_CALL
	C.AKEYCODE_ENDCALL
	C.AKEYCODE_STAR
	C.AKEYCODE_POUND
	C.AKEYCODE_DPAD_UP
	C.AKEYCODE_DPAD_DOWN
	C.AKEYCODE_DPAD_LEFT
	C.AKEYCODE_DPAD_RIGHT
	C.AKEYCODE_DPAD_CENTER
	C.AKEYCODE_POWER
	C.AKEYCODE_CAMERA
	C.AKEYCODE_CLEAR
	C.AKEYCODE_SYM
	C.AKEYCODE_EXPLORER
	C.AKEYCODE_ENVELOPE
	C.AKEYCODE_AT
	C.AKEYCODE_NUM
	C.AKEYCODE_HEADSETHOOK
	C.AKEYCODE_FOCUS
	C.AKEYCODE_PLUS
	C.AKEYCODE_MENU
	C.AKEYCODE_NOTIFICATION
	C.AKEYCODE_SEARCH
	C.AKEYCODE_MEDIA_PLAY_PAUSE
	C.AKEYCODE_MEDIA_STOP
	C.AKEYCODE_MEDIA_NEXT
	C.AKEYCODE_MEDIA_PREVIOUS
	C.AKEYCODE_MEDIA_REWIND
	C.AKEYCODE_MEDIA_FAST_FORWARD
	C.AKEYCODE_MUTE
	C.AKEYCODE_PICTSYMBOLS
	C.AKEYCODE_SWITCH_CHARSET
	C.AKEYCODE_BUTTON_A
	C.AKEYCODE_BUTTON_B
	C.AKEYCODE_BUTTON_C
	C.AKEYCODE_BUTTON_X
	C.AKEYCODE_BUTTON_Y
	C.AKEYCODE_BUTTON_Z
	C.AKEYCODE_BUTTON_L1
	C.AKEYCODE_BUTTON_R1
	C.AKEYCODE_BUTTON_L2
	C.AKEYCODE_BUTTON_R2
	C.AKEYCODE_BUTTON_THUMBL
	C.AKEYCODE_BUTTON_THUMBR
	C.AKEYCODE_BUTTON_START
	C.AKEYCODE_BUTTON_SELECT
	C.AKEYCODE_BUTTON_MODE
	C.AKEYCODE_SCROLL_LOCK
	C.AKEYCODE_FUNCTION
	C.AKEYCODE_SYSRQ
	C.AKEYCODE_BREAK
	C.AKEYCODE_MOVE_HOME
	C.AKEYCODE_MOVE_END
	C.AKEYCODE_FORWARD
	C.AKEYCODE_MEDIA_PLAY
	C.AKEYCODE_MEDIA_PAUSE
	C.AKEYCODE_MEDIA_CLOSE
	C.AKEYCODE_MEDIA_EJECT
	C.AKEYCODE_MEDIA_RECORD
	C.AKEYCODE_NUMPAD_COMMA
	C.AKEYCODE_NUMPAD_LEFT_PAREN
	C.AKEYCODE_NUMPAD_RIGHT_PAREN
	C.AKEYCODE_INFO
	C.AKEYCODE_CHANNEL_UP
	C.AKEYCODE_CHANNEL_DOWN
	C.AKEYCODE_ZOOM_IN
	C.AKEYCODE_ZOOM_OUT
	C.AKEYCODE_TV
	C.AKEYCODE_WINDOW
	C.AKEYCODE_GUIDE
	C.AKEYCODE_DVR
	C.AKEYCODE_BOOKMARK
	C.AKEYCODE_CAPTIONS
	C.AKEYCODE_SETTINGS
	C.AKEYCODE_TV_POWER
	C.AKEYCODE_TV_INPUT
	C.AKEYCODE_STB_POWER
	C.AKEYCODE_STB_INPUT
	C.AKEYCODE_AVR_POWER
	C.AKEYCODE_AVR_INPUT
	C.AKEYCODE_PROG_RED
	C.AKEYCODE_PROG_GREEN
	C.AKEYCODE_PROG_YELLOW
	C.AKEYCODE_PROG_BLUE
	C.AKEYCODE_APP_SWITCH
	C.AKEYCODE_BUTTON_1
	C.AKEYCODE_BUTTON_2
	C.AKEYCODE_BUTTON_3
	C.AKEYCODE_BUTTON_4
	C.AKEYCODE_BUTTON_5
	C.AKEYCODE_BUTTON_6
	C.AKEYCODE_BUTTON_7
	C.AKEYCODE_BUTTON_8
	C.AKEYCODE_BUTTON_9
	C.AKEYCODE_BUTTON_10
	C.AKEYCODE_BUTTON_11
	C.AKEYCODE_BUTTON_12
	C.AKEYCODE_BUTTON_13
	C.AKEYCODE_BUTTON_14
	C.AKEYCODE_BUTTON_15
	C.AKEYCODE_BUTTON_16
	C.AKEYCODE_LANGUAGE_SWITCH
	C.AKEYCODE_MANNER_MODE
	C.AKEYCODE_3D_MODE
	C.AKEYCODE_CONTACTS
	C.AKEYCODE_CALENDAR
	C.AKEYCODE_MUSIC
	C.AKEYCODE_CALCULATOR

	Defined in an NDK API version beyond what we use today:
	C.AKEYCODE_ASSIST
	C.AKEYCODE_BRIGHTNESS_DOWN
	C.AKEYCODE_BRIGHTNESS_UP
	C.AKEYCODE_RO
	C.AKEYCODE_YEN
	C.AKEYCODE_ZENKAKU_HANKAKU
*/
