package efl

// #cgo pkg-config: eina evas ecore-evas
// #cgo CFLAGS: -DEFL_BETA_API_SUPPORT=master-compatibility-hack
// #include <Eina.h>
// #include <Evas.h>
// #include <Ecore.h>
// #include <Ecore_Evas.h>
//
// void onObjectMouseDown_cgo(Evas_Object *, void *);
// void onExit_cgo(Ecore_Event_Signal_Exit *);
import "C"

import "log"
import "math"
import "runtime"
import "sync"
import "time"
import "unsafe"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/input"
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

	ev := new(ui.MouseEvent)
	ev.Position = pos
	ev.Button = input.MouseButton(int(info.button))

	switch obj := co.(type) {
	case ui.ClickableObject:
		obj.OnMouseDown(ev)
	case widget.FocussableWidget:
		if canvas.focussed != nil {
			if canvas.focussed == obj {
				return
			}

			canvas.focussed.OnFocusLost()
		}
		canvas.focussed = obj
		obj.OnFocusGained()
	}
}

type eflCanvas struct {
	ui.Canvas
	evas  *C.Evas
	size  ui.Size
	scale float32

	content  ui.CanvasObject
	window   *window
	focussed widget.FocussableWidget

	onKeyDown func(*ui.KeyEvent)

	objects map[*C.Evas_Object]ui.CanvasObject
	native  map[ui.CanvasObject]*C.Evas_Object
}

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

func stepRenderer() {
	C.ecore_main_loop_iterate()
}

func hookIntoEFL() {
	C.ecore_event_handler_add(C.ECORE_EVENT_SIGNAL_EXIT, (C.Ecore_Event_Handler_Cb)(unsafe.Pointer(C.onExit_cgo)), nil)

	tick := time.NewTicker(time.Millisecond * 10)
	go func() {
		for {
			<-tick.C
			mainfunc <- func() {
				stepRenderer()
			}
		}
	}()
}

// initEFL runs our mainthread loop to execute UI functions for EFL
func initEFL() {
	hookIntoEFL()

	for {
		if f, ok := <-mainfunc; ok {
			f()
		} else {
			break
		}
	}
}

// Quit will cause the render loop to end and the application to exit
//export Quit
func Quit() {
	close(mainfunc)
}

// queue of work to run in main thread.
var mainfunc = make(chan func())

// do runs f on the main thread.
func queueRender(c *eflCanvas, co ui.CanvasObject) {
	done := make(chan bool, 1)
	mainfunc <- func() {
		c.doRefresh(co)

		done <- true
	}
	<-done
}

func (c *eflCanvas) buildObject(o ui.CanvasObject, target ui.CanvasObject, size ui.Size) *C.Evas_Object {
	var obj *C.Evas_Object

	switch o.(type) {
	case *canvas.Text:
		obj = C.evas_object_text_add(c.evas)

		to, _ := o.(*canvas.Text)
		C.evas_object_text_text_set(obj, C.CString(to.Text))
		C.evas_object_color_set(obj, C.int(to.Color.R), C.int(to.Color.G),
			C.int(to.Color.B), C.int(to.Color.A))

		updateFont(obj, c, to)
	case *canvas.Rectangle:
		obj = C.evas_object_rectangle_add(c.evas)
		ro, _ := o.(*canvas.Rectangle)

		C.evas_object_color_set(obj, C.int(ro.FillColor.R), C.int(ro.FillColor.G),
			C.int(ro.FillColor.B), C.int(ro.FillColor.A))
	case *canvas.Image:
		obj = C.evas_object_image_add(c.evas)
		img, _ := o.(*canvas.Image)
		C.evas_object_image_alpha_set(obj, C.EINA_FALSE)
		C.evas_object_image_filled_set(obj, C.EINA_TRUE)

		if img.File != "" {
			C.evas_object_image_file_set(obj, C.CString(img.File), nil)
		}
	case *canvas.Line:
		obj = C.evas_object_line_add(c.evas)
		lo, _ := o.(*canvas.Line)

		C.evas_object_color_set(obj, C.int(lo.StrokeColor.R), C.int(lo.StrokeColor.G),
			C.int(lo.StrokeColor.B), C.int(lo.StrokeColor.A))
	default:
		log.Printf("Unrecognised Object %#v\n", o)
		return nil
	}

	c.native[o] = obj
	c.objects[obj] = target
	C.evas_object_event_callback_add(obj, C.EVAS_CALLBACK_MOUSE_DOWN,
		(C.Evas_Object_Event_Cb)(unsafe.Pointer(C.onObjectMouseDown_cgo)),
		nil)

	C.evas_object_show(obj)
	return obj
}

func (c *eflCanvas) buildContainer(objs []ui.CanvasObject, target ui.CanvasObject, size ui.Size) {
	obj := C.evas_object_rectangle_add(c.evas)
	bg := theme.BackgroundColor()
	C.evas_object_color_set(obj, C.int(bg.R), C.int(bg.G), C.int(bg.B), C.int(bg.A))

	C.evas_object_show(obj)
	c.native[target] = obj

	for _, child := range objs {
		switch child.(type) {
		case *ui.Container:
			container := child.(*ui.Container)

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

func (c *eflCanvas) refreshObject(o, o2 ui.CanvasObject, pos ui.Position, size ui.Size) {
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
}

func (c *eflCanvas) refreshContainer(objs []ui.CanvasObject, target ui.CanvasObject, pos ui.Position, size ui.Size) {
	containerPos := pos
	containerSize := size
	if target == c.content {
		containerPos = ui.NewPos(0, 0)
		containerSize = ui.NewSize(size.Width+2*theme.Padding(), size.Height+2*theme.Padding())
	}

	obj := c.native[target]
	bg := theme.BackgroundColor()
	C.evas_object_color_set(obj, C.int(bg.R), C.int(bg.G), C.int(bg.B), C.int(bg.A))
	C.evas_object_geometry_set(obj, C.Evas_Coord(scaleInt(c, containerPos.X)), C.Evas_Coord(scaleInt(c, containerPos.Y)),
		C.Evas_Coord(scaleInt(c, containerSize.Width)), C.Evas_Coord(scaleInt(c, containerSize.Height)))

	for _, child := range objs {
		switch typed := child.(type) {
		case *ui.Container:
			container := child.(*ui.Container)

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

func (c *eflCanvas) Size() ui.Size {
	return c.size
}

func (c *eflCanvas) setup(o ui.CanvasObject) {
	C.ecore_thread_main_loop_begin()

	switch o.(type) {
	case *ui.Container:
		container := o.(*ui.Container)

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

func (c *eflCanvas) Refresh(o ui.CanvasObject) {
	go queueRender(c, o)
}

func (c *eflCanvas) doRefresh(o ui.CanvasObject) {
	c.fitContent()
	position := o.CurrentPosition()
	if o != c.content {
		position = position.Add(ui.NewPos(theme.Padding(), theme.Padding()))
	}

	switch o.(type) {
	case *ui.Container:
		container := o.(*ui.Container)
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

func (c *eflCanvas) Contains(obj ui.CanvasObject) bool {
	return c.native[obj] != nil
}

func (c *eflCanvas) fitContent() {
	var w, h C.int
	C.ecore_evas_geometry_get(c.window.ee, nil, nil, &w, &h)

	min := c.content.MinSize()
	minWidth := scaleInt(c, min.Width+theme.Padding()*2)
	minHeight := scaleInt(c, min.Height+theme.Padding()*2)

	width := ui.Max(minWidth, int(w))
	height := ui.Max(minHeight, int(h))

	C.ecore_evas_size_min_set(c.window.ee, C.int(minWidth), C.int(minHeight))
	C.ecore_evas_resize(c.window.ee, C.int(width), C.int(height))

	c.content.Move(ui.NewPos(theme.Padding(), theme.Padding()))
	c.content.Resize(ui.NewSize(unscaleInt(c, width)-theme.Padding()*2, unscaleInt(c, height)-theme.Padding()*2))
}

func (c *eflCanvas) Content() ui.CanvasObject {
	return c.content
}

func (c *eflCanvas) SetContent(o ui.CanvasObject) {
	canvases[C.ecore_evas_get(c.window.ee)] = c
	c.objects = make(map[*C.Evas_Object]ui.CanvasObject)
	c.native = make(map[ui.CanvasObject]*C.Evas_Object)
	c.content = o

	c.setup(o)
	c.Refresh(o)
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
	log.Println("FYNE_SCALE", scale)

	ratio := scale / c.scale
	c.scale = scale

	var w, h C.int
	C.ecore_evas_geometry_get(c.window.ee, nil, nil, &w, &h)
	width := int(float32(w) * ratio)
	height := int(float32(h) * ratio)
	C.ecore_evas_resize(c.window.ee, C.int(width), C.int(height))

	c.content.Move(ui.NewPos(theme.Padding(), theme.Padding()))
	c.content.Resize(ui.NewSize(unscaleInt(c, width)-theme.Padding()*2, unscaleInt(c, height)-theme.Padding()*2))
}

func (c *eflCanvas) SetOnKeyDown(keyDown func(*ui.KeyEvent)) {
	c.onKeyDown = keyDown
}
