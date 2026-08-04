//go:build windows

// The run loop and the main-thread work queue.
//
// Win32 requires that a window is only touched from the thread that created
// it, so everything funnels through here: the loop pumps the message queue,
// runs queued functions, ticks animations and repaints.

package directx

import (
	"runtime"
	"sync/atomic"
	"time"

	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2"
	intapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/driver/common"
	paint "fyne.io/fyne/v2/internal/painter"
)

// Arrange that main.main runs on the main thread, as Win32 requires all window
// and message calls for a window to happen on its creating thread.
func init() {
	runtime.LockOSThread()
	async.SetMainGoroutine()
}

type funcData struct {
	f    func()
	done chan struct{}
}

var (
	funcQueue        = async.NewUnboundedChan[funcData]()
	running, drained atomic.Bool
	// wakeTarget is any live window handle, used to nudge the message loop when
	// work is queued from another goroutine.
	wakeTarget windows.Handle
)

// Arrange that main.main runs on the main thread, as Win32 requires all window
// and message calls for a window to happen on its creating thread.

func (d *dxDriver) init() {
	if d.initialized {
		return
	}
	d.initialized = true
	enableDpiAwareness()
}

func runOnMain(f func()) {
	runOnMainWithWait(f, true)
}

func runOnMainWithWait(f func(), wait bool) {
	if (!running.Load() && async.IsMainGoroutine()) || drained.Load() {
		f()
		return
	}

	if wait {
		done := common.DonePool.Get()
		defer common.DonePool.Put(done)

		funcQueue.In() <- funcData{f: f, done: done}
		wake()
		<-done
		return
	}
	funcQueue.In() <- funcData{f: f}
	wake()
}

// wake posts a no-op message so a message loop blocked in GetMessage returns.
func wake() {
	if wakeTarget != 0 {
		postMessage(wakeTarget, wmFyneDo, 0, 0)
	}
}

func (d *dxDriver) runLoop() {
	if !running.CompareAndSwap(false, true) {
		return // Run was called twice
	}
	d.init()
	if d.trayStart != nil {
		d.trayStart()
	}

	fyne.CurrentApp().Settings().AddListener(func(set fyne.Settings) {
		paint.ClearFontCache()
		cache.ResetThemeCaches()
		intapp.ApplySettingsWithCallback(set, fyne.CurrentApp(), func(w fyne.Window) {
			if win, ok := w.(*window); ok {
				win.setDarkMode()
			}
			c, ok := w.Canvas().(*dxCanvas)
			if !ok {
				return
			}
			c.applyThemeOutOfTreeObjects()
			c.reloadScale()
		})
	})

	if f := fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle).OnStarted(); f != nil {
		f()
	}

	eventTick := time.NewTicker(time.Second / 60)
	for {
		select {
		case <-d.done:
			eventTick.Stop()
			l := fyne.CurrentApp().Lifecycle().(*intapp.Lifecycle)
			if f := l.OnStopped(); f != nil {
				l.QueueEvent(f)
			}
			for len(funcQueue.Out()) > 0 {
				f := <-funcQueue.Out()
				if f.done != nil {
					f.done <- struct{}{}
				}
			}
			drained.Store(true)
			funcQueue.Close()
			return

		case f := <-funcQueue.Out():
			f.f()
			if f.done != nil {
				f.done <- struct{}{}
			}

		case <-eventTick.C:
			d.pollEvents()
			for i := 0; i < len(d.windows); i++ {
				w := d.windows[i].(*window)
				if w.closing {
					d.destroyWindow(w, i)
					i--
				}
			}
			d.animation.TickAnimations()
			d.drawSingleFrame()
		}
	}
}

// pollEvents drains the Win32 message queue without blocking, so the ticker
// keeps driving repaints and the func queue stays responsive.
func (*dxDriver) pollEvents() {
	var m msg
	for peekMessage(&m) {
		if m.message == wmQuit {
			return
		}
		translateDispatch(&m)
	}
}

func (d *dxDriver) destroyWindow(w *window, index int) {
	w.visible = false
	w.destroy()

	if index < len(d.windows)-1 {
		copy(d.windows[index:], d.windows[index+1:])
	}
	d.windows[len(d.windows)-1] = nil
	d.windows = d.windows[:len(d.windows)-1]

	// An app with a system tray menu outlives its windows - quitting is then the
	// tray's Quit item, which addMissingQuitForMenu guarantees exists.
	if w.master || (len(d.windows) == 0 && d.systrayMenu == nil) {
		d.Quit()
	}
}

func (d *dxDriver) drawSingleFrame() {
	refreshed := false
	for _, win := range d.windows {
		w := win.(*window)
		if w.closing || !w.visible || w.gpu == nil {
			continue
		}
		if !w.canvas.CheckDirtyAndClear() {
			w.markCacheAlive()
			continue
		}
		if d.repaintWindow(w) {
			refreshed = true
		}
	}
	cache.Clean(refreshed)
}

func (d *dxDriver) repaintWindow(w *window) bool {
	if w.painting {
		return false
	}
	w.painting = true
	defer func() { w.painting = false }()

	// Last chance to catch a client area that moved without a WM_SIZE reaching us.
	w.syncSurface()

	canvas := w.canvas
	if canvas.EnsureMinSize() {
		w.fitContent()
	}
	freed := canvas.FreeDirtyTextures() > 0

	// Drive the viewport from the swap chain rather than from the canvas size: the
	// back buffer is what is actually being drawn into, and mid-resize the canvas
	// has not necessarily caught up with it.
	fbWidth, fbHeight := w.gpu.Size()
	canvas.Painter().SetFrameBufferScale(canvas.texScale)
	canvas.Painter().SetOutputSize(int(fbWidth), int(fbHeight))

	canvas.paint(canvas.Size())
	if w.gpu.Present() {
		w.recoverDevice()
	}

	w.lastWalkedTime = time.Now()
	return freed
}

func (w *window) markCacheAlive() {
	threshold := time.Now().Add(10*time.Second - cache.ValidDuration)
	if w.lastWalkedTime.Before(threshold) {
		w.canvas.WalkTrees(nil, func(node *common.RenderCacheNode, _ fyne.Position) {
			_ = cache.GetCanvasForObject(node.Obj())
			if wid, ok := node.Obj().(fyne.Widget); ok {
				_, _ = cache.CachedRenderer(wid)
			}
		})
		w.lastWalkedTime = time.Now()
	}
}
