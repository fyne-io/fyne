package glfw

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/scale"
)

// for queuing functions on the main thread
var funcQueue = newGlfwFuncQueue()
var running atomic.Bool
var initOnce = &sync.Once{}

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
	async.SetMainGoroutine()
}

// force a function f to run on the main thread
func runOnMain(f func()) {
	runOnMainWithWait(f, true)
}

// force a function f to run on the main thread and specify if we should wait for it to return
func runOnMainWithWait(f func(), wait bool) {
	// If we are on main just execute - otherwise add it to the main queue and wait.
	// The "running" variable is normally false when we are on the main thread.
	if !running.Load() {
		f()
		return
	}

	funcQueue.Push(f, wait)
}

// Preallocate to avoid allocations on every drawSingleFrame.
// Note that the capacity of this slice can only grow,
// but its length will never be longer than the total number of
// window canvases that are dirty on a single frame.
// So its memory impact should be negligible and does not
// need periodic shrinking.
var refreshingCanvases []fyne.Canvas

func (d *gLDriver) drawSingleFrame() {
	refreshed := false
	for _, win := range d.windowList() {
		w := win.(*window)
		if w.closing {
			continue
		}

		canvas := w.canvas

		// CheckDirtyAndClear must be checked after visibility,
		// because when a window becomes visible, it could be
		// showing old content without a dirty flag set to true.
		// Do the clear if and only if the window is visible.
		if !w.visible || !canvas.CheckDirtyAndClear() {
			// Window hidden or not being redrawn, mark canvasForObject
			// cache alive if it hasn't been done recently
			// n.b. we need to make sure threshold is a bit *after*
			// time.Now() - CacheDuration()
			threshold := time.Now().Add(time.Second - cache.ValidDuration)
			if w.lastWalkedTime.Before(threshold) {
				w.canvas.WalkTrees(nil, func(node *common.RenderCacheNode, _ fyne.Position) {
					// marks canvas for widget cache entry alive
					_ = cache.GetCanvasForObject(node.Obj())
				})
				w.lastWalkedTime = time.Now()
			}
			continue
		}

		refreshed = refreshed || d.repaintWindow(w)
		refreshingCanvases = append(refreshingCanvases, canvas)
	}
	cache.CleanCanvases(refreshingCanvases)
	cache.Clean(refreshed)

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

	d.initGLFW()
	if d.trayStart != nil {
		d.trayStart()
	}

	fyne.CurrentApp().Settings().AddListener(func(set fyne.Settings) {
		painter.ClearFontCache()
		cache.ResetThemeCaches()
		app.ApplySettingsWithCallback(set, fyne.CurrentApp(), func(w fyne.Window) {
			c, ok := w.Canvas().(*glCanvas)
			if !ok {
				return
			}
			c.applyThemeOutOfTreeObjects()
			c.reloadScale()
		})
	})

	if f := fyne.CurrentApp().Lifecycle().(*app.Lifecycle).OnStarted(); f != nil {
		f()
	}

	for {
		// idle until events (or d.wakeUp()) received
		d.waitEvents()

		// run active loop until we are exiting, or idle again
		if exit := d.runActiveLoop(); exit {
			return
		}
	}
}

func (d *gLDriver) destroyWindow(w *window, index int) {
	w.visible = false
	w.viewport.Destroy()
	w.destroy(d)

	if index < len(d.windows)-1 {
		copy(d.windows[index:], d.windows[index+1:])
	}
	d.windows[len(d.windows)-1] = nil
	d.windows = d.windows[:len(d.windows)-1]

	if len(d.windows) == 0 {
		d.Quit()
	}
}

// runs a loop that invokes runSingleFrame at 60 fps until there are no more
// animations or events. uses pollEvents to check for events every frame.
func (d *gLDriver) runActiveLoop() bool {
	t := time.NewTicker(time.Second / 60)
	defer t.Stop()
	for range t.C {
		d.pollEvents()
		exit, idle := d.runSingleFrame()
		if exit {
			return exit
		} else if idle {
			break
		}
	}
	return false
}

func (d *gLDriver) runSingleFrame() (exit, idle bool) {
	// check if we're shutting down
	if !running.Load() {
		d.Terminate()
		l := fyne.CurrentApp().Lifecycle().(*app.Lifecycle)
		if f := l.OnStopped(); f != nil {
			l.QueueEvent(f)
		}
		return true, false
	}

	// process received glfw events
	noEvents := true
	for {
		evt, ok := eventQueue.Pull()
		if !ok {
			break
		}
		noEvents = false
		processEvent(evt)
	}

	// run funcs queued to main
	var buf [16]funcData
	funcsDone := false
	for !funcsDone {
		n := funcQueue.PullN(buf[:])
		for i := 0; i < n; i++ {
			buf[i].f()
			if buf[i].done != nil {
				buf[i].done <- struct{}{}
			}
			buf[i] = funcData{}
		}
		funcsDone = n == 0
	}

	for i := 0; i < len(d.windows); i++ {
		w := d.windows[i].(*window)
		if w.viewport == nil {
			continue
		}

		if w.viewport.ShouldClose() {
			d.destroyWindow(w, i)
			i-- // Trailing windows are moved forward one position.
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
	}
	animationsDone := d.animation.TickAnimations()
	d.drawSingleFrame()

	return false, animationsDone && noEvents
}

func (d *gLDriver) repaintWindow(w *window) bool {
	canvas := w.canvas
	freed := false
	w.RunWithContext(func() {
		if canvas.EnsureMinSize() {
			w.shouldExpand = true
		}
		freed = canvas.FreeDirtyTextures() > 0

		updateGLContext(w)
		canvas.paint(canvas.Size())

		view := w.viewport
		visible := w.visible

		if view != nil && visible {
			view.SwapBuffers()
		}

		// mark that we have walked the window and don't
		// need to walk it again to mark caches alive
		w.lastWalkedTime = time.Now()
	})
	return freed
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
	canvas.Painter().SetOutputSize(int(winWidth), int(winHeight))
}
