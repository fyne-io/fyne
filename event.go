package fyne

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	Name KeyName
}

// PointEvent describes a pointer input event. The position is relative to the
// top-left of the CanvasObject this is triggered on.
type PointEvent struct {
	Position Position // The position of the event
}

// ScrollEvent defines the parameters of a pointer or other scroll event.
// The DeltaX and DeltaY represent how large the scroll was in two dimensions.
type ScrollEvent struct {
	PointEvent
	DeltaX, DeltaY int
}
