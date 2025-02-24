//go:build wasm || test_web_driver

package glfw

import (
	"fyne.io/fyne/v2"
	"github.com/fyne-io/glfw-js"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = clipboard{}

func NewClipboard() fyne.Clipboard {
	return clipboard{}
}

// clipboard represents the system clipboard
type clipboard struct{}

// Content returns the clipboard content
func (c clipboard) Content() string {
	return glfw.GetClipboardString()
}

// SetContent sets the clipboard content
func (c clipboard) SetContent(content string) {
	glfw.SetClipboardString(content)
}
