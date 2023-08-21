package mobile

import (
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/animation"
	intapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/driver/mobile/app"
	"fyne.io/fyne/v2/internal/driver/mobile/event/key"
	"fyne.io/fyne/v2/internal/driver/mobile/event/lifecycle"
	"fyne.io/fyne/v2/internal/driver/mobile/event/paint"
	"fyne.io/fyne/v2/internal/driver/mobile/event/size"
	"fyne.io/fyne/v2/internal/driver/mobile/event/touch"
	"fyne.io/fyne/v2/internal/driver/mobile/gl"
	"fyne.io/fyne/v2/internal/painter"
	pgl "fyne.io/fyne/v2/internal/painter/gl"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/theme"
)

const (
	tapMoveThreshold  = 4.0                    // how far can we move before it is a drag
	tapSecondaryDelay = 300 * time.Millisecond // how long before secondary tap
)

// Configuration is the system information about the current device
type Configuration struct {
	SystemTheme fyne.ThemeVariant
}

// ConfiguredDriver is a simple type that allows packages to hook into configuration changes of this driver.
type ConfiguredDriver interface {
	SetOnConfigurationChanged(func(*Configuration))
}

type mobileDriver struct {
	app   app.App
	glctx gl.Context

	windows     []fyne.Window
	device      *device
	animation   *animation.Runner
	currentSize size.Event

	theme           fyne.ThemeVariant
	onConfigChanged func(*Configuration)
	painting        bool
}

// Declare conformity with Driver
var _ fyne.Driver = (*mobileDriver)(nil)
var _ ConfiguredDriver = (*mobileDriver)(nil)

func init() {
	runtime.LockOSThread()
}

func (d *mobileDriver) CreateWindow(title string) fyne.Window {
	c := NewCanvas().(*mobileCanvas) // silence lint
	ret := &window{title: title, canvas: c, isChild: len(d.windows) > 0}
	ret.InitEventQueue()
	go ret.RunEventQueue()
	c.setContent(&canvas.Rectangle{FillColor: theme.BackgroundColor()})
	c.SetPainter(pgl.NewPainter(c, ret))
	d.windows = append(d.windows, ret)
	return ret
}

func (d *mobileDriver) AllWindows() []fyne.Window {
	return d.windows
}

// currentWindow returns the most recently opened window - we can only show one at a time.
func (d *mobileDriver) currentWindow() *window {
	if len(d.windows) == 0 {
		return nil
	}

	var last *window
	for i := len(d.windows) - 1; i >= 0; i-- {
		last = d.windows[i].(*window)
		if last.visible {
			return last
		}
	}

	return last
}

func (d *mobileDriver) RenderedTextSize(text string, textSize float32, style fyne.TextStyle) (size fyne.Size, baseline float32) {
	return painter.RenderedTextSize(text, textSize, style)
}

func (d *mobileDriver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	if len(d.windows) == 0 {
		return nil
	}

	// TODO figure out how we handle multiple windows...
	return d.currentWindow().Canvas()
}

func (d *mobileDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return fyne.NewPos(0, 0)
	}

	mc := c.(*mobileCanvas)
	pos := driver.AbsolutePositionForObject(co, mc.ObjectTrees())
	inset, _ := c.InteractiveArea()

	if mc.windowHead != nil {
		if len(mc.windowHead.(*fyne.Container).Objects) > 1 {
			topHeight := mc.windowHead.MinSize().Height
			pos = pos.Subtract(fyne.NewSize(0, topHeight))
		}
	}
	return pos.Subtract(inset)
}

func (d *mobileDriver) GoBack() {
	app.GoBack()
}

func (d *mobileDriver) Quit() {
	// Android and iOS guidelines say this should not be allowed!
}

func (d *mobileDriver) Run() {
	app.Main(func(a app.App) {
		d.app = a
		settingsChange := make(chan fyne.Settings)
		fyne.CurrentApp().Settings().AddChangeListener(settingsChange)
		draw := time.NewTicker(time.Second / 60)

		for {
			select {
			case <-draw.C:
				d.sendPaintEvent()
			case set := <-settingsChange:
				painter.ClearFontCache()
				cache.ResetThemeCaches()
				intapp.ApplySettingsWithCallback(set, fyne.CurrentApp(), func(w fyne.Window) {
					c, ok := w.Canvas().(*mobileCanvas)
					if !ok {
						return
					}
					c.applyThemeOutOfTreeObjects()
				})
			case e, ok := <-a.Events():
				if !ok {
					return // events channel closed, app done
				}
				current := d.currentWindow()
				if current == nil {
					continue
				}
				c := current.Canvas().(*mobileCanvas)

				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					d.handleLifecycle(e, current)
				case size.Event:
					if e.WidthPx <= 0 {
						continue
					}
					d.currentSize = e
					currentOrientation = e.Orientation
					currentDPI = e.PixelsPerPt * 72
					d.setTheme(e.DarkMode)

					dev := d.device
					dev.safeTop = e.InsetTopPx
					dev.safeLeft = e.InsetLeftPx
					dev.safeHeight = e.HeightPx - e.InsetTopPx - e.InsetBottomPx
					dev.safeWidth = e.WidthPx - e.InsetLeftPx - e.InsetRightPx
					c.scale = fyne.CurrentDevice().SystemScaleForWindow(nil)
					c.Painter().SetFrameBufferScale(1.0)

					// make sure that we paint on the next frame
					c.Content().Refresh()
				case paint.Event:
					d.handlePaint(e, current)
				case touch.Event:
					switch e.Type {
					case touch.TypeBegin:
						d.tapDownCanvas(current, e.X, e.Y, e.Sequence)
					case touch.TypeMove:
						d.tapMoveCanvas(current, e.X, e.Y, e.Sequence)
					case touch.TypeEnd:
						d.tapUpCanvas(current, e.X, e.Y, e.Sequence)
					}
				case key.Event:
					if e.Direction == key.DirPress {
						d.typeDownCanvas(c, e.Rune, e.Code, e.Modifiers)
					} else if e.Direction == key.DirRelease {
						d.typeUpCanvas(c, e.Rune, e.Code, e.Modifiers)
					}
				}
			}
		}
	})
}

func (d *mobileDriver) handleLifecycle(e lifecycle.Event, w fyne.Window) {
	c := w.Canvas().(*mobileCanvas)
	switch e.Crosses(lifecycle.StageVisible) {
	case lifecycle.CrossOn:
		d.glctx, _ = e.DrawContext.(gl.Context)
		d.onStart()

		// this is a fix for some android phone to prevent the app from being drawn as a blank screen after being pushed in the background
		c.Content().Refresh()

		d.sendPaintEvent()
	case lifecycle.CrossOff:
		d.onStop()
		d.glctx = nil
	}
	switch e.Crosses(lifecycle.StageFocused) {
	case lifecycle.CrossOn: // foregrounding
		fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerEnteredForeground()
	case lifecycle.CrossOff: // will enter background
		if runtime.GOOS == "darwin" {
			if d.glctx == nil {
				return
			}

			s := fyne.NewSize(float32(d.currentSize.WidthPx)/c.scale, float32(d.currentSize.HeightPx)/c.scale)
			d.paintWindow(w, s)
			d.app.Publish()
		}
		fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerExitedForeground()
	}
}

func (d *mobileDriver) handlePaint(e paint.Event, w fyne.Window) {
	c := w.Canvas().(*mobileCanvas)
	d.painting = false
	if d.glctx == nil || e.External {
		return
	}
	if !c.inited {
		c.inited = true
		c.Painter().Init() // we cannot init until the context is set above
	}

	canvasNeedRefresh := c.FreeDirtyTextures() > 0 || c.CheckDirtyAndClear()
	if canvasNeedRefresh {
		newSize := fyne.NewSize(float32(d.currentSize.WidthPx)/c.scale, float32(d.currentSize.HeightPx)/c.scale)

		if c.EnsureMinSize() {
			c.sizeContent(newSize) // force resize of content
		} else { // if screen changed
			w.Resize(newSize)
		}

		d.paintWindow(w, newSize)
		d.app.Publish()
	}
	cache.Clean(canvasNeedRefresh)
}

func (d *mobileDriver) onStart() {
	fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerStarted()
}

func (d *mobileDriver) onStop() {
	fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerStopped()
}

func (d *mobileDriver) paintWindow(window fyne.Window, size fyne.Size) {
	clips := &internal.ClipStack{}
	c := window.Canvas().(*mobileCanvas)

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	d.glctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	d.glctx.Clear(gl.ColorBufferBit)

	draw := func(node *common.RenderCacheNode, pos fyne.Position) {
		obj := node.Obj()
		if _, ok := obj.(fyne.Scrollable); ok {
			inner := clips.Push(pos, obj.Size())
			c.Painter().StartClipping(inner.Rect())
		}

		if size.Width <= 0 || size.Height <= 0 { // iconifying on Windows can do bad things
			return
		}
		c.Painter().Paint(obj, pos, size)
	}
	afterDraw := func(node *common.RenderCacheNode, pos fyne.Position) {
		if _, ok := node.Obj().(fyne.Scrollable); ok {
			c.Painter().StopClipping()
			clips.Pop()
			if top := clips.Top(); top != nil {
				c.Painter().StartClipping(top.Rect())
			}
		}

		if c.debug {
			c.DrawDebugOverlay(node.Obj(), pos, size)
		}
	}

	c.WalkTrees(draw, afterDraw)
}

func (d *mobileDriver) sendPaintEvent() {
	if d.painting {
		return
	}
	d.app.Send(paint.Event{})
	d.painting = true
}

func (d *mobileDriver) setTheme(dark bool) {
	var mode fyne.ThemeVariant
	if dark {
		mode = theme.VariantDark
	} else {
		mode = theme.VariantLight
	}

	if d.theme != mode && d.onConfigChanged != nil {
		d.onConfigChanged(&Configuration{SystemTheme: mode})
	}
	d.theme = mode
}

func (d *mobileDriver) tapDownCanvas(w *window, x, y float32, tapID touch.Sequence) {
	tapX := scale.ToFyneCoordinate(w.canvas, int(x))
	tapY := scale.ToFyneCoordinate(w.canvas, int(y))
	pos := fyne.NewPos(tapX, tapY+tapYOffset)

	w.canvas.tapDown(pos, int(tapID))
}

func (d *mobileDriver) tapMoveCanvas(w *window, x, y float32, tapID touch.Sequence) {
	tapX := scale.ToFyneCoordinate(w.canvas, int(x))
	tapY := scale.ToFyneCoordinate(w.canvas, int(y))
	pos := fyne.NewPos(tapX, tapY+tapYOffset)

	w.canvas.tapMove(pos, int(tapID), func(wid fyne.Draggable, ev *fyne.DragEvent) {
		w.QueueEvent(func() { wid.Dragged(ev) })
	})
}

func (d *mobileDriver) tapUpCanvas(w *window, x, y float32, tapID touch.Sequence) {
	tapX := scale.ToFyneCoordinate(w.canvas, int(x))
	tapY := scale.ToFyneCoordinate(w.canvas, int(y))
	pos := fyne.NewPos(tapX, tapY+tapYOffset)

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

var keyCodeMap = map[key.Code]fyne.KeyName{
	// non-printable
	key.CodeEscape:          fyne.KeyEscape,
	key.CodeReturnEnter:     fyne.KeyReturn,
	key.CodeTab:             fyne.KeyTab,
	key.CodeDeleteBackspace: fyne.KeyBackspace,
	key.CodeInsert:          fyne.KeyInsert,
	key.CodePageUp:          fyne.KeyPageUp,
	key.CodePageDown:        fyne.KeyPageDown,
	key.CodeHome:            fyne.KeyHome,
	key.CodeEnd:             fyne.KeyEnd,

	key.CodeF1:  fyne.KeyF1,
	key.CodeF2:  fyne.KeyF2,
	key.CodeF3:  fyne.KeyF3,
	key.CodeF4:  fyne.KeyF4,
	key.CodeF5:  fyne.KeyF5,
	key.CodeF6:  fyne.KeyF6,
	key.CodeF7:  fyne.KeyF7,
	key.CodeF8:  fyne.KeyF8,
	key.CodeF9:  fyne.KeyF9,
	key.CodeF10: fyne.KeyF10,
	key.CodeF11: fyne.KeyF11,
	key.CodeF12: fyne.KeyF12,

	key.CodeKeypadEnter: fyne.KeyEnter,

	// printable
	key.CodeA:       fyne.KeyA,
	key.CodeB:       fyne.KeyB,
	key.CodeC:       fyne.KeyC,
	key.CodeD:       fyne.KeyD,
	key.CodeE:       fyne.KeyE,
	key.CodeF:       fyne.KeyF,
	key.CodeG:       fyne.KeyG,
	key.CodeH:       fyne.KeyH,
	key.CodeI:       fyne.KeyI,
	key.CodeJ:       fyne.KeyJ,
	key.CodeK:       fyne.KeyK,
	key.CodeL:       fyne.KeyL,
	key.CodeM:       fyne.KeyM,
	key.CodeN:       fyne.KeyN,
	key.CodeO:       fyne.KeyO,
	key.CodeP:       fyne.KeyP,
	key.CodeQ:       fyne.KeyQ,
	key.CodeR:       fyne.KeyR,
	key.CodeS:       fyne.KeyS,
	key.CodeT:       fyne.KeyT,
	key.CodeU:       fyne.KeyU,
	key.CodeV:       fyne.KeyV,
	key.CodeW:       fyne.KeyW,
	key.CodeX:       fyne.KeyX,
	key.CodeY:       fyne.KeyY,
	key.CodeZ:       fyne.KeyZ,
	key.Code0:       fyne.Key0,
	key.CodeKeypad0: fyne.Key0,
	key.Code1:       fyne.Key1,
	key.CodeKeypad1: fyne.Key1,
	key.Code2:       fyne.Key2,
	key.CodeKeypad2: fyne.Key2,
	key.Code3:       fyne.Key3,
	key.CodeKeypad3: fyne.Key3,
	key.Code4:       fyne.Key4,
	key.CodeKeypad4: fyne.Key4,
	key.Code5:       fyne.Key5,
	key.CodeKeypad5: fyne.Key5,
	key.Code6:       fyne.Key6,
	key.CodeKeypad6: fyne.Key6,
	key.Code7:       fyne.Key7,
	key.CodeKeypad7: fyne.Key7,
	key.Code8:       fyne.Key8,
	key.CodeKeypad8: fyne.Key8,
	key.Code9:       fyne.Key9,
	key.CodeKeypad9: fyne.Key9,

	key.CodeSemicolon: fyne.KeySemicolon,
	key.CodeEqualSign: fyne.KeyEqual,

	key.CodeSpacebar:           fyne.KeySpace,
	key.CodeApostrophe:         fyne.KeyApostrophe,
	key.CodeComma:              fyne.KeyComma,
	key.CodeHyphenMinus:        fyne.KeyMinus,
	key.CodeKeypadHyphenMinus:  fyne.KeyMinus,
	key.CodeFullStop:           fyne.KeyPeriod,
	key.CodeKeypadFullStop:     fyne.KeyPeriod,
	key.CodeSlash:              fyne.KeySlash,
	key.CodeLeftSquareBracket:  fyne.KeyLeftBracket,
	key.CodeBackslash:          fyne.KeyBackslash,
	key.CodeRightSquareBracket: fyne.KeyRightBracket,
	key.CodeGraveAccent:        fyne.KeyBackTick,

	key.CodeBackButton: mobile.KeyBack,
}

func keyToName(code key.Code) fyne.KeyName {
	ret, ok := keyCodeMap[code]
	if !ok {
		return ""
	}

	return ret
}

func runeToPrintable(r rune) rune {
	if strconv.IsPrint(r) {
		return r
	}

	return 0
}

func (d *mobileDriver) typeDownCanvas(canvas *mobileCanvas, r rune, code key.Code, mod key.Modifiers) {
	keyName := keyToName(code)
	switch keyName {
	case fyne.KeyTab:
		capture := false
		if ent, ok := canvas.Focused().(fyne.Tabbable); ok {
			capture = ent.AcceptsTab()
		}
		if !capture {
			switch mod {
			case 0:
				canvas.FocusNext()
				return
			case key.ModShift:
				canvas.FocusPrevious()
				return
			}
		}
	}

	r = runeToPrintable(r)
	keyEvent := &fyne.KeyEvent{Name: keyName}

	if canvas.Focused() != nil {
		if keyName != "" {
			canvas.Focused().TypedKey(keyEvent)
		}
		if r > 0 {
			canvas.Focused().TypedRune(r)
		}
	} else {
		if keyName != "" {
			if canvas.onTypedKey != nil {
				canvas.onTypedKey(keyEvent)
			} else if keyName == mobile.KeyBack {
				d.GoBack()
			}
		}
		if r > 0 && canvas.onTypedRune != nil {
			canvas.onTypedRune(r)
		}
	}
}

func (d *mobileDriver) typeUpCanvas(_ *mobileCanvas, _ rune, _ key.Code, _ key.Modifiers) {
}

func (d *mobileDriver) Device() fyne.Device {
	if d.device == nil {
		d.device = &device{}
	}

	return d.device
}

func (d *mobileDriver) SetOnConfigurationChanged(f func(*Configuration)) {
	d.onConfigChanged = f
}

// NewGoMobileDriver sets up a new Driver instance implemented using the Go
// Mobile extension and OpenGL bindings.
func NewGoMobileDriver() fyne.Driver {
	d := &mobileDriver{
		theme:     fyne.ThemeVariant(2), // unspecified
		animation: &animation.Runner{},
	}

	registerRepository(d)
	return d
}
