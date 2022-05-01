//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package glfw

import (
	"runtime"

	"fyne.io/systray"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
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
	d.trayStart, d.trayStop = systray.RunWithExternalLoop(func() {
		if fyne.CurrentApp().Icon() != nil {
			img, err := toOSIcon(fyne.CurrentApp().Icon())
			if err == nil {
				systray.SetIcon(img)
			}
		} else {
			img, err := toOSIcon(theme.FyneLogo())
			if err == nil {
				systray.SetIcon(img)
			}
		}

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
	}, func() {
		// anything required for tear-down
	})
}

func (d *gLDriver) SetSystemTrayIcon(resource fyne.Resource) {
	systray.SetIcon(resource.Content())
}
