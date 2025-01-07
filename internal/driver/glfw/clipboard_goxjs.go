//go:build wasm || test_web_driver

package glfw

import (
	"fyne.io/fyne/v2"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = (*clipboard)(nil)

// clipboard represents the system clipboard
type clipboard struct {
}

// Content returns the clipboard content
func (c *clipboard) Content() string {
	content := ""
	runOnMain(func() {
		win := fyne.CurrentApp().Driver().AllWindows()[0].(*window).viewport
		content, _ = win.GetClipboardString()
	})
	return content
}

// SetContent sets the clipboard content
func (c *clipboard) SetContent(content string) {
	runOnMain(func() {
		win := fyne.CurrentApp().Driver().AllWindows()[0].(*window).viewport
		win.SetClipboardString(content)
	})
}
