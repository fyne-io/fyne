//go:build linux || freebsd || openbsd || netbsd

package glfw

import "C"

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"

	"fyne.io/fyne/v2"
)

const portalInhibitIdle = 8 // See https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.impl.portal.Inhibit.html#org-freedesktop-impl-portal-inhibit-inhibit

var inhibitCookie uint32
var requestPath dbus.ObjectPath

func setDisableScreenBlank(disable bool) {
	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		fyne.LogError("Unable to connect to session D-Bus", err)
		return
	}

	if !disable {
		if requestPath != "" {
			obj := conn.Object("org.freedesktop.portal.Desktop", requestPath)
			call := obj.Call("org.freedesktop.portal.Request.Close", 0)
			if call.Err != nil {
				fyne.LogError("Failed to close portal inhibit request", call.Err)
			}
			requestPath = ""
		}

		// Shouldn't hit this if the portal is in use
		if inhibitCookie != 0 {
			obj := conn.Object("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver")
			call := obj.Call("org.freedesktop.ScreenSaver.UnInhibit", 0, inhibitCookie)
			if call.Err != nil {
				fyne.LogError("Failed to send message to bus", call.Err)
			}
			inhibitCookie = 0
		}
		return
	}

	if requestPath != "" || inhibitCookie != 0 {
		return // Idle was already inhibited, so don't worry about anything.
	}

	// Try the portal first
	portal := conn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	call := portal.Call("org.freedesktop.portal.Inhibit.Inhibit", 0, "",
		uint32(portalInhibitIdle),
		// meh, don't bother passing anything
		map[string]dbus.Variant{})
	if call.Err == nil {
		var handle dbus.ObjectPath
		if err := call.Store(&handle); err == nil && handle != "" {
			requestPath = handle
			return
		}
		fyne.LogError(fmt.Sprintf("Unexpected portal result: %#v", call.Body), nil)
	} else {
		fyne.LogError("Portal inhibit failed, falling back to org.freedesktop.ScreenSaver", call.Err)
	}

	// Portal either wasn't available or failed, so fall back to ScreenSaver
	obj := conn.Object("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver")
	call := obj.Call("org.freedesktop.ScreenSaver.Inhibit", 0, fyne.CurrentApp().Metadata().ID,
		"App disabled screensaver")
	if call.Err == nil {
		var ok bool
		callResult := call.Body[0]
		if inhibitCookie, ok = callResult.(uint32); !ok {
			fyne.LogError(fmt.Sprintf("Unexpected D-Bus call result: %#v", callResult), nil)
		}
	} else {
		fyne.LogError("Failed to send message to bus", call.Err)
	}
}

func (*gLDriver) DoubleTapDelay() time.Duration {
	return desktopDefaultDoubleTapDelay
}
