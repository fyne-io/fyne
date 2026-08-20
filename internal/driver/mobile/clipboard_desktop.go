//go:build !ios && !android

package mobile

import "fyne.io/fyne/v2"

// Content returns the clipboard content for mobile simulator runs
func (mobileClipboard) Content() string {
	fyne.LogError("Clipboard is not supported in mobile simulation", nil)
	return ""
}

// SetContent sets the clipboard content for mobile simulator runs
func (mobileClipboard) SetContent(string) {
	fyne.LogError("Clipboard is not supported in mobile simulation", nil)
}
