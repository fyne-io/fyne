//go:build !ci && (!android || !ios || !mobile) && (js || wasm || test_web_driver)
// +build !ci
// +build !android !ios !mobile
// +build js wasm test_web_driver

package app

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"syscall/js"

	"fyne.io/fyne/v2"
)

func (app *fyneApp) SendNotification(n *fyne.Notification) {
	window := js.Global().Get("window")
	if window.IsUndefined() {
		fyne.LogError("Current browser does not support notifications.", nil)
		return
	}
	notification := window.Get("Notification")
	if window.IsUndefined() {
		fyne.LogError("Current browser does not support notifications.", nil)
		return
	}
	// check permission
	permission := notification.Get("permission")
	showNotification := func() {
		icon := app.icon.Content()
		base64Str := base64.StdEncoding.EncodeToString(icon)
		mimeType := http.DetectContentType(icon)
		base64Img := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
		notification.New(n.Title, map[string]interface{}{
			"body": n.Content,
			"icon": base64Img,
		})
		fyne.LogError("done show...", nil)
	}
	if permission.Type() != js.TypeString || permission.String() != "granted" {
		// need to request for permission
		notification.Call("requestPermission", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 && args[0].Type() == js.TypeString && args[0].String() == "granted" {
				showNotification()
			} else {
				fyne.LogError("User rejected the request for notifications.", nil)
			}
			return nil
		}))
	} else {
		showNotification()
	}
}

func rootConfigDir() string {
	return "/data/"
}
