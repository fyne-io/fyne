package glfw

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/internal"
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
		d.startDrawThread()
	})
}

func (d *gLDriver) runGL() {
	eventTick := time.NewTicker(time.Second / 60)
	runMutex.Lock()
	runFlag = true
	runMutex.Unlock()

	d.initGLFW()

	for {
		select {
		case <-d.done:
			eventTick.Stop()
			d.drawDone <- nil // wait for draw thread to stop
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
					w.viewLock.Lock()
					w.visible = false
					v := w.viewport
					w.viewport = nil
					w.viewLock.Unlock()

					// remove window from window list
					v.Destroy()
					go w.destroy(d)
					continue
				}

				w.viewLock.RLock()
				expand := w.shouldExpand
				w.viewLock.RUnlock()

				if expand {
					w.fitContent()
					w.viewLock.Lock()
					w.shouldExpand = false
					view := w.viewport
					w.viewLock.Unlock()
					view.SetSize(w.width, w.height)
				}

				newWindows = append(newWindows, win)
			}
			if reassign {
				d.windowLock.Lock()
				d.windows = newWindows
				d.windowLock.Unlock()

				if len(newWindows) == 0 {
					d.Quit()
				}
			}
		}
	}
}

func (d *gLDriver) repaintWindow(w *window) {
	canvas := w.canvas
	w.RunWithContext(func() {
		if w.canvas.ensureMinSize() {
			w.viewLock.Lock()
			w.shouldExpand = true
			w.viewLock.Unlock()
		}
		freeDirtyTextures(canvas)

		updateGLContext(w)
		canvas.paint(canvas.Size())

		w.viewLock.RLock()
		view := w.viewport
		visible := w.visible
		w.viewLock.RUnlock()

		if view != nil && visible {
			view.SwapBuffers()
		}
	})
}

func (d *gLDriver) startDrawThread() {
	settingsChange := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(settingsChange)
	draw := time.NewTicker(time.Second / 60)

	go func() {
		runtime.LockOSThread()

		for {
			select {
			case <-d.drawDone:
				return
			case f := <-drawFuncQueue:
				f.win.RunWithContext(f.f)
				if f.done != nil {
					f.done <- true
				}
			case <-settingsChange:
				painter.ClearFontCache()
			case <-draw.C:
				for _, win := range d.windowList() {
					w := win.(*window)
					w.viewLock.RLock()
					canvas := w.canvas
					view := w.viewport
					visible := w.visible
					w.viewLock.RUnlock()
					if view == nil || !canvas.isDirty() || !visible {
						continue
					}

					d.repaintWindow(w)
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
	w.canvas.setDirty(true)
}

func updateGLContext(w *window) {
	canvas := w.Canvas().(*glCanvas)
	size := canvas.Size()

	// w.width and w.height are not correct if we are maximised, so figure from canvas
	winWidth := float32(internal.ScaleInt(canvas, size.Width)) * canvas.texScale
	winHeight := float32(internal.ScaleInt(canvas, size.Height)) * canvas.texScale

	canvas.painter.SetFrameBufferScale(canvas.texScale)
	w.canvas.painter.SetOutputSize(int(winWidth), int(winHeight))
}
