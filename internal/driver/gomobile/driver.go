package gomobile

import (
	"runtime"
	"strconv"
	"time"

	"github.com/fyne-io/mobile/app"
	"github.com/fyne-io/mobile/event/key"
	"github.com/fyne-io/mobile/event/lifecycle"
	"github.com/fyne-io/mobile/event/paint"
	"github.com/fyne-io/mobile/event/size"
	"github.com/fyne-io/mobile/event/touch"
	"github.com/fyne-io/mobile/gl"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/animation"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/painter"
	pgl "fyne.io/fyne/v2/internal/painter/gl"
	"fyne.io/fyne/v2/theme"
)

const (
	tapMoveThreshold  = 4.0                    // how far can we move before it is a drag
	tapSecondaryDelay = 300 * time.Millisecond // how long before secondary tap
)

type mobileDriver struct {
	app   app.App
	glctx gl.Context

	windows   []fyne.Window
	device    *device
	animation *animation.Runner
}

// Declare conformity with Driver
var _ fyne.Driver = (*mobileDriver)(nil)

func init() {
	runtime.LockOSThread()
}

func (d *mobileDriver) CreateWindow(title string) fyne.Window {
	c := NewCanvas().(*mobileCanvas) // silence lint
	ret := &window{title: title, canvas: c, isChild: len(d.windows) > 0}
	c.content = &canvas.Rectangle{FillColor: theme.BackgroundColor()}
	c.painter = pgl.NewPainter(c, ret)
	d.windows = append(d.windows, ret)
	return ret
}

func (d *mobileDriver) AllWindows() []fyne.Window {
	return d.windows
}

// currentWindow returns the most recently opened window - we can only show one at a time.
func (d *mobileDriver) currentWindow() fyne.Window {
	if len(d.windows) == 0 {
		return nil
	}

	var last fyne.Window
	for i := len(d.windows) - 1; i >= 0; i-- {
		last = d.windows[i]
		if last.(*window).visible {
			return last
		}
	}

	return last
}

func (d *mobileDriver) RenderedTextSize(text string, size float32, style fyne.TextStyle) fyne.Size {
	return painter.RenderedTextSize(text, size, style)
}

func (d *mobileDriver) CanvasForObject(fyne.CanvasObject) fyne.Canvas {
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
	pos := driver.AbsolutePositionForObject(co, mc.objectTrees())
	inset, _ := c.InteractiveArea()

	if mc.windowHead != nil {
		if len(mc.windowHead.(*fyne.Container).Objects) > 1 {
			topHeight := mc.windowHead.MinSize().Height
			pos = pos.Subtract(fyne.NewSize(0, topHeight))
		}
	}
	return pos.Subtract(inset)
}

func (d *mobileDriver) Quit() {
	// Android and iOS guidelines say this should not be allowed!
}

func (d *mobileDriver) Run() {
	app.Main(func(a app.App) {
		d.app = a

		var currentSize size.Event
		for e := range a.Events() {
			current := d.currentWindow()
			if current == nil {
				continue
			}
			canvas := current.Canvas().(*mobileCanvas)

			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					d.glctx, _ = e.DrawContext.(gl.Context)
					d.onStart()

					// this is a fix for some android phone to prevent the app from being drawn as a blank screen after being pushed in the background
					canvas.Content().Refresh()

					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					d.onStop()
					d.glctx = nil
				}
				switch e.Crosses(lifecycle.StageFocused) {
				case lifecycle.CrossOff: // will enter background
					if runtime.GOOS == "darwin" {
						if d.glctx == nil {
							continue
						}

						size := fyne.NewSize(float32(currentSize.WidthPx)/canvas.scale, float32(currentSize.HeightPx)/canvas.scale)
						d.paintWindow(current, size)
						a.Publish()
					}
				}
			case size.Event:
				if e.WidthPx <= 0 {
					continue
				}
				currentSize = e
				currentOrientation = e.Orientation
				currentDPI = e.PixelsPerPt * 72

				dev := d.device
				dev.safeTop = e.InsetTopPx
				dev.safeLeft = e.InsetLeftPx
				dev.safeHeight = e.HeightPx - e.InsetTopPx - e.InsetBottomPx
				dev.safeWidth = e.WidthPx - e.InsetLeftPx - e.InsetRightPx
				canvas.scale = fyne.CurrentDevice().SystemScaleForWindow(nil)
				canvas.painter.SetFrameBufferScale(1.0)

				// make sure that we paint on the next frame
				canvas.Content().Refresh()
			case paint.Event:
				if d.glctx == nil || e.External {
					continue
				}
				if !canvas.inited {
					canvas.inited = true
					canvas.painter.Init() // we cannot init until the context is set above
				}

				if d.freeDirtyTextures(canvas) {
					newSize := fyne.NewSize(float32(currentSize.WidthPx)/canvas.scale, float32(currentSize.HeightPx)/canvas.scale)

					if canvas.minSizeChanged() {
						canvas.ensureMinSize()

						canvas.sizeContent(newSize) // force resize of content
					} else { // if screen changed
						current.Resize(newSize)
					}

					d.paintWindow(current, newSize)
					a.Publish()
				}

				time.Sleep(time.Millisecond * 10)
				a.Send(paint.Event{})
			case touch.Event:
				switch e.Type {
				case touch.TypeBegin:
					d.tapDownCanvas(canvas, e.X, e.Y, e.Sequence)
				case touch.TypeMove:
					d.tapMoveCanvas(canvas, e.X, e.Y, e.Sequence)
				case touch.TypeEnd:
					d.tapUpCanvas(canvas, e.X, e.Y, e.Sequence)
				}
			case key.Event:
				if e.Direction == key.DirPress {
					d.typeDownCanvas(canvas, e.Rune, e.Code)
				} else if e.Direction == key.DirRelease {
					d.typeUpCanvas(canvas, e.Rune, e.Code)
				}
			}
		}
	})
}

func (d *mobileDriver) onStart() {
}

func (d *mobileDriver) onStop() {
}

func (d *mobileDriver) paintWindow(window fyne.Window, size fyne.Size) {
	clips := &internal.ClipStack{}
	canvas := window.Canvas().(*mobileCanvas)

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	d.glctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	d.glctx.Clear(gl.COLOR_BUFFER_BIT)

	paint := func(obj fyne.CanvasObject, pos fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if _, ok := obj.(fyne.Scrollable); ok {
			inner := clips.Push(pos, obj.Size())
			canvas.painter.StartClipping(inner.Rect())
		}
		canvas.painter.Paint(obj, pos, size)
		return false
	}
	afterPaint := func(obj, _ fyne.CanvasObject) {
		if _, ok := obj.(fyne.Scrollable); ok {
			canvas.painter.StopClipping()
			clips.Pop()
			if top := clips.Top(); top != nil {
				canvas.painter.StartClipping(top.Rect())
			}
		}
	}

	canvas.walkTree(paint, afterPaint)
}

func (d *mobileDriver) tapDownCanvas(canvas *mobileCanvas, x, y float32, tapID touch.Sequence) {
	tapX := internal.UnscaleInt(canvas, int(x))
	tapY := internal.UnscaleInt(canvas, int(y))
	pos := fyne.NewPos(tapX, tapY+tapYOffset)

	canvas.tapDown(pos, int(tapID))
}

func (d *mobileDriver) tapMoveCanvas(canvas *mobileCanvas, x, y float32, tapID touch.Sequence) {
	tapX := internal.UnscaleInt(canvas, int(x))
	tapY := internal.UnscaleInt(canvas, int(y))
	pos := fyne.NewPos(tapX, tapY+tapYOffset)

	canvas.tapMove(pos, int(tapID), func(wid fyne.Draggable, ev *fyne.DragEvent) {
		go wid.Dragged(ev)
	})
}

func (d *mobileDriver) tapUpCanvas(canvas *mobileCanvas, x, y float32, tapID touch.Sequence) {
	tapX := internal.UnscaleInt(canvas, int(x))
	tapY := internal.UnscaleInt(canvas, int(y))
	pos := fyne.NewPos(tapX, tapY+tapYOffset)

	canvas.tapUp(pos, int(tapID), func(wid fyne.Tappable, ev *fyne.PointEvent) {
		go wid.Tapped(ev)
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
		go wid.TappedSecondary(ev)
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		go wid.DoubleTapped(ev)
	}, func(wid fyne.Draggable) {
		go wid.DragEnd()
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

func (d *mobileDriver) typeDownCanvas(canvas *mobileCanvas, r rune, code key.Code) {
	keyName := keyToName(code)
	r = runeToPrintable(r)
	keyEvent := &fyne.KeyEvent{Name: keyName}

	if canvas.Focused() != nil {
		if keyName != "" {
			canvas.Focused().TypedKey(keyEvent)
		}
		if r > 0 {
			canvas.Focused().TypedRune(r)
		}
	} else if canvas.onTypedKey != nil {
		if keyName != "" {
			canvas.onTypedKey(keyEvent)
		}
		if r > 0 {
			canvas.onTypedRune(r)
		}
	}
}

func (d *mobileDriver) typeUpCanvas(canvas *mobileCanvas, r rune, code key.Code) {

}

func (d *mobileDriver) freeDirtyTextures(canvas *mobileCanvas) bool {
	freed := false
	for {
		select {
		case object := <-canvas.refreshQueue:
			freed = true
			freeWalked := func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
				canvas.painter.Free(obj)
				return false
			}
			driver.WalkCompleteObjectTree(object, freeWalked, nil)
		default:
			return freed
		}
	}
}

func (d *mobileDriver) Device() fyne.Device {
	if d.device == nil {
		d.device = &device{}
	}

	return d.device
}

// NewGoMobileDriver sets up a new Driver instance implemented using the Go
// Mobile extension and OpenGL bindings.
func NewGoMobileDriver() fyne.Driver {
	d := new(mobileDriver)
	d.animation = &animation.Runner{}

	registerRepository(d)
	return d
}
