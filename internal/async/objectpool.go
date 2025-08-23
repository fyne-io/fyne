//go:build !migrated_fynedo

package async

import "fyne.io/fyne/v2"

// ObjectPool represents a pool of reusable canvas objects.
type ObjectPool struct {
	Pool[fyne.CanvasObject]
}

func (p *ObjectPool) Clear() {
	for p.Get() != nil {
	}
}
