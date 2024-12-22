//go:build wasm || test_web_driver

package glfw

import (
	glfw "github.com/fyne-io/glfw-js"

	"fyne.io/fyne/v2"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = clipboard{}

func NewClipboard() fyne.Clipboard {
	return clipboard{}
}

// clipboard represents the system clipboard
type clipboard struct {
	window *glfw.Window
}

// Content returns the clipboard content
func (c clipboard) Content() string {
	content, _ := c.window.GetClipboardString()
	return content
}

// SetContent sets the clipboard content
func (c clipboard) SetContent(content string) {
	c.window.SetClipboardString(content)
}
