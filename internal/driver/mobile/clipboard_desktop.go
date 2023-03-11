//go:build !ios && !android
// +build !ios,!android

package mobile

import "fyne.io/fyne/v2"

// Content returns the clipboard content for mobile simulator runs
func (c *mobileClipboard) Content() string {
	fyne.LogError("Clipboard is not supported in mobile simulation", nil)
	return ""
}

// SetContent sets the clipboard content for mobile simulator runs
func (c *mobileClipboard) SetContent(content string) {
	fyne.LogError("Clipboard is not supported in mobile simulation", nil)
}
