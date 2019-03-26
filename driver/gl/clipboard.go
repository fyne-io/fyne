package gl

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// clipboard represents the system clipboard
type clipboard struct {
	window *glfw.Window
}

// Content returns the clipboard content
func (c *clipboard) Content() string {
	content, err := c.window.GetClipboardString()
	if err != nil {
		return ""
	}
	return content
}

// SetContent sets the clipboard content
func (c *clipboard) SetContent(content string) {
	c.window.SetClipboardString(content)
}
