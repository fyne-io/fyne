//go:build linux || freebsd || openbsd || netbsd

package glfw

import "C"
import (
	"fyne.io/fyne/v2"
	"github.com/godbus/dbus/v5"
)

var inhibitCookie uint32

func setDisableScreenBlank(disable bool) {
	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		fyne.LogError("Unable to connect to session D-Bus", err)
		return
	}

	if !disable {
		if inhibitCookie != 0 {
			obj := conn.Object("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver")
			call := obj.Call("org.freedesktop.ScreenSaver.Uninhibit", 0, inhibitCookie)
			if call.Err != nil {
				fyne.LogError("Failed to send message to bus", call.Err)
			}
			inhibitCookie = 0
		}
		return
	}

	obj := conn.Object("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver")
	call := obj.Call("org.freedesktop.ScreenSaver.Inhibit", 0, fyne.CurrentApp().Metadata().Name,
		"App disabled screensaver", &inhibitCookie)
	if call.Err != nil {
		fyne.LogError("Failed to send message to bus", call.Err)
	}
}
