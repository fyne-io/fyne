package test

import "fyne.io/fyne/v2"

type clipboard struct {
	content string
}

func (c *clipboard) Content() string {
	return c.content
}

func (c *clipboard) SetContent(content string) {
	c.content = content
}

// NewClipboard returns a single use in-memory clipboard used for testing
func NewClipboard() fyne.Clipboard {
	return &clipboard{}
}
