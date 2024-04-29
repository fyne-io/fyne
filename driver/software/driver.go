package software

import (
	"fmt"
	"image"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/mobile/event/touch"

	"fyne.io/fyne/v2"
	intapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/painter/software"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/storage/repository"
)

type softwareDriver struct {
	device       fyne.Device
	painter      painter.Painter
	windows      []fyne.Window
	windowsMutex sync.RWMutex

	Output func(image.Image)
	Events chan any

	running  atomic.Bool
	painting atomic.Bool
}

// Declare conformity with Driver
var _ fyne.Driver = (*softwareDriver)(nil)

type device struct {
}

// Declare conformity with Device
var _ fyne.Device = (*device)(nil)

func (d *device) Orientation() fyne.DeviceOrientation {
	return fyne.OrientationVertical
}

func (d *device) HasKeyboard() bool {
	return false
}

func (d *device) SystemScale() float32 {
	return d.SystemScaleForWindow(nil)
}

func (d *device) SystemScaleForWindow(fyne.Window) float32 {
	return 1
}

func (d *device) Locale() fyne.Locale {
	return "en"
}

func (*device) IsBrowser() bool {
	return runtime.GOARCH == "js" || runtime.GOOS == "js"
}

func (d *device) IsMobile() bool {
	return true
}

// NewDriver sets up and registers a new dummy driver for test purpose
func NewDriver(painter func(image.Image), events chan any) fyne.Driver {
	drv := &softwareDriver{
		windowsMutex: sync.RWMutex{},
		Output:       painter,
		device:       &device{},
		Events:       events,
	}
	repository.Register("file", intRepo.NewFileRepository())

	// make a single dummy window for rendering tests
	drv.CreateWindow("")

	return drv
}

// NewDriverWithPainter creates a new dummy driver that will pass the given
// painter to all canvases created
func NewDriverWithPainter(painter painter.Painter) fyne.Driver {
	return &softwareDriver{
		painter:      painter,
		windowsMutex: sync.RWMutex{},
	}
}

func (d *softwareDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return fyne.NewPos(0, 0)
	}

	tc := c.(*softwareCanvas)
	return driver.AbsolutePositionForObject(co, tc.objectTrees())
}

func (d *softwareDriver) AllWindows() []fyne.Window {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	return d.windows
}

func (d *softwareDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	// cheating: probably the last created window is meant
	return d.windows[len(d.windows)-1].Canvas()
}

func (d *softwareDriver) CreateWindow(string) fyne.Window {
	canvas := NewCanvas().(*softwareCanvas)
	if d.painter != nil {
		canvas.SetPainter(d.painter)
	} else {
		canvas.SetPainter(software.NewPainter())
	}

	canvas.Initialize(canvas, canvas.SetDirty)

	window := &softwareWindow{canvas: canvas, driver: d}
	// window.clipboard = &testClipboard{}
	window.InitEventQueue()
	go window.RunEventQueue()

	d.windowsMutex.Lock()
	d.windows = append(d.windows, window)
	d.windowsMutex.Unlock()
	return window
}

func (d *softwareDriver) CurrentWindow() *softwareWindow {
	d.windowsMutex.RLock()
	defer d.windowsMutex.RUnlock()
	if len(d.windows) == 0 {
		return nil
	}
	return d.windows[len(d.windows)-1].(*softwareWindow)
}

func (d *softwareDriver) Device() fyne.Device {
	return d.device
}

func (d *softwareDriver) handlePaint(w *softwareWindow) {
	if !d.painting.CompareAndSwap(false, true) {
		return
	}
	defer d.painting.Store(false)

	c := w.Canvas().(*softwareCanvas)
	// d.painting = false
	if c.Painter() == nil {
		c.SetPainter(software.NewPainter())
	}

	canvasNeedRefresh := c.FreeDirtyTextures() > 0 || c.CheckDirtyAndClear()
	if canvasNeedRefresh {
		fmt.Println("painting")
		t := time.Now()

		// 	newSize := fyne.NewSize(float32(d.currentSize.WidthPx)/c.scale, float32(d.currentSize.HeightPx)/c.scale)

		// if c.EnsureMinSize() {
		// 	c.sizeContent(newSize) // force resize of content
		// } else { // if screen changed
		// 	w.Resize(newSize)
		// }

		d.Output(c.Capture())
		fmt.Println("painting took", time.Since(t))
	}
	cache.Clean(canvasNeedRefresh)
}

// RenderedTextSize looks up how bit a string would be if drawn on screen
func (d *softwareDriver) RenderedTextSize(text string, size float32, style fyne.TextStyle) (fyne.Size, float32) {
	return painter.RenderedTextSize(text, size, style)
}

func (d *softwareDriver) Run() {
	if !d.running.CompareAndSwap(false, true) {
		return // Run was called twice.
	}

	settingsChange := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(settingsChange)
	draw := time.NewTicker(time.Second / 15)

	// make sure we have a theme set
	painter.ClearFontCache()
	cache.ResetThemeCaches()
	intapp.ApplyThemeTo(d.CurrentWindow().Canvas().Content(), d.CurrentWindow().Canvas())

	for d.running.Load() {
		select {
		case <-draw.C:
			w := d.CurrentWindow()
			go d.handlePaint(w)
		case <-settingsChange:
			painter.ClearFontCache()
			cache.ResetThemeCaches()
			intapp.ApplyThemeTo(d.CurrentWindow().Canvas().Content(), d.CurrentWindow().Canvas())
		case e, ok := <-d.Events:
			if !ok {
				return // events channel closed, app done
			}
			current := d.CurrentWindow()
			if current == nil {
				continue
			}
			// c := current.Canvas().(*softwareCanvas)

			switch e := e.(type) {
			// case lifecycle.Event:
			// 	d.handleLifecycle(e, current)
			// case size.Event:
			// 	if e.WidthPx <= 0 {
			// 		continue
			// 	}
			// 	d.currentSize = e
			// 	currentOrientation = e.Orientation
			// 	currentDPI = e.PixelsPerPt * 72
			// 	d.setTheme(e.DarkMode)
			//
			// 	dev := &d.device
			// 	insetChange := dev.safeTop != e.InsetTopPx || dev.safeBottom != e.InsetBottomPx ||
			// 		dev.safeLeft != e.InsetLeftPx || dev.safeRight != e.InsetRightPx
			// 	dev.safeTop = e.InsetTopPx
			// 	dev.safeLeft = e.InsetLeftPx
			// 	dev.safeBottom = e.InsetBottomPx
			// 	dev.safeRight = e.InsetRightPx
			// 	c.scale = fyne.CurrentDevice().SystemScaleForWindow(nil)
			// 	c.Painter().SetFrameBufferScale(1.0)
			//
			// 	if insetChange {
			// 		current.canvas.sizeContent(current.canvas.size) // even if size didn't change we invalidate
			// 	}
			// 	// make sure that we paint on the next frame
			// 	c.Content().Refresh()
			case touch.Event:
				switch e.Type {
				case touch.TypeBegin:
					d.tapDownCanvas(current, e.X, e.Y, e.Sequence)
				// case touch.TypeMove:
				// 	d.tapMoveCanvas(current, e.X, e.Y, e.Sequence)
				case touch.TypeEnd:
					d.tapUpCanvas(current, e.X, e.Y, e.Sequence)
				}
				// case key.Event:
				// 	if e.Direction == key.DirPress {
				// 		d.typeDownCanvas(c, e.Rune, e.Code, e.Modifiers)
				// 	} else if e.Direction == key.DirRelease {
				// 		d.typeUpCanvas(c, e.Rune, e.Code, e.Modifiers)
				// 	}
			}
		}
	}
}

func (d *softwareDriver) StartAnimation(a *fyne.Animation) {
	// currently no animations in test app, we just initialise it and leave
	a.Tick(1.0)
}

func (d *softwareDriver) StopAnimation(a *fyne.Animation) {
	// currently no animations in test app, do nothing
}

func (d *softwareDriver) Quit() {
	d.running.Store(false)
}

func (d *softwareDriver) removeWindow(w *softwareWindow) {
	d.windowsMutex.Lock()
	i := 0
	for _, window := range d.windows {
		if window == w {
			break
		}
		i++
	}

	d.windows = append(d.windows[:i], d.windows[i+1:]...)
	d.windowsMutex.Unlock()
}

func (d *softwareDriver) DoubleTapDelay() time.Duration {
	return 300 * time.Millisecond
}

func (d *softwareDriver) SetDisableScreenBlanking(_ bool) {
	// no-op for test
}
func (d *softwareDriver) tapDownCanvas(w *softwareWindow, x, y float32, tapID touch.Sequence) {
	tapX := scale.ToFyneCoordinate(w.canvas, int(x))
	tapY := scale.ToFyneCoordinate(w.canvas, int(y))
	pos := fyne.NewPos(tapX, tapY)

	fmt.Println("tap down", pos)

	w.canvas.tapDown(pos, int(tapID))
}

//	func (d *softwareDriver) tapMoveCanvas(w *softwareWindow, x, y float32, tapID touch.Sequence) {
//		tapX := scale.ToFyneCoordinate(w.canvas, int(x))
//		tapY := scale.ToFyneCoordinate(w.canvas, int(y))
//		pos := fyne.NewPos(tapX, tapY+tapYOffset)
//
//		w.canvas.tapMove(pos, int(tapID), func(wid fyne.Draggable, ev *fyne.DragEvent) {
//			w.QueueEvent(func() { wid.Dragged(ev) })
//		})
//	}
func (d *softwareDriver) tapUpCanvas(w *softwareWindow, x, y float32, tapID touch.Sequence) {
	tapX := scale.ToFyneCoordinate(w.canvas, int(x))
	tapY := scale.ToFyneCoordinate(w.canvas, int(y))
	pos := fyne.NewPos(tapX, tapY)

	fmt.Println("tap up", pos)

	w.canvas.tapUp(pos, int(tapID), func(wid fyne.Tappable, ev *fyne.PointEvent) {
		w.QueueEvent(func() { wid.Tapped(ev) })
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
		w.QueueEvent(func() { wid.TappedSecondary(ev) })
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		w.QueueEvent(func() { wid.DoubleTapped(ev) })
	}, func(wid fyne.Draggable) {
		w.QueueEvent(wid.DragEnd)
	})
}
