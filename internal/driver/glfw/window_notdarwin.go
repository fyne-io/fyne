//go:build !darwin

package glfw

func (w *window) doSetFullScreen(full bool) {
	monitor := w.getMonitorForWindow()
	mode := monitor.GetVideoMode()

	if full {
		w.viewport.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		if w.width == 0 && w.height == 0 { // if we were fullscreen on creation...
			s := w.canvas.Size().Max(w.canvas.MinSize())
			w.width, w.height = w.screenSize(s)
		}
		w.viewport.SetMonitor(nil, w.xpos, w.ypos, w.width, w.height, 0)
	}
}
