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
	content := ""
	runOnMain(func() {
		content = glfw.GetClipboardString()
	})
	return content
}

// SetContent sets the clipboard content
func (c *clipboard) SetContent(content string) {
	runOnMain(func() {
		defer func() {
			if r := recover(); r != nil {
				fyne.LogError("GLFW clipboard error (details above)", nil)
			}
		}()

		glfw.SetClipboardString(content)
	})
}
