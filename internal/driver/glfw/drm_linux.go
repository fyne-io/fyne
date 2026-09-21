//go:build linux && !wasm && !test_web_driver

package glfw

import (
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/internal/cache"

	"github.com/go-gl/glfw/v3.4/glfw"
)

// guardDRM runs on the GLFW thread before PollEvents.
//
// KVM often leaves status=connected and only flips enabled=disabled, so both
// are watched. The crash is glfw.PollEvents dispatching wl_registry.global_remove.
// On loss we Terminate immediately and do not poll until the output is back;
// flushing those events later SIGSEGVs even with the window hidden.
func (d *gLDriver) guardDRM() {
	cur := drmSnapshot()
	if !drmDown && drmLostOutput(drmPrev, cur) {
		drmLost = cloneDRM(drmPrev)
		d.dropDisplay()
		drmDown = true
	}
	if drmDown && !drmStillMissing(drmLost, cur) {
		if d.restoreDisplay() {
			drmDown = false
			drmPrev = cur
		}
	}
	if !drmDown {
		drmPrev = cur
	}
}

func (d *gLDriver) dropDisplay() {
	for _, win := range d.AllWindows() {
		w, ok := win.(*window)
		if !ok {
			continue
		}
		w.visible = false
		w.viewport = nil
		w.created = false
		if w.frame != nil {
			w.frame.resetAfterDisplayLoss()
		}
		if w.canvas != nil {
			cache.DropTexturesFor(w.canvas)
		}
	}
	glfw.Terminate()
}

func (d *gLDriver) restoreDisplay() bool {
	d.rebindActive = true
	defer func() { d.rebindActive = false }()

	if err := d.initGLFW(); err != nil {
		return false
	}
	for _, win := range d.AllWindows() {
		w, ok := win.(*window)
		if !ok || w.closing {
			continue
		}
		w.remapAfterHotplug()
	}
	d.rebindPaint = true
	return true
}

func (w *window) remapAfterHotplug() {
	w.create()
	if w.viewport == nil {
		return
	}
	w.created = true
	if w.frame != nil {
		w.frame.resetAfterDisplayLoss()
	}
	if w.canvas != nil {
		w.canvas.SetDirty()
		if content := w.canvas.Content(); content != nil {
			content.Refresh()
		}
	}
	view := w.viewport
	view.SetTitle(w.title)
	if !build.IsWayland && w.centered {
		w.doCenterOnScreen()
	}
	w.visible = true
	view.Show()
	if !build.IsWayland {
		w.xpos, w.ypos = view.GetPos()
	}
	if w.fullScreenSecondary {
		w.doSetFullScreen2(true)
	} else if w.fullScreen {
		w.doSetFullScreen(true)
	}
}

func (*gLDriver) shouldSkipPoll() bool {
	return drmDown
}

type drmConn struct {
	Status  string
	Enabled string
}

var (
	drmPrev map[string]drmConn
	drmLost map[string]drmConn
	drmDown bool
)

func drmSnapshot() map[string]drmConn {
	out := make(map[string]drmConn)
	matches, err := filepath.Glob("/sys/class/drm/card*-*/status")
	if err != nil {
		return out
	}
	for _, p := range matches {
		dir := filepath.Dir(p)
		name := filepath.Base(dir)
		out[name] = drmConn{
			Status:  drmRead(p),
			Enabled: drmRead(filepath.Join(dir, "enabled")),
		}
	}
	return out
}

func drmRead(p string) string {
	b, err := os.ReadFile(p)
	if err != nil {
		return "?"
	}
	return strings.TrimSpace(string(b))
}

func drmLostOutput(prev, cur map[string]drmConn) bool {
	if len(prev) == 0 {
		return false
	}
	for name, p := range prev {
		c := cur[name]
		if p.Status == "connected" && c.Status != "connected" {
			return true
		}
		if p.Enabled == "enabled" && c.Enabled != "enabled" {
			return true
		}
	}
	return false
}

func drmStillMissing(lost, cur map[string]drmConn) bool {
	for name, p := range lost {
		if p.Status != "connected" && p.Enabled != "enabled" {
			continue
		}
		c := cur[name]
		if p.Status == "connected" && c.Status != "connected" {
			return true
		}
		if p.Enabled == "enabled" && c.Enabled != "enabled" {
			return true
		}
	}
	return false
}

func cloneDRM(m map[string]drmConn) map[string]drmConn {
	out := make(map[string]drmConn, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
