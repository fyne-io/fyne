package fyne

// OverlayStack is a stack of CanvasObjects intended to be used as overlays of a Canvas.
type OverlayStack interface {
	// Overlays returns the overlays currently on the overlay stack.
	Overlays() []CanvasObject
	// PopOverlay removes and returns the top-most object of the overlay stack.
	PopOverlay() CanvasObject
	// PushOverlay pushes an overlay on the top of the overlay stack.
	PushOverlay(overlay CanvasObject)
	// RemoveOverlay removes the given object and all objects above it from the overlay stack.
	RemoveOverlay(overlay CanvasObject)
	// TopOverlay returns the top-most object of the overlay stack.
	TopOverlay() CanvasObject
}
