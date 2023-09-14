package glfw

import (
	"runtime"
	"sync"
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

type drawData struct {
	f    func()
	win  *window
	done chan struct{} // Zero allocation signalling channel
}

type runFlag struct {
	sync.Cond
	flag bool
}

// channel for queuing functions on the main thread
var funcQueue = make(chan funcData)
var drawFuncQueue = make(chan drawData)
var run = &runFlag{Cond: sync.Cond{L: &sync.Mutex{}}}
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
	run.L.Lock()
	running := !run.flag
	run.L.Unlock()

	if running {
		f()
	} else {
		done := common.DonePool.Get().(chan struct{})
		defer common.DonePool.Put(done)

		funcQueue <- funcData{f: f, done: done}

		<-done
	}
}

// force a function f to run on the draw thread
func runOnDraw(w *window, f func()) {
	if drawOnMainThread {
		runOnMain(func() { w.RunWithContext(f) })
		return
	}
	done := common.DonePool.Get().(chan struct{})
	defer common.DonePool.Put(done)

	drawFuncQueue <- drawData{f: f, win: w, done: done}
	<-done
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
		w.viewLock.RLock()
		canvas := w.canvas
		closing := w.closing
		visible := w.visible
		w.viewLock.RUnlock()

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
	eventTick := time.NewTicker(time.Second / 60)

	run.L.Lock()
	run.flag = true
	run.L.Unlock()
	run.Broadcast()

	d.initGLFW()
	if d.trayStart != nil {
		d.trayStart()
	}
	fyne.CurrentApp().Lifecycle().(*app.Lifecycle).TriggerStarted()
	for {
		select {
		case <-d.done:
			eventTick.Stop()
			d.drawDone <- nil // wait for draw thread to stop
			d.Terminate()
			fyne.CurrentApp().Lifecycle().(*app.Lifecycle).TriggerStopped()
			return
		case f := <-funcQueue:
			f.f()
			if f.done != nil {
				f.done <- struct{}{}
			}
		case <-eventTick.C:
			d.tryPollEvents()
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

				w.viewLock.RLock()
				expand := w.shouldExpand
				fullScreen := w.fullScreen
				w.viewLock.RUnlock()

				if expand && !fullScreen {
					w.fitContent()
					w.viewLock.Lock()
					shouldExpand := w.shouldExpand
					w.shouldExpand = false
					view := w.viewport
					w.viewLock.Unlock()
					if shouldExpand {
						view.SetSize(w.shouldWidth, w.shouldHeight)
					}
				}

				if drawOnMainThread {
					d.drawSingleFrame()
				}
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
						w.viewLock.Lock()
						w.visible = false
						v := w.viewport
						w.viewLock.Unlock()

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
		}
	}
}

func (d *gLDriver) repaintWindow(w *window) {
	canvas := w.canvas
	w.RunWithContext(func() {
		if canvas.EnsureMinSize() {
			w.viewLock.Lock()
			w.shouldExpand = true
			w.viewLock.Unlock()
		}
		canvas.FreeDirtyTextures()

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
	var drawCh <-chan time.Time
	if drawOnMainThread {
		drawCh = make(chan time.Time) // don't tick when on M1
	} else {
		drawCh = time.NewTicker(time.Second / 60).C
	}

	go func() {
		runtime.LockOSThread()

		for {
			select {
			case <-d.drawDone:
				return
			case f := <-drawFuncQueue:
				f.win.RunWithContext(f.f)
				if f.done != nil {
					f.done <- struct{}{}
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
			case <-drawCh:
				d.drawSingleFrame()
			}
		}
	}()
}

// refreshWindow requests that the specified window be redrawn
func refreshWindow(w *window) {
	w.canvas.SetDirty()
}

func updateGLContext(w *window) {
	canvas := w.Canvas().(*glCanvas)
	size := canvas.Size()

	// w.width and w.height are not correct if we are maximised, so figure from canvas
	winWidth := float32(scale.ToScreenCoordinate(canvas, size.Width)) * canvas.texScale
	winHeight := float32(scale.ToScreenCoordinate(canvas, size.Height)) * canvas.texScale

	canvas.Painter().SetFrameBufferScale(canvas.texScale)
	w.canvas.Painter().SetOutputSize(int(winWidth), int(winHeight))
}
