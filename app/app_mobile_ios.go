// +build !ci

// +build darwin,ios

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework UserNotifications

#include <stdlib.h>

void openURL(char *urlStr);
void sendNotification(char *title, char *content);
*/
import "C"
import (
	"net/url"
	"os"
	"path/filepath"
	"unsafe"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func defaultTheme() fyne.Theme {
	// TODO read the iOS setting when 10.13 arrives in 2019
	return theme.LightTheme()
}

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir() // TODO this should be <Application Home>

	desktopConfig := filepath.Join(filepath.Join(homeDir, "Library"), "Preferences")
	return filepath.Join(desktopConfig, "fyne")
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	urlStr := C.CString(url.String())
	C.openURL(urlStr)
	C.free(unsafe.Pointer(urlStr))

	return nil
}

func (app *fyneApp) SendNotification(n *fyne.Notification) {
	titleStr := C.CString(n.Title)
	defer C.free(unsafe.Pointer(titleStr))
	contentStr := C.CString(n.Content)
	defer C.free(unsafe.Pointer(contentStr))

	C.sendNotification(titleStr, contentStr)
}
