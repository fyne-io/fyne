//go:build !windows || !ci

package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

// countingWidget counts how many times the visible-object-tree walk asks its
// renderer for its objects. One walk is one such request, so the counter is a
// probe for whether tapMove re-walked the tree.
type countingWidget struct {
	widget.BaseWidget
	objects int
}

func newCountingWidget() *countingWidget {
	w := &countingWidget{}
	w.ExtendBaseWidget(w)
	w.Resize(fyne.NewSize(40, 40))
	return w
}

func (w *countingWidget) CreateRenderer() fyne.WidgetRenderer {
	return &countingRenderer{w: w}
}

type countingRenderer struct {
	w *countingWidget
}

func (*countingRenderer) Destroy()           {}
func (*countingRenderer) Layout(fyne.Size)   {}
func (*countingRenderer) Refresh()           {}
func (*countingRenderer) MinSize() fyne.Size { return fyne.NewSize(40, 40) }

func (r *countingRenderer) Objects() []fyne.CanvasObject {
	r.w.objects++
	return nil
}

// TestCanvasDragSkipsHitTestWalk reproduces the scroll/text hot path from the
// reported lag: once a drag is in flight the target is locked, so a move must not
// walk the visible object tree again. The counting widget records one hit for
// each walk.
func TestCanvasDragSkipsHitTestWalk(t *testing.T) {
	slider := widget.NewSlider(0, 100)
	slider.Resize(fyne.NewSize(160, 40))
	slider.Move(fyne.NewPos(10, 10))

	probe := newCountingWidget()
	probe.Move(fyne.NewPos(10, 120))

	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(container.NewWithoutLayout(slider, probe))
	c.Resize(fyne.NewSize(220, 220))

	dragged := 0
	drag := func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged++
	}

	c.tapDown(fyne.NewPos(30, 30), 0)
	require.Positive(t, probe.objects, "the press should walk the tree")

	c.tapMove(fyne.NewPos(90, 30), 0, drag) // establishes the drag and walks once
	assert.Equal(t, 1, dragged, "the first move should start dragging the slider")
	walks := probe.objects

	c.tapMove(fyne.NewPos(150, 30), 0, drag) // the drag is locked, nothing to re-resolve
	assert.Equal(t, 2, dragged, "the drag should continue")
	assert.Equal(t, walks, probe.objects, "a move while dragging must not re-walk the object tree")
}

// touchDrag is both touchable and draggable, so the move has to keep hit testing
// to decide whether the captured touch target must be cancelled.
type touchDrag struct {
	widget.BaseWidget

	drags   int
	cancels int
}

func newTouchDrag() *touchDrag {
	w := &touchDrag{}
	w.ExtendBaseWidget(w)
	return w
}

func (*touchDrag) CreateRenderer() fyne.WidgetRenderer { return &touchDragRenderer{} }

func (w *touchDrag) Dragged(*fyne.DragEvent) { w.drags++ }
func (*touchDrag) DragEnd()                  {}

func (*touchDrag) TouchDown(*mobile.TouchEvent) {}
func (*touchDrag) TouchUp(*mobile.TouchEvent)   {}
func (w *touchDrag) TouchCancel(*mobile.TouchEvent) {
	w.cancels++
}

type touchDragRenderer struct{}

func (*touchDragRenderer) Destroy()           {}
func (*touchDragRenderer) Layout(fyne.Size)   {}
func (*touchDragRenderer) Refresh()           {}
func (*touchDragRenderer) MinSize() fyne.Size { return fyne.NewSize(100, 100) }
func (*touchDragRenderer) Objects() []fyne.CanvasObject {
	return nil
}

// TestCanvasDragCancelsTouchTarget checks the guard added with the skip: when the
// dragged object is also the captured touch target, the tree is still walked so a
// move off it cancels the touch.
func TestCanvasDragCancelsTouchTarget(t *testing.T) {
	target := newTouchDrag()
	target.Resize(fyne.NewSize(60, 60))
	target.Move(fyne.NewPos(10, 10))

	c := newCanvas(fyne.CurrentDevice()).(*canvas)
	c.SetContent(container.NewWithoutLayout(target))
	c.Resize(fyne.NewSize(200, 200))

	drag := func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
	}

	c.tapDown(fyne.NewPos(30, 30), 0)
	require.NotNil(t, c.touched[0], "the touch target should be captured")

	c.tapMove(fyne.NewPos(40, 40), 0, drag) // starts the drag, touch still on target
	assert.NotNil(t, c.touched[0], "the touch target is still under the finger")

	c.tapMove(fyne.NewPos(160, 160), 0, drag) // dragged away from the touch target
	assert.Nil(t, c.touched[0], "moving away from the touch target must cancel the touch")
	assert.Equal(t, 1, target.cancels)
	assert.Positive(t, target.drags)
}
