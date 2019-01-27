package gl

import (
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"
)

// Clipboard represents the system clipboard
type Clipboard struct {
	window *glfw.Window
}

// Content returns the clipboard content
func (c *Clipboard) Content() string {
	content, err := c.window.GetClipboardString()
	if err != nil {
		log.Printf("unable to get clipboard string: %v", err)
		return ""
	}
	return content
}

// SetContent sets the clipboard content
func (c *Clipboard) SetContent(content string) {
	c.window.SetClipboardString(content)
}
