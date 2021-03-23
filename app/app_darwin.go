// +build !ci

// +build !ios

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#include <AppKit/AppKit.h>

bool isBundled();
bool isDarkMode();
void sendNotification(const char *title, const char *content);
void watchTheme();
*/
import "C"
import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func defaultVariant() fyne.ThemeVariant {
	if fyne.CurrentDevice().IsMobile() { // this is called in mobile simulate mode
		return theme.VariantLight
	}
	if C.isDarkMode() {
		return theme.VariantDark
	}
	return theme.VariantLight
}

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "Library"), "Preferences")
	return filepath.Join(desktopConfig, "fyne")
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := app.exec("open", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func (app *fyneApp) SendNotification(n *fyne.Notification) {
	if C.isBundled() {
		title := C.CString(n.Title)
		defer C.free(unsafe.Pointer(title))
		content := C.CString(n.Content)
		defer C.free(unsafe.Pointer(content))

		C.sendNotification(title, content)
		return
	}

	title := escapeNotificationString(n.Title)
	content := escapeNotificationString(n.Content)
	template := `display notification "%s" with title "%s"`
	script := fmt.Sprintf(template, content, title)

	err := exec.Command("osascript", "-e", script).Start()
	if err != nil {
		fyne.LogError("Failed to launch darwin notify script", err)
	}
}

func escapeNotificationString(in string) string {
	noSlash := strings.ReplaceAll(in, "\\", "\\\\")
	return strings.ReplaceAll(noSlash, "\"", "\\\"")
}

//export themeChanged
func themeChanged() {
	fyne.CurrentApp().Settings().(*settings).setupTheme()
}

func watchTheme() {
	C.watchTheme()
}
