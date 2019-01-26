package fyne

// Clipboard represents the system clipboard
type Clipboard struct {
	Window
}

// Content returns the clipboard content
func (c *Clipboard) Content() string {
	content, err := c.Window.GetClipboardString()
	if err != nil {
		return ""
	}
	return content
}

// SetContent sets the clipboard content
func (c *Clipboard) SetContent(content string) {
	c.Window.SetClipboardString(content)
}
