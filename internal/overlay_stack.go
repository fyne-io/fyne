package internal

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/widget"
)

// OverlayStack implements fyne.OverlayStack
type OverlayStack struct {
	OnChange      func()
	Canvas        fyne.Canvas
	focusManagers []*app.FocusManager
	overlays      []fyne.CanvasObject
	propertyLock  sync.RWMutex
}

var _ fyne.OverlayStack = (*OverlayStack)(nil)

// Add puts an overlay on the stack.
//
// Implements: fyne.OverlayStack
func (s *OverlayStack) Add(overlay fyne.CanvasObject) {
	if overlay == nil {
		return
	}

	if s.OnChange != nil {
		defer s.OnChange()
	}

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.overlays = append(s.overlays, overlay)

	// TODO this should probably apply to all once #707 is addressed
	if _, ok := overlay.(*widget.OverlayContainer); ok {
		safePos, safeSize := s.Canvas.InteractiveArea()

		overlay.Resize(safeSize)
		overlay.Move(safePos)
	}

	s.focusManagers = append(s.focusManagers, app.NewFocusManager(overlay))
}

// List returns all overlays on the stack from bottom to top.
//
// Implements: fyne.OverlayStack
func (s *OverlayStack) List() []fyne.CanvasObject {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	return s.overlays
}

// ListFocusManagers returns all focus managers on the stack from bottom to top.
func (s *OverlayStack) ListFocusManagers() []*app.FocusManager {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	return s.focusManagers
}

// Remove deletes an overlay and all overlays above it from the stack.
//
// Implements: fyne.OverlayStack
func (s *OverlayStack) Remove(overlay fyne.CanvasObject) {
	if s.OnChange != nil {
		defer s.OnChange()
	}

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	overlayIdx := -1
	for i, o := range s.overlays {
		if o == overlay {
			overlayIdx = i
			break
		}
	}
	if overlayIdx == -1 {
		return
	}
	// set removed elements in backing array to nil to release memory references
	for i := overlayIdx; i < len(s.overlays); i++ {
		s.overlays[i] = nil
		s.focusManagers[i] = nil
	}
	s.overlays = s.overlays[:overlayIdx]
	s.focusManagers = s.focusManagers[:overlayIdx]
}

// Top returns the top-most overlay of the stack.
//
// Implements: fyne.OverlayStack
func (s *OverlayStack) Top() fyne.CanvasObject {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	if len(s.overlays) == 0 {
		return nil
	}
	return s.overlays[len(s.overlays)-1]
}

// TopFocusManager returns the app.FocusManager assigned to the top-most overlay of the stack.
func (s *OverlayStack) TopFocusManager() *app.FocusManager {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	return s.topFocusManager()
}

func (s *OverlayStack) topFocusManager() *app.FocusManager {
	if len(s.focusManagers) == 0 {
		return nil
	}

	return s.focusManagers[len(s.focusManagers)-1]
}
