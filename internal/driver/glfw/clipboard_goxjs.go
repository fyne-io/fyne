//go:build js || wasm || test_web_driver
// +build js wasm test_web_driver

package glfw

import (
	"fyne.io/fyne/v2"
	glfw "github.com/fyne-io/glfw-js"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = (*clipboard)(nil)

// clipboard represents the system clipboard
type clipboard struct {
	window *glfw.Window
}

// Content returns the clipboard content
func (c *clipboard) Content() string {
	content := ""
	runOnMain(func() {
		content, _ = c.window.GetClipboardString()
	})
	return content
}

// SetContent sets the clipboard content
func (c *clipboard) SetContent(content string) {
	runOnMain(func() {
		c.window.SetClipboardString(content)
	})
}
