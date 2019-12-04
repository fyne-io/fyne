package gomobile

import (
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/internal"
	"fyne.io/fyne/widget"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter"
	pgl "fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
)

const tapSecondaryDelay = 300 * time.Millisecond

type mobileDriver struct {
	app   app.App
	glctx gl.Context

	windows []fyne.Window
	device  fyne.Device
}

// Declare conformity with Driver
var _ fyne.Driver = (*mobileDriver)(nil)

func init() {
	runtime.LockOSThread()
}

func (d *mobileDriver) CreateWindow(title string) fyne.Window {
	canvas := NewCanvas().(*mobileCanvas) // silence lint
	ret := &window{title: title, canvas: canvas, isChild: len(d.windows) > 0}
	canvas.painter = pgl.NewPainter(canvas, ret)

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

	return d.windows[len(d.windows)-1]
}

func (d *mobileDriver) RenderedTextSize(text string, size int, style fyne.TextStyle) fyne.Size {
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
	var pos fyne.Position
	c := fyne.CurrentApp().Driver().CanvasForObject(co).(*mobileCanvas)

	c.walkTree(func(o fyne.CanvasObject, p fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		if o == co {
			pos = p
			return true
		}
		return false
	}, nil)

	return pos
}

func (d *mobileDriver) Quit() {
	if d.app == nil {
		return
	}

	// TODO? often mobile apps should not allow this...
	d.app.Send(lifecycle.Event{From: lifecycle.StageAlive, To: lifecycle.StageDead, DrawContext: nil})
}

func (d *mobileDriver) Run() {
	app.Main(func(a app.App) {
		d.app = a
		quit := false

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
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					d.onStop()
					d.glctx = nil
				}
				if e.Crosses(lifecycle.StageAlive) == lifecycle.CrossOff {
					quit = true
				}
			case size.Event:
				currentSize = e
			case paint.Event:
				if d.glctx == nil || e.External {
					continue
				}
				if !canvas.inited {
					canvas.inited = true
					canvas.painter.Init() // we cannot init until the context is set above
				}

				d.freeDirtyTextures(canvas)
				d.paintWindow(current, currentSize)
				a.Publish()
				a.Send(paint.Event{})
			case touch.Event:
				switch e.Type {
				case touch.TypeBegin:
					d.tapDownCanvas(canvas, e.X, e.Y)
				case touch.TypeMove:
					d.tapMoveCanvas(canvas, e.X, e.Y)
				case touch.TypeEnd:
					d.tapUpCanvas(canvas, e.X, e.Y)
				}
			case key.Event:
				if e.Direction == key.DirPress {
					d.typeDownCanvas(canvas, e.Rune, e.Code)
				} else if e.Direction == key.DirRelease {
					d.typeUpCanvas(canvas, e.Rune, e.Code)
				}
			}

			if quit {
				break
			}
		}
	})
}

func (d *mobileDriver) onStart() {
	for _, win := range d.AllWindows() {
		win.Canvas().(*mobileCanvas).painter.Init() // we cannot init until the context is set above
	}
}

func (d *mobileDriver) onStop() {
}

func (d *mobileDriver) paintWindow(window fyne.Window, sz size.Event) {
	canvas := window.Canvas().(*mobileCanvas)
	currentOrientation = sz.Orientation

	r, g, b, a := theme.BackgroundColor().RGBA()
	max16bit := float32(255 * 255)
	d.glctx.ClearColor(float32(r)/max16bit, float32(g)/max16bit, float32(b)/max16bit, float32(a)/max16bit)
	d.glctx.Clear(gl.COLOR_BUFFER_BIT)

	newSize := fyne.NewSize(int(float32(sz.WidthPx)/canvas.scale), int(float32(sz.HeightPx)/canvas.scale))
	window.Resize(newSize)

	paint := func(obj fyne.CanvasObject, pos fyne.Position, _ fyne.Position, _ fyne.Size) bool {
		// TODO should this be somehow not scroll container specific?
		if _, ok := obj.(*widget.ScrollContainer); ok {
			canvas.painter.StartClipping(
				fyne.NewPos(pos.X, canvas.Size().Height-pos.Y-obj.Size().Height),
				obj.Size(),
			)
		}
		canvas.painter.Paint(obj, pos, newSize)
		return false
	}
	afterPaint := func(obj, _ fyne.CanvasObject) {
		if _, ok := obj.(*widget.ScrollContainer); ok {
			canvas.painter.StopClipping()
		}
	}

	canvas.walkTree(paint, afterPaint)
}

func (d *mobileDriver) tapDownCanvas(canvas *mobileCanvas, x, y float32) {
	tapX := internal.UnscaleInt(canvas, int(x))
	tapY := internal.UnscaleInt(canvas, int(y))
	pos := fyne.NewPos(tapX, tapY)

	canvas.tapDown(pos)
}

func (d *mobileDriver) tapMoveCanvas(canvas *mobileCanvas, x, y float32) {
	tapX := internal.UnscaleInt(canvas, int(x))
	tapY := internal.UnscaleInt(canvas, int(y))
	pos := fyne.NewPos(tapX, tapY)

	canvas.tapMove(pos, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		go wid.Dragged(ev)
	})
}

func (d *mobileDriver) tapUpCanvas(canvas *mobileCanvas, x, y float32) {
	tapX := internal.UnscaleInt(canvas, int(x))
	tapY := internal.UnscaleInt(canvas, int(y))
	pos := fyne.NewPos(tapX, tapY)

	canvas.tapUp(pos, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		go wid.Tapped(ev)
	}, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		go wid.TappedSecondary(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		go wid.DragEnd()
	})
}

func keyToName(code key.Code) fyne.KeyName {
	switch code {
	// non-printable
	case key.CodeEscape:
		return fyne.KeyEscape
	case key.CodeReturnEnter:
		return fyne.KeyReturn
	case key.CodeTab:
		return fyne.KeyTab
	case key.CodeDeleteBackspace:
		return fyne.KeyBackspace
	case key.CodeInsert:
		return fyne.KeyInsert
	case key.CodePageUp:
		return fyne.KeyPageUp
	case key.CodePageDown:
		return fyne.KeyPageDown
	case key.CodeHome:
		return fyne.KeyHome
	case key.CodeEnd:
		return fyne.KeyEnd

	case key.CodeF1:
		return fyne.KeyF1
	case key.CodeF2:
		return fyne.KeyF2
	case key.CodeF3:
		return fyne.KeyF3
	case key.CodeF4:
		return fyne.KeyF4
	case key.CodeF5:
		return fyne.KeyF5
	case key.CodeF6:
		return fyne.KeyF6
	case key.CodeF7:
		return fyne.KeyF7
	case key.CodeF8:
		return fyne.KeyF8
	case key.CodeF9:
		return fyne.KeyF9
	case key.CodeF10:
		return fyne.KeyF10
	case key.CodeF11:
		return fyne.KeyF11
	case key.CodeF12:
		return fyne.KeyF12

	case key.CodeKeypadEnter:
		return fyne.KeyEnter

	// printable
	case key.CodeA:
		return fyne.KeyA
	case key.CodeB:
		return fyne.KeyB
	case key.CodeC:
		return fyne.KeyC
	case key.CodeD:
		return fyne.KeyD
	case key.CodeE:
		return fyne.KeyE
	case key.CodeF:
		return fyne.KeyF
	case key.CodeG:
		return fyne.KeyG
	case key.CodeH:
		return fyne.KeyH
	case key.CodeI:
		return fyne.KeyI
	case key.CodeJ:
		return fyne.KeyJ
	case key.CodeK:
		return fyne.KeyK
	case key.CodeL:
		return fyne.KeyL
	case key.CodeM:
		return fyne.KeyM
	case key.CodeN:
		return fyne.KeyN
	case key.CodeO:
		return fyne.KeyO
	case key.CodeP:
		return fyne.KeyP
	case key.CodeQ:
		return fyne.KeyQ
	case key.CodeR:
		return fyne.KeyR
	case key.CodeS:
		return fyne.KeyS
	case key.CodeT:
		return fyne.KeyT
	case key.CodeU:
		return fyne.KeyU
	case key.CodeV:
		return fyne.KeyV
	case key.CodeW:
		return fyne.KeyW
	case key.CodeX:
		return fyne.KeyX
	case key.CodeY:
		return fyne.KeyY
	case key.CodeZ:
		return fyne.KeyZ
	case key.Code0, key.CodeKeypad0:
		return fyne.Key0
	case key.Code1, key.CodeKeypad1:
		return fyne.Key1
	case key.Code2, key.CodeKeypad2:
		return fyne.Key2
	case key.Code3, key.CodeKeypad3:
		return fyne.Key3
	case key.Code4, key.CodeKeypad4:
		return fyne.Key4
	case key.Code5, key.CodeKeypad5:
		return fyne.Key5
	case key.Code6, key.CodeKeypad6:
		return fyne.Key6
	case key.Code7, key.CodeKeypad7:
		return fyne.Key7
	case key.Code8, key.CodeKeypad8:
		return fyne.Key8
	case key.Code9, key.CodeKeypad9:
		return fyne.Key9

	case key.CodeSemicolon:
		return fyne.KeySemicolon
	case key.CodeEqualSign:
		return fyne.KeyEqual

	case key.CodeSpacebar:
		return fyne.KeySpace
	case key.CodeApostrophe:
		return fyne.KeyApostrophe
	case key.CodeComma:
		return fyne.KeyComma
	case key.CodeHyphenMinus, key.CodeKeypadHyphenMinus:
		return fyne.KeyMinus
	case key.CodeFullStop, key.CodeKeypadFullStop:
		return fyne.KeyPeriod
	case key.CodeSlash:
		return fyne.KeySlash
	case key.CodeLeftSquareBracket:
		return fyne.KeyLeftBracket
	case key.CodeBackslash:
		return fyne.KeyBackslash
	case key.CodeRightSquareBracket:
		return fyne.KeyRightBracket
	}
	return ""
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

func (d *mobileDriver) freeDirtyTextures(canvas *mobileCanvas) {
	for {
		select {
		case object := <-canvas.refreshQueue:
			freeWalked := func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
				canvas.painter.Free(obj)
				return false
			}
			driver.WalkCompleteObjectTree(object, freeWalked, nil)
		default:
			return
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
	driver := new(mobileDriver)

	return driver
}
