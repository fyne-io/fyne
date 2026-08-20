//go:build !darwin

package glfw

import (
	"time"

	"fyne.io/fyne/v2/internal"
)

const desktopDefaultDoubleTapDelay = 300 * time.Millisecond

func (w *window) doSetFullScreen(full bool) {
	monitor := w.getMonitorForWindow()

	w.doApplyFullScreen(monitor, full)
}

func (w *window) doSetFullScreen2(full bool) {
	monitor := w.getSecondaryMonitor()

	w.doApplyFullScreen(monitor, full)
}

func (w *window) doApplyFullScreen(monitor *monitor, full bool) {
	if full {
		mode := monitor.GetVideoMode()
		if mode == nil { // monitor was disconnected
			return
		}
		w.viewport.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		if w.width == 0 && w.height == 0 { // if we were fullscreen on creation...
			s := internal.MaxSizes(w.canvas.Size(), w.canvas.MinSize())
			w.width, w.height = w.screenSize(s)
		}
		w.viewport.SetMonitor(nil, w.xpos, w.ypos, w.width, w.height, 0)
	}
}
