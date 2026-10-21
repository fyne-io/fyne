//go:build !darwin

package glfw

import (
	"time"

	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/build"
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
		if build.IsWayland && w.width != 0 && w.height != 0 {
			// At its full screen size GLFW commits no smaller window before the
			// compositor's configure brings back the size it had. A compositor
			// that sends no size leaves it full screen sized.
			width, height := w.viewport.GetSize()
			w.viewport.SetMonitor(nil, 0, 0, width, height, 0)
			return
		}
		if w.width == 0 && w.height == 0 { // if we were fullscreen on creation...
			s := internal.MaxSizes(w.canvas.Size(), w.canvas.MinSize())
			w.width, w.height = w.screenSize(s)
		}
		w.viewport.SetMonitor(nil, w.xpos, w.ypos, w.width, w.height, 0)
		w.processResized(w.width, w.height) // GLFW on Wayland sends no size event for a size it sets
	}
}
