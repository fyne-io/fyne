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
	canvas ui.Canvas
}

func (o *canvasobject) SetColor(c color.RGBA) {
	C.evas_object_color_set(o.obj, C.int(c.R), C.int(c.G), C.int(c.B), C.int(c.A))
}

func (o *canvasobject) Canvas() ui.Canvas {
	return o.canvas
}

type canvas struct {
	evas *C.Evas
	scale float32
}

func (c *canvas) NewRectangle(r image.Rectangle) ui.CanvasObject {
	o := &canvasobject{
		obj: C.evas_object_rectangle_add(c.evas),
		canvas: c,
	}

	C.evas_object_geometry_set(o.obj, C.Evas_Coord(scaleInt(c, r.Min.X)), C.Evas_Coord(scaleInt(c, r.Min.Y)),
		C.Evas_Coord(scaleInt(c, r.Max.X-r.Min.X)), C.Evas_Coord(scaleInt(c, r.Max.Y-r.Min.Y)))
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

	C.evas_object_text_font_set(t.canvasobject.obj, C.CString(font), C.Evas_Font_Size(scaleInt(t.Canvas(), t.size)))
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
			canvas: c,
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

func scaleInt(c ui.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) * c.Scale())
	}
}

func (c *canvas) Scale() float32 {
	return c.scale
}

func (c *canvas) SetScale(scale float32) {
	c.scale = scale
	log.Println("TODO Update all our objects")
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

func scaleByDPI(w *window) float32 {
	xdpi := C.int(0)

	C.ecore_evas_screen_dpi_get(w.ee, &xdpi, nil)
	if (xdpi > 96) {
		log.Println("High DPI", xdpi, "- scaling to 1.5")
		return float32(1.5)
	}

	return float32(1.0)
}

func (d EFLDriver) CreateWindow(title string) ui.Window {
	engine := findEngineName()
	size := image.Pt(300, 200)

	C.evas_init()
	C.ecore_init()
	C.ecore_evas_init()

	w := &window{
		ee: C.ecore_evas_new(C.CString(engine), 0, 0, 100, 100, nil),
	}
	c := &canvas{
		evas: C.ecore_evas_get(w.ee),
		scale: scaleByDPI(w),
	}
	w.canvas = c
	C.ecore_evas_resize(w.ee, C.int(scaleInt(c, size.X)), C.int(scaleInt(c, size.Y)))

	if engine == WaylandEngineName() {
		WaylandWindowInit(w)
	} else {
		X11WindowInit(w)
	}

	bg := c.NewRectangle(image.Rect(0, 0, size.X, size.Y))
	bg.SetColor(theme.BackgroundColor())

	w.SetTitle(title)
	w.Show()
	return w
}

func (d EFLDriver) Run() {
	C.ecore_main_loop_begin()
}
