package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
)

type base struct {
	hidden bool
	pos    fyne.Position
	size   fyne.Size
}

// Move satisfies the fyne.Widget interface.
func (b *base) Move(pos fyne.Position) {
	b.pos = pos
}

// Position satisfies the fyne.Widget interface.
func (b *base) Position() fyne.Position {
	return b.pos
}

// Size satisfies the fyne.Widget interface.
func (b *base) Size() fyne.Size {
	return b.size
}

// Visible satisfies the fyne.Widget interface.
func (b *base) Visible() bool {
	return !b.hidden
}

func (b *base) hide(w fyne.Widget) {
	if b.hidden {
		return
	}

	b.hidden = true
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
	if b.size == size {
		return
	}

	b.size = size
	r := cache.Renderer(w)
	if r == nil {
		return
	}

	r.Layout(size)
}

func (b *base) show(w fyne.Widget) {
	if !b.hidden {
		return
	}

	b.hidden = false
	w.Refresh()
}
