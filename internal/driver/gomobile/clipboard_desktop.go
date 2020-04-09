// +build !ios,!android

package gomobile

import "fyne.io/fyne"

// Content returns the clipboard content for mobile simulator runs
func (c *mobileClipboard) Content() string {
	fyne.LogError("Clipboard is not supported in mobile simulation", nil)
	return ""
}

// SetContent sets the clipboard content for mobile simulator runs
func (c *mobileClipboard) SetContent(content string) {
	fyne.LogError("Clipboard is not supported in mobile simulation", nil)
}
