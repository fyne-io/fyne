package gomobile

import (
	"fyne.io/fyne"
)

// Declare conformity with Clipboard interface
var _ fyne.Clipboard = (*mobileClipboard)(nil)

// mobileClipboard represents the system mobileClipboard
type mobileClipboard struct {
}

// Content returns the mobileClipboard content
func (c *mobileClipboard) Content() string {
	return "" // TODO implement ticket #414
}

// SetContent sets the mobileClipboard content
func (c *mobileClipboard) SetContent(content string) {
	// TODO implement ticket #414
}
