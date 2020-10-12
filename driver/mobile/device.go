// Package mobile provides mobile specific driver functionality.
package mobile

import "fyne.io/fyne"

// Device describes functionality only available on mobile
type Device interface {
	// ScreenInsets returns the space around a mobile screen that should not be drawn used for interactive elements
	ScreenInsets() (topLeft, bottomRight fyne.Size)

	// Request that the mobile device show the touch screen keyboard (standard layout)
	ShowVirtualKeyboard()
	// Request that the mobile device show the touch screen keyboard (custom layout)
	ShowVirtualKeyboardType(KeyboardType)
	// Request that the mobile device dismiss the touch screen keyboard
	HideVirtualKeyboard()
}
