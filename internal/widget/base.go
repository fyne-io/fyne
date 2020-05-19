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

func (b *base) Move(pos fyne.Position) {
	b.setFieldsAndRefresh(func() {
		b.pos = pos
	}, nil)
}

func (b *base) Position() fyne.Position {
	b.propertyLock.RLock()
	defer b.propertyLock.RUnlock()

	return b.pos
}

func (b *base) Size() fyne.Size {
	b.propertyLock.RLock()
	defer b.propertyLock.RUnlock()

	return b.size
}

func (b *base) Visible() bool {
	b.propertyLock.RLock()
	defer b.propertyLock.RUnlock()

	return !b.hidden
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

// setFieldsAndRefresh helps to make changes to a widget that should be followed by a refresh.
// This method is a guaranteed thread-safe way of directly manipulating widget fields.
func (b *base) setFieldsAndRefresh(f func(), w fyne.Widget) {
	b.propertyLock.Lock()
	f()
	b.propertyLock.Unlock()

	if w != nil { // the wrapping function didn't tell us what to refresh
		b.refresh(w)
	}
}

func (b *base) show(w fyne.Widget) {
	if b.Visible() {
		return
	}

	b.setFieldsAndRefresh(func() {
		b.hidden = false
	}, w)
}
