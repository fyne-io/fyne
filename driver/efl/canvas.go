// +build !ci,efl

package efl

// #cgo pkg-config: eina evas ecore-evas
// #include <Eina.h>
// #include <Evas.h>
// #include <Ecore.h>
// #include <Ecore_Evas.h>
// #include <stdlib.h>
//
// void onObjectMouseDown_cgo(Evas_Object *, void *);
// void onObjectMouseWheel_cgo(Evas_Object *, void *);
// void onObjectMouseIn_cgo(Evas_Object *, void *);
// void onObjectMouseOut_cgo(Evas_Object *, void *);
import "C"

import (
	"image/color"
	"log"
	"math"
	"sync"
	"unsafe"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

var canvasMutex sync.RWMutex
var canvases = make(map[*C.Evas]*eflCanvas)

func getMouseObject(obj *C.Evas_Object, cx, cy C.int) (fyne.CanvasObject, fyne.Canvas, fyne.Position) {
	canvasMutex.RLock()
	current := canvases[C.evas_object_evas_get(obj)]
	canvasMutex.RUnlock()
	co := current.objects[obj]
	target := current.native[co]

	var x, y C.int
	C.evas_object_geometry_get(target, &x, &y, nil, nil)
	pos := fyne.NewPos(unscaleInt(current, int(cx-x)), unscaleInt(current, int(cy-y)))

	return co, current, pos
}

//export onObjectMouseDown
func onObjectMouseDown(obj *C.Evas_Object, info *C.Evas_Event_Mouse_Down) {
	co, current, pos := getMouseObject(obj, info.canvas.x, info.canvas.y)

	ev := new(fyne.PointEvent)
	ev.Position = pos

	switch w := co.(type) {
	case fyne.Tappable:
		if int(info.button) == 3 {
			go w.TappedSecondary(ev)
		} else {
			go w.Tapped(ev)
		}
	case fyne.Focusable:
		current.Focus(w)
	}
}

//export onObjectMouseWheel
func onObjectMouseWheel(obj *C.Evas_Object, info *C.Evas_Event_Mouse_Wheel) {
	co, _, pos := getMouseObject(obj, info.canvas.x, info.canvas.y)

	ev := new(fyne.ScrollEvent)
	ev.Position = pos
	ev.DeltaY = int(-info.z)

	switch w := co.(type) {
	case fyne.Scrollable:
		w.Scrolled(ev)
	}
}

//export onObjectMouseIn
func onObjectMouseIn(obj *C.Evas_Object, info *C.Evas_Event_Mouse_In) {
	co, current, _ := getMouseObject(obj, info.canvas.x, info.canvas.y)

	setCursor(current.(*eflCanvas).window, co)
}

//export onObjectMouseOut
func onObjectMouseOut(obj *C.Evas_Object, info *C.Evas_Event_Mouse_Out) {
	co, current, _ := getMouseObject(obj, info.canvas.x, info.canvas.y)

	unsetCursor(current.(*eflCanvas).window, co)
}

type eflCanvas struct {
	fyne.Canvas
	evas  *C.Evas
	size  fyne.Size
	scale float32

	content fyne.CanvasObject
	window  *window
	focused fyne.Focusable

	onTypedRune func(rune)
	onTypedKey  func(*fyne.KeyEvent)

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

func (c *eflCanvas) buildObject(o fyne.CanvasObject, target fyne.CanvasObject, offset fyne.Position,
	clip *C.Evas_Object) *C.Evas_Object {
	var obj *C.Evas_Object

	switch co := o.(type) {
	case *canvas.Text:
		obj = C.evas_object_text_add(c.evas)

		cstr := C.CString(co.Text)
		C.evas_object_text_text_set(obj, cstr)
		C.free(unsafe.Pointer(cstr))
		setColor(obj, co.Color)

		updateFont(obj, c, co.TextSize, co.TextStyle)
	case *canvas.Rectangle:
		obj = C.evas_object_rectangle_add(c.evas)

		setColor(obj, co.FillColor)
	case *canvas.Image:
		obj = C.evas_object_image_filled_add(c.evas)
		C.evas_object_image_alpha_set(obj, C.EINA_TRUE)
		alpha := C.int(float64(255) * co.Alpha())
		C.evas_object_color_set(obj, alpha, alpha, alpha, alpha) // premul ffffff*alpha

		if co.File != "" {
			c.loadImage(co, obj)
		}
	case *canvas.Line:
		obj = C.evas_object_line_add(c.evas)

		setColor(obj, co.StrokeColor)
	case *canvas.Circle:
		// TODO - this isnt all there yet, but at least this stops lots of debug output
		obj = C.evas_object_rectangle_add(c.evas)

		if co.FillColor != nil {
			setColor(obj, co.FillColor)
		}
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
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_WHEEL,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseWheel_cgo)),
		nil)
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_IN,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseIn_cgo)),
		nil)
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_OUT,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseOut_cgo)),
		nil)

	if clip != nil {
		C.evas_object_clip_set(obj, clip)
	}
	C.evas_object_show(obj)
	return obj
}

func (c *eflCanvas) buildContainer(parent fyne.CanvasObject, target fyne.CanvasObject, objs []fyne.CanvasObject,
	size fyne.Size, pos, offset fyne.Position, clip *C.Evas_Object) {

	obj := C.evas_object_rectangle_add(c.evas)
	r, g, b, a := theme.BackgroundColor().RGBA()
	C.evas_object_color_set(obj, C.int(r), C.int(g), C.int(b), C.int(a))

	C.evas_object_show(obj)
	c.native[parent] = obj
	c.offsets[parent] = offset
	c.objects[obj] = target
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_DOWN,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseDown_cgo)),
		nil)

	childOffset := offset.Add(pos)
	for _, child := range objs {
		switch co := child.(type) {
		case *fyne.Container:
			c.buildContainer(co, co, co.Objects, child.Size(),
				child.Position(), childOffset, clip)
		case *widget.ScrollContainer:
			C.evas_object_color_set(obj, 255, 255, 255, 255)
			click := child
			if _, ok := parent.(*widget.Entry); ok {
				click = parent
			}
			c.buildContainer(co, click, widget.Renderer(co).Objects(),
				child.Size(), child.Position(), childOffset, obj)
		case fyne.Widget:
			click := child
			if _, ok := parent.(*widget.Entry); ok {
				click = parent
			}
			c.buildContainer(co, click, widget.Renderer(co).Objects(),
				child.Size(), child.Position(), childOffset, clip)
		default:
			if target == nil {
				target = parent
				if target == nil {
					target = child
				}
			}

			if !ignoreObject(child) {
				c.buildObject(child, target, childOffset, clip)
			}
		}
	}

	if themed, ok := parent.(fyne.Themeable); ok {
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

			// convert from the 16 bit-per-channel RGBA to an 8 bit-per-channel ARGB
			r, g, b, a := col.RGBA()
			rgba := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
			pixels[i] = (uint32)(((uint32)(rgba.A) << 24) | ((uint32)(rgba.R) << 16) |
				((uint32)(rgba.G) << 8) | (uint32)(rgba.B))
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
	size := img.Size()

	C.evas_object_image_load_size_set(obj, C.int(scaleInt(c, size.Width)), C.int(scaleInt(c, size.Height)))

	var file string
	if img.Resource != nil {
		file = img.Resource.CachePath()
	} else {
		file = img.File
	}
	cstr := C.CString(file)
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
		obj = c.buildObject(o, o2, c.offsets[o], nil)
	}
	pos := c.offsets[o].Add(o.Position())
	size := o.Size()

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
		width := scaleInt(c, size.Width)
		height := scaleInt(c, size.Height)

		if co.FillMode == canvas.ImageFillStretch {
			C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
				C.Evas_Coord(width), C.Evas_Coord(height))
		} else {
			viewAspect := float32(size.Width) / float32(size.Height)

			var iw, ih C.int
			C.evas_object_image_size_get(obj, &iw, &ih)
			aspect := float32(iw) / float32(ih)

			// if the image specifies it should be original size we need at least that many pixels on screen
			if co.FillMode == canvas.ImageFillOriginal {
				pixSize := fyne.NewSize(unscaleInt(c, int(iw)), unscaleInt(c, int(ih)))
				co.SetMinSize(pixSize)
			}

			widthPad, heightPad := 0, 0
			if viewAspect > aspect {
				newWidth := int(float32(height) * aspect)
				widthPad = (width - newWidth) / 2
				width = newWidth
			} else {
				newHeight := int(float32(width) / aspect)
				heightPad = (height - newHeight) / 2
				height = newHeight
			}

			C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)+widthPad), C.Evas_Coord(scaleInt(c, pos.Y)+heightPad),
				C.Evas_Coord(width), C.Evas_Coord(height))
		}
		C.evas_object_image_alpha_set(obj, C.EINA_TRUE)
		alpha := C.int(float64(255) * co.Alpha())
		C.evas_object_color_set(obj, alpha, alpha, alpha, alpha) // premul ffffff*alpha

		if co.File != "" || co.Resource != nil {
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
		if co.FillColor != nil {
			setColor(obj, co.FillColor)
		}
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	default:
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	}

	if o.Visible() != (C.evas_object_visible_get(obj) != 0) {
		if o.Visible() {
			C.evas_object_show(obj)
		} else {
			C.evas_object_hide(obj)
		}
	}
}

func (c *eflCanvas) refreshContainer(objs []fyne.CanvasObject, target fyne.CanvasObject, pos fyne.Position, size fyne.Size) {
	position := c.offsets[target].Add(pos)

	obj := c.native[target]
	switch target.(type) {
	case *widget.ScrollContainer:
		C.evas_object_color_set(obj, 255, 255, 255, 255)
	default:
		bg := theme.BackgroundColor()
		if wid, ok := target.(fyne.Widget); ok {
			bg = widget.Renderer(wid).BackgroundColor()
		}
		setColor(obj, bg)
	}

	if target == c.content {
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, 0)), C.Evas_Coord(scaleInt(c, 0)),
			C.Evas_Coord(scaleInt(c, size.Width+theme.Padding()*2)), C.Evas_Coord(scaleInt(c, size.Height+theme.Padding()*2)))
	} else {
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, position.X)), C.Evas_Coord(scaleInt(c, position.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	}

	if target.Visible() != (C.evas_object_visible_get(obj) != 0) {
		if target.Visible() {
			C.evas_object_show(obj)
		} else {
			C.evas_object_hide(obj)
		}
	}

	for _, child := range objs {
		c.offsets[child] = position

		switch typed := child.(type) {
		case *fyne.Container:
			c.refreshContainer(typed.Objects, child, child.Position(), child.Size())
		case fyne.Widget:
			c.refreshContainer(widget.Renderer(typed).Objects(), child,
				child.Position(), child.Size())
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
			c.buildContainer(set, set, set.Objects, set.MinSize(), o.Position(), offset, nil)
		case fyne.Widget:
			c.buildContainer(set, set, widget.Renderer(set).Objects(),
				set.MinSize(), o.Position(), offset, nil)
		default:
			if !ignoreObject(o) {
				c.buildObject(o, o, offset, nil)
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
		c.refreshContainer(ref.Objects, o, o.Position(),
			o.Size())
	case fyne.Widget:
		c.refreshContainer(widget.Renderer(ref).Objects(), o,
			o.Position(), o.Size())
	default:
		c.refreshObject(o, o)
	}
}

func (c *eflCanvas) Focus(obj fyne.Focusable) {
	if c.focused != nil {
		if c.focused == obj {
			return
		}

		c.focused.FocusLost()
	}

	c.focused = obj
	obj.FocusGained()
}

func (c *eflCanvas) Focused() fyne.Focusable {
	return c.focused
}

// fitContent is thread safe - its only every called from the main loop
func (c *eflCanvas) fitContent() {
	var w, h C.int
	C.ecore_evas_geometry_get(c.window.ee, nil, nil, &w, &h)

	pad := theme.Padding()
	if !c.window.Padded() {
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
	if !c.window.Padded() {
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
	canvasMutex.Lock()
	canvases[C.ecore_evas_get(c.window.ee)] = c
	canvasMutex.Unlock()

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

func (c *eflCanvas) OnTypedRune() func(rune) {
	return c.onTypedRune
}

func (c *eflCanvas) SetOnTypedRune(typed func(rune)) {
	c.onTypedRune = typed
}

func (c *eflCanvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.onTypedKey
}

func (c *eflCanvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.onTypedKey = typed
}
