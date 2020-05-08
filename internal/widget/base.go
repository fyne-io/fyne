package widget

import (
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
)

type base struct {
	hidden bool
	pos    fyne.Position
	size   fyne.Size

	propertyLock sync.RWMutex
}

// Move satisfies the fyne.Widget interface.
func (b *base) Move(pos fyne.Position) {
	b.setFieldsAndRefresh(nil, func() {
		b.pos = pos
	})
}

// Position satisfies the fyne.Widget interface.
func (b *base) Position() fyne.Position {
	var pos fyne.Position
	b.readFields(func() {
		pos = b.pos
	})
	return pos
}

// readFields provides a guaranteed thread safe way to access widget fields.
func (b *base) readFields(f func()) {
	b.propertyLock.RLock()
	defer b.propertyLock.RUnlock()

	f()
}

// SetFieldsAndRefresh helps to make changes to a widget that should be followed by a refresh.
// This method is a guaranteed thread-safe way of directly manipulating widget fields.
func (b *base) setFieldsAndRefresh(w fyne.Widget, f func()) {
	b.propertyLock.Lock()
	f()
	b.propertyLock.Unlock()

	if w != nil { // the wrapping function didn't tell us what to refresh
		b.refresh(w)
	}
}

// Size satisfies the fyne.Widget interface.
func (b *base) Size() fyne.Size {
	var size fyne.Size
	b.readFields(func() {
		size = b.size
	})
	return size
}

// Visible satisfies the fyne.Widget interface.
func (b *base) Visible() bool {
	var hidden bool
	b.readFields(func() {
		hidden = b.hidden
	})
	return !hidden
}

func (b *base) hide(w fyne.Widget) {
	if !b.Visible() {
		return
	}

	b.propertyLock.Lock()
	b.hidden = true
	b.propertyLock.Unlock()
	canvas.Refresh(w)
}

func (b *base) minSize(w fyne.Widget) fyne.Size {
	r := cache.Renderer(w)
	if r == nil {
		return fyne.NewSize(0, 0)
	}

	return r.MinSize()
}

func (b *base) refresh(w fyne.Widget) {
	r := cache.Renderer(w)
	if r == nil {
		return
	}

	r.Refresh()
}

func (b *base) resize(size fyne.Size, w fyne.Widget) {
	var baseSize fyne.Size
	b.readFields(func() {
		baseSize = b.size
	})
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

func (b *base) show(w fyne.Widget) {
	if b.Visible() {
		return
	}

	b.setFieldsAndRefresh(w, func() {
		b.hidden = false
	})
}
