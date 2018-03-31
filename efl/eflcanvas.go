package efl

// #cgo pkg-config: evas ecore-evas
// #include <Evas.h>
// #include <Ecore_Evas.h>
//
// void onObjectMouseDown_cgo(Evas_Object *, void *);
import "C"

import "log"
import "math"
import "unsafe"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/event"
import "github.com/fyne-io/fyne/ui/layout"
import "github.com/fyne-io/fyne/ui/theme"
import "github.com/fyne-io/fyne/ui/widget"

var canvases = make(map[*C.Evas]*eflCanvas)

const vectorPad = 10

//export onObjectMouseDown
func onObjectMouseDown(obj *C.Evas_Object, info *C.Evas_Event_Mouse_Down) {
	canvas := canvases[C.evas_object_evas_get(obj)]
	co := canvas.objects[obj]

	var x, y C.int
	C.evas_object_geometry_get(obj, &x, &y, nil, nil)
	pos := ui.NewPos(unscaleInt(canvas, int(info.canvas.x-x)), unscaleInt(canvas, int(info.canvas.y-y)))

	ev := new(event.MouseEvent)
	ev.Position = pos
	ev.Button = event.MouseButton(int(info.button))

	switch co.(type) {
	case *widget.Button:
		co.(*widget.Button).OnClicked(ev)
	}
}

type eflCanvas struct {
	ui.Canvas
	evas    *C.Evas
	size    ui.Size
	scale   float32
	content ui.CanvasObject
	window  *window

	objects map[*C.Evas_Object]ui.CanvasObject
}

func nativeTextBounds(obj *C.Evas_Object) ui.Size {
	width, height := 0, 0
	var w, h C.Evas_Coord
	length := int(C.strlen(C.evas_object_text_text_get(obj)))

	for i := 0; i < length; i++ {
		C.evas_object_text_char_pos_get(obj, C.int(i), nil, nil, &w, &h)
		width += int(w) + 2
		if int(h) > height {
			height = int(h)
		}
	}

	return ui.Size{width, height}
}

func buildCanvasObject(c *eflCanvas, o ui.CanvasObject, target ui.CanvasObject, size ui.Size) (*C.Evas_Object, *ui.Size) {
	var obj *C.Evas_Object
	var min *ui.Size

	switch o.(type) {
	case *canvas.TextObject:
		obj = C.evas_object_text_add(c.evas)

		to, _ := o.(*canvas.TextObject)
		C.evas_object_color_set(obj, C.int(to.Color.R), C.int(to.Color.G),
			C.int(to.Color.B), C.int(to.Color.A))

		updateFont(obj, c, to)
		C.evas_object_text_text_set(obj, C.CString(to.Text))

		native := nativeTextBounds(obj)
		min = &ui.Size{unscaleInt(c, native.Width), unscaleInt(c, native.Height)}
		to.SetMinSize(*min)
	case *canvas.RectangleObject:
		obj = C.evas_object_rectangle_add(c.evas)

		ro, _ := o.(*canvas.RectangleObject)
		C.evas_object_color_set(obj, C.int(ro.Color.R), C.int(ro.Color.G),
			C.int(ro.Color.B), C.int(ro.Color.A))
	case *canvas.LineObject:
		obj = C.evas_object_vg_add(c.evas)
		lo, _ := o.(*canvas.LineObject)

		shape := C.evas_vg_shape_add(C.evas_object_vg_root_node_get(obj))
		C.evas_vg_shape_append_move_to(shape, vectorPad, vectorPad)
		C.evas_vg_shape_append_line_to(shape, C.double(vectorPad+scaleInt(c, size.Width)), vectorPad+C.double(scaleInt(c, size.Height)))
		C.evas_vg_shape_stroke_color_set(shape, C.int(lo.Color.R), C.int(lo.Color.G),
			C.int(lo.Color.B), C.int(lo.Color.A))
		C.evas_vg_shape_stroke_width_set(shape, C.double(lo.Width*c.Scale()))
	case *canvas.CircleObject:
		obj = C.evas_object_vg_add(c.evas)
		lo, _ := o.(*canvas.CircleObject)

		shape := C.evas_vg_shape_add(C.evas_object_vg_root_node_get(obj))
		C.evas_vg_shape_append_circle(shape, C.double(vectorPad+scaleInt(c, size.Width/2)), vectorPad+C.double(scaleInt(c, size.Height/2)), C.double(scaleInt(c, size.Width/2)))
		C.evas_vg_shape_stroke_color_set(shape, C.int(lo.Color.R), C.int(lo.Color.G),
			C.int(lo.Color.B), C.int(lo.Color.A))
		C.evas_vg_shape_stroke_width_set(shape, C.double(lo.Width*c.Scale()))
	default:
		log.Printf("Unrecognised Object %#v\n", o)
	}

	c.objects[obj] = target
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_DOWN,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseDown_cgo)),
		nil)
	return obj, min
}

func (c *eflCanvas) setupObj(o, o2 ui.CanvasObject, pos ui.Position, size ui.Size) {
	obj, min := buildCanvasObject(c, o, o2, size)
	if min != nil {
		pos = ui.NewPos(pos.X+(size.Width-min.Width)/2, pos.Y+(size.Height-min.Height)/2)
	}

	switch o.(type) {
	case *canvas.LineObject, *canvas.CircleObject:
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)-vectorPad), C.Evas_Coord(scaleInt(c, pos.Y)-vectorPad),
			C.Evas_Coord(scaleInt(c, size.Width)+vectorPad*2), C.Evas_Coord(scaleInt(c, size.Height)+vectorPad*2))
	default:
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	}
	C.evas_object_show(obj)
}

func (c *eflCanvas) setupContainer(objs []ui.CanvasObject, target ui.CanvasObject, pos ui.Position, size ui.Size) {
	for _, child := range objs {
		switch child.(type) {
		case *ui.Container:
			container := child.(*ui.Container)

			if container.Layout != nil {
				container.Layout.Layout(container.Objects, child.CurrentSize())
			} else {
				layout.NewMaxLayout().Layout(container.Objects, child.CurrentSize())
			}
			c.setupContainer(container.Objects, nil, child.CurrentPosition().Add(pos), child.CurrentSize())
		case widget.Widget:
			c.setupContainer(child.(widget.Widget).Layout(), child, child.CurrentPosition().Add(pos), child.CurrentSize())

		default:
			if target == nil {
				target = child
			}

			childPos := child.CurrentPosition().Add(pos)
			c.setupObj(child, target, childPos, child.CurrentSize())
		}
	}
}

func (c *eflCanvas) Refresh(o ui.CanvasObject) {
	inner := c.size.Add(ui.NewSize(theme.Padding()*-2, theme.Padding()*-2))
	switch o.(type) {
	case *ui.Container:
		r := canvas.NewRectangle(theme.BackgroundColor())
		obj, _ := buildCanvasObject(c, r, r, inner)
		C.evas_object_geometry_set(obj, 0, 0, C.Evas_Coord(scaleInt(c, c.size.Width)), C.Evas_Coord(scaleInt(c, c.size.Height)))
		C.evas_object_show(obj)

		container := o.(*ui.Container)
		container.Move(ui.NewPos(theme.Padding(), theme.Padding()))
		container.Resize(inner)
		// TODO should this move into container like widget?
		if container.Layout != nil {
			container.Layout.Layout(container.Objects, inner)
		} else {
			layout.NewMaxLayout().Layout(container.Objects, inner)
		}

		c.setupContainer(container.Objects, nil, ui.NewPos(theme.Padding(), theme.Padding()), inner)
	case widget.Widget:
		widget := o.(widget.Widget)
		c.setupContainer(widget.Layout(), o, ui.NewPos(theme.Padding(), theme.Padding()), inner)
	default:
		c.setupObj(o, o, ui.NewPos(theme.Padding(), theme.Padding()), inner)
	}

	c.content = o
}

func (c *eflCanvas) SetContent(o ui.CanvasObject) {
	canvases[C.ecore_evas_get(c.window.ee)] = c
	c.objects = make(map[*C.Evas_Object]ui.CanvasObject)
	c.Refresh(o)

	min := o.MinSize()
	minWidth := scaleInt(c, min.Width+theme.Padding()*2)
	minHeight := scaleInt(c, min.Height+theme.Padding()*2)

	C.ecore_evas_size_min_set(c.window.ee, C.int(minWidth), C.int(minHeight))

	var w, h C.int
	C.ecore_evas_geometry_get(c.window.ee, nil, nil, &w, &h)
	if int(w) < minWidth || int(h) < minHeight {
		C.ecore_evas_resize(c.window.ee, C.int(minWidth), C.int(minHeight))
	}
}

func updateFont(obj *C.Evas_Object, c *eflCanvas, t *canvas.TextObject) {
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
		return int(math.Round(float64(v) * float64(c.Scale())))
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
