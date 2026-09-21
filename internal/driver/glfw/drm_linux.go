//go:build linux && !wasm && !test_web_driver

package glfw

import (
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/internal/cache"

	"github.com/go-gl/glfw/v3.4/glfw"
)

// hideIfDRMOutputLost runs on the GLFW thread immediately before PollEvents.
// KVM often leaves status=connected and only flips enabled=disabled, so both
// are watched. While a previously live connector is down we skip PollEvents
// so GLFW never dispatches wl_registry.global_remove on a mapped surface.
//
// On resume we must not PollEvents: the socket still holds global_remove
// from the skip period, which SIGSEGVs even with the window hidden.
// Terminate drops the connection; Init + create() is a clean compositor.
func (d *gLDriver) hideIfDRMOutputLost() {
	cur := drmSnapshot()
	if drmLostOutput(drmPrev, cur) {
		d.hideAllNative()
		drmSkipPoll.Store(true)
		drmAtSkip = cloneDRM(drmPrev)
	}
	if drmSkipPoll.Load() && !drmStillMissing(drmAtSkip, cur) {
		drmSkipPoll.Store(false)
		d.rebindGLFWAfterHotplug()
	}
	drmPrev = cur
}

func (d *gLDriver) hideAllNative() {
	for _, win := range d.AllWindows() {
		w, ok := win.(*window)
		if !ok || w.closing || w.viewport == nil {
			continue
		}
		w.visible = false
		w.viewport.Hide()
	}
}

func (d *gLDriver) rebindGLFWAfterHotplug() {
	d.rebindActive = true
	defer func() { d.rebindActive = false }()

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
	d.initGLFW()
	for _, win := range d.AllWindows() {
		w, ok := win.(*window)
		if !ok || w.closing {
			continue
		}
		w.remapAfterHotplug()
	}
	d.rebindPaint = true
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
	return drmSkipPoll.Load()
}

type drmConn struct {
	Status  string
	Enabled string
}

var (
	drmPrev     map[string]drmConn
	drmAtSkip   map[string]drmConn
	drmSkipPoll atomic.Bool
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
		st := drmRead(p)
		en := drmRead(filepath.Join(dir, "enabled"))
		out[name] = drmConn{Status: st, Enabled: en}
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

func drmStillMissing(atSkip, cur map[string]drmConn) bool {
	for name, p := range atSkip {
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
