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

	"github.com/fyne-io/fyne"
)

type window struct {
	ee     *C.Ecore_Evas
	canvas fyne.Canvas
	icon   fyne.Resource

	master    bool
	focused   bool
	fixedSize bool

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

func (w *window) doShow(warn bool) {
	runOnMain(func() {
		C.ecore_evas_show(w.ee)
	})

	if len(windows) == 1 {
		w.master = true

		if warn {
			log.Println("window.Show() no longer blocks to run the application.")
			log.Println("If this program quit immediately try window.ShowAndRun().")
		}
	}
}

func (w *window) Show() {
	w.doShow(true)
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
	w.doShow(false)
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
		canFocus.focused.OnFocusGained()
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
		canFocus.focused.OnFocusLost()
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

	if canvas.focused == nil && canvas.onKeyDown == nil {
		return
	}

	ev := new(fyne.KeyEvent)
	ev.String = C.GoString(info.string)
	switch C.GoString(info.keyname) {
	case "Shift_L":
		fallthrough
	case "Shift_R":
		ev.Name = fyne.KeyShift
	case "Control_L":
		fallthrough
	case "Control_R":
		ev.Name = fyne.KeyControl
	case "Alt_L":
		fallthrough
	case "Alt_R":
		ev.Name = fyne.KeyAlt
	case "Super_L":
		fallthrough
	case "Super_R":
		ev.Name = fyne.KeySuper

	default:
		ev.Name = fyne.KeyName(C.GoString(info.keyname))
	}

	if (info.modifiers & C.ECORE_EVENT_MODIFIER_SHIFT) != 0 {
		ev.Modifiers |= fyne.ShiftModifier
	}
	if (info.modifiers & C.ECORE_EVENT_MODIFIER_CTRL) != 0 {
		ev.Modifiers |= fyne.ControlModifier
	}
	if (info.modifiers & C.ECORE_EVENT_MODIFIER_ALT) != 0 {
		ev.Modifiers |= fyne.AltModifier
	}

	if canvas.focused != nil {
		canvas.focused.OnKeyDown(ev)
	}
	if canvas.onKeyDown != nil {
		canvas.onKeyDown(ev)
	}
}

// CreateWindowWithEngine will create a new efl backed window using the specified
// engine name. The possible options for EFL engines is out of scope of this
// documentation and can be found on the http://enlightenment.org website.
// USE OF THIS METHOD IS NOT RECOMMENDED
func CreateWindowWithEngine(engine string) fyne.Window {
	cstr := C.CString(engine)
	w := &window{}

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
