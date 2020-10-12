// Package mobile provides mobile specific driver functionality.
package mobile

import "fyne.io/fyne"

// Device describes functionality only available on mobile
type Device interface {
	// ScreenInteractiveArea returns the position and size of the central interactive area.
	// Operating system elements may overlap the portions outside this area.
	ScreenInteractiveArea() (fyne.Position, fyne.Size)

	// Request that the mobile device show the touch screen keyboard (standard layout)
	ShowVirtualKeyboard()
	// Request that the mobile device show the touch screen keyboard (custom layout)
	ShowVirtualKeyboardType(KeyboardType)
	// Request that the mobile device dismiss the touch screen keyboard
	HideVirtualKeyboard()
}
