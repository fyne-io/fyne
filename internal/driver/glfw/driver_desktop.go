//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package glfw

import (
	"sync"

	"fyne.io/systray"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var setup sync.Once

func (d *gLDriver) Run() {
	if goroutineID() != mainGoroutineID {
		panic("Run() or ShowAndRun() must be called from main goroutine")
	}
	d.runGL()
}

func (d *gLDriver) SetSystemTrayMenu(m *fyne.Menu) {
	setup.Do(func() {
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
	systray.SetIcon(resource.Content())
}

func (d *gLDriver) SystemTrayMenu() *fyne.Menu {
	return d.systrayMenu
}
