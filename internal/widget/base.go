package widget

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
)

// Base is a simple base widget for internal use.
type Base struct {
	hidden bool
	pos    fyne.Position
	size   fyne.Size

	propertyLock sync.RWMutex
}

// Position returns the position of the widget relative to its parent.
//
// Implements: fyne.Widget
func (b *Base) Position() fyne.Position {
	b.propertyLock.RLock()
	defer b.propertyLock.RUnlock()

	return b.pos
}

// Size returns the current size of the widget.
//
// Implements: fyne.Widget
func (b *Base) Size() fyne.Size {
	b.propertyLock.RLock()
	defer b.propertyLock.RUnlock()

	return b.size
}

// Visible returns whether the widget is visible or not.
//
// Implements: fyne.Widget
func (b *Base) Visible() bool {
	b.propertyLock.RLock()
	defer b.propertyLock.RUnlock()

	return !b.hidden
}

// HideWidget is a helper method to hide and refresh a widget.
func HideWidget(b *Base, w fyne.Widget) {
	if !b.Visible() {
		return
	}

	b.propertyLock.Lock()
	b.hidden = true
	b.propertyLock.Unlock()
	canvas.Refresh(w)
}

// MinSizeOf is a helper method to get the minSize of a widget.
func MinSizeOf(w fyne.Widget) fyne.Size {
	r := cache.Renderer(w)
	if r == nil {
		return fyne.NewSize(0, 0)
	}

	return r.MinSize()
}

// MoveWidget is a helper method to set the position of a widget relative to its parent.
func MoveWidget(b *Base, w fyne.Widget, pos fyne.Position) {
	b.propertyLock.Lock()
	b.pos = pos
	b.propertyLock.Unlock()

	RefreshWidget(w)
}

// RefreshWidget is a helper method to refresh a widget.
func RefreshWidget(w fyne.Widget) {
	r := cache.Renderer(w)
	if r == nil {
		return
	}

	r.Refresh()
}

// ResizeWidget is a helper method to resize a widget.
func ResizeWidget(b *Base, w fyne.Widget, size fyne.Size) {
	b.propertyLock.RLock()
	baseSize := b.size
	b.propertyLock.RUnlock()
	if baseSize == size {
		return
	}

	b.propertyLock.Lock()
	b.size = size
	b.propertyLock.Unlock()
	r := cache.Renderer(w)
	if r == nil {
		return
	}

	r.Layout(size)
}

// ShowWidget is a helper method to show and refresh a widget.
func ShowWidget(b *Base, w fyne.Widget) {
	if b.Visible() {
		return
	}

	b.propertyLock.Lock()
	b.hidden = false
	b.propertyLock.Unlock()

	RefreshWidget(w)
}
