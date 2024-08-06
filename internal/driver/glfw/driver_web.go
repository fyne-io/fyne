//go:build wasm || test_web_driver

package glfw

import (
	"time"

	"fyne.io/fyne/v2"
)

const webDefaultDoubleTapDelay = 300 * time.Millisecond

func (d *gLDriver) SetSystemTrayMenu(m *fyne.Menu) {
	// no-op for mobile apps using this driver
}

func (d *gLDriver) catchTerm() {
}

func setDisableScreenBlank(disable bool) {
	// awaiting complete support for WakeLock
}

func (g *gLDriver) DoubleTapDelay() time.Duration {
	return webDefaultDoubleTapDelay
}
