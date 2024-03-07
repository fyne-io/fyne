package desktop

import "fyne.io/fyne/v2"

// App defines the desktop specific extensions to a fyne.App.
//
// Since: 2.2
type App interface {
	SetSystemTrayMenu(menu *fyne.Menu)
	// SetSystemTrayIcon sets the icon to be used in system tray.
	// If you pass a `ThemedResource` then any OS that adjusts look to match theme will adapt the icon.
	SetSystemTrayIcon(icon fyne.Resource)
}
