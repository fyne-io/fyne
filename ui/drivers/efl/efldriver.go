package efl

// #cgo pkg-config: ecore ecore-evas ecore-wl2 evas
// #cgo CFLAGS: -DEFL_BETA_API_SUPPORT=1
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <Ecore_Wl2.h>
// #include <Evas.h>
import "C"

import "fmt"
import "image"
import "image/color"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/theme"

type window struct {
	ee *C.Ecore_Evas
	canvas ui.Canvas
}

func (w *window) Init(ee *C.Ecore_Evas) {
	w.ee = ee
}

func (w *window) Title() string {
	return C.GoString(C.ecore_evas_title_get(w.ee))
}

func (w *window) SetTitle(title string) {
	C.ecore_evas_title_set(w.ee, C.CString(title))
}

func (w *window) Show() {
	C.ecore_evas_show(w.ee)
}

func (w *window) Hide() {
	C.ecore_evas_hide(w.ee)
}

func (w *window) Canvas() (ui.Canvas) {
	return w.canvas
}

type canvasobject struct {
	obj *C.Evas_Object
}

func (o *canvasobject) SetColor(c color.RGBA) {
	C.evas_object_color_set(o.obj, C.int(c.R), C.int(c.G), C.int(c.B), C.int(c.A))
}

type canvas struct {
	evas *C.Evas
}

func (c *canvas) NewRectangle(r image.Rectangle) (ui.CanvasObject) {
	o := &canvasobject{
		obj: C.evas_object_rectangle_add(c.evas),
	}

	C.evas_object_geometry_set(o.obj, C.Evas_Coord(r.Min.X), C.Evas_Coord(r.Min.Y),
				   C.Evas_Coord(r.Max.X - r.Min.X), C.Evas_Coord(r.Max.Y - r.Min.Y))
	C.evas_object_show(o.obj)
	return o
}

type EFLDriver struct {
}

func (d EFLDriver) CreateWindow(title string) ui.Window {
	engine := "wayland_shm"
	env := C.getenv(C.CString("WAYLAND_DISPLAY"))

	if env == nil {
		fmt.Println("Unable to connect to Wayland - attempting X")
		engine = "software_x11"
	}

	C.evas_init()
	C.ecore_init()
	C.ecore_evas_init()

	w := &window{
		ee: C.ecore_evas_new(C.CString(engine), 10, 10, 300, 200, nil),
	}
	c := &canvas{
		evas: C.ecore_evas_get(w.ee),
	}
	w.canvas = c

	if engine == "wayland_shm" {
		win := C.ecore_evas_wayland2_window_get(w.ee)
		C.ecore_wl2_window_type_set(win, C.ECORE_WL2_WINDOW_TYPE_TOPLEVEL)
	}

	bg := c.NewRectangle(image.Rect(0, 0, 300, 200))
	bg.SetColor(theme.BackgroundColor())

	w.SetTitle(title)
	w.Show()
	return w
}

func (d EFLDriver) Run() {
	C.ecore_main_loop_begin()
}
