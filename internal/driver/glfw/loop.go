package glfw

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/scale"
)

type funcData struct {
	f    func()
	done chan struct{} // Zero allocation signalling channel
}

// channel for queuing functions on the main thread
var funcQueue = make(chan funcData)
var running atomic.Bool
var initOnce = &sync.Once{}

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
	mainGoroutineID = goroutineID()
}

// force a function f to run on the main thread
func runOnMain(f func()) {
	// If we are on main just execute - otherwise add it to the main queue and wait.
	// The "running" variable is normally false when we are on the main thread.
	if !running.Load() {
		f()
		return
	}

	done := common.DonePool.Get()
	defer common.DonePool.Put(done)

	funcQueue <- funcData{f: f, done: done}

	<-done
}

// force a function f to run on the draw thread
func runOnMainWithContext(w *window, f func()) {
	runOnMain(func() { w.RunWithContext(f) }) // TODO remove this completely
}

// Preallocate to avoid allocations on every drawSingleFrame.
// Note that the capacity of this slice can only grow,
// but its length will never be longer than the total number of
// window canvases that are dirty on a single frame.
// So its memory impact should be negligible and does not
// need periodic shrinking.
var refreshingCanvases []fyne.Canvas

func (d *gLDriver) drawSingleFrame() {
	for _, win := range d.windowList() {
		w := win.(*window)
		canvas := w.canvas
		closing := w.closing
		visible := w.visible

		// CheckDirtyAndClear must be checked after visibility,
		// because when a window becomes visible, it could be
		// showing old content without a dirty flag set to true.
		// Do the clear if and only if the window is visible.
		if closing || !visible || !canvas.CheckDirtyAndClear() {
			continue
		}

		d.repaintWindow(w)
		refreshingCanvases = append(refreshingCanvases, canvas)
	}
	cache.CleanCanvases(refreshingCanvases)

	// cleanup refreshingCanvases slice
	for i := 0; i < len(refreshingCanvases); i++ {
		refreshingCanvases[i] = nil
	}
	refreshingCanvases = refreshingCanvases[:0]
}

func (d *gLDriver) runGL() {
	if !running.CompareAndSwap(false, true) {
		return // Run was called twice.
	}
	close(d.waitForStart) // Signal that execution can continue.

	d.initGLFW()
	if d.trayStart != nil {
		d.trayStart()
	}
	if f := fyne.CurrentApp().Lifecycle().(*app.Lifecycle).OnStarted(); f != nil {
		go f() // don't block main, we don't have window event queue
	}

	settingsChange := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(settingsChange)

	eventTick := time.NewTicker(time.Second / 60)
	for {
		select {
		case <-d.done:
			eventTick.Stop()
			d.Terminate()
			l := fyne.CurrentApp().Lifecycle().(*app.Lifecycle)
			if f := l.OnStopped(); f != nil {
				l.QueueEvent(f)
			}
			return
		case f := <-funcQueue:
			f.f()
			f.done <- struct{}{}
		case <-eventTick.C:
			d.pollEvents()
			windowsToRemove := 0
			for _, win := range d.windowList() {
				w := win.(*window)
				if w.viewport == nil {
					continue
				}

				if w.viewport.ShouldClose() {
					windowsToRemove++
					continue
				}

				expand := w.shouldExpand
				fullScreen := w.fullScreen

				if expand && !fullScreen {
					w.fitContent()
					shouldExpand := w.shouldExpand
					w.shouldExpand = false
					view := w.viewport

					if shouldExpand && runtime.GOOS != "js" {
						view.SetSize(w.shouldWidth, w.shouldHeight)
					}
				}

				d.animation.TickAnimations()
				d.drawSingleFrame()
			}
			if windowsToRemove > 0 {
				oldWindows := d.windowList()
				newWindows := make([]fyne.Window, 0, len(oldWindows)-windowsToRemove)

				for _, win := range oldWindows {
					w := win.(*window)
					if w.viewport == nil {
						continue
					}

					if w.viewport.ShouldClose() {
						w.visible = false
						v := w.viewport

						// remove window from window list
						v.Destroy()
						w.destroy(d)
						continue
					}

					newWindows = append(newWindows, win)
				}

				d.windowLock.Lock()
				d.windows = newWindows
				d.windowLock.Unlock()

				if len(newWindows) == 0 {
					d.Quit()
				}
			}
		case set := <-settingsChange:
			painter.ClearFontCache()
			cache.ResetThemeCaches()
			app.ApplySettingsWithCallback(set, fyne.CurrentApp(), func(w fyne.Window) {
				c, ok := w.Canvas().(*glCanvas)
				if !ok {
					return
				}
				c.applyThemeOutOfTreeObjects()
				go c.reloadScale()
			})

		}
	}
}

func (d *gLDriver) repaintWindow(w *window) {
	canvas := w.canvas
	w.RunWithContext(func() {
		if canvas.EnsureMinSize() {
			w.shouldExpand = true
		}
		canvas.FreeDirtyTextures()

		updateGLContext(w)
		canvas.paint(canvas.Size())

		view := w.viewport
		visible := w.visible

		if view != nil && visible {
			view.SwapBuffers()
		}
	})
}

// refreshWindow requests that the specified window be redrawn
func refreshWindow(w *window) {
	w.canvas.SetDirty()
}

func updateGLContext(w *window) {
	canvas := w.canvas
	size := canvas.Size()

	// w.width and w.height are not correct if we are maximised, so figure from canvas
	winWidth := float32(scale.ToScreenCoordinate(canvas, size.Width)) * canvas.texScale
	winHeight := float32(scale.ToScreenCoordinate(canvas, size.Height)) * canvas.texScale

	canvas.Painter().SetFrameBufferScale(canvas.texScale)
	w.canvas.Painter().SetOutputSize(int(winWidth), int(winHeight))
}
