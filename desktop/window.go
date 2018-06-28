// +build !ci

package desktop

// #cgo pkg-config: ecore ecore-evas ecore-input evas
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <Ecore_Input.h>
// #include <Evas.h>
//
// void onWindowResize_cgo(Ecore_Evas *);
// void onWindowMove_cgo(Ecore_Evas *);
// void onWindowFocusIn_cgo(Ecore_Evas *);
// void onWindowFocusOut_cgo(Ecore_Evas *);
// void onWindowClose_cgo(Ecore_Evas *);
import "C"

import "log"
import "os"
import "strconv"
import "unsafe"

import "github.com/fyne-io/fyne"

type window struct {
	ee      *C.Ecore_Evas
	canvas  fyne.Canvas
	master  bool
	focused bool
}

func init() {
	C.evas_init()
	C.ecore_init()
	C.ecore_evas_init()
}

var windows = make(map[*C.Ecore_Evas]*window)

func (w *window) Title() string {
	return C.GoString(C.ecore_evas_title_get(w.ee))
}

func (w *window) SetTitle(title string) {
	C.ecore_evas_title_set(w.ee, C.CString(title))
}

func (w *window) Fullscreen() bool {
	return C.ecore_evas_fullscreen_get(w.ee) != 0
}

func (w *window) SetFullscreen(full bool) {
	if full {
		C.ecore_evas_fullscreen_set(w.ee, 1)
	} else {
		C.ecore_evas_fullscreen_set(w.ee, 0)
	}
}

func (w *window) Show() {
	C.ecore_evas_show(w.ee)

	if len(windows) == 1 {
		w.master = true
		initEFL()
	}
}

func (w *window) Hide() {
	C.ecore_evas_hide(w.ee)
}

func (w *window) Close() {
	w.Hide()

	if w.master || len(windows) == 1 {
		DoQuit()
	} else {
		delete(windows, w.ee)
	}
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
	if xdpi > 250 {
		return float32(1.5)
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

	canvas := w.canvas.(*eflCanvas)
	canvas.size = fyne.NewSize(int(float32(ww)/canvas.Scale()), int(float32(hh)/canvas.Scale()))
	canvas.resizeContent()
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

	canvas := w.canvas.(*eflCanvas)

	if canvas.focused == nil {
		return
	}

	canvas.focused.OnFocusGained()
}

//export onWindowFocusLost
func onWindowFocusLost(ee *C.Ecore_Evas) {
	w := windows[ee]

	// we may be closing the window
	if w == nil {
		return
	}

	canvas := w.canvas.(*eflCanvas)
	if canvas.focused == nil {
		return
	}

	canvas.focused.OnFocusLost()
}

//export onWindowClose
func onWindowClose(ee *C.Ecore_Evas) {
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
	ev.Name = C.GoString(info.keyname)
	ev.Code = fyne.KeyCode(int(info.keycode))
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
	evas := C.ecore_evas_new(C.CString(engine), 0, 0, 10, 10, nil)
	if evas == nil {
		log.Fatalln("Unable to create canvas, perhaps missing module for", engine)
	}

	w := &window{
		ee: evas,
	}
	oSWindowInit(w)
	windows[w.ee] = w

	w.canvas = &eflCanvas{
		evas:   C.ecore_evas_get(evas),
		scale:  1.0,
		window: w,
	}

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
