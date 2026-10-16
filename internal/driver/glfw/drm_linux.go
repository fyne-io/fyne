//go:build linux && !wasm && !test_web_driver

package glfw

import (
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/internal/cache"

	"github.com/go-gl/glfw/v3.4/glfw"
)

// skipPollForDRM runs on the GLFW thread before PollEvents.
//
// The crash is glfw.PollEvents dispatching wl_registry.global_remove for an
// output that just went away. KVM often leaves status=connected and only
// flips enabled=disabled, so both are watched. GLFW's own monitor callback
// runs inside PollEvents, which is the call that crashes, so this looks at
// sysfs first.
//
// On loss we Terminate without polling. The socket still holds global_remove,
// and reading it later crashes even if the window is hidden. When the output
// is back we Init and Show the same windows again.
func (d *gLDriver) skipPollForDRM() bool {
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
	return drmDown
}

func (d *gLDriver) dropDisplay() {
	for _, win := range d.AllWindows() {
		w, ok := win.(*window)
		if !ok {
			continue
		}
		w.visible = false
		// Display is still connected. Free the frame callback now; after
		// Terminate it would point at a dead wl_display.
		if w.frame != nil {
			w.frame.free()
		}
		w.frame = newPresentGate(w)
		if w.canvas != nil {
			cache.DropTexturesFor(w.canvas)
		}
		w.viewport = nil
		w.created = false
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
		w.Show()
	}
	d.rebindPaint = true
	return true
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
		out[filepath.Base(dir)] = drmConn{
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
