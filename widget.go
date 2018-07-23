package fyne

// Widget defines the standard behaviours of any widget. This extends the
// CanvasObject - a widget behaves in the same basic way but will encapsulate
// many child objects to create the rendered widget.
type Widget interface {
	CurrentSize() Size
	Resize(Size)
	CurrentPosition() Position
	Move(Position)

	MinSize() Size
	ApplyTheme()
	CanvasObjects() []CanvasObject
}
