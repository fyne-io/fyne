//go:build !ci && !ios && !wasm && !test_web_driver && !mobile

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <AppKit/AppKit.h>

bool isBundled();
void watchTheme();
*/
import "C"

import (
	"net/url"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
)

func (a *fyneApp) OpenURL(url *url.URL) error {
	cmd := exec.Command("open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	a.Driver().(systrayDriver).SetSystemTrayIcon(icon)
}

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default, this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	if desk, ok := a.Driver().(systrayDriver); ok {
		desk.SetSystemTrayMenu(menu)
	}
}

// SetSystemTrayWindow assigns a window to be shown with the system tray menu is tapped.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayWindow(w fyne.Window) {
	a.Driver().(systrayDriver).SetSystemTrayWindow(w)
}

//export themeChanged
func themeChanged() {
	fyne.CurrentApp().Settings().(*settings).setupTheme()
}

func watchTheme(_ *settings) {
	C.watchTheme()
}

func (a *fyneApp) registerRepositories() {
	// no-op
}
