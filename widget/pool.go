package widget

import (
	"sync"

	"fyne.io/fyne/v2"
)

type pool interface {
	Obtain() fyne.CanvasObject
	Release(fyne.CanvasObject)
}

var _ pool = (*syncPool)(nil)

type syncPool struct {
	sync.Pool
}

// Obtain returns an item from the pool for use
func (p *syncPool) Obtain() (item fyne.CanvasObject) {
	o := p.Get()
	if o != nil {
		item = o.(fyne.CanvasObject)
	}
	return
}

// Release adds an item into the pool to be used later
func (p *syncPool) Release(item fyne.CanvasObject) {
	p.Put(item)
}
