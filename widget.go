package fyne

import "image/color"

// Widget defines the standard behaviours of any widget. This extends the
// CanvasObject - a widget behaves in the same basic way but will encapsulate
// many child objects to create the rendered widget.
type Widget interface {
	CanvasObject

	// CreateRenderer must return the same renderer if called twice.
	CreateRenderer() WidgetRenderer
}

// WidgetRenderer defines the behaviour of a widget's implementation.
// This is returned from a widget's declarative object through the Render()
// function and should be exactly one instance per widget in memory.
type WidgetRenderer interface {
	Layout(Size)
	MinSize() Size

	Refresh()
	ApplyTheme()
	BackgroundColor() color.Color
	Objects() []CanvasObject
	Destroy()
}
