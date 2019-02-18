// +build !ci,efl

package efl

// #cgo pkg-config: ecore ecore-evas ecore-input evas
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <Ecore_Input.h>
// #include <Evas.h>
// #include <stdlib.h>
//
// void onWindowResize_cgo(Ecore_Evas *);
// void onWindowMove_cgo(Ecore_Evas *);
// void onWindowFocusIn_cgo(Ecore_Evas *);
// void onWindowFocusOut_cgo(Ecore_Evas *);
// void onWindowClose_cgo(Ecore_Evas *);
//
// void force_render();
import "C"

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"unsafe"

	"fyne.io/fyne"
)

type window struct {
	ee     *C.Ecore_Evas
	canvas fyne.Canvas
	icon   fyne.Resource

	master    bool
	focused   bool
	fixedSize bool
	padded    bool

	onClosed func()
}

func init() {
	// This is a workaround for a logged issue, phab.enlightenment.org/T7099
	os.Setenv("EVAS_DRM_BUFFERS", "5")

	C.evas_init()
	C.ecore_init()
	C.ecore_evas_init()
}

var windows = make(map[*C.Ecore_Evas]*window)

func (w *window) Title() string {
	return C.GoString(C.ecore_evas_title_get(w.ee))
}

func (w *window) SetTitle(title string) {
	runOnMain(func() {
		cstr := C.CString(title)
		C.ecore_evas_title_set(w.ee, cstr)
		C.free(unsafe.Pointer(cstr))
	})
}

func (w *window) FullScreen() bool {
	return C.ecore_evas_fullscreen_get(w.ee) != 0
}

func (w *window) SetFullScreen(full bool) {
	runOnMain(func() {
		if full {
			C.ecore_evas_fullscreen_set(w.ee, 1)
		} else {
			C.ecore_evas_fullscreen_set(w.ee, 0)
		}
	})
}

func (w *window) CenterOnScreen() {
	winW, winH := w.sizeOnScreen()

	runOnMain(func() {
		var screenW, screenH C.int
		C.ecore_evas_screen_geometry_get(w.ee, nil, nil, &screenW, &screenH)

		C.ecore_evas_move(w.ee, screenW/2-C.int(winW)/2, screenH/2-C.int(winH)/2)
	})
}

// sizeOnScreen gets the size of a window content in screen pixels
func (w *window) sizeOnScreen() (int, int) {
	// get current size of content inside the window
	winContentSize := w.Content().MinSize()
	// content size can be scaled, so factor that in to determining window size
	scale := w.canvas.Scale()

	// calculate how many pixels will be used at this scale
	viewWidth := int(float32(winContentSize.Width) * scale)
	viewHeight := int(float32(winContentSize.Height) * scale)

	return viewWidth, viewHeight
}

func (w *window) Resize(size fyne.Size) {
	scale := w.canvas.Scale()
	runOnMain(func() {
		C.ecore_evas_resize(w.ee, C.int(float32(size.Width)*scale), C.int(float32(size.Height)*scale))
	})
}

func (w *window) FixedSize() bool {
	return w.fixedSize
}

func (w *window) SetFixedSize(fixed bool) {
	w.fixedSize = fixed
}

func (w *window) Padded() bool {
	return w.padded
}

func (w *window) SetPadded(padded bool) {
	w.padded = padded
}

func (w *window) Icon() fyne.Resource {
	if w.icon == nil {
		return fyne.CurrentApp().Icon()
	}

	return w.icon
}

func (w *window) SetIcon(icon fyne.Resource) {
	w.icon = icon
}

func (w *window) SetOnClosed(closed func()) {
	w.onClosed = closed
}

func (w *window) Show() {
	runOnMain(func() {
		C.ecore_evas_show(w.ee)
	})

	if len(windows) == 1 {
		w.master = true
	}
}

func (w *window) Hide() {
	runOnMain(func() {
		C.ecore_evas_hide(w.ee)
	})
}

func (w *window) Close() {
	w.Hide()
	if w.onClosed != nil {
		w.onClosed()
	}

	if w.master || len(windows) == 1 {
		DoQuit()
	} else {
		delete(windows, w.ee)
	}
}

func (w *window) ShowAndRun() {
	w.Show()
	runEFL()
}

func (w *window) Content() fyne.CanvasObject {
	return w.Canvas().Content()
}

func (w *window) SetContent(obj fyne.CanvasObject) {
	w.Canvas().SetContent(obj)
}

func (w *window) Canvas() fyne.Canvas {
	return w.canvas
}

func scaleByDPI(w *window) float32 {
	xdpi := C.int(0)

	env := os.Getenv("FYNE_SCALE")
	if env != "" {
		scale, _ := strconv.ParseFloat(env, 32)
		return float32(scale)
	}
	C.ecore_evas_screen_dpi_get(w.ee, &xdpi, nil)
	if xdpi > 1000 { // assume that this is a mistake and bail
		return float32(1.0)
	}

	if xdpi > 192 {
		return float32(1.5)
	} else if xdpi > 144 {
		return float32(1.35)
	} else if xdpi > 120 {
		return float32(1.2)
	}

	return float32(1.0)
}

//export onWindowResize
func onWindowResize(ee *C.Ecore_Evas) {
	w := windows[ee]
	if w == nil {
		return
	}

	var ww, hh C.int
	C.ecore_evas_geometry_get(ee, nil, nil, &ww, &hh)

	current := w.canvas.(*eflCanvas)
	current.size = fyne.NewSize(int(float32(ww)/current.Scale()), int(float32(hh)/current.Scale()))
	current.resizeContent()

	if runtime.GOOS == "darwin" {
		// due to NSRunLoop freezing during window resize we need to force a refresh
		runOnMain(func() {
			current.dirty[current.Content()] = true
			renderCycle()
		})
		C.force_render()
	}
}

//export onWindowMove
func onWindowMove(ee *C.Ecore_Evas) {
	w := windows[ee]
	if w == nil {
		return
	}

	canvas := w.canvas.(*eflCanvas)

	scale := scaleByDPI(w)
	if scale != canvas.Scale() {
		canvas.SetScale(scaleByDPI(w))
	}
}

//export onWindowFocusGained
func onWindowFocusGained(ee *C.Ecore_Evas) {
	w := windows[ee]
	if w == nil {
		return
	}

	if canFocus, ok := w.canvas.(*eflCanvas); ok && canFocus.focused != nil {
		canFocus.focused.FocusGained()
	}
}

//export onWindowFocusLost
func onWindowFocusLost(ee *C.Ecore_Evas) {
	w := windows[ee]

	// we may be closing the window
	if w == nil {
		return
	}

	if canFocus, ok := w.canvas.(*eflCanvas); ok && canFocus.focused != nil {
		canFocus.focused.FocusLost()
	}
}

//export onWindowClose
func onWindowClose(ee *C.Ecore_Evas) {
	// TODO notify any onclose...
	windows[ee].Close()
}

//export onWindowKeyDown
func onWindowKeyDown(ew C.Ecore_Window, info *C.Ecore_Event_Key) {
	if ew == 0 {
		log.Println("Keystroke missing window")
		return
	}

	var w *window
	for _, win := range windows {
		if C.ecore_evas_window_get(win.ee) == ew {
			w = win
		}
	}

	if w == nil {
		log.Println("Window not found")
		return
	}
	canvas := w.canvas.(*eflCanvas)

	if canvas.focused == nil && canvas.onTypedRune == nil && canvas.onTypedKey == nil {
		return
	}

	ev := new(fyne.KeyEvent)
	str := C.GoString(info.string)
	if str != "" && []rune(str)[0] < ' ' {
		str = ""
	}
	ev.Name = fyne.KeyName(C.GoString(info.keyname))

	if canvas.focused != nil {
		if str != "" {
			canvas.focused.TypedRune([]rune(str)[0])
		} else {
			canvas.focused.TypedKey(ev)
		}
	}
	if str != "" {
		if canvas.onTypedRune != nil {
			canvas.onTypedRune([]rune(str)[0])
		}
	} else {
		if canvas.onTypedKey != nil {
			canvas.onTypedKey(ev)
		}
	}
}

// CreateWindowWithEngine will create a new efl backed window using the specified
// engine name. The possible options for EFL engines is out of scope of this
// documentation and can be found on the http://enlightenment.org website.
// USE OF THIS METHOD IS NOT RECOMMENDED
func CreateWindowWithEngine(engine string) fyne.Window {
	cstr := C.CString(engine)
	w := &window{padded: true}

	runOnMain(func() {
		evas := C.ecore_evas_new(cstr, 0, 0, 10, 10, nil)
		C.free(unsafe.Pointer(cstr))
		if evas == nil {
			log.Fatalln("Unable to create canvas, perhaps missing module for", engine)
		}

		w.ee = evas
		oSWindowInit(w)
		windows[w.ee] = w

		w.canvas = &eflCanvas{
			evas:   C.ecore_evas_get(evas),
			scale:  scaleByDPI(w),
			window: w,
		}
	})

	return w
}

func (d *eFLDriver) CreateWindow(title string) fyne.Window {
	win := CreateWindowWithEngine(oSEngineName())
	win.SetTitle(title)

	w := win.(*window)
	C.ecore_evas_callback_resize_set(w.ee, (C.Ecore_Evas_Event_Cb)(unsafe.Pointer(C.onWindowResize_cgo)))
	C.ecore_evas_callback_move_set(w.ee, (C.Ecore_Evas_Event_Cb)(unsafe.Pointer(C.onWindowMove_cgo)))
	C.ecore_evas_callback_focus_in_set(w.ee, (C.Ecore_Evas_Event_Cb)(unsafe.Pointer(C.onWindowFocusIn_cgo)))
	C.ecore_evas_callback_focus_out_set(w.ee, (C.Ecore_Evas_Event_Cb)(unsafe.Pointer(C.onWindowFocusOut_cgo)))
	C.ecore_evas_callback_delete_request_set(w.ee, (C.Ecore_Evas_Event_Cb)(unsafe.Pointer(C.onWindowClose_cgo)))

	w.SetContent(new(fyne.Container))
	return w
}

func (d *eFLDriver) AllWindows() []fyne.Window {
	wins := make([]fyne.Window, 0, len(windows))

	for _, win := range windows {
		wins = append(wins, win)
	}

	return wins
}
