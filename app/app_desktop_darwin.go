//go:build !ci && !ios && !js && !wasm && !test_web_driver
// +build !ci,!ios,!js,!wasm,!test_web_driver

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <AppKit/AppKit.h>

bool isBundled();
bool isDarkMode();
void watchTheme();
*/
import "C"
import (
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/sys/execabs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	if desk, ok := a.Driver().(systrayDriver); ok {
		desk.SetSystemTrayMenu(menu)
	}
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	a.Driver().(systrayDriver).SetSystemTrayIcon(icon)
}

func defaultVariant() fyne.ThemeVariant {
	if C.isDarkMode() {
		return theme.VariantDark
	}
	return theme.VariantLight
}

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "Library"), "Preferences")
	return filepath.Join(desktopConfig, "fyne")
}

func (a *fyneApp) OpenURL(url *url.URL) error {
	cmd := execabs.Command("open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

//export themeChanged
func themeChanged() {
	fyne.CurrentApp().Settings().(*settings).setupTheme()
}

func watchTheme() {
	C.watchTheme()
}
