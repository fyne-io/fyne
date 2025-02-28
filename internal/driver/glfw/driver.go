// Package glfw provides a full Fyne desktop driver that uses the system OpenGL libraries.
// This supports Windows, Mac OS X and Linux using the gl and glfw packages from go-gl.
package glfw

import (
	"bytes"
	"image"
	"os"
	"runtime"

	"fyne.io/fyne/v2/internal/async"
	"github.com/fyne-io/image/ico"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/animation"
	intapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/driver"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage/repository"
)

var curWindow *window

// Declare conformity with Driver
var _ fyne.Driver = (*gLDriver)(nil)

type gLDriver struct {
	windows     []fyne.Window
	initialized bool
	done        chan struct{}

	animation animation.Runner

	currentKeyModifiers fyne.KeyModifier // desktop driver only

	trayStart, trayStop func()     // shut down the system tray, if used
	systrayMenu         *fyne.Menu // cache the menu set so we know when to refresh
}

func (d *gLDriver) init() {
	if !d.initialized {
		d.initialized = true
		d.initGLFW()
	}
}

func toOSIcon(icon []byte) ([]byte, error) {
	if runtime.GOOS != "windows" {
		return icon, nil
	}

	img, _, err := image.Decode(bytes.NewReader(icon))
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	err = ico.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *gLDriver) DoFromGoroutine(f func(), wait bool) {
	if wait {
		async.EnsureNotMain(func() {
			runOnMainWithWait(f, true)
		})
	} else {
		runOnMainWithWait(f, false)
	}
}

func (d *gLDriver) RenderedTextSize(text string, textSize float32, style fyne.TextStyle, source fyne.Resource) (size fyne.Size, baseline float32) {
	return painter.RenderedTextSize(text, textSize, style, source)
}

func (d *gLDriver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	return common.CanvasForObject(obj)
}

func (d *gLDriver) AbsolutePositionForObject(co fyne.CanvasObject) fyne.Position {
	c := d.CanvasForObject(co)
	if c == nil {
		return fyne.NewPos(0, 0)
	}

	glc := c.(*glCanvas)
	return driver.AbsolutePositionForObject(co, glc.ObjectTrees())
}

func (d *gLDriver) Device() fyne.Device {
	return &glDevice{}
}

func (d *gLDriver) Quit() {
	if curWindow != nil {
		if f := fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle).OnExitedForeground(); f != nil {
			f()
		}
		curWindow = nil
		if d.trayStop != nil {
			d.trayStop()
		}
	}

	// Only call close once to avoid panic.
	if running.CompareAndSwap(true, false) {
		close(d.done)
	}
}

func (d *gLDriver) addWindow(w *window) {
	d.windows = append(d.windows, w)
}

// a trivial implementation of "focus previous" - return to the most recently opened, or master if set.
// This may not do the right thing if your app has 3 or more windows open, but it was agreed this was not much
// of an issue, and the added complexity to track focus was not needed at this time.
func (d *gLDriver) focusPreviousWindow() {
	var chosen *window
	for _, w := range d.windows {
		win := w.(*window)
		if !win.visible {
			continue
		}
		chosen = win
		if win.master {
			break
		}
	}

	if chosen == nil || chosen.view() == nil {
		return
	}
	chosen.RequestFocus()
}

func (d *gLDriver) windowList() []fyne.Window {
	return d.windows
}

func (d *gLDriver) initFailed(msg string, err error) {
	logError(msg, err)

	if !running.Load() {
		d.Quit()
	} else {
		os.Exit(1)
	}
}

func (d *gLDriver) Run() {
	if !async.IsMainGoroutine() {
		panic("Run() or ShowAndRun() must be called from main goroutine")
	}

	go d.catchTerm()
	d.runGL()

	// Ensure lifecycle events run to completion before the app exits
	l := fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle)
	l.WaitForEvents()
	l.DestroyEventQueue()
}

func (d *gLDriver) SetDisableScreenBlanking(disable bool) {
	setDisableScreenBlank(disable)
}

// NewGLDriver sets up a new Driver instance implemented using the GLFW Go library and OpenGL bindings.
func NewGLDriver() *gLDriver {
	repository.Register("file", intRepo.NewFileRepository())

	return &gLDriver{
		done: make(chan struct{}),
	}
}
