//go:build !ci && !js && !wasm && !test_web_driver && (linux || openbsd || freebsd || netbsd) && !android
// +build !ci
// +build !js
// +build !wasm
// +build !test_web_driver
// +build linux openbsd freebsd netbsd
// +build !android

package app

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/godbus/dbus/v5"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func defaultVariant() fyne.ThemeVariant {
	return theme.VariantDark
}

func (a *fyneApp) OpenURL(url *url.URL) error {
	cmd := a.exec("xdg-open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Start()
}

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		fyne.LogError("Unable to connect to session D-Bus", err)
		return
	}

	appName := fyne.CurrentApp().UniqueID()
	appIcon := a.cachedIconPath()
	timeout := int32(0) // we don't support this yet

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, appName, uint32(0),
		appIcon, n.Title, n.Content, []string{}, map[string]dbus.Variant{}, timeout)
	if call.Err != nil {
		fyne.LogError("Failed to send message to bus", call.Err)
	}
}

func (a *fyneApp) cachedIconPath() string {
	metadata := a.Metadata()
	if metadata.Icon == nil || len(metadata.Icon.Content()) == 0 {
		return ""
	}

	dirPath := filepath.Join(rootCacheDir(), a.uniqueID)
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		fyne.LogError("Unable to create application cache directory", err)
		return ""
	}

	filepath := filepath.Join(dirPath, "icon.png")
	file, err := os.Create(filepath)
	if err != nil {
		fyne.LogError("Unable to create icon file", err)
		return ""
	}

	defer file.Close()

	_, err = file.Write(meta.Icon.Content())
	if err != nil {
		fyne.LogError("Unable to write icon contents", err)
		return ""
	}

	return filepath
}

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	a.Driver().(systrayDriver).SetSystemTrayMenu(menu)
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	a.Driver().(systrayDriver).SetSystemTrayIcon(icon)
}

func rootConfigDir() string {
	desktopConfig, _ := os.UserConfigDir()
	return filepath.Join(desktopConfig, "fyne")
}

func rootCacheDir() string {
	desktopCache, _ := os.UserCacheDir()
	return filepath.Join(desktopCache, "fyne")
}

func watchTheme() {
	// no-op, not able to read linux theme in a standard way
}
