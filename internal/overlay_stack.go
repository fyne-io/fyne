package internal

import (
	"sync"

	"fyne.io/fyne"
)

// OverlayStack implements fyne.OverlayStack
type OverlayStack struct {
	overlays     []fyne.CanvasObject
	propertyLock sync.RWMutex
}

var _ fyne.OverlayStack = (*OverlayStack)(nil)

// Add puts an overlay on the stack.
// Implements: fyne.OverlayStack
func (s *OverlayStack) Add(overlay fyne.CanvasObject) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	if overlay == nil {
		return
	}
	s.overlays = append(s.overlays, overlay)
}

// List returns all overlays on the stack from bottom to top.
// Implements: fyne.OverlayStack
func (s *OverlayStack) List() []fyne.CanvasObject {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	return s.overlays
}

// Remove deletes an overlay and all overlays above it from the stack.
// Implements: fyne.OverlayStack
func (s *OverlayStack) Remove(overlay fyne.CanvasObject) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	for i, o := range s.overlays {
		if o == overlay {
			s.overlays = s.overlays[:i]
			break
		}
	}
}

// Top returns the top-most overlay of the stack.
// Implements: fyne.OverlayStack
func (s *OverlayStack) Top() fyne.CanvasObject {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	if len(s.overlays) == 0 {
		return nil
	}
	return s.overlays[len(s.overlays)-1]
}
