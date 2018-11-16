// +build !ci

package efl

// #cgo pkg-config: eina evas ecore-evas
// #include <Eina.h>
// #include <Evas.h>
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <stdlib.h>
//
// void onObjectMouseDown_cgo(Evas_Object *, void *);
import "C"

import "log"
import "image/color"
import "math"
import "sync"
import "unsafe"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"
import "github.com/fyne-io/fyne/layout"
import "github.com/fyne-io/fyne/widget"

var canvases = make(map[*C.Evas]*eflCanvas)

//export onObjectMouseDown
func onObjectMouseDown(obj *C.Evas_Object, info *C.Evas_Event_Mouse_Down) {
	current := canvases[C.evas_object_evas_get(obj)]
	co := current.objects[obj]

	var x, y C.int
	C.evas_object_geometry_get(obj, &x, &y, nil, nil)
	pos := fyne.NewPos(unscaleInt(current, int(info.canvas.x)), unscaleInt(current, int(info.canvas.y)))
	pos = pos.Subtract(current.offsets[co])

	ev := new(fyne.MouseEvent)
	ev.Position = pos
	ev.Button = fyne.MouseButton(int(info.button))

	switch w := co.(type) {
	case fyne.ClickableObject:
		w.OnMouseDown(ev)
	case fyne.FocusableObject:
		current.Focus(w)
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

func ignoreObject(o fyne.CanvasObject) bool {
	if _, ok := o.(layout.SpacerObject); ok {
		return true
	}

	return false
}

func setColor(obj *C.Evas_Object, col color.Color) {
	r, g, b, a := col.RGBA()

	C.evas_object_color_set(obj, C.int((uint8)(r)), C.int((uint8)(g)), C.int((uint8)(b)), C.int((uint8)(a)))
}

func (c *eflCanvas) buildObject(o fyne.CanvasObject, target fyne.CanvasObject, offset fyne.Position) *C.Evas_Object {
	var obj *C.Evas_Object
	var opts canvas.Options

	switch co := o.(type) {
	case *canvas.Text:
		obj = C.evas_object_text_add(c.evas)

		cstr := C.CString(co.Text)
		C.evas_object_text_text_set(obj, cstr)
		C.free(unsafe.Pointer(cstr))
		setColor(obj, co.Color)

		updateFont(obj, c, co.TextSize, co.TextStyle)
		opts = co.Options
	case *canvas.Rectangle:
		obj = C.evas_object_rectangle_add(c.evas)

		setColor(obj, co.FillColor)
		opts = co.Options
	case *canvas.Image:
		obj = C.evas_object_image_filled_add(c.evas)
		C.evas_object_image_alpha_set(obj, C.EINA_TRUE)
		alpha := C.int(float64(255) * co.Alpha())
		C.evas_object_color_set(obj, alpha, alpha, alpha, alpha) // premul ffffff*alpha

		if co.File != "" {
			c.loadImage(co, obj)
		}
		opts = co.Options
	case *canvas.Line:
		obj = C.evas_object_line_add(c.evas)

		setColor(obj, co.StrokeColor)
		opts = co.Options
	case *canvas.Circle:
		// TODO - this isnt all there yet, but at least this stops lots of debug output
		obj = C.evas_object_rectangle_add(c.evas)

		setColor(obj, co.FillColor)
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

func (c *eflCanvas) buildContainer(parent fyne.CanvasObject, target fyne.CanvasObject, objs []fyne.CanvasObject,
	size fyne.Size, pos, offset fyne.Position) {

	obj := C.evas_object_rectangle_add(c.evas)
	r, g, b, a := theme.BackgroundColor().RGBA()
	C.evas_object_color_set(obj, C.int(r), C.int(g), C.int(b), C.int(a))

	C.evas_object_show(obj)
	c.native[parent] = obj
	c.offsets[parent] = offset

	childOffset := offset.Add(pos)
	for _, child := range objs {
		switch co := child.(type) {
		case *fyne.Container:
			c.buildContainer(co, co, co.Objects, child.CurrentSize(),
				child.CurrentPosition(), childOffset)
		case fyne.Widget:
			click := child
			if _, ok := parent.(*widget.Entry); ok {
				click = parent
			}
			c.buildContainer(co, click, co.Renderer().Objects(),
				child.CurrentSize(), child.CurrentPosition(), childOffset)
		default:
			if target == nil {
				target = parent
				if target == nil {
					target = child
				}
			}

			if !ignoreObject(child) {
				c.buildObject(child, target, childOffset)
			}
		}
	}

	if themed, ok := parent.(fyne.ThemedObject); ok {
		themed.ApplyTheme()
	}
}

func renderImagePortion(img *canvas.Image, pixels []uint32, wg *sync.WaitGroup,
	startx, starty, width, height, imgWidth, imgHeight int) {
	defer wg.Done()

	// calculate image pixels
	i := startx + starty*imgWidth
	for y := starty; y < starty+height; y++ {
		for x := startx; x < startx+width; x++ {
			col := img.PixelColor(x, y, imgWidth, imgHeight)
			// TODO support other color models
			if rgba, ok := col.(color.RGBA); ok {
				pixels[i] = (uint32)(((uint32)(rgba.A) << 24) | ((uint32)(rgba.R) << 16) |
					((uint32)(rgba.G) << 8) | (uint32)(rgba.B))
			}
			i++
		}
		i += imgWidth - width
	}
}

func (c *eflCanvas) renderImage(img *canvas.Image, x, y, width, height int) {
	if width <= 0 || height <= 0 {
		return
	}

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

func (c *eflCanvas) loadImage(img *canvas.Image, obj *C.Evas_Object) {
	size := img.CurrentSize()

	C.evas_object_image_load_size_set(obj, C.int(scaleInt(c, size.Width)), C.int(scaleInt(c, size.Height)))
	cstr := C.CString(img.File)
	C.evas_object_image_file_set(obj, cstr, nil)
	C.free(unsafe.Pointer(cstr))
}

func (c *eflCanvas) refreshObject(o, o2 fyne.CanvasObject) {
	obj := c.native[o]

	// TODO a better solution here as objects are added to the UI
	if obj == nil {
		if ignoreObject(o) {
			return
		}
		obj = c.buildObject(o, o2, c.offsets[o])
	}
	pos := c.offsets[o].Add(o.CurrentPosition())
	size := o.CurrentSize()

	switch co := o.(type) {
	case *canvas.Text:
		cstr := C.CString(co.Text)
		C.evas_object_text_text_set(obj, cstr)
		C.free(unsafe.Pointer(cstr))

		setColor(obj, co.Color)
		updateFont(obj, c, co.TextSize, co.TextStyle)
		pos = getTextPosition(co, pos, size)

		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	case *canvas.Rectangle:
		setColor(obj, co.FillColor)
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	case *canvas.Image:
		var oldWidth, oldHeight C.int
		C.evas_object_geometry_get(obj, nil, nil, &oldWidth, &oldHeight)

		width := scaleInt(c, size.Width)
		height := scaleInt(c, size.Height)
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(width), C.Evas_Coord(height))

		C.evas_object_image_alpha_set(obj, C.EINA_TRUE)
		alpha := C.int(float64(255) * co.Alpha())
		C.evas_object_color_set(obj, alpha, alpha, alpha, alpha) // premul ffffff*alpha

		if co.File != "" {
			c.loadImage(co, obj)
		}
		if co.PixelColor != nil {
			C.evas_object_image_size_set(obj, C.int(width), C.int(height))

			c.renderImage(co, 0, 0, width, height)
		}
	case *canvas.Line:
		width := co.Position2.X - co.Position1.X
		height := co.Position2.Y - co.Position1.Y

		setColor(obj, co.StrokeColor)

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
	case *canvas.Circle:
		setColor(obj, co.FillColor)
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	default:
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	}

	if o.IsVisible() != (C.evas_object_visible_get(obj) != 0) {
		if o.IsVisible() {
			C.evas_object_show(obj)
		} else {
			C.evas_object_hide(obj)
		}
	}
}

func (c *eflCanvas) refreshContainer(objs []fyne.CanvasObject, target fyne.CanvasObject, pos fyne.Position, size fyne.Size) {
	position := c.offsets[target].Add(pos)

	obj := c.native[target]
	bg := theme.BackgroundColor()
	if _, ok := target.(*widget.Toolbar); ok { // TODO don't make this a special case
		bg = theme.ButtonColor()
	}
	setColor(obj, bg)

	if target == c.content {
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, 0)), C.Evas_Coord(scaleInt(c, 0)),
			C.Evas_Coord(scaleInt(c, size.Width+theme.Padding()*2)), C.Evas_Coord(scaleInt(c, size.Height+theme.Padding()*2)))
	} else {
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, position.X)), C.Evas_Coord(scaleInt(c, position.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	}

	if target.IsVisible() != (C.evas_object_visible_get(obj) != 0) {
		if target.IsVisible() {
			C.evas_object_show(obj)
		} else {
			C.evas_object_hide(obj)
		}
	}

	for _, child := range objs {
		c.offsets[child] = position

		switch typed := child.(type) {
		case *fyne.Container:
			c.refreshContainer(typed.Objects, child, child.CurrentPosition(), child.CurrentSize())
		case fyne.Widget:
			c.refreshContainer(typed.Renderer().Objects(), child,
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
	runOnMain(func() {
		switch set := o.(type) {
		case *fyne.Container:
			c.buildContainer(set, set, set.Objects, set.MinSize(), o.CurrentPosition(), offset)
		case fyne.Widget:
			c.buildContainer(set, set, set.Renderer().Objects(),
				set.MinSize(), o.CurrentPosition(), offset)
		default:
			if !ignoreObject(o) {
				c.buildObject(o, o, offset)
			}
		}
	})
}

func (c *eflCanvas) Refresh(o fyne.CanvasObject) {
	queueRender(c, o)
}

func (c *eflCanvas) doRefresh(o fyne.CanvasObject) {
	switch ref := o.(type) {
	case *fyne.Container:
		c.refreshContainer(ref.Objects, o, o.CurrentPosition(),
			o.CurrentSize())
	case fyne.Widget:
		c.refreshContainer(ref.Renderer().Objects(), o,
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

		if c.window.FixedSize() {
			C.ecore_evas_size_max_set(c.window.ee, C.int(minWidth), C.int(minHeight))
		}
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
