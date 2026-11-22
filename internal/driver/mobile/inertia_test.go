package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

type inertiaRecorder struct {
	drags int
	ends  int
	last  fyne.Delta
}

func (r *inertiaRecorder) Dragged(ev *fyne.DragEvent) {
	r.drags++
	r.last = ev.Dragged
}

func (r *inertiaRecorder) DragEnd() {
	r.ends++
}

func Test_driver_InertiaDecaysAndEndsOnce(t *testing.T) {
	d := &driver{}
	rec := &inertiaRecorder{}
	d.startInertia([]inertiaTarget{{wid: rec, event: *dragEvent(0, 20)}})

	for i := 0; i < 100 && d.inertia != nil; i++ {
		d.stepInertia()
	}

	assert.Nil(t, d.inertia)
	assert.Greater(t, rec.drags, 1)
	assert.Equal(t, 1, rec.ends)
}

func Test_driver_InertiaStopsWhenTouchBegins(t *testing.T) {
	d := &driver{}
	rec := &inertiaRecorder{}
	d.startInertia([]inertiaTarget{{wid: rec, event: *dragEvent(0, 20)}})
	d.stepInertia()
	assert.Equal(t, 1, rec.drags)

	// A new touch cancels the old flick and ends it exactly once.
	d.stopInertia()
	assert.Equal(t, 1, rec.ends)
	d.stepInertia()
	assert.Equal(t, 1, rec.drags)
	assert.Equal(t, 1, rec.ends)

	// A following flick in the other direction is the only one still running.
	d.startInertia([]inertiaTarget{{wid: rec, event: *dragEvent(0, -10)}})
	d.stepInertia()
	assert.Equal(t, 2, rec.drags)
	assert.InDelta(t, -10*tapMoveDecay, rec.last.DY, 0.0001)
	assert.Equal(t, 1, rec.ends)
}

func dragEvent(dx, dy float32) *fyne.DragEvent {
	return &fyne.DragEvent{
		PointEvent: fyne.PointEvent{Position: fyne.NewPos(5, 5), AbsolutePosition: fyne.NewPos(5, 5)},
		Dragged:    fyne.NewDelta(dx, dy),
	}
}
