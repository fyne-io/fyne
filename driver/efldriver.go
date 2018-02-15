package driver

// #cgo pkg-config: ecore ecore-evas evas
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <Evas.h>
import "C"

import "log"
import "image"
import "image/color"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/theme"

type EFLDriver struct {
}

type window struct {
	ee     *C.Ecore_Evas
	canvas ui.Canvas
	driver *EFLDriver
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
	evas  *C.Evas
	w, h  int
	scale float32
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

func (c *eflCanvas) addObject(o ui.CanvasObject) {
	obj := buildCanvasObject(c, o)

	C.evas_object_geometry_set(obj, 0, 0, C.Evas_Coord(scaleInt(c, c.w)), C.Evas_Coord(scaleInt(c, c.h)))
	C.evas_object_show(obj)
}

func (c *eflCanvas) SetContent(o ui.CanvasObject) {
	switch o.(type) {
	case *ui.Container:
		c.addObject(ui.NewRectangle(theme.BackgroundColor()))
		container := o.(*ui.Container)

		for _, child := range container.Objects {
			c.addObject(child)
		}

		if container.Layout != nil {
			container.Layout.Layout(container)
		}
	default:
		c.addObject(o)
	}
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
	env := C.getenv(C.CString("WAYLAND_DISPLAY"))

	if env == nil {
		log.Println("Unable to connect to Wayland - attempting X")
		return X11EngineName()
	}

	return WaylandEngineName()
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

func (d *EFLDriver) CreateWindow(title string) ui.Window {
	engine := findEngineName()
	size := image.Pt(300, 200)

	C.evas_init()
	C.ecore_init()
	C.ecore_evas_init()

	w := &window{
		ee:     C.ecore_evas_new(C.CString(engine), 0, 0, 100, 100, nil),
		driver: d,
	}
	c := &eflCanvas{
		evas:  C.ecore_evas_get(w.ee),
		w:     size.X,
		h:     size.Y,
		scale: scaleByDPI(w),
	}
	w.canvas = c
	C.ecore_evas_resize(w.ee, C.int(scaleInt(c, size.X)), C.int(scaleInt(c, size.Y)))

	if engine == WaylandEngineName() {
		WaylandWindowInit(w)
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
