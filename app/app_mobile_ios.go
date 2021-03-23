// +build !ci

// +build ios

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework UserNotifications

#include <stdlib.h>

char *documentsPath(void);
void openURL(char *urlStr);
void sendNotification(char *title, char *content);
*/
import "C"
import (
	"net/url"
	"path/filepath"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func defaultVariant() fyne.ThemeVariant {
	// TODO read the iOS setting when 10.13 arrives in 2019
	return theme.VariantLight
}

func rootConfigDir() string {
	root := C.documentsPath()
	return filepath.Join(C.GoString(root), "fyne")
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
