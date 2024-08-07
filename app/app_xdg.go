//go:build !ci && !wasm && !test_web_driver && !android && !ios && !mobile && (linux || openbsd || freebsd || netbsd)

package app

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/notification"
	"github.com/rymdport/portal/openuri"
	portalSettings "github.com/rymdport/portal/settings"

	"fyne.io/fyne/v2"
	internalapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/theme"
)

func (a *fyneApp) OpenURL(url *url.URL) error {
	if build.IsFlatpak {
		err := openuri.OpenURI("", url.String(), nil)
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
func findFreedesktopColorScheme() fyne.ThemeVariant {
	dbusConn, err := dbus.SessionBus()
	if err != nil {
		fyne.LogError("Unable to connect to session D-Bus", err)
		return theme.VariantDark
	}

	dbusObj := dbusConn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	call := dbusObj.Call(
		"org.freedesktop.portal.Settings.Read",
		dbus.FlagNoAutoStart,
		"org.freedesktop.appearance",
		"color-scheme",
	)
	if call.Err != nil {
		// many desktops don't have this exported yet
		return theme.VariantDark
	}

	var value uint8
	if err = call.Store(&value); err != nil {
		fyne.LogError("failed to read theme variant from D-Bus", err)
		return theme.VariantDark
	}

	// See: https://github.com/flatpak/xdg-desktop-portal/blob/1.16.0/data/org.freedesktop.impl.portal.Settings.xml#L32-L46
	// 0: No preference
	// 1: Prefer dark appearance
	// 2: Prefer light appearance
	switch value {
	case 2:
		return theme.VariantLight
	case 1:
		return theme.VariantDark
	default:
		// Default to light theme to support Gnome's default see https://github.com/fyne-io/fyne/pull/3561
		return theme.VariantLight
	}
}

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	if build.IsFlatpak {
		err := a.sendNotificationThroughPortal(n)
		if err != nil {
			fyne.LogError("Sending notification using portal failed", err)
		}
		return
	}

	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		fyne.LogError("Unable to connect to session D-Bus", err)
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
var notificationID atomic.Uint64

// See https://flatpak.github.io/xdg-desktop-portal/docs/#gdbus-org.freedesktop.portal.Notification.
func (a *fyneApp) sendNotificationThroughPortal(n *fyne.Notification) error {
	return notification.Add(
		uint(notificationID.Add(1)),
		notification.Content{
			Title: n.Title,
			Body:  n.Content,
			Icon:  a.uniqueID,
		},
	)
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

func watchTheme() {
	go func() {
		// with portal this may not be immediate, so we update a cache instead
		internalapp.CurrentVariant.Store(uint64(findFreedesktopColorScheme()))
		// ensure initial variant setting is applied at startup
		fyne.CurrentApp().Settings().(*settings).setupTheme()

		portalSettings.OnSignalSettingChanged(func(changed portalSettings.Changed) {
			if changed.Namespace == "org.freedesktop.appearance" && changed.Key == "color-scheme" {
				internalapp.CurrentVariant.Store(uint64(findFreedesktopColorScheme()))
				fyne.CurrentApp().Settings().(*settings).setupTheme()
			}
		})
	}()
}
