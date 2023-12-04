//go:build !js && !wasm && !test_web_driver
// +build !js,!wasm,!test_web_driver

package glfw

import (
	"bytes"
	"image/png"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/svg"
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
				d.SetSystemTrayIcon(theme.BrokenImageIcon())
			}

			// Some XDG systray crash without a title (See #3678)
			if runtime.GOOS == "linux" || runtime.GOOS == "openbsd" || runtime.GOOS == "freebsd" || runtime.GOOS == "netbsd" {
				app := fyne.CurrentApp()
				title := app.Metadata().Name
				if title == "" {
					title = app.UniqueID()
				}

				systray.SetTitle(title)
			}

			// it must be refreshed after init, so an earlier call would have been ineffective
			d.refreshSystray(m)
		}, func() {
			// anything required for tear-down
		})

		// the only way we know the app was asked to quit is if this window is asked to close...
		w := d.CreateWindow("SystrayMonitor")
		w.(*window).create()
		w.SetCloseIntercept(d.Quit)
		w.SetOnClosed(systray.Quit)
	})

	d.refreshSystray(m)
}

func itemForMenuItem(i *fyne.MenuItem, parent *systray.MenuItem) *systray.MenuItem {
	if i.IsSeparator {
		if parent != nil {
			parent.AddSeparator()
		} else {
			systray.AddSeparator()
		}
		return nil
	}

	var item *systray.MenuItem
	if i.Checked {
		if parent != nil {
			item = parent.AddSubMenuItemCheckbox(i.Label, i.Label, true)
		} else {
			item = systray.AddMenuItemCheckbox(i.Label, i.Label, true)
		}
	} else {
		if parent != nil {
			item = parent.AddSubMenuItem(i.Label, i.Label)
		} else {
			item = systray.AddMenuItem(i.Label, i.Label)
		}
	}
	if i.Disabled {
		item.Disable()
	}
	if i.Icon != nil {
		data := i.Icon.Content()
		if svg.IsResourceSVG(i.Icon) {
			b := &bytes.Buffer{}
			res := i.Icon
			if runtime.GOOS == "windows" && isDark() { // windows menus don't match dark mode so invert icons
				res = theme.NewInvertedThemedResource(i.Icon)
			}
			img := painter.PaintImage(canvas.NewImageFromResource(res), nil, 64, 64)
			err := png.Encode(b, img)
			if err != nil {
				fyne.LogError("Failed to encode SVG icon for menu", err)
			} else {
				data = b.Bytes()
			}
		}

		img, err := toOSIcon(data)
		if err != nil {
			fyne.LogError("Failed to convert systray icon", err)
		} else {
			item.SetIcon(img)
		}
	}
	return item
}

func (d *gLDriver) refreshSystray(m *fyne.Menu) {
	d.systrayMenu = m
	systray.ResetMenu()
	d.refreshSystrayMenu(m, nil)

	addMissingQuitForMenu(m, d)
}

func (d *gLDriver) refreshSystrayMenu(m *fyne.Menu, parent *systray.MenuItem) {
	for _, i := range m.Items {
		item := itemForMenuItem(i, parent)
		if item == nil {
			continue // separator
		}
		if i.ChildMenu != nil {
			d.refreshSystrayMenu(i.ChildMenu, item)
		}

		fn := i.Action
		go func() {
			for range item.ClickedCh {
				if fn != nil {
					fn()
				}
			}
		}()
	}
}

func (d *gLDriver) SetSystemTrayIcon(resource fyne.Resource) {
	systrayIcon = resource // in case we need it later

	img, err := toOSIcon(resource.Content())
	if err != nil {
		fyne.LogError("Failed to convert systray icon", err)
		return
	}

	systray.SetIcon(img)
}

func (d *gLDriver) SystemTrayMenu() *fyne.Menu {
	return d.systrayMenu
}

func (d *gLDriver) CurrentKeyModifiers() fyne.KeyModifier {
	return d.currentKeyModifiers
}

func (d *gLDriver) catchTerm() {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	d.Quit()
}

func addMissingQuitForMenu(menu *fyne.Menu, d *gLDriver) {
	var lastItem *fyne.MenuItem
	if len(menu.Items) > 0 {
		lastItem = menu.Items[len(menu.Items)-1]
		if lastItem.Label == "Quit" {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit { // make sure the menu always has a quit option
		quitItem := fyne.NewMenuItem("Quit", nil)
		quitItem.IsQuit = true
		menu.Items = append(menu.Items, fyne.NewMenuItemSeparator(), quitItem)
	}
	for _, item := range menu.Items {
		if item.IsQuit && item.Action == nil {
			item.Action = d.Quit
		}
	}
}
