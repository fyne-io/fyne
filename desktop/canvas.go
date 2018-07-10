// +build !ci

package desktop

// #cgo pkg-config: eina evas ecore-evas
// #include <Eina.h>
// #include <Evas.h>
// #include <Ecore.h>
// #include <Ecore_Evas.h>
//
// void onObjectMouseDown_cgo(Evas_Object *, void *);
import "C"

import "log"
import "math"
import "sync"
import "unsafe"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"
import "github.com/fyne-io/fyne/widget"

var canvases = make(map[*C.Evas]*eflCanvas)

const vectorPad = 10

//export onObjectMouseDown
func onObjectMouseDown(obj *C.Evas_Object, info *C.Evas_Event_Mouse_Down) {
	canvas := canvases[C.evas_object_evas_get(obj)]
	co := canvas.objects[obj]

	var x, y C.int
	C.evas_object_geometry_get(obj, &x, &y, nil, nil)
	pos := fyne.NewPos(unscaleInt(canvas, int(info.canvas.x-x)), unscaleInt(canvas, int(info.canvas.y-y)))

	ev := new(fyne.MouseEvent)
	ev.Position = pos
	ev.Button = fyne.MouseButton(int(info.button))

	switch obj := co.(type) {
	case fyne.ClickableObject:
		obj.OnMouseDown(ev)
	case fyne.FocusableObject:
		canvas.Focus(obj)
	}
}

type eflCanvas struct {
	fyne.Canvas
	evas  *C.Evas
	size  fyne.Size
	scale float32

	content fyne.CanvasObject
	window  *window
	focused fyne.FocusableObject

	onKeyDown func(*fyne.KeyEvent)

	objects map[*C.Evas_Object]fyne.CanvasObject
	native  map[fyne.CanvasObject]*C.Evas_Object
	offsets map[fyne.CanvasObject]fyne.Position
	dirty   map[fyne.CanvasObject]bool
}

func (c *eflCanvas) buildObject(o fyne.CanvasObject, target fyne.CanvasObject, offset fyne.Position) *C.Evas_Object {
	var obj *C.Evas_Object
	var opts canvas.Options

	switch co := o.(type) {
	case *canvas.Text:
		obj = C.evas_object_text_add(c.evas)

		C.evas_object_text_text_set(obj, C.CString(co.Text))
		C.evas_object_color_set(obj, C.int(co.Color.R), C.int(co.Color.G),
			C.int(co.Color.B), C.int(co.Color.A))

		updateFont(obj, c, co)
		opts = co.Options
	case *canvas.Rectangle:
		obj = C.evas_object_rectangle_add(c.evas)

		C.evas_object_color_set(obj, C.int(co.FillColor.R), C.int(co.FillColor.G),
			C.int(co.FillColor.B), C.int(co.FillColor.A))
		opts = co.Options
	case *canvas.Image:
		obj = C.evas_object_image_filled_add(c.evas)
		C.evas_object_image_alpha_set(obj, C.EINA_TRUE)

		if co.File != "" {
			c.loadImage(co)
		}
		opts = co.Options
	case *canvas.Line:
		obj = C.evas_object_line_add(c.evas)

		C.evas_object_color_set(obj, C.int(co.StrokeColor.R), C.int(co.StrokeColor.G),
			C.int(co.StrokeColor.B), C.int(co.StrokeColor.A))
		opts = co.Options
	default:
		log.Printf("Unrecognised Object %#v\n", o)
		return nil
	}

	c.native[o] = obj
	c.offsets[o] = offset
	c.objects[obj] = target
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_DOWN,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseDown_cgo)),
		nil)
	if opts.RepeatEvents {
		C.evas_object_repeat_events_set(obj, 1)
	}

	C.evas_object_show(obj)
	return obj
}

func (c *eflCanvas) buildContainer(objs []fyne.CanvasObject,
	target fyne.CanvasObject, size fyne.Size, pos, offset fyne.Position) {

	obj := C.evas_object_rectangle_add(c.evas)
	bg := theme.BackgroundColor()
	C.evas_object_color_set(obj, C.int(bg.R), C.int(bg.G), C.int(bg.B), C.int(bg.A))

	C.evas_object_show(obj)
	c.native[target] = obj
	c.offsets[target] = offset

	childOffset := offset.Add(pos)
	for _, child := range objs {
		switch co := child.(type) {
		case *fyne.Container:
			c.buildContainer(co.Objects, child, child.CurrentSize(),
				child.CurrentPosition(), childOffset)
		case widget.Widget:
			c.buildContainer(co.CanvasObjects(), child,
				child.CurrentSize(), child.CurrentPosition(), childOffset)
		default:
			if target == nil {
				target = child
			}

			c.buildObject(child, target, childOffset)
		}
	}
}

func renderImagePortion(img *canvas.Image, pixels []uint32, wg *sync.WaitGroup,
	startx, starty, width, height, imgWidth, imgHeight int) {
	defer wg.Done()

	// calculate image pixels
	i := startx + starty*imgWidth
	for y := starty; y < starty+height; y++ {
		for x := startx; x < startx+width; x++ {
			color := img.PixelColor(x, y, imgWidth, imgHeight)
			pixels[i] = (uint32)(((uint32)(color.A) << 24) | ((uint32)(color.R) << 16) |
				((uint32)(color.G) << 8) | (uint32)(color.B))
			i++
		}
		i += imgWidth - width
	}
}

func (c *eflCanvas) renderImage(img *canvas.Image, x, y, width, height int) {
	pixels := make([]uint32, width*height)

	// Spawn 4 threads each calculating the pixels for a quadrant of the image
	halfWidth := width / 2
	halfHeight := height / 2

	// use a WaitGroup so we don't render our pixels before they are ready
	var wg sync.WaitGroup
	wg.Add(4)
	go renderImagePortion(img, pixels, &wg, 0, 0, halfWidth, halfHeight, width, height)
	go renderImagePortion(img, pixels, &wg, halfWidth, 0, width-halfWidth, halfHeight, width, height)
	go renderImagePortion(img, pixels, &wg, 0, halfHeight, halfWidth, height-halfHeight, width, height)
	go renderImagePortion(img, pixels, &wg, halfWidth, halfHeight, width-halfWidth, height-halfHeight, width, height)
	wg.Wait()

	// write pixels to canvas
	obj := c.native[img]
	C.evas_object_image_data_set(obj, unsafe.Pointer(&pixels[0]))
	C.evas_object_image_data_update_add(obj, 0, 0, C.int(width), C.int(height))
}

func (c *eflCanvas) loadImage(img *canvas.Image) {
	size := img.CurrentSize()
	obj := c.native[img]

	C.evas_object_image_load_size_set(obj, C.int(scaleInt(c, size.Width)), C.int(scaleInt(c, size.Height)))
	C.evas_object_image_file_set(obj, C.CString(img.File), nil)
}

func (c *eflCanvas) refreshObject(o, o2 fyne.CanvasObject) {
	obj := c.native[o]

	// TODO a better solution here as objects are added to the UI
	if obj == nil {
		obj = c.buildObject(o, o2, fyne.NewPos(0, 0)) // TODO fix offset
	}
	pos := c.offsets[o].Add(o.CurrentPosition())
	size := o.CurrentSize()

	switch co := o.(type) {
	case *canvas.Text:
		C.evas_object_text_text_set(obj, C.CString(co.Text))
		C.evas_object_color_set(obj, C.int(co.Color.R), C.int(co.Color.G),
			C.int(co.Color.B), C.int(co.Color.A))

		updateFont(obj, c, co)
		pos = getTextPosition(co, pos, size)

		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	case *canvas.Rectangle:
		C.evas_object_color_set(obj, C.int(co.FillColor.R), C.int(co.FillColor.G),
			C.int(co.FillColor.B), C.int(co.FillColor.A))
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	case *canvas.Image:
		var oldWidth, oldHeight C.int
		C.evas_object_geometry_get(obj, nil, nil, &oldWidth, &oldHeight)

		width := scaleInt(c, size.Width)
		height := scaleInt(c, size.Height)
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(width), C.Evas_Coord(height))

		if co.File != "" {
			c.loadImage(co)
		}
		if co.PixelColor != nil {
			C.evas_object_image_size_set(obj, C.int(width), C.int(height))

			c.renderImage(co, 0, 0, width, height)
		}
	case *canvas.Line:
		width := co.Position2.X - co.Position1.X
		height := co.Position2.Y - co.Position1.Y

		C.evas_object_color_set(obj, C.int(co.StrokeColor.R), C.int(co.StrokeColor.G),
			C.int(co.StrokeColor.B), C.int(co.StrokeColor.A))

		if width >= 0 {
			if height >= 0 {
				C.evas_object_line_xy_set(obj, C.Evas_Coord(scaleInt(c, 0)), C.Evas_Coord(scaleInt(c, 0)),
					C.Evas_Coord(scaleInt(c, width)), C.Evas_Coord(scaleInt(c, height)))
			} else {
				C.evas_object_line_xy_set(obj, C.Evas_Coord(scaleInt(c, 0)), C.Evas_Coord(scaleInt(c, -height)),
					C.Evas_Coord(scaleInt(c, width)), C.Evas_Coord(scaleInt(c, 0)))
			}
		} else {
			if height >= 0 {
				C.evas_object_line_xy_set(obj, C.Evas_Coord(scaleInt(c, -width)), C.Evas_Coord(scaleInt(c, 0)),
					C.Evas_Coord(scaleInt(c, 0)), C.Evas_Coord(scaleInt(c, height)))
			} else {
				C.evas_object_line_xy_set(obj, C.Evas_Coord(scaleInt(c, 0)), C.Evas_Coord(scaleInt(c, 0)),
					C.Evas_Coord(scaleInt(c, -width)), C.Evas_Coord(scaleInt(c, -height)))
			}
		}
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, int(math.Abs(float64(width))))+1), C.Evas_Coord(scaleInt(c, int(math.Abs(float64(height))))+1))
	default:
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	}
}

func (c *eflCanvas) refreshContainer(objs []fyne.CanvasObject, target fyne.CanvasObject, pos fyne.Position, size fyne.Size) {
	position := c.offsets[target].Add(pos)

	obj := c.native[target]
	bg := theme.BackgroundColor()
	C.evas_object_color_set(obj, C.int(bg.R), C.int(bg.G), C.int(bg.B), C.int(bg.A))

	if target == c.content {
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, 0)), C.Evas_Coord(scaleInt(c, 0)),
			C.Evas_Coord(scaleInt(c, size.Width+theme.Padding()*2)), C.Evas_Coord(scaleInt(c, size.Height+theme.Padding()*2)))
	} else {
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, position.X)), C.Evas_Coord(scaleInt(c, position.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	}

	for _, child := range objs {
		c.offsets[child] = position

		switch typed := child.(type) {
		case *fyne.Container:
			c.refreshContainer(typed.Objects, child, child.CurrentPosition(), child.CurrentSize())
		case widget.Widget:
			typed.ApplyTheme()
			c.refreshContainer(typed.CanvasObjects(), child,
				child.CurrentPosition(), child.CurrentSize())
		default:
			if target == nil {
				target = child
			}

			c.refreshObject(child, target)
		}
	}
}

func (c *eflCanvas) Size() fyne.Size {
	return c.size
}

func (c *eflCanvas) setup(o fyne.CanvasObject, offset fyne.Position) {
	C.ecore_thread_main_loop_begin()

	switch set := o.(type) {
	case *fyne.Container:
		c.buildContainer(set.Objects, o, set.MinSize(), o.CurrentPosition(), offset)
	case widget.Widget:
		c.buildContainer(set.CanvasObjects(), o,
			set.MinSize(), o.CurrentPosition(), offset)
	default:
		c.buildObject(o, o, offset)
	}

	C.ecore_thread_main_loop_end()
}

func (c *eflCanvas) Refresh(o fyne.CanvasObject) {
	queueRender(c, o)
}

func (c *eflCanvas) doRefresh(o fyne.CanvasObject) {
	switch ref := o.(type) {
	case *fyne.Container:
		c.refreshContainer(ref.Objects, o, o.CurrentPosition(),
			o.CurrentSize())
	case widget.Widget:
		c.refreshContainer(ref.CanvasObjects(), o,
			o.CurrentPosition(), o.CurrentSize())
	default:
		c.refreshObject(o, o)
	}
}

func (c *eflCanvas) Contains(obj fyne.CanvasObject) bool {
	return c.native[obj] != nil
}

func (c *eflCanvas) Focus(obj fyne.FocusableObject) {
	if c.focused != nil {
		if c.focused == obj {
			return
		}

		c.focused.OnFocusLost()
	}

	c.focused = obj
	obj.OnFocusGained()
}

func (c *eflCanvas) Focused() fyne.FocusableObject {
	return c.focused
}

func (c *eflCanvas) fitContent() {
	var w, h C.int
	C.ecore_evas_geometry_get(c.window.ee, nil, nil, &w, &h)

	pad := theme.Padding()
	if c.window.Fullscreen() {
		pad = 0
	}
	min := c.content.MinSize()
	minWidth := scaleInt(c, min.Width+pad*2)
	minHeight := scaleInt(c, min.Height+pad*2)

	width := fyne.Max(minWidth, int(w))
	height := fyne.Max(minHeight, int(h))

	if width != int(w) || height != int(h) {
		C.ecore_evas_size_min_set(c.window.ee, C.int(minWidth), C.int(minHeight))
		C.ecore_evas_resize(c.window.ee, C.int(width), C.int(height))
	}
}

func (c *eflCanvas) resizeContent() {
	var w, h C.int
	C.ecore_evas_geometry_get(c.window.ee, nil, nil, &w, &h)

	pad := theme.Padding()
	if c.window.Fullscreen() {
		pad = 0
	}
	width := unscaleInt(c, int(w)) - pad*2
	height := unscaleInt(c, int(h)) - pad*2

	c.content.Resize(fyne.NewSize(width, height))
	c.content.Move(fyne.NewPos(pad, pad))
	queueRender(c, c.content)
}

func (c *eflCanvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *eflCanvas) SetContent(o fyne.CanvasObject) {
	c.objects = make(map[*C.Evas_Object]fyne.CanvasObject)
	c.native = make(map[fyne.CanvasObject]*C.Evas_Object)
	c.offsets = make(map[fyne.CanvasObject]fyne.Position)
	c.dirty = make(map[fyne.CanvasObject]bool)
	c.content = o
	canvases[C.ecore_evas_get(c.window.ee)] = c

	c.resizeContent()
	c.setup(o, fyne.NewPos(0, 0))
}

func scaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(math.Round(float64(v) * float64(c.Scale())))
	}
}

func unscaleInt(c fyne.Canvas, v int) int {
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
	log.Println("FYNE_SCALE", scale)

	ratio := scale / c.scale
	c.scale = scale

	var w, h C.int
	C.ecore_evas_geometry_get(c.window.ee, nil, nil, &w, &h)
	width := int(float32(w) * ratio)
	height := int(float32(h) * ratio)
	C.ecore_evas_resize(c.window.ee, C.int(width), C.int(height))
}

func (c *eflCanvas) SetOnKeyDown(keyDown func(*fyne.KeyEvent)) {
	c.onKeyDown = keyDown
}
