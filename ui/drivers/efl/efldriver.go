package efl

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

type window struct {
	ee     *C.Ecore_Evas
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

func (w *window) Canvas() ui.Canvas {
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

func (c *canvas) NewRectangle(r image.Rectangle) ui.CanvasObject {
	o := &canvasobject{
		obj: C.evas_object_rectangle_add(c.evas),
	}

	C.evas_object_geometry_set(o.obj, C.Evas_Coord(r.Min.X), C.Evas_Coord(r.Min.Y),
		C.Evas_Coord(r.Max.X-r.Min.X), C.Evas_Coord(r.Max.Y-r.Min.Y))
	C.evas_object_show(o.obj)
	return o
}

type canvasTextObject struct {
	*canvasobject
	size         int
	italic, bold bool
}

func (t *canvasTextObject) updateFont() {
	font := theme.TextFont()

	if t.bold {
		if t.italic {
			font = theme.TextBoldItalicFont()
		} else {
			font = theme.TextBoldFont()
		}
	} else if t.italic {
		font = theme.TextItalicFont()
	}

	C.evas_object_text_font_set(t.canvasobject.obj, C.CString(font), C.Evas_Font_Size(t.size))
}

func (t *canvasTextObject) FontSize() int {
	return t.size
}

func (t *canvasTextObject) SetFontSize(size int) {
	t.size = size
	t.updateFont()
}

func (t *canvasTextObject) Bold() bool {
	return t.bold
}

func (t *canvasTextObject) SetBold(bold bool) {
	t.bold = bold
	t.updateFont()
}

func (t *canvasTextObject) Italic() bool {
	return t.italic
}

func (t *canvasTextObject) SetItalic(italic bool) {
	t.italic = italic
	t.updateFont()
}

func (c *canvas) NewText(text string) ui.CanvasTextObject {
	o := &canvasTextObject{
		&canvasobject{
			obj: C.evas_object_text_add(c.evas),
		},
		theme.TextSize(),
		false,
		false,
	}

	C.evas_object_text_text_set(o.obj, C.CString(text))
	o.SetColor(theme.TextColor())
	o.updateFont()

	C.evas_object_show(o.obj)
	return o
}

type EFLDriver struct {
}

func findEngineName() string {
	env := C.getenv(C.CString("WAYLAND_DISPLAY"))

	if env == nil {
		log.Println("Unable to connect to Wayland - attempting X")
		return X11EngineName()
	}

	return WaylandEngineName()
}

func (d EFLDriver) CreateWindow(title string) ui.Window {
	engine := findEngineName()

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

	if engine == WaylandEngineName() {
		WaylandWindowInit(w)
	} else {
		X11WindowInit(w)
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
