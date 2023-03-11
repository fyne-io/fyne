package desktop

import "fyne.io/fyne/v2"

// App defines the desktop specific extensions to a fyne.App.
//
// Since: 2.2
type App interface {
	SetSystemTrayMenu(menu *fyne.Menu)
	SetSystemTrayIcon(icon fyne.Resource)
}
