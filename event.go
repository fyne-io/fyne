package fyne

import (
	"math"
)

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	Name KeyName
}

// PointerID tries to identify events created by the same finger or pointer device.
type PointerID uint32

// PointerEvent describes a pointer input of a specific input source (mouse, stylus or finger).
// The object is persistent for a given finger, mouse or stylus and the platform driver
// mutates it as it tracks its movement
type PointerEvent struct {
	// Device type reporting the event
	Type PointerDeviceType

	// ID tries to identify events created by the same finger or pointer device.
	// The same ID for cannot be used for a finger and a mouse pointer at the same time but
	// if the finger is released and pressed down again, it can be re-assigned the same
	// or different ID. If the same finger touches down and moves, it should be identified
	// by the same ID in those events.
	ID PointerID

	// Primary pointer of a device type (e.g. the finger which touched the multitouch-capable
	// display first, while reporing input changes for other fingers as well). Note there may
	// be another primary pointer for another device type (e.g. mouse) at the same time.
	Primary bool

	// The position of the event relative to the top-left of the CanvasObject
	Position Position

	// A bit map of current states of individual buttons (bit set if pressed).
	Buttons PointerButton

	// A bit map of buttons which changed state since the previous event. To determine
	// whether the button is now pressed or not, use Buttons attribute.
	ButtonsChange PointerButton

	// Pressure if reported by the device in range [0.0 - 1.0]. Set to 0.5 if device
	// doesn't support pressure.
	Pressure float64

	// Angle of tilt in degress of the pointer along the x-axis in a range of (-90, +90),
	// with a positive value indicating a tilt to the right.
	// See https://www.w3.org/TR/pointerevents/tiltX_600px.png
	XTilt int

	// Angle of tilt in degress of the pointer along the y-axis in a range of (-90, +90),
	// with a positive value indicating tilt toward the user
	// See https://www.w3.org/TR/pointerevents/tiltY_600px.png
	YTilt int

	// The clockwise rotation (in degrees, in the range of [0,359]) of a transducer
	// (e.g. pen stylus) around its own major axis. 0 for hardware and platforms
	// that do not report twist
	Twist int
}

// PointEvent has been renamed to PointerEvent.
// Please update the code to use PointerEvent!
type PointEvent PointerEvent

// XTiltYTiltRad returns stylus XTilt and YTilt in radians
func (ev *PointerEvent) XTiltYTiltRad() (float64, float64) {
	return float64(ev.XTilt) * math.Pi / 180, float64(ev.YTilt) * math.Pi / 180
}

// AzimuthAltitudeRad returns stylus orientation alternatively
// converted to azimuth and altitude
func (ev *PointerEvent) AzimuthAltitudeRad() (float64, float64) {
	xTiltRad, yTiltRad := ev.XTiltYTiltRad()
	azimuthRad := 0.0

	if ev.XTilt != 0 {
		azimuthRad = math.Pi/2 - math.Atan2(-math.Cos(xTiltRad)*math.Sin(yTiltRad), math.Cos(yTiltRad)*math.Sin(xTiltRad))
		if azimuthRad < 0 { // fix range to [0, 2*pi)
			azimuthRad += 2 * math.Pi
		}
	}

	return azimuthRad, math.Pi/2 - math.Acos(math.Cos(xTiltRad)*math.Cos(yTiltRad))
}

// ScrollEvent defines the parameters of a pointer or other scroll event.
// The DeltaX and DeltaY represent how large the scroll was in two dimensions.
type ScrollEvent struct {
	DeltaX, DeltaY int
}
