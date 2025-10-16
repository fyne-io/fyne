//go:build !wasm && !test_web_driver

package glfw

import (
	"bytes"
	"image/png"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/internal/svg"
	"fyne.io/fyne/v2/lang"
	"fyne.io/systray"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var (
	systrayIcon    fyne.Resource
	systrayRunning bool
)

func (d *gLDriver) SetSystemTrayMenu(m *fyne.Menu) {
	if !systrayRunning {
		systrayRunning = true
		d.runSystray(m)
	}

	d.refreshSystray(m)
}

func (d *gLDriver) runSystray(m *fyne.Menu) {
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

		if m != nil {
			// it must be refreshed after init, so an earlier call would have been ineffective
			runOnMain(func() {
				d.refreshSystray(m)
			})
		}
	}, func() {
		// anything required for tear-down
	})

	// the only way we know the app was asked to quit is if this window is asked to close...
	w := d.CreateWindow("SystrayMonitor")
	w.(*window).create()
	w.SetCloseIntercept(d.Quit)
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
			item = parent.AddSubMenuItemCheckbox(i.Label, "", true)
		} else {
			item = systray.AddMenuItemCheckbox(i.Label, "", true)
		}
	} else {
		if parent != nil {
			item = parent.AddSubMenuItem(i.Label, "")
		} else {
			item = systray.AddMenuItem(i.Label, "")
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
			if _, ok := i.Icon.(*theme.ThemedResource); ok {
				item.SetTemplateIcon(img, img)
			} else {
				item.SetIcon(img)
			}
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
	if m == nil {
		return
	}

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
					runOnMain(fn)
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

	if _, ok := resource.(*theme.ThemedResource); ok {
		systray.SetTemplateIcon(img, img)
	} else {
		systray.SetIcon(img)
	}
}

func (d *gLDriver) SetSystemTrayWindow(w fyne.Window) {
	if !systrayRunning {
		systrayRunning = true
		d.runSystray(nil)
	}

	w.SetCloseIntercept(w.Hide)
	glw := w.(*window)
	if glw.decorate {
		systray.SetOnTapped(glw.Show)
	} else {
		systray.SetOnTapped(glw.toggleVisible)
	}
}

func (d *gLDriver) SystemTrayMenu() *fyne.Menu {
	return d.systrayMenu
}

func (d *gLDriver) CurrentKeyModifiers() fyne.KeyModifier {
	return d.currentKeyModifiers
}

// this function should be invoked from a goroutine
func (d *gLDriver) catchTerm() {
	terminateSignal := make(chan os.Signal, 1)
	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminateSignal
	fyne.Do(d.Quit)
}

func addMissingQuitForMenu(menu *fyne.Menu, d *gLDriver) {
	localQuit := lang.L("Quit")
	var lastItem *fyne.MenuItem
	if len(menu.Items) > 0 {
		lastItem = menu.Items[len(menu.Items)-1]
		if lastItem.Label == localQuit {
			lastItem.IsQuit = true
		}
	}
	if lastItem == nil || !lastItem.IsQuit { // make sure the menu always has a quit option
		quitItem := fyne.NewMenuItem(localQuit, nil)
		quitItem.IsQuit = true
		menu.Items = append(menu.Items, fyne.NewMenuItemSeparator(), quitItem)
	}
	for _, item := range menu.Items {
		if item.IsQuit && item.Action == nil {
			item.Action = d.Quit
		}
	}
}
