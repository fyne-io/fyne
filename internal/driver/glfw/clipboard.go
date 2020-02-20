package glfw

import (
	"fyne.io/fyne"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = (*clipboard)(nil)

// clipboard represents the system clipboard
type clipboard struct {
	window *glfw.Window
}

// Content returns the clipboard content
func (c *clipboard) Content() string {
	return c.window.GetClipboardString()
}

// SetContent sets the clipboard content
func (c *clipboard) SetContent(content string) {
	c.window.SetClipboardString(content)
}
