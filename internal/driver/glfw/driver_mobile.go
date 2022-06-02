//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package glfw

import "fyne.io/fyne/v2"

func (d *gLDriver) SetSystemTrayMenu(m *fyne.Menu) {
	// no-op for mobile apps using this driver
}
