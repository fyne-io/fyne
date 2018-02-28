package event

import "github.com/fyne-io/fyne/ui"

type MouseEvent struct {
	Position ui.Position
	Button   MouseButton
}

type MouseButton int

const (
	LeftMouseButton MouseButton = iota + 1
	MiddleMouseButton
	RightMouseButton
)
