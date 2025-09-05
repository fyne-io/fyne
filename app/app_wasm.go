//go:build !ci && (!android || !ios || !mobile) && (wasm || test_web_driver)

package app

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"syscall/js"

	"fyne.io/fyne/v2"
	intRepo "fyne.io/fyne/v2/internal/repository"
	"fyne.io/fyne/v2/storage/repository"
)

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	notification := js.Global().Get("window").Get("Notification")
	if notification.IsUndefined() {
		fyne.LogError("Current browser does not support notifications.", nil)
		return
	}

	permission := notification.Get("permission")
	if permission.Type() != js.TypeString || permission.String() != "granted" {
		request := js.FuncOf(func(this js.Value, args []js.Value) any {
			if len(args) > 0 && args[0].Type() == js.TypeString && args[0].String() == "granted" {
				a.showNotification(n, &notification)
			} else {
				fyne.LogError("User rejected the request for notifications.", nil)
			}
			return nil
		})
		defer request.Release()
		notification.Call("requestPermission", request)
		return
	}

	a.showNotification(n, &notification)
}

func (a *fyneApp) showNotification(data *fyne.Notification, notification *js.Value) {
	icon := a.icon.Content()
	base64Str := base64.StdEncoding.EncodeToString(icon)
	mimeType := http.DetectContentType(icon)
	base64Img := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Str)
	notification.New(data.Title, map[string]any{
		"body": data.Content,
		"icon": base64Img,
	})
}

var themeChanged = js.FuncOf(func(this js.Value, args []js.Value) any {
	if len(args) > 0 && args[0].Type() == js.TypeObject {
		fyne.CurrentApp().Settings().(*settings).setupTheme()
	}
	return nil
})

func watchTheme(_ *settings) {
	js.Global().Call("matchMedia", "(prefers-color-scheme: dark)").Call("addEventListener", "change", themeChanged)
}

func stopWatchingTheme() {
	js.Global().Call("matchMedia", "(prefers-color-scheme: dark)").Call("removeEventListener", "change", themeChanged)
}

func (a *fyneApp) registerRepositories() {
	repo, err := intRepo.NewIndexDBRepository()
	if err != nil {
		fyne.LogError("failed to create repository: %v", err)
		return
	}

	repository.Register("idbfile", repo)
}
