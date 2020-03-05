package internal

import "fyne.io/fyne"

// OverlayStack implements fyne.OverlayStack
type OverlayStack struct {
	overlay fyne.CanvasObject
}

var _ fyne.OverlayStack = (*OverlayStack)(nil)

// AddOverlay satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) AddOverlay(overlay fyne.CanvasObject) {
	if overlay == nil {
		return
	}
	s.overlay = overlay
}

// Overlays satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) Overlays() []fyne.CanvasObject {
	if s.overlay == nil {
		return nil
	}
	return []fyne.CanvasObject{s.overlay}
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
