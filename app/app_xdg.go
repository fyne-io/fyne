//go:build !ci && !wasm && !test_web_driver && (linux || openbsd || freebsd || netbsd) && !android

package app

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/notification"
	"github.com/rymdport/portal/openuri"
	portalSettings "github.com/rymdport/portal/settings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/theme"
)

var once sync.Once

func defaultVariant() fyne.ThemeVariant {
	return findFreedestktopColorScheme()
}

func (a *fyneApp) OpenURL(url *url.URL) error {
	if build.IsFlatpak {
		err := openuri.OpenURI("", url.String())
		if err != nil {
			fyne.LogError("Opening url in portal failed", err)
		}
		return err
	}

	cmd := exec.Command("xdg-open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Start()
}

// fetch color variant from dbus portal desktop settings.
func findFreedestktopColorScheme() fyne.ThemeVariant {
	colourScheme, err := portalSettings.GetColorScheme()
	if err != nil {
		return theme.VariantDark
	}

	switch colourScheme {
	case portalSettings.Light:
		return theme.VariantLight
	case portalSettings.Dark:
		return theme.VariantDark
	default:
		// Default to light theme to support Gnome's default see https://github.com/fyne-io/fyne/pull/3561
		return theme.VariantLight
	}
}

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		fyne.LogError("Unable to connect to session D-Bus", err)
		return
	}

	if build.IsFlatpak {
		err := a.sendNotificationThroughPortal(conn, n)
		if err != nil {
			fyne.LogError("Sending notification using portal failed", err)
		}
		return
	}

	appIcon := a.cachedIconPath()
	timeout := int32(0) // we don't support this yet

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, a.uniqueID, uint32(0),
		appIcon, n.Title, n.Content, []string{}, map[string]dbus.Variant{}, timeout)
	if call.Err != nil {
		fyne.LogError("Failed to send message to bus", call.Err)
	}
}

// Sending with same ID replaces the old notification.
var notificationID uint = 0

// See https://flatpak.github.io/xdg-desktop-portal/docs/#gdbus-org.freedesktop.portal.Notification.
func (a *fyneApp) sendNotificationThroughPortal(_ *dbus.Conn, n *fyne.Notification) error {
	err := notification.Add(notificationID,
		&notification.Content{
			Title: n.Title,
			Body:  n.Content,
			Icon:  a.uniqueID,
		},
	)
	if err != nil {
		return err
	}

	notificationID++
	return nil
}

func (a *fyneApp) saveIconToCache(dirPath, filePath string) error {
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		fyne.LogError("Unable to create application cache directory", err)
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		fyne.LogError("Unable to create icon file", err)
		return err
	}

	defer file.Close()

	if icon := a.Icon(); icon != nil {
		_, err = file.Write(icon.Content())
		if err != nil {
			fyne.LogError("Unable to write icon contents", err)
			return err
		}
	}

	return nil
}

func (a *fyneApp) cachedIconPath() string {
	if a.Icon() == nil {
		return ""
	}

	dirPath := filepath.Join(rootCacheDir(), a.UniqueID())
	filePath := filepath.Join(dirPath, "icon.png")
	once.Do(func() {
		err := a.saveIconToCache(dirPath, filePath)
		if err != nil {
			filePath = ""
		}
	})

	return filePath
}

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	if desk, ok := a.Driver().(systrayDriver); ok { // don't use this on mobile tag
		desk.SetSystemTrayMenu(menu)
	}
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	if desk, ok := a.Driver().(systrayDriver); ok { // don't use this on mobile tag
		desk.SetSystemTrayIcon(icon)
	}
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
	go portalSettings.WatchSettingsChange(func(body []any) {
		for _, v := range body {
			if v == "color-scheme" {
				fyne.CurrentApp().Settings().(*settings).setupTheme()
				break
			}
		}
	})
}
