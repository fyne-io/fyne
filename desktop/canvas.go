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
import "github.com/fyne-io/fyne/layout"
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
	dirty   map[fyne.CanvasObject]bool
}

func (c *eflCanvas) buildObject(o fyne.CanvasObject, target fyne.CanvasObject, size fyne.Size) *C.Evas_Object {
	var obj *C.Evas_Object
	var opts canvas.Options

	switch o.(type) {
	case *canvas.Text:
		obj = C.evas_object_text_add(c.evas)

		to, _ := o.(*canvas.Text)
		C.evas_object_text_text_set(obj, C.CString(to.Text))
		C.evas_object_color_set(obj, C.int(to.Color.R), C.int(to.Color.G),
			C.int(to.Color.B), C.int(to.Color.A))

		updateFont(obj, c, to)
		opts = to.Options
	case *canvas.Rectangle:
		obj = C.evas_object_rectangle_add(c.evas)
		ro, _ := o.(*canvas.Rectangle)

		C.evas_object_color_set(obj, C.int(ro.FillColor.R), C.int(ro.FillColor.G),
			C.int(ro.FillColor.B), C.int(ro.FillColor.A))
		opts = ro.Options
	case *canvas.Image:
		obj = C.evas_object_image_add(c.evas)
		img, _ := o.(*canvas.Image)
		C.evas_object_image_alpha_set(obj, C.EINA_TRUE)
		C.evas_object_image_filled_set(obj, C.EINA_TRUE)

		if img.File != "" {
			C.evas_object_image_file_set(obj, C.CString(img.File), nil)
		}
		opts = img.Options
	case *canvas.Line:
		obj = C.evas_object_line_add(c.evas)
		lo, _ := o.(*canvas.Line)

		C.evas_object_color_set(obj, C.int(lo.StrokeColor.R), C.int(lo.StrokeColor.G),
			C.int(lo.StrokeColor.B), C.int(lo.StrokeColor.A))
		opts = lo.Options
	default:
		log.Printf("Unrecognised Object %#v\n", o)
		return nil
	}

	c.native[o] = obj
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

func (c *eflCanvas) buildContainer(objs []fyne.CanvasObject, target fyne.CanvasObject, size fyne.Size) {
	obj := C.evas_object_rectangle_add(c.evas)
	bg := theme.BackgroundColor()
	C.evas_object_color_set(obj, C.int(bg.R), C.int(bg.G), C.int(bg.B), C.int(bg.A))

	C.evas_object_show(obj)
	c.native[target] = obj

	for _, child := range objs {
		switch child.(type) {
		case *fyne.Container:
			container := child.(*fyne.Container)

			c.buildContainer(container.Objects, child, child.CurrentSize())
		case widget.Widget:
			c.buildContainer(child.(widget.Widget).Layout(child.CurrentSize()),
				child, child.CurrentSize())

		default:
			if target == nil {
				target = child
			}

			c.buildObject(child, target, child.CurrentSize())
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

func (c *eflCanvas) refreshObject(o, o2 fyne.CanvasObject, pos fyne.Position, size fyne.Size) {
	obj := c.native[o]

	// TODO a better solution here as objects are added to the UI
	if obj == nil {
		obj = c.buildObject(o, o2, size)
	}

	switch o.(type) {
	case *canvas.Text:
		to, _ := o.(*canvas.Text)
		C.evas_object_text_text_set(obj, C.CString(to.Text))
		C.evas_object_color_set(obj, C.int(to.Color.R), C.int(to.Color.G),
			C.int(to.Color.B), C.int(to.Color.A))

		updateFont(obj, c, to)
		pos = getTextPosition(to, pos, size)

		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	case *canvas.Rectangle:
		ro, _ := o.(*canvas.Rectangle)

		C.evas_object_color_set(obj, C.int(ro.FillColor.R), C.int(ro.FillColor.G),
			C.int(ro.FillColor.B), C.int(ro.FillColor.A))
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(scaleInt(c, size.Width)), C.Evas_Coord(scaleInt(c, size.Height)))
	case *canvas.Image:
		img, _ := o.(*canvas.Image)
		var oldWidth, oldHeight C.int
		C.evas_object_geometry_get(obj, nil, nil, &oldWidth, &oldHeight)

		width := scaleInt(c, size.Width)
		height := scaleInt(c, size.Height)
		C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, pos.X)), C.Evas_Coord(scaleInt(c, pos.Y)),
			C.Evas_Coord(width), C.Evas_Coord(height))

		if img.PixelColor != nil {
			C.evas_object_image_size_set(obj, C.int(width), C.int(height))

			c.renderImage(img, 0, 0, width, height)
		}
	case *canvas.Line:
		lo, _ := o.(*canvas.Line)
		width := lo.Position2.X - lo.Position1.X
		height := lo.Position2.Y - lo.Position1.Y

		C.evas_object_color_set(obj, C.int(lo.StrokeColor.R), C.int(lo.StrokeColor.G),
			C.int(lo.StrokeColor.B), C.int(lo.StrokeColor.A))

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
	containerPos := pos
	containerSize := size
	if target == c.content {
		containerPos = fyne.NewPos(0, 0)
		containerSize = fyne.NewSize(size.Width+2*theme.Padding(), size.Height+2*theme.Padding())
	}

	obj := c.native[target]
	bg := theme.BackgroundColor()
	C.evas_object_color_set(obj, C.int(bg.R), C.int(bg.G), C.int(bg.B), C.int(bg.A))
	C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, containerPos.X)), C.Evas_Coord(scaleInt(c, containerPos.Y)),
		C.Evas_Coord(scaleInt(c, containerSize.Width)), C.Evas_Coord(scaleInt(c, containerSize.Height)))

	for _, child := range objs {
		switch typed := child.(type) {
		case *fyne.Container:
			container := child.(*fyne.Container)

			if container.Layout != nil {
				container.Layout.Layout(container.Objects, child.CurrentSize())
			} else {
				layout.NewMaxLayout().Layout(container.Objects, child.CurrentSize())
			}
			c.refreshContainer(container.Objects, nil, child.CurrentPosition().Add(pos), child.CurrentSize())
		case widget.Widget:
			typed.ApplyTheme()
			c.refreshContainer(child.(widget.Widget).Layout(child.CurrentSize()),
				child, child.CurrentPosition().Add(pos), child.CurrentSize())

		default:
			if target == nil {
				target = child
			}

			childPos := child.CurrentPosition().Add(pos)
			c.refreshObject(child, target, childPos, child.CurrentSize())
		}
	}
}

func (c *eflCanvas) Size() fyne.Size {
	return c.size
}

func (c *eflCanvas) setup(o fyne.CanvasObject) {
	C.ecore_thread_main_loop_begin()

	switch o.(type) {
	case *fyne.Container:
		container := o.(*fyne.Container)

		c.buildContainer(container.Objects, o, container.CurrentSize())
	case widget.Widget:
		widget := o.(widget.Widget)
		c.buildContainer(widget.Layout(widget.CurrentSize()), o,
			widget.CurrentSize())
	default:
		c.buildObject(o, o, o.CurrentSize())
	}

	C.ecore_thread_main_loop_end()
}

func (c *eflCanvas) Refresh(o fyne.CanvasObject) {
	queueRender(c, o)
}

func (c *eflCanvas) doRefresh(o fyne.CanvasObject) {
	position := o.CurrentPosition()
	if o != c.content {
		position = position.Add(fyne.NewPos(theme.Padding(), theme.Padding()))
	}

	switch o.(type) {
	case *fyne.Container:
		container := o.(*fyne.Container)
		// TODO should this move into container like widget?
		if container.Layout != nil {
			container.Layout.Layout(container.Objects, container.CurrentSize())
		} else {
			layout.NewMaxLayout().Layout(container.Objects, container.CurrentSize())
		}

		c.refreshContainer(container.Objects, o, position, container.CurrentSize())
	case widget.Widget:
		widget := o.(widget.Widget)
		c.refreshContainer(widget.Layout(widget.CurrentSize()), o,
			position, widget.CurrentSize())
	default:
		c.refreshObject(o, o, position, o.CurrentSize())
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
	if fyne.GetWindow(c).Fullscreen() {
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
	if fyne.GetWindow(c).Fullscreen() {
		pad = 0
	}
	width := unscaleInt(c, int(w)) - pad*2
	height := unscaleInt(c, int(h)) - pad*2

	c.content.Resize(fyne.NewSize(width, height))
	queueRender(c, c.content)
}

func (c *eflCanvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *eflCanvas) SetContent(o fyne.CanvasObject) {
	c.objects = make(map[*C.Evas_Object]fyne.CanvasObject)
	c.native = make(map[fyne.CanvasObject]*C.Evas_Object)
	c.dirty = make(map[fyne.CanvasObject]bool)
	c.content = o
	canvases[C.ecore_evas_get(c.window.ee)] = c

	c.setup(o)
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

	c.content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	c.content.Resize(fyne.NewSize(unscaleInt(c, width)-theme.Padding()*2, unscaleInt(c, height)-theme.Padding()*2))
}

func (c *eflCanvas) SetOnKeyDown(keyDown func(*fyne.KeyEvent)) {
	c.onKeyDown = keyDown
}
