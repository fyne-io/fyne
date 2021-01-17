package fyne

// Widget defines the standard behaviours of any widget. This extends the
// CanvasObject - a widget behaves in the same basic way but will encapsulate
// many child objects to create the rendered widget.
type Widget interface {
	CanvasObject

	// CreateRenderer returns a new WidgetRenderer for this widget.
	// This should not be called by regular code, it is used internally to render a widget.
	CreateRenderer() WidgetRenderer
}

// WidgetRenderer defines the behaviour of a widget's implementation.
// This is returned from a widget's declarative object through the CreateRenderer()
// function and should be exactly one instance per widget in memory.
type WidgetRenderer interface {
	// Destroy is for internal use.
	Destroy()
	// Layout is a hook that is called if the widget needs to be laid out.
	// This should never call Refresh.
	Layout(Size)
	// MinSize returns the minimum size of the widget that is rendered by this renderer.
	MinSize() Size
	// Objects returns all objects that should be drawn.
	Objects() []CanvasObject
	// Refresh is a hook that is called if the widget has updated and needs to be redrawn.
	// This might trigger a Layout.
	Refresh()
}
