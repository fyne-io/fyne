//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package glfw

import (
	"github.com/fyne-io/systray"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

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
