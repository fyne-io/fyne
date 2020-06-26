// Package desktop provides desktop specific driver functionality.
package desktop

import "fyne.io/fyne"

// Driver represents the extended capabilities of a desktop driver
type Driver interface {
	// Create a new borderless window that is centered on screen
	CreateSplashWindow() fyne.Window
}
