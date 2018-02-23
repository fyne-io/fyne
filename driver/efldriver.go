package driver

// #cgo pkg-config: ecore ecore-evas evas
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <Evas.h>
//
// void onWindowResize_cgo(Ecore_Evas *);
import "C"

import "log"
import "image/color"
import "unsafe"
import "runtime"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"

type EFLDriver struct {
}

type window struct {
	ee     *C.Ecore_Evas
	canvas ui.Canvas
	driver *EFLDriver
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
	w.driver.Run()
}

func (w *window) Hide() {
	C.ecore_evas_hide(w.ee)
}

func (w *window) Canvas() ui.Canvas {
	return w.canvas
}

type canvasobject struct {
	obj    *C.Evas_Object
	canvas ui.Canvas
}

func (o *canvasobject) SetColor(c color.RGBA) {
	C.evas_object_color_set(o.obj, C.int(c.R), C.int(c.G), C.int(c.B), C.int(c.A))
}

func (o *canvasobject) Canvas() ui.Canvas {
	return o.canvas
}

type eflCanvas struct {
	ui.Canvas
	evas    *C.Evas
	size    ui.Size
	scale   float32
	content ui.CanvasObject
}

func buildCanvasObject(c *eflCanvas, o ui.CanvasObject) *C.Evas_Object {
	var obj *C.Evas_Object

	switch o.(type) {
	case *ui.TextObject:
		obj = C.evas_object_text_add(c.evas)

		to, _ := o.(*ui.TextObject)
		C.evas_object_color_set(obj, C.int(to.Color.R), C.int(to.Color.G),
			C.int(to.Color.B), C.int(to.Color.A))

		C.evas_object_text_text_set(obj, C.CString(to.Text))
		updateFont(obj, c, to)
	case *ui.RectangleObject:
		obj = C.evas_object_rectangle_add(c.evas)

		ro, _ := o.(*ui.RectangleObject)
		C.evas_object_color_set(obj, C.int(ro.Color.R), C.int(ro.Color.G),
			C.int(ro.Color.B), C.int(ro.Color.A))
	default:
		log.Println("Unrecognised Object", o)
	}

	return obj
}

func (c *eflCanvas) SetContent(o ui.CanvasObject) {
	switch o.(type) {
	case *ui.Container:
		obj := buildCanvasObject(c, ui.NewRectangle(theme.BackgroundColor()))
		C.evas_object_geometry_set(obj, 0, 0, C.Evas_Coord(scaleInt(c, c.size.Width)), C.Evas_Coord(scaleInt(c, c.size.Height)))
		C.evas_object_show(obj)

		container := o.(*ui.Container)

		var objs = make([]*C.Evas_Object, len(container.Objects))
		for i, child := range container.Objects {
			obj := buildCanvasObject(c, child)
			objs[i] = obj
		}

		if container.Layout != nil {
			container.Layout.Layout(container, c.size)
		} else {
			layout.NewMaxLayout().Layout(container, c.size)
		}

		for i, child := range container.Objects {
			size := child.CurrentSize()
			pos := child.CurrentPosition()
			C.evas_object_geometry_set(objs[i], C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
				C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
			C.evas_object_show(objs[i])
		}
	default:
		obj := buildCanvasObject(c, o)
		C.evas_object_geometry_set(obj, 0, 0, C.Evas_Coord(scaleInt(c, c.size.Width)), C.Evas_Coord(scaleInt(c, c.size.Height)))
		C.evas_object_show(obj)
	}

	c.content = o
}

func updateFont(obj *C.Evas_Object, c *eflCanvas, t *ui.TextObject) {
	font := theme.TextFont()

	if t.Bold {
		if t.Italic {
			font = theme.TextBoldItalicFont()
		} else {
			font = theme.TextBoldFont()
		}
	} else if t.Italic {
		font = theme.TextItalicFont()
	}

	C.evas_object_text_font_set(obj, C.CString(font), C.Evas_Font_Size(scaleInt(c, t.FontSize)))
}

func scaleInt(c ui.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) * c.Scale())
	}
}

func (c *eflCanvas) Scale() float32 {
	return c.scale
}

func (c *eflCanvas) SetScale(scale float32) {
	c.scale = scale
	log.Println("TODO Update all our objects")
}

func findEngineName() string {
	if runtime.GOOS == "darwin" {
		return CocoaEngineName()
	}

	env := C.getenv(C.CString("WAYLAND_DISPLAY"))

	if env != nil {
		log.Println("Wayland support is currently disabled - attempting XWayland")
	}

	return X11EngineName()
}

func scaleByDPI(w *window) float32 {
	xdpi := C.int(0)

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
	canvas.SetContent(canvas.content)
}

func (d *EFLDriver) CreateWindow(title string) ui.Window {
	engine := findEngineName()
	size := ui.NewSize(300, 200)

	C.evas_init()
	C.ecore_init()
	C.ecore_evas_init()

	w := &window{
		ee:     C.ecore_evas_new(C.CString(engine), 0, 0, 100, 100, nil),
		driver: d,
	}
	c := &eflCanvas{
		evas:  C.ecore_evas_get(w.ee),
		size:  size,
		scale: scaleByDPI(w),
	}
	w.canvas = c
	windows[w.ee] = w
	C.ecore_evas_resize(w.ee, C.int(scaleInt(c, size.Width)), C.int(scaleInt(c, size.Height)))
	C.ecore_evas_callback_resize_set(w.ee, (C.Ecore_Evas_Event_Cb)(unsafe.Pointer(C.onWindowResize_cgo)))

	if engine == WaylandEngineName() {
		WaylandWindowInit(w)
	} else if engine == CocoaEngineName() {
		CocoaWindowInit(w)
	} else {
		X11WindowInit(w)
	}

	c.SetContent(new(ui.Container))

	w.SetTitle(title)
	return w
}

func (d *EFLDriver) Run() {
	C.ecore_main_loop_begin()
}
