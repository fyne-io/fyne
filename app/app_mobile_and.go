// +build !ci

// +build android

package app

/*
#cgo LDFLAGS: -landroid -llog

#include <stdlib.h>

void sendNotification(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *title, char *content);
*/
import "C"
import (
	"log"
	"net/url"
	"os"
	"unsafe"

	mobileApp "github.com/fyne-io/mobile/app"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func defaultTheme() fyne.Theme {
	return theme.LightTheme()
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	cmd := app.exec("/system/bin/am", "start", "-a", "android.intent.action.VIEW", "--user", "0",
		"-d", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func (app *fyneApp) SendNotification(n *fyne.Notification) {
	titleStr := C.CString(n.Title)
	defer C.free(unsafe.Pointer(titleStr))
	contentStr := C.CString(n.Content)
	defer C.free(unsafe.Pointer(contentStr))

	mobileApp.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.sendNotification(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), titleStr, contentStr)
		return nil
	})
}

func rootConfigDir() string {
	filesDir := os.Getenv("FILESDIR")
	if filesDir == "" {
		log.Println("FILESDIR env was not set by android native code")
		return "/data/data" // probably won't work, but we can't make a better guess
	}

	return filesDir
}
