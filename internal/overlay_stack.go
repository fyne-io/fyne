package internal

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/internal/widget"
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
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	if overlay == nil {
		return
	}
	s.overlays = append(s.overlays, overlay)

	// TODO this should probably apply to all once #707 is addressed
	if _, ok := overlay.(*widget.OverlayContainer); ok {
		safePos, safeSize := s.Canvas.InteractiveArea()

		overlay.Resize(safeSize)
		overlay.Move(safePos)
	}

	s.focusManagers = append(s.focusManagers, app.NewFocusManager(overlay))
	if s.OnChange != nil {
		s.OnChange()
	}
}

// List returns all overlays on the stack from bottom to top.
//
// Implements: fyne.OverlayStack
func (s *OverlayStack) List() []fyne.CanvasObject {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()

	return s.overlays
}

// Remove deletes an overlay and all overlays above it from the stack.
//
// Implements: fyne.OverlayStack
func (s *OverlayStack) Remove(overlay fyne.CanvasObject) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()

	for i, o := range s.overlays {
		if o == overlay {
			s.overlays = s.overlays[:i]
			s.focusManagers = s.focusManagers[:i]
			break
		}
	}
	if s.OnChange != nil {
		s.OnChange()
	}
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
	var fm *app.FocusManager
	if len(s.focusManagers) > 0 {
		fm = s.focusManagers[len(s.focusManagers)-1]
	}
	return fm
}
