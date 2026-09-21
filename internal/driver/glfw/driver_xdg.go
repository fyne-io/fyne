//go:build linux || freebsd || openbsd || netbsd

package glfw

import "C"

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"

	"fyne.io/fyne/v2"
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
			call := obj.Call("org.freedesktop.ScreenSaver.UnInhibit", 0, inhibitCookie)
			if call.Err != nil {
				fyne.LogError("Failed to send message to bus", call.Err)
			}
			inhibitCookie = 0
		}
		return
	}

	if inhibitCookie != 0 {
		return
	}

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
