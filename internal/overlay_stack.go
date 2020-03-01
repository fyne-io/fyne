package internal

import "fyne.io/fyne"

// OverlayStack is a mix-in for providing an overlay stack to canvases.
type OverlayStack struct {
	overlay fyne.CanvasObject
}

// Overlay satisfies the fyne.Canvas interface.
// Deprecated: Use Overlays() instead.
func (s *OverlayStack) Overlay() fyne.CanvasObject {
	return s.overlay
}

// Overlays satisfies the fyne.Canvas interface.
func (s *OverlayStack) Overlays() []fyne.CanvasObject {
	if s.overlay == nil {
		return nil
	}
	return []fyne.CanvasObject{s.overlay}
}

// PopOverlay satisfies the fyne.Canvas interface.
func (s *OverlayStack) PopOverlay() fyne.CanvasObject {
	overlay := s.overlay
	s.overlay = nil
	return overlay
}

// PushOverlay satisfies the fyne.Canvas interface.
func (s *OverlayStack) PushOverlay(overlay fyne.CanvasObject) {
	if overlay == nil {
		return
	}
	s.overlay = overlay
}

// RemoveOverlay satisfies the fyne.Canvas interface.
func (s *OverlayStack) RemoveOverlay(overlay fyne.CanvasObject) {
	if s.overlay != overlay {
		return
	}
	s.overlay = nil
}

// SetOverlay satisfies the fyne.Canvas interface.
// Deprecated: Use PushOverlay() instead.
func (s *OverlayStack) SetOverlay(overlay fyne.CanvasObject) {
	s.overlay = overlay
}

// TopOverlay satisfies the fyne.Canvas interface.
func (s *OverlayStack) TopOverlay() fyne.CanvasObject {
	return s.overlay
}
