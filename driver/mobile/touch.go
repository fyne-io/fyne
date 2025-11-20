package mobile

import "fyne.io/fyne/v2"

// TouchEvent contains data relating to mobile touch events
type TouchEvent struct {
	fyne.PointEvent

	// ID represents the current ID of this touch, used to differentiate multiple fingers during a gesture.
	// The ID value may be re-used after that touch is released from the device (via TouchUp or TouchCancel).
	//
	// Since: 2.8
	ID int
}

// Touchable represents mobile touch events that can be sent to CanvasObjects
type Touchable interface {
	TouchDown(*TouchEvent)
	TouchUp(*TouchEvent)
	TouchCancel(*TouchEvent)
}
