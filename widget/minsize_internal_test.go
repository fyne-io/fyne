package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
)

type countingWidget struct {
	BaseWidget
	calls int
}

func (w *countingWidget) CreateRenderer() fyne.WidgetRenderer {
	return &countingRenderer{w: w}
}

type countingRenderer struct {
	w *countingWidget
}

func (r *countingRenderer) MinSize() fyne.Size {
	r.w.calls++
	return fyne.NewSize(10, 10)
}

func (*countingRenderer) Layout(fyne.Size)             {}
func (*countingRenderer) Refresh()                     {}
func (*countingRenderer) Objects() []fyne.CanvasObject { return nil }
func (*countingRenderer) Destroy()                     {}

func TestBaseWidget_MinSizeOncePerPass(t *testing.T) {
	w := &countingWidget{}
	w.ExtendBaseWidget(w)

	w.MinSize()
	w.MinSize()
	assert.Equal(t, 2, w.calls, "outside a pass every call asks the renderer")

	cache.BeginMinSizePass()
	defer cache.EndMinSizePass()
	w.MinSize()
	w.MinSize()
	assert.Equal(t, 3, w.calls, "in a pass the renderer is asked once")
	w.Resize(fyne.NewSize(20, 20))
	w.MinSize()
	assert.Equal(t, 4, w.calls, "a resize can change the min size")
	w.Refresh()
	w.MinSize()
	assert.Equal(t, 5, w.calls, "so can a refresh")

	cache.EndMinSizePass()
	w.MinSize()
	assert.Equal(t, 6, w.calls)
	cache.BeginMinSizePass()
	w.MinSize()
	assert.Equal(t, 7, w.calls, "a new pass computes anew")
}
