package internal

import "fyne.io/fyne"

// OverlayStack implements fyne.OverlayStack
type OverlayStack struct {
	overlay fyne.CanvasObject
}

var _ fyne.OverlayStack = (*OverlayStack)(nil)

// Add satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) Add(overlay fyne.CanvasObject) {
	if overlay == nil {
		return
	}
	s.overlay = overlay
}

// List satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) List() []fyne.CanvasObject {
	if s.overlay == nil {
		return nil
	}
	return []fyne.CanvasObject{s.overlay}
}

// Remove satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) Remove(overlay fyne.CanvasObject) {
	if s.overlay != overlay {
		return
	}
	s.overlay = nil
}

// Top satisfies the fyne.OverlayStack interface.
func (s *OverlayStack) Top() fyne.CanvasObject {
	return s.overlay
}
