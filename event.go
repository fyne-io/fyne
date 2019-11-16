package fyne

// KeyEvent describes a keyboard input event.
type KeyEvent struct {
	Name KeyName
}

// PointEvent describes a pointer input event. The position is relative to the
// top-left of the CanvasObject this is triggered on.
type PointEvent struct {
	AbsolutePosition Position // The absolute position of the event
	Position         Position // The relative position of the event
}

// ScrollEvent defines the parameters of a pointer or other scroll event.
// The DeltaX and DeltaY represent how large the scroll was in two dimensions.
type ScrollEvent struct {
	PointEvent
	DeltaX, DeltaY int
}

// DragEvent defines the parameters of a pointer or other drag event.
// The DraggedX and DraggedY fields show how far the item was dragged since the last event.
type DragEvent struct {
	PointEvent
	DraggedX, DraggedY int
}
