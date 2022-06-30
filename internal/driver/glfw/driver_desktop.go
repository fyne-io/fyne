//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package glfw

import (
	"runtime"
	"sync"

	"fyne.io/systray"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	systrayIcon fyne.Resource
	setup       sync.Once
)

func goroutineID() (id uint64) {
	var buf [30]byte
	runtime.Stack(buf[:], false)
	for i := 10; buf[i] != ' '; i++ {
		id = id*10 + uint64(buf[i]&15)
	}
	return id
}

func (d *gLDriver) SetSystemTrayMenu(m *fyne.Menu) {
	setup.Do(func() {
		d.trayStart, d.trayStop = systray.RunWithExternalLoop(func() {
			if systrayIcon != nil {
				d.SetSystemTrayIcon(systrayIcon)
			} else if fyne.CurrentApp().Icon() != nil {
				d.SetSystemTrayIcon(fyne.CurrentApp().Icon())
			} else {
				d.SetSystemTrayIcon(theme.FyneLogo())
			}

			// it must be refreshed after init, so an earlier call would have been ineffective
			d.refreshSystray(m)
		}, func() {
			// anything required for tear-down
		})
	})

	d.refreshSystray(m)
}

func (d *gLDriver) refreshSystray(m *fyne.Menu) {
	d.systrayMenu = m
	systray.ResetMenu()
	for _, i := range m.Items {
		if i.IsSeparator {
			systray.AddSeparator()
			continue
		}

		var item *systray.MenuItem
		fn := i.Action

		if i.Checked {
			item = systray.AddMenuItemCheckbox(i.Label, i.Label, true)
		} else {
			item = systray.AddMenuItem(i.Label, i.Label)
		}
		if i.Disabled {
			item.Disable()
		}

		go func() {
			for range item.ClickedCh {
				fn()
			}
		}()
	}

	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Quit application")
	go func() {
		<-quit.ClickedCh
		d.Quit()
	}()
}

func (d *gLDriver) SetSystemTrayIcon(resource fyne.Resource) {
	systrayIcon = resource // in case we need it later

	img, err := toOSIcon(resource)
	if err != nil {
		fyne.LogError("Failed to convert systray icon", err)
		return
	}

	systray.SetIcon(img)
}

func (d *gLDriver) SystemTrayMenu() *fyne.Menu {
	return d.systrayMenu
}
