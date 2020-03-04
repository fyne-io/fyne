package internal

import "fyne.io/fyne"

// OverlayStack implements fyne.OverlayStack
type OverlayStack struct {
	overlay fyne.CanvasObject
}

var _ fyne.OverlayStack = (*OverlayStack)(nil)

// Overlays satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) Overlays() []fyne.CanvasObject {
	if s.overlay == nil {
		return nil
	}
	return []fyne.CanvasObject{s.overlay}
}

// PopOverlay satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) PopOverlay() fyne.CanvasObject {
	overlay := s.overlay
	s.overlay = nil
	return overlay
}

// PushOverlay satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) PushOverlay(overlay fyne.CanvasObject) {
	if overlay == nil {
		return
	}
	s.overlay = overlay
}

// RemoveOverlay satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) RemoveOverlay(overlay fyne.CanvasObject) {
	if s.overlay != overlay {
		return
	}
	s.overlay = nil
}

// TopOverlay satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) TopOverlay() fyne.CanvasObject {
	return s.overlay
}
