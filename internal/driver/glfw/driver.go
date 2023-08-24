// Package glfw provides a full Fyne desktop driver that uses the system OpenGL libraries.
// This supports Windows, Mac OS X and Linux using the gl and glfw packages from go-gl.
package glfw

import (
	"bytes"
	"image"
	"os"
	"runtime"
	"sync"

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

// mainGoroutineID stores the main goroutine ID.
// This ID must be initialized in main.init because
// a main goroutine may not equal to 1 due to the
// influence of a garbage collector.
var mainGoroutineID uint64

var curWindow *window

// Declare conformity with Driver
var _ fyne.Driver = (*gLDriver)(nil)

// A workaround on Apple M1/M2, just use 1 thread until fixed upstream.
const drawOnMainThread bool = runtime.GOOS == "darwin" && runtime.GOARCH == "arm64"

type gLDriver struct {
	windowLock sync.RWMutex
	windows    []fyne.Window
	device     *glDevice
	done       chan interface{}
	drawDone   chan interface{}

	animation *animation.Runner

	currentKeyModifiers fyne.KeyModifier // desktop driver only

	trayStart, trayStop func()     // shut down the system tray, if used
	systrayMenu         *fyne.Menu // cache the menu set so we know when to refresh
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

func (d *gLDriver) RenderedTextSize(text string, textSize float32, style fyne.TextStyle) (size fyne.Size, baseline float32) {
	return painter.RenderedTextSize(text, textSize, style)
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
	if d.device == nil {
		d.device = &glDevice{}
	}

	return d.device
}

func (d *gLDriver) Quit() {
	if curWindow != nil {
		curWindow = nil
		if d.trayStop != nil {
			d.trayStop()
		}
		fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle).TriggerExitedForeground()
	}
	defer func() {
		recover() // we could be called twice - no safe way to check if d.done is closed
	}()
	close(d.done)
}

func (d *gLDriver) addWindow(w *window) {
	d.windowLock.Lock()
	defer d.windowLock.Unlock()
	d.windows = append(d.windows, w)
}

// a trivial implementation of "focus previous" - return to the most recently opened, or master if set.
// This may not do the right thing if your app has 3 or more windows open, but it was agreed this was not much
// of an issue, and the added complexity to track focus was not needed at this time.
func (d *gLDriver) focusPreviousWindow() {
	d.windowLock.RLock()
	wins := d.windows
	d.windowLock.RUnlock()

	var chosen fyne.Window
	for _, w := range wins {
		if !w.(*window).visible {
			continue
		}
		chosen = w
		if w.(*window).master {
			break
		}
	}

	if chosen == nil || chosen.(*window).view() == nil {
		return
	}
	chosen.RequestFocus()
}

func (d *gLDriver) windowList() []fyne.Window {
	d.windowLock.RLock()
	defer d.windowLock.RUnlock()
	return d.windows
}

func (d *gLDriver) initFailed(msg string, err error) {
	logError(msg, err)

	run.L.Lock()
	running := !run.flag
	run.L.Unlock()

	if running {
		d.Quit()
	} else {
		os.Exit(1)
	}
}

func (d *gLDriver) Run() {
	if goroutineID() != mainGoroutineID {
		panic("Run() or ShowAndRun() must be called from main goroutine")
	}

	go d.catchTerm()
	d.runGL()
}

// NewGLDriver sets up a new Driver instance implemented using the GLFW Go library and OpenGL bindings.
func NewGLDriver() fyne.Driver {
	repository.Register("file", intRepo.NewFileRepository())

	return &gLDriver{
		done:      make(chan interface{}),
		drawDone:  make(chan interface{}),
		animation: &animation.Runner{},
	}
}
