package mobile

// Device describes functionality only available on mobile
type Device interface {
	ShowVirtualKeyboard() // Request that the mobile device show the touch screen keyboard (standard layout)
	HideVirtualKeyboard() // Request that the mobile device dismiss the touch screen keyboard
}
