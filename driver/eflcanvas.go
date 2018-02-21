package driver

// #cgo pkg-config: evas
// #include <Evas.h>
//
// void onObjectMouseDown_cgo(Evas_Object *, void *);
import "C"

import "log"
import "unsafe"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"
import "github.com/fyne-io/fyne/ui/widget"

type canvasobject struct {
	obj    *C.Evas_Object
	canvas ui.Canvas
}

var objects = make(map[*C.Evas_Object]ui.CanvasObject)

func (o *canvasobject) Canvas() ui.Canvas {
	return o.canvas
}

//export onObjectMouseDown
func onObjectMouseDown(obj *C.Evas_Object) {
	co := objects[obj]

	switch co.(type) {
	case *widget.Button:
		co.(*widget.Button).OnClicked()
	}
}

type eflCanvas struct {
	ui.Canvas
	evas    *C.Evas
	size    ui.Size
	scale   float32
	content ui.CanvasObject
}

func nativeTextBounds(obj *C.Evas_Object) ui.Size {
	width, height := 0, 0
	var w, h C.Evas_Coord
	length := int(C.strlen(C.evas_object_text_text_get(obj)))

	for i := 0; i < length; i++ {
		C.evas_object_text_char_pos_get(obj, C.int(i), nil, nil, &w, &h)
		width += int(w) + 1
		if int(h) > height {
			height = int(h)
		}
	}

	return ui.Size{width, height}
}

func buildCanvasObject(c *eflCanvas, o ui.CanvasObject, target ui.CanvasObject) (*C.Evas_Object, *ui.Size) {
	var obj *C.Evas_Object
	var min *ui.Size

	switch o.(type) {
	case *ui.TextObject:
		obj = C.evas_object_text_add(c.evas)

		to, _ := o.(*ui.TextObject)
		C.evas_object_color_set(obj, C.int(to.Color.R), C.int(to.Color.G),
			C.int(to.Color.B), C.int(to.Color.A))

		updateFont(obj, c, to)
		C.evas_object_text_text_set(obj, C.CString(to.Text))

		native := nativeTextBounds(obj)
		min = &ui.Size{unscaleInt(c, native.Width), unscaleInt(c, native.Height)}
	case *ui.RectangleObject:
		obj = C.evas_object_rectangle_add(c.evas)

		ro, _ := o.(*ui.RectangleObject)
		C.evas_object_color_set(obj, C.int(ro.Color.R), C.int(ro.Color.G),
			C.int(ro.Color.B), C.int(ro.Color.A))
	default:
		log.Printf("Unrecognised Object %#v\n", o)
	}

	objects[obj] = target
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_DOWN,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseDown_cgo)),
		nil)
	return obj, min
}

func (c *eflCanvas) setupObj(o, o2 ui.CanvasObject, pos ui.Position, size ui.Size) {
	obj, min := buildCanvasObject(c, o, o2)
	if min != nil {
		pos = ui.NewPos(pos.X+(size.Width-min.Width)/2, pos.Y+(size.Height-min.Height)/2)
	}

	C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
		C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	C.evas_object_show(obj)
}

func (c *eflCanvas) setupContainer(objs []ui.CanvasObject, target ui.CanvasObject, pos ui.Position, size ui.Size) {
	for _, child := range objs {
		switch child.(type) {
		case *ui.Container:
			container := child.(*ui.Container)

			if container.Layout != nil {
				container.Layout.Layout(container.Objects, c.size)
			} else {
				layout.NewMaxLayout().Layout(container.Objects, c.size)
			}
			c.setupContainer(container.Objects, nil, child.CurrentPosition(), child.CurrentSize())
		case ui.Widget:
			c.setupContainer(child.(ui.Widget).Layout(), child, child.CurrentPosition(), child.CurrentSize())

		default:
			if target == nil {
				target = child
			}

			childPos := child.CurrentPosition().Add(pos)
			c.setupObj(child, target, childPos, child.CurrentSize())
		}
	}
}

func (c *eflCanvas) SetContent(o ui.CanvasObject) {
	switch o.(type) {
	case *ui.Container:
		r := ui.NewRectangle(theme.BackgroundColor())
		obj, _ := buildCanvasObject(c, r, r)
		C.evas_object_geometry_set(obj, 0, 0, C.Evas_Coord(scaleInt(c, c.size.Width)), C.Evas_Coord(scaleInt(c, c.size.Height)))
		C.evas_object_show(obj)

		container := o.(*ui.Container)
		container.Resize(c.size)
		// TODO should this move into container like widget?
		if container.Layout != nil {
			container.Layout.Layout(container.Objects, c.size)
		} else {
			layout.NewMaxLayout().Layout(container.Objects, c.size)
		}

		c.setupContainer(container.Objects, nil, ui.NewPos(0, 0), container.CurrentSize())
	case ui.Widget:
		widget := o.(ui.Widget)
		c.setupContainer(widget.Layout(), o, widget.CurrentPosition(), widget.CurrentSize())
	default:
		c.setupObj(o, o, ui.NewPos(0, 0), c.size)
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

func unscaleInt(c ui.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) / c.Scale())
	}
}

func (c *eflCanvas) Scale() float32 {
	return c.scale
}

func (c *eflCanvas) SetScale(scale float32) {
	c.scale = scale
	log.Println("TODO Update all our objects")
}
