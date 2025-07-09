//go:build !ci && !wasm && !test_web_driver && !android && !ios && !mobile && (linux || openbsd || freebsd || netbsd)

package app

import (
	"net/url"
	"os"
	"os/exec"
	"sync/atomic"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/notification"
	"github.com/rymdport/portal/openuri"
	portalSettings "github.com/rymdport/portal/settings"
	"github.com/rymdport/portal/settings/appearance"

	"fyne.io/fyne/v2"
	internalapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/theme"
)

const systemTheme = fyne.ThemeVariant(99)

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
	colorScheme, err := appearance.GetColorScheme()
	if err != nil {
		return systemTheme
	}

	return colorSchemeToThemeVariant(colorScheme)
}

func colorSchemeToThemeVariant(colorScheme appearance.ColorScheme) fyne.ThemeVariant {
	switch colorScheme {
	case appearance.Light:
		return theme.VariantLight
	case appearance.Dark:
		return theme.VariantDark
	}

	// Default to light theme to support Gnome's default see https://github.com/fyne-io/fyne/pull/3561
	return theme.VariantLight
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

func watchTheme(s *settings) {
	go func() {
		// Theme lookup hangs on some desktops. Update theme variant cache from within goroutine.
		themeVariant := findFreedesktopColorScheme()
		if themeVariant != systemTheme {
			internalapp.CurrentVariant.Store(uint64(themeVariant))
			fyne.Do(func() { s.applyVariant(themeVariant) })
		}

		portalSettings.OnSignalSettingChanged(func(changed portalSettings.Changed) {
			if changed.Namespace == appearance.Namespace && changed.Key == "color-scheme" {
				themeVariant := colorSchemeToThemeVariant(appearance.ColorScheme(changed.Value.(uint32)))
				internalapp.CurrentVariant.Store(uint64(themeVariant))
				fyne.Do(func() { s.applyVariant(themeVariant) })
			}
		})
	}()
}

func (a *fyneApp) registerRepositories() {
	// no-op
}
