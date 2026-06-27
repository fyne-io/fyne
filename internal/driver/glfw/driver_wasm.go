//go:build wasm || test_web_driver

package glfw

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const webDefaultDoubleTapDelay = 300 * time.Millisecond

func (d *gLDriver) CreateWindow(title string) (win fyne.Window) {
	// handling multiple windows by overlaying on the root for web
	var root fyne.Window
	hasVisible := false
	for _, w := range d.windows {
		if w.(*window).visible {
			hasVisible = true
			root = w
			break
		}
	}

	if !hasVisible {
		return d.newWindow(title, true)
	}

	c := root.Canvas().(*glCanvas)
	multi := c.webExtraWindows
	if multi == nil {
		multi = container.NewMultipleWindows()
		multi.Resize(c.Size())
		c.webExtraWindows = multi
	}
	inner := container.NewInnerWindow(title, canvas.NewRectangle(color.Transparent))
	multi.Add(inner)

	return wrapInnerWindow(inner, root, d)
}

func (d *gLDriver) HasSecondaryDisplay() bool {
	return false
}

func (d *gLDriver) SetSystemTrayMenu(m *fyne.Menu) {
	// no-op for wasm apps using this driver
}

func (d *gLDriver) catchTerm() {}

func setDisableScreenBlank(disable bool) {
	// awaiting complete support for WakeLock
}

func (d *gLDriver) DoubleTapDelay() time.Duration {
	return webDefaultDoubleTapDelay
}
