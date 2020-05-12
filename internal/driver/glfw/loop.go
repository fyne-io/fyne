package glfw

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/internal/painter"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type funcData struct {
	f    func()
	done chan bool
}

type drawData struct {
	f    func()
	win  *window
	done chan bool
}

// channel for queuing functions on the main thread
var funcQueue = make(chan funcData)
var drawFuncQueue = make(chan drawData)
var windowQueue = make(chan *window, 16)
var runFlag = false
var runMutex = &sync.Mutex{}
var initOnce = &sync.Once{}

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

func running() bool {
	runMutex.Lock()
	defer runMutex.Unlock()
	return runFlag
}

// force a function f to run on the main thread
func runOnMain(f func()) {
	// If we are on main just execute - otherwise add it to the main queue and wait.
	// The "running" variable is normally false when we are on the main thread.
	if !running() {
		f()
	} else {
		done := make(chan bool)

		funcQueue <- funcData{f: f, done: done}
		<-done
	}
}

// force a function f to run on the draw thread
func runOnDraw(w *window, f func()) {
	done := make(chan bool)

	drawFuncQueue <- drawData{f: f, win: w, done: done}
	<-done
}

func (d *gLDriver) initGLFW() {
	initOnce.Do(func() {
		err := glfw.Init()
		if err != nil {
			fyne.LogError("failed to initialise GLFW", err)
			return
		}

		initCursors()
	})
}

func (d *gLDriver) runGL() {
	eventTick := time.NewTicker(time.Second / 10)
	runMutex.Lock()
	runFlag = true
	runMutex.Unlock()

	d.initGLFW()
	d.startDrawThread()
	d.startRedrawTimer()

	for {
		select {
		case <-d.done:
			eventTick.Stop()
			glfw.Terminate()
			return
		case f := <-funcQueue:
			f.f()
			if f.done != nil {
				f.done <- true
			}
		case <-eventTick.C:
			d.tryPollEvents()
			newWindows := []fyne.Window{}
			reassign := false
			for _, win := range d.windowList() {
				w := win.(*window)
				if w.viewport == nil {
					continue
				}

				if w.viewport.ShouldClose() {
					reassign = true
					// remove window from window list
					w.viewport.Destroy()

					go w.destroy(d)
					continue
				}
				newWindows = append(newWindows, win)
			}
			if reassign {
				d.windowLock.Lock()
				d.windows = newWindows
				d.windowLock.Unlock()
			}
		}
	}
}

func (d *gLDriver) repaintWindow(w *window) {
	canvas := w.canvas
	w.RunWithContext(func() {
		freeDirtyTextures(canvas)

		updateGLContext(w)
		if canvas.ensureMinSize() {
			// TODO we can no longer run fitContent on this thread as it can impact viewport so must be on main
			// TODO but if we run it on main from here we will block the redraw queue on resize ... so need another signal :(
			//runOnMain(func() {
			//	w.fitContent()
			//})
		}
		canvas.paint(canvas.Size())

		w.viewport.SwapBuffers()
	})
}

func (d *gLDriver) startDrawThread() {
	settingsChange := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(settingsChange)

	go func() {
		runtime.LockOSThread()

		for {
			select {
			case f := <-drawFuncQueue:
				f.win.RunWithContext(f.f)
				if f.done != nil {
					f.done <- true
				}
			case <-settingsChange:
				painter.ClearFontCache()
			case w := <-windowQueue:
				d.repaintWindow(w)
			}
		}
	}()
}

func (d *gLDriver) startRedrawTimer() {
	draw := time.NewTicker(time.Second / 60)
	go func() {
		for {
			select {
			case <-draw.C:
				for _, win := range d.windowList() {
					w := win.(*window)
					canvas := w.canvas
					if w.viewport == nil || !canvas.isDirty() || !w.visible {
						continue
					}

					windowQueue <- w
				}
			}
		}
	}()
}

func (d *gLDriver) tryPollEvents() {
	defer func() {
		if r := recover(); r != nil {
			fyne.LogError(fmt.Sprint("GLFW poll event error: ", r), nil)
		}
	}()

	glfw.PollEvents() // This call blocks while window is being resized, which prevents freeDirtyTextures from being called
}

func freeDirtyTextures(canvas *glCanvas) {
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

// refreshWindow requests that the specified window be redrawn
func refreshWindow(w *window) {
	windowQueue <- w
}
