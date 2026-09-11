//go:build !wasm && !test_web_driver

package glfw

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"fyne.io/systray"
	"github.com/go-gl/glfw/v3.4/glfw"

	"fyne.io/fyne/v2"
	intsystray "fyne.io/fyne/v2/internal/driver/systray"
	"fyne.io/fyne/v2/internal/goos"
)

const systrayIconSize = 64

var systrayTray *intsystray.Tray

func (d *gLDriver) systray() *intsystray.Tray {
	if systrayTray == nil {
		systrayTray = &intsystray.Tray{
			IconSize:    systrayIconSize,
			RunOnMain:   runOnMain,
			Quit:        d.Quit,
			ToOSIcon:    toOSIcon,
			SupportsSVG: runtime.GOOS == goos.Darwin, // only macOS takes SVG tray icons
			// Windows menus don't match dark mode so icons are inverted there
			InvertMenuIcons: func() bool { return runtime.GOOS == goos.Windows && isDark() },
		}
	}
	return systrayTray
}

func (*gLDriver) HasSecondaryDisplay() bool {
	monitors := glfw.GetMonitors()
	if len(monitors) == 1 {
		return false
	}

	primaryTop, primaryLeft := monitors[0].GetPos()
	for _, m := range monitors[1:] {
		top, left := m.GetPos()
		if top != primaryTop || left != primaryLeft {
			return true
		}
	}

	return false // all the monitors had same origin, thus mirroring
}

func (d *gLDriver) SetSystemTrayMenu(m *fyne.Menu) {
	d.systrayMenu = m
	if !d.systray().Running() {
		d.runSystray(m)
	}

	d.systray().Refresh(m)
}

func (d *gLDriver) runSystray(m *fyne.Menu) {
	d.trayStart, d.trayStop = d.systray().Start(m, func() {
		// Some XDG systray crash without a title (See #3678)
		if runtime.GOOS == goos.Linux || goos.IsBSD(runtime.GOOS) {
			app := fyne.CurrentApp()
			title := app.Metadata().Name
			if title == "" {
				title = app.UniqueID()
			}

			systray.SetTitle(title)
		}
	})

	// the only way we know the app was asked to quit is if this window is asked to close...
	w := d.CreateWindow("SystrayMonitor")
	w.(*window).create()
	w.SetCloseIntercept(d.Quit)
}

func (d *gLDriver) SetSystemTrayIcon(resource fyne.Resource) {
	d.systray().SetIcon(resource)
}

func (d *gLDriver) SetSystemTrayWindow(w fyne.Window) {
	if !d.systray().Running() {
		d.runSystray(nil)
	}

	w.SetCloseIntercept(w.Hide)
	glw, _ := w.(*window)
	if glw.decorate {
		systray.SetOnTapped(func() { fyne.Do(glw.Show) })
	} else {
		systray.SetOnTapped(func() { fyne.Do(glw.toggleVisible) })
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
