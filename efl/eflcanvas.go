package efl

// #cgo pkg-config: evas ecore-evas
// #include <Evas.h>
// #include <Ecore.h>
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
	native  map[ui.CanvasObject]*C.Evas_Object
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

func buildCanvasObject(c *eflCanvas, o ui.CanvasObject, target ui.CanvasObject, size ui.Size) *C.Evas_Object {
	var obj *C.Evas_Object

	switch o.(type) {
	case *canvas.Text:
		obj = C.evas_object_text_add(c.evas)

		to, _ := o.(*canvas.Text)
		C.evas_object_color_set(obj, C.int(to.Color.R), C.int(to.Color.G),
			C.int(to.Color.B), C.int(to.Color.A))

		updateFont(obj, c, to)
	case *canvas.Rectangle:
		obj = C.evas_object_rectangle_add(c.evas)
		ro, _ := o.(*canvas.Rectangle)

		C.evas_object_color_set(obj, C.int(ro.FillColor.R), C.int(ro.FillColor.G),
			C.int(ro.FillColor.B), C.int(ro.FillColor.A))
	case *canvas.Line:
		obj = C.evas_object_line_add(c.evas)
		lo, _ := o.(*canvas.Line)

		C.evas_object_color_set(obj, C.int(lo.StrokeColor.R), C.int(lo.StrokeColor.G),
			C.int(lo.StrokeColor.B), C.int(lo.StrokeColor.A))
	case *canvas.Circle:
		obj = C.evas_object_vg_add(c.evas)
		co, _ := o.(*canvas.Circle)

		shape := C.evas_vg_shape_add(C.evas_object_vg_root_node_get(obj))
		C.evas_vg_shape_append_circle(shape, C.double(scaleInt(c, vectorPad+size.Width/2)), C.double(scaleInt(c, vectorPad+size.Height/2)), C.double(scaleInt(c, size.Width/2)))
		C.evas_vg_shape_stroke_color_set(shape, C.int(co.StrokeColor.R), C.int(co.StrokeColor.G),
			C.int(co.StrokeColor.B), C.int(co.StrokeColor.A))
		if co.FillColor.A != 0 {
			C.evas_vg_node_color_set(shape, C.int(co.FillColor.R), C.int(co.FillColor.G),
				C.int(co.FillColor.B), C.int(co.FillColor.A))
		}
		C.evas_vg_shape_stroke_width_set(shape, C.double(co.StrokeWidth*c.Scale()))
	default:
		log.Printf("Unrecognised Object %#v\n", o)
	}

	c.objects[obj] = target
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_DOWN,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseDown_cgo)),
		nil)
	return obj
}

func (c *eflCanvas) setupObj(o, o2 ui.CanvasObject, pos ui.Position, size ui.Size) {
	var obj *C.Evas_Object
	if c.native[o] == nil {
		obj = buildCanvasObject(c, o, o2, size)
		c.native[o] = obj
	} else {
		obj = c.native[o]
	}

	switch o.(type) {
	case *canvas.Text:
		to, _ := o.(*canvas.Text)
		C.evas_object_text_text_set(obj, C.CString(to.Text))

		native := nativeTextBounds(obj)
		min := ui.Size{unscaleInt(c, native.Width), unscaleInt(c, native.Height)}
		to.SetMinSize(min)

		pos = ui.NewPos(pos.X+(size.Width-min.Width)/2, pos.Y+(size.Height-min.Height)/2)
	}

	switch o.(type) {
	case *canvas.Circle:
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X-vectorPad)), C.Evas_Coord(scaleInt(c, pos.Y-vectorPad)),
			C.Evas_Coord(scaleInt(c, int(math.Abs(float64(size.Width)))+vectorPad*2)), C.Evas_Coord(scaleInt(c, int(math.Abs(float64(size.Height)))+vectorPad*2)))
	case *canvas.Line:
		lo, _ := o.(*canvas.Line)
		width := lo.Position2.X - lo.Position1.X
		height := lo.Position2.Y - lo.Position1.Y

		if width >= 0 {
			if height >= 0 {
				C.evas_object_line_xy_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
					C.Evas_Coord(scaleInt(c, pos.X+size.Width)), C.Evas_Coord(scaleInt(c, pos.Y+size.Height)))
			} else {
				C.evas_object_line_xy_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y-size.Height)),
					C.Evas_Coord(scaleInt(c, pos.X+size.Width)), C.Evas_Coord(scaleInt(c, pos.Y)))
			}
		} else {
			if height >= 0 {
				C.evas_object_line_xy_set(obj, C.Evas_Coord(scaleInt(c, pos.X-size.Width)), C.Evas_Coord(scaleInt(c, pos.Y)),
					C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y+size.Height)))
			} else {
				C.evas_object_line_xy_set(obj, C.Evas_Coord(scaleInt(c, pos.X-size.Width)), C.Evas_Coord(scaleInt(c, pos.Y-size.Height)),
					C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)))
			}
		}
	default:
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	}
	C.evas_object_show(obj)
}

func (c *eflCanvas) setupContainer(objs []ui.CanvasObject, target ui.CanvasObject, pos ui.Position, size ui.Size) {
	containerPos := pos
	containerSize := size
	if target == c.content {
		containerPos = ui.NewPos(0, 0)
		containerSize = ui.NewSize(size.Width+2*theme.Padding(), size.Height+2*theme.Padding())
	}

	if c.native[target] == nil {
		obj := C.evas_object_rectangle_add(c.evas)
		bg := theme.BackgroundColor()
		C.evas_object_color_set(obj, C.int(bg.R), C.int(bg.G), C.int(bg.B), C.int(bg.A))

		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, containerPos.X)), C.Evas_Coord(scaleInt(c, containerPos.Y)),
			C.Evas_Coord(scaleInt(c, containerSize.Width)), C.Evas_Coord(scaleInt(c, containerSize.Height)))
		C.evas_object_show(obj)

		c.native[target] = obj
	} else {
		obj := c.native[target]
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, containerPos.X)), C.Evas_Coord(scaleInt(c, containerPos.Y)),
			C.Evas_Coord(scaleInt(c, containerSize.Width)), C.Evas_Coord(scaleInt(c, containerSize.Height)))
	}

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
			c.setupContainer(child.(widget.Widget).Layout(child.CurrentSize()), child, child.CurrentPosition().Add(pos), child.CurrentSize())

		default:
			if target == nil {
				target = child
			}

			childPos := child.CurrentPosition().Add(pos)
			c.setupObj(child, target, childPos, child.CurrentSize())
		}
	}
}

func (c *eflCanvas) Size() ui.Size {
	return c.size
}

// TODO let's not draw over the top for every refresh!
func (c *eflCanvas) Refresh(o ui.CanvasObject) {
	C.ecore_thread_main_loop_begin()
	inner := c.size.Add(ui.NewSize(theme.Padding()*-2, theme.Padding()*-2))
	switch o.(type) {
	case *ui.Container:
		container := o.(*ui.Container)
		container.Move(ui.NewPos(theme.Padding(), theme.Padding()))
		container.Resize(inner)
		// TODO should this move into container like widget?
		if container.Layout != nil {
			container.Layout.Layout(container.Objects, inner)
		} else {
			layout.NewMaxLayout().Layout(container.Objects, inner)
		}

		c.setupContainer(container.Objects, o, ui.NewPos(theme.Padding(), theme.Padding()), inner)
	case widget.Widget:
		widget := o.(widget.Widget)
		c.setupContainer(widget.Layout(inner), o, ui.NewPos(theme.Padding(), theme.Padding()), inner)
	default:
		c.setupObj(o, o, ui.NewPos(theme.Padding(), theme.Padding()), inner)
	}

	C.ecore_thread_main_loop_end()
}

func (c *eflCanvas) fitContent() {
	min := c.content.MinSize()
	minWidth := scaleInt(c, min.Width+theme.Padding()*2)
	minHeight := scaleInt(c, min.Height+theme.Padding()*2)

	C.ecore_evas_size_min_set(c.window.ee, C.int(minWidth), C.int(minHeight))
	C.ecore_evas_resize(c.window.ee, C.int(minWidth), C.int(minHeight))
}

func (c *eflCanvas) SetContent(o ui.CanvasObject) {
	canvases[C.ecore_evas_get(c.window.ee)] = c
	c.objects = make(map[*C.Evas_Object]ui.CanvasObject)
	c.native = make(map[ui.CanvasObject]*C.Evas_Object)
	c.content = o

	c.Refresh(o)
	c.fitContent()
}

func updateFont(obj *C.Evas_Object, c *eflCanvas, t *canvas.Text) {
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

	c.Refresh(c.content)
	c.fitContent()
}
