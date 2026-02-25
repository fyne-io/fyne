package embedded

import "fyne.io/fyne/v2"

// TouchDownEvent is for indicating that an embedded device touch screen or pointing device was pressed.
//
// Since: 2.7
type TouchDownEvent struct {
	Position fyne.Position
	ID       int
}

func (t *TouchDownEvent) isEvent() {}

// TouchMoveEvent is for indicating that an embedded device touch screen or pointing device was moved whilst being pressed.
//
// Since: 2.7
type TouchMoveEvent struct {
	Position fyne.Position
	ID       int
}

func (t *TouchMoveEvent) isEvent() {}

// TouchUpEvent is for indicating that an embedded device touch screen or pointing device was released.
//
// Since: 2.7
type TouchUpEvent struct {
	Position fyne.Position
	ID       int
}

func (t *TouchUpEvent) isEvent() {}
