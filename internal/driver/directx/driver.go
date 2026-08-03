//go:build windows

// Package directx implements a Fyne desktop driver for Windows on native Win32,
// rendering through the Direct3D 11 painter in internal/painter/dx. It needs
// neither cgo nor GLFW - every Win32 call here goes through syscall.
//
// The split mirrors the OpenGL driver: this package owns windows, input, menus,
// the clipboard and the run loop, the way internal/driver/glfw does, while all
// drawing lives in the painter package.
//
// Select it with the `directx` build tag:
//
//	go build -tags directx
package directx

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/animation"
	intapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	paint "fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/painter/dx"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage/repository"
)

// Declare conformity with Driver.
var _ fyne.Driver = (*dxDriver)(nil)

// curWindow tracks the focused window for lifecycle events.
var curWindow *window

type dxDriver struct {
	windows     []fyne.Window
	initialized bool
	done        chan struct{}

	animation animation.Runner

	trayStart, trayStop func()     // start/stop the system tray, when one is used
	systrayMenu         *fyne.Menu // the tray menu, nil when there is none
}

// NewDriver creates a Fyne driver backed by native Win32 windows and Direct3D 11.
func NewDriver() fyne.Driver {
	repository.Register(fyne.URISchemeFile, intRepo.NewFileRepository())
	return &dxDriver{done: make(chan struct{})}
}

func (*dxDriver) DoFromGoroutine(f func(), wait bool) {
	if wait {
		async.EnsureNotMain(func() {
			runOnMainWithWait(f, true)
		})
		return
	}
	runOnMainWithWait(f, false)
}

func (d *dxDriver) CreateWindow(title string) fyne.Window {
	return d.createWindow(title, true)
}

// CreateSplashWindow makes a borderless, unpadded window centred on screen.
func (d *dxDriver) CreateSplashWindow() fyne.Window {
	w := d.createWindow("", false)
	w.SetPadded(false)
	w.CenterOnScreen()
	return w
}

func (d *dxDriver) createWindow(title string, decorate bool) fyne.Window {
	d.init()

	w := &window{
		driver:   d,
		title:    title,
		decorate: decorate,
		width:    fallbackWidth,
		height:   fallbackHeight,
	}
	w.canvas = newCanvas()
	w.canvas.win = w

	// CreateWindowEx binds the window to the calling thread, so everything from
	// here on must happen on main.
	async.EnsureMain(func() {
		if err := registerWindowClass(); err != nil {
			fyne.LogError("directx: registering window class", err)
			return
		}

		pendingWindow = w
		style := w.style()
		r := rect{Right: int32(w.width), Bottom: int32(w.height)}
		adjustWindowRect(&r, style, 96, false)
		hwnd, err := createWindowEx(0, utf16Ptr(windowClassName), utf16Ptr(title), style,
			cwUseDefault, cwUseDefault, r.Right-r.Left, r.Bottom-r.Top,
			0, 0, getModuleHandle(), nil)
		pendingWindow = nil
		if err != nil {
			fyne.LogError("directx: creating window", err)
			return
		}
		w.hwnd = hwnd
		windowsByHandle[hwnd] = w
		if wakeTarget == 0 {
			wakeTarget = hwnd
		}
		w.setDarkMode()

		client := getClientRect(hwnd)
		g, err := dx.NewGPU(hwnd, uint32(client.Right), uint32(client.Bottom))
		if err != nil {
			fyne.LogError("directx: initialising Direct3D 11", err)
			return
		}
		w.gpu = g
		w.paint = dx.NewPainter(w.canvas, g)
		w.canvas.SetPainter(w.paint)
		w.paint.Init()

		w.canvas.scale = w.calculatedScale()
		w.resized(client.Right, client.Bottom)
	})

	d.addWindow(w)
	return w
}

func (d *dxDriver) AllWindows() []fyne.Window { return d.windows }

func (d *dxDriver) addWindow(w *window) {
	d.windows = append(d.windows, w)
}

func (*dxDriver) RenderedTextSize(text string, textSize float32, style fyne.TextStyle,
	source fyne.Resource) (fyne.Size, float32) {
	return paint.RenderedTextSize(text, textSize, style, source)
}

func (*dxDriver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	return common.CanvasForObject(obj)
}

func (d *dxDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return fyne.NewPos(0, 0)
	}
	dc, ok := c.(*dxCanvas)
	if !ok {
		return fyne.NewPos(0, 0)
	}
	return driver.AbsolutePositionForObject(co, dc.ObjectTrees())
}

func (*dxDriver) Device() fyne.Device { return &dxDevice{} }

// DoubleTapDelay reports the system double-click time so taps match the rest of
// the desktop rather than a hard coded constant.
func (*dxDriver) DoubleTapDelay() time.Duration {
	return time.Duration(doubleClickTime()) * time.Millisecond
}

// SetDisableScreenBlanking keeps the display awake while disable is true. The
// execution state is per-thread, so it is set on the main thread, which lives
// for the whole app.
func (*dxDriver) SetDisableScreenBlanking(disable bool) {
	async.EnsureMain(func() {
		state := uintptr(esContinuous)
		if disable {
			state |= esDisplayRequired
		}
		procSetThreadExecutionState.Call(state)
	})
}

// catchTerm shuts down through Quit on Ctrl+C or a terminate signal, so the
// lifecycle's OnStopped still runs rather than the process just vanishing.
func (d *dxDriver) catchTerm() {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	fyne.Do(d.Quit)
}

func (d *dxDriver) Quit() {
	if curWindow != nil {
		if f := fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle).OnExitedForeground(); f != nil {
			f()
		}
		curWindow = nil
	}
	if d.trayStop != nil {
		d.trayStop()
	}
	if running.CompareAndSwap(true, false) {
		close(d.done)
	}
}

func (d *dxDriver) Run() {
	if !async.IsMainGoroutine() {
		panic("Run() or ShowAndRun() must be called from main goroutine")
	}
	go d.catchTerm()
	d.runLoop()

	l := fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle)
	l.WaitForEvents()
	l.DestroyEventQueue()
}

// CurrentKeyModifiers reports the modifier keys held right now. Windows keeps
// modifier state in the keyboard rather than in each message, so this samples it
// on demand.
func (*dxDriver) CurrentKeyModifiers() fyne.KeyModifier {
	return currentModifiers()
}

// HasSecondaryDisplay reports whether more than one monitor is attached.
func (*dxDriver) HasSecondaryDisplay() bool {
	n, _, _ := procGetSystemMetrics.Call(smCMonitors)
	return n > 1
}

// The app and widget packages reach the driver's extended capabilities through
// type assertions, several of which are unchecked and so fail at runtime rather
// than at build time. These declarations turn a missing method back into a
// compile error.
//
// app.systrayDriver is unexported, so it is mirrored here by method set.
var (
	_ desktop.Driver = (*dxDriver)(nil)

	// mirrors the unexported app.systrayDriver
	_ interface {
		SetSystemTrayMenu(*fyne.Menu)
		SetSystemTrayIcon(fyne.Resource)
		SetSystemTrayWindow(fyne.Window)
	} = (*dxDriver)(nil)

	// mirrors the unexported fyne.systemTrayDriver used by fyne.Menu.Refresh
	_ interface {
		fyne.Driver
		SetSystemTrayMenu(*fyne.Menu)
		SystemTrayMenu() *fyne.Menu
	} = (*dxDriver)(nil)
)
