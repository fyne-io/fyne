//go:build !ci && !js && !wasm && !test_web_driver && (linux || openbsd || freebsd || netbsd) && !android
// +build !ci
// +build !js
// +build !wasm
// +build !test_web_driver
// +build linux openbsd freebsd netbsd
// +build !android

package app

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var once sync.Once

func defaultVariant() fyne.ThemeVariant {
	return findThemeVariant()
}

func (a *fyneApp) OpenURL(url *url.URL) error {
	cmd := a.exec("xdg-open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Start()
}

func findCurrentWM() string {

	wm := os.Getenv("XDG_CURRENT_DESKTOP")
	if wm == "" {
		wm = os.Getenv("DESKTOP_SESSION")
	}
	wm = strings.ToLower(wm)
	return wm
}

func findThemeVariant() fyne.ThemeVariant {
	wm := findCurrentWM()
	switch wm {
	case "gnome", "xfce", "unity", "gnome-shell", "gnome-classic", "mate", "gnome-mate":
		return findGnomeThemeVariant()
	case "kde", "kde-plasma", "plasma":
		return findKdeThemeVariant()
	default:
		return theme.VariantDark
	}
}

// find the current KDE theme variant. At this time, no solution.
func findKdeThemeVariant() fyne.ThemeVariant {
	homedir, err := os.UserHomeDir()
	if err != nil || homedir == "" {
		// there is a problem, fallback to dark theme
		return theme.VariantDark
	}

	// check if the user has a .config/kdeglobals file
	kdeGlobals := filepath.Join(homedir, ".config/kdeglobals")
	if _, err := os.Stat(kdeGlobals); os.IsNotExist(err) {
		// no kdeglobals file, fallback to dark theme
		return theme.VariantDark
	}

	// find the LookAndFeelPackage key in the kdeglobals file
	content, err := ioutil.ReadFile(kdeGlobals)
	if err != nil {
		// there is a problem, fallback to dark theme
		return theme.VariantDark
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "LookAndFeelPackage=") {
			// there is 2 possible values for the LookAndFeelPackage key
			// =org.kde.breeze.desktop for the light theme
			// =org.kde.breezedark.desktop for the dark theme
			if strings.HasSuffix(line, "org.kde.breeze.desktop") {
				return theme.VariantLight
			} else if strings.HasSuffix(line, "org.kde.breezedark.desktop") {
				return theme.VariantDark
			}
		}
	}

	// If we reach this point, it means that the LookAndFeelPackage key is not present in the kdeglobals file
	// we can try to calculate the theme variant from the current KDE theme using the WM activeBackground key
	for _, line := range lines {
		if strings.HasPrefix(line, "activeBackground=") {
			// read the value, it's a color in the form of r,g,b
			col := strings.Split(line, "=")[1]
			// convert the color to a hex string
			cols := strings.Split(col, ",")
			// convert the string to int
			r, _ := strconv.Atoi(cols[0])
			g, _ := strconv.Atoi(cols[1])
			b, _ := strconv.Atoi(cols[2])

			// calculate the luminance of the color
			brightness := (float32(r)/255*299 + float32(g)/255*587 + float32(b)/255*114) / 1000
			if brightness > 0.5 {
				return theme.VariantLight
			} else {
				return theme.VariantDark
			}

		}
	}

	return theme.VariantDark
}

// fetch org.gnome.desktop.interface color-scheme 'prefer-dark' or 'prefer-light' from gsettings
func findGnomeThemeVariant() fyne.ThemeVariant {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	out, err := cmd.CombinedOutput()
	if err == nil {
		w := strings.TrimSpace(string(out))
		w = strings.Trim(w, "'")
		switch w {
		case "prefer-light", "default":
			return theme.VariantLight
		case "prefer-dark":
			return theme.VariantDark
		}
	}
	return theme.VariantDark
}

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	conn, err := dbus.SessionBus() // shared connection, don't close
	if err != nil {
		fyne.LogError("Unable to connect to session D-Bus", err)
		return
	}

	appName := fyne.CurrentApp().UniqueID()
	appIcon := a.cachedIconPath()
	timeout := int32(0) // we don't support this yet

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0, appName, uint32(0),
		appIcon, n.Title, n.Content, []string{}, map[string]dbus.Variant{}, timeout)
	if call.Err != nil {
		fyne.LogError("Failed to send message to bus", call.Err)
	}
}

func (a *fyneApp) saveIconToCache(dirPath, filePath string) error {
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		fyne.LogError("Unable to create application cache directory", err)
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		fyne.LogError("Unable to create icon file", err)
		return err
	}

	defer file.Close()

	if icon := a.Icon(); icon != nil {
		_, err = file.Write(icon.Content())
		if err != nil {
			fyne.LogError("Unable to write icon contents", err)
			return err
		}
	}

	return nil
}

func (a *fyneApp) cachedIconPath() string {
	if a.Icon() == nil {
		return ""
	}

	dirPath := filepath.Join(rootCacheDir(), a.UniqueID())
	filePath := filepath.Join(dirPath, "icon.png")
	once.Do(func() {
		err := a.saveIconToCache(dirPath, filePath)
		if err != nil {
			filePath = ""
		}
	})

	return filePath

}

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	if desk, ok := a.Driver().(systrayDriver); ok { // don't use this on mobile tag
		desk.SetSystemTrayMenu(menu)
	}
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	if desk, ok := a.Driver().(systrayDriver); ok { // don't use this on mobile tag
		desk.SetSystemTrayIcon(icon)
	}
}

func rootConfigDir() string {
	desktopConfig, _ := os.UserConfigDir()
	return filepath.Join(desktopConfig, "fyne")
}

func rootCacheDir() string {
	desktopCache, _ := os.UserCacheDir()
	return filepath.Join(desktopCache, "fyne")
}

func watchTheme() {
	wm := findCurrentWM()
	switch wm {
	case "gnome":
		go watchGnomeTheme()
	case "kde":
		// no-op, not able to read linux theme in a standard way
	}
}

func watchGnomeTheme() {
	// connect to dbus to detect color-schem theme changes
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Println(err)
		return
	}

	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath("/org/freedesktop/portal/desktop"),
		dbus.WithMatchInterface("org.freedesktop.portal.Settings"),
		dbus.WithMatchMember("SettingChanged"),
	); err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	dbusChan := make(chan *dbus.Signal, 10)
	conn.Signal(dbusChan)
	currentTheme := fyne.CurrentApp().Settings().Theme()
	for {
		select {
		case sig := <-dbusChan:
			for _, v := range sig.Body {
				switch v {
				case "color-scheme":
					if currentTheme == theme.DefaultTheme() {
						variant := findGnomeThemeVariant()
						switch variant {
						case theme.VariantLight:
							fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
						case theme.VariantDark:
							fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
						}
					}
				}
			}
		}
	}

}
