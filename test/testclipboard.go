package test

import "fyne.io/fyne/v2"

type testClipboard struct {
	content string
}

func (c *testClipboard) Content() string {
	return c.content
}

func (c *testClipboard) SetContent(content string) {
	c.content = content
}

// NewClipboard returns a single use in-memory clipboard used for testing
func NewClipboard() fyne.Clipboard {
	return &testClipboard{}
}
