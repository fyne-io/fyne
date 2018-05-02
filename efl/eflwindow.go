package efl

// #cgo pkg-config: ecore ecore-evas evas
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <Evas.h>
//
// void onWindowResize_cgo(Ecore_Evas *);
// void onWindowClose_cgo(Ecore_Evas *);
import "C"

import "log"
import "os"
import "strconv"
import "unsafe"

import "github.com/fyne-io/fyne/ui"

type window struct {
	ee     *C.Ecore_Evas
	canvas ui.Canvas
	driver *eFLDriver
	master bool
}

var windows = make(map[*C.Ecore_Evas]*window)

func (w *window) Title() string {
	return C.GoString(C.ecore_evas_title_get(w.ee))
}

func (w *window) SetTitle(title string) {
	C.ecore_evas_title_set(w.ee, C.CString(title))
}

func (w *window) Show() {
	C.ecore_evas_show(w.ee)

	w.master = len(windows) == 1
	if !w.driver.running {
		w.driver.Run()
	}
}

func (w *window) Hide() {
	C.ecore_evas_hide(w.ee)
}

func (w *window) Close() {
	w.Hide()

	if w.master || len(windows) == 1 {
		w.driver.Quit()
	} else {
		delete(windows, w.ee)
	}
}

func (w *window) Canvas() ui.Canvas {
	return w.canvas
}

func scaleByDPI(w *window) float32 {
	xdpi := C.int(0)

	env := os.Getenv("FYNE_SCALE")
	if env != "" {
		scale, _ := strconv.ParseFloat(env, 32)
		log.Println("Scale specified, rendering at", scale)
		return float32(scale)
	}
	C.ecore_evas_screen_dpi_get(w.ee, &xdpi, nil)
	if xdpi > 96 {
		log.Println("High DPI", xdpi, "- scaling to 1.5")
		return float32(1.5)
	}

	return float32(1.0)
}

//export onWindowResize
func onWindowResize(ee *C.Ecore_Evas) {
	var ww, hh C.int
	C.ecore_evas_geometry_get(ee, nil, nil, &ww, &hh)

	w := windows[ee]

	canvas := w.canvas.(*eflCanvas)
	canvas.size = ui.NewSize(int(float32(ww)/canvas.Scale()), int(float32(hh)/canvas.Scale()))
	canvas.Refresh(canvas.content)
}

//export onWindowClose
func onWindowClose(ee *C.Ecore_Evas) {
	windows[ee].Close()
}

func (d *eFLDriver) CreateWindow(title string) ui.Window {
	engine := oSEngineName()

	C.evas_init()
	C.ecore_init()
	C.ecore_evas_init()

	evas := C.ecore_evas_new(C.CString(engine), 0, 0, 10, 10, nil)
	if evas == nil {
		log.Fatalln("Unable to create canvas, perhaps missing module for", engine)
	}

	w := &window{
		ee:     evas,
		driver: d,
	}
	w.SetTitle(title)
	oSWindowInit(w)
	c := &eflCanvas{
		scale:  scaleByDPI(w),
		evas:   C.ecore_evas_get(evas),
		window: w,
	}
	w.canvas = c
	windows[w.ee] = w
	C.ecore_evas_callback_resize_set(w.ee, (C.Ecore_Evas_Event_Cb)(unsafe.Pointer(C.onWindowResize_cgo)))
	C.ecore_evas_callback_delete_request_set(w.ee, (C.Ecore_Evas_Event_Cb)(unsafe.Pointer(C.onWindowClose_cgo)))

	c.SetContent(new(ui.Container))
	return w
}

func (d *eFLDriver) AllWindows() []ui.Window {
	wins := make([]ui.Window, 0, len(windows))

	for _, win := range windows {
		wins = append(wins, win)
	}

	return wins
}
