//go:build !ci && !js && !wasm && !test_web_driver
// +build !ci,!js,!wasm,!test_web_driver

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <stdbool.h>
#include <stdlib.h>

bool isBundled();
void sendNotification(char *title, char *content);
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"fyne.io/fyne/v2"
	"golang.org/x/sys/execabs"
)

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	if C.isBundled() {
		titleStr := C.CString(n.Title)
		defer C.free(unsafe.Pointer(titleStr))
		contentStr := C.CString(n.Content)
		defer C.free(unsafe.Pointer(contentStr))

		C.sendNotification(titleStr, contentStr)
		return
	}

	fallbackNotification(n.Title, n.Content)
}

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	if desk, ok := a.Driver().(systrayDriver); ok {
		desk.SetSystemTrayMenu(menu)
	}
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	a.Driver().(systrayDriver).SetSystemTrayIcon(icon)
}

func escapeNotificationString(in string) string {
	noSlash := strings.ReplaceAll(in, "\\", "\\\\")
	return strings.ReplaceAll(noSlash, "\"", "\\\"")
}

//export fallbackSend
func fallbackSend(cTitle, cContent *C.char) {
	title := C.GoString(cTitle)
	content := C.GoString(cContent)
	fallbackNotification(title, content)
}

func fallbackNotification(title, content string) {
	template := `display notification "%s" with title "%s"`
	script := fmt.Sprintf(template, escapeNotificationString(content), escapeNotificationString(title))

	err := execabs.Command("osascript", "-e", script).Start()
	if err != nil {
		fyne.LogError("Failed to launch darwin notify script", err)
	}
}
