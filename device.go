package fyne

// DeviceOrientation represents the different ways that a mobile device can be held
type DeviceOrientation int

const (
	// OrientationVertical is the default vertical orientation
	OrientationVertical DeviceOrientation = iota
	// OrientationVerticalUpsideDown is the portrait orientation held upside down
	OrientationVerticalUpsideDown
	// OrientationHorizontalLeft is used to indicate a landscape orientation with the top to the left
	OrientationHorizontalLeft
	// OrientationHorizontalRight is used to indicate a landscape orientation with the top to the right
	OrientationHorizontalRight
)

// IsVertical is a helper utility that determines if a passed orientation is vertical
func IsVertical(orient DeviceOrientation) bool {
	return orient == OrientationVertical || orient == OrientationVerticalUpsideDown
}

// IsHorizontal is a helper utility that determines if a passed orientation is horizontal
func IsHorizontal(orient DeviceOrientation) bool {
	return !IsVertical(orient)
}

// Device provides information about the devices the code is running on
type Device interface {
	Orientation() DeviceOrientation
	IsMobile() bool
	HasKeyboard() bool
	SystemScale() float32
}

// CurrentDevice returns the device information for the current hardware (via the driver)
func CurrentDevice() Device {
	return CurrentApp().Driver().Device()
}
