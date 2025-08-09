package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestMultipleWindows_Add(t *testing.T) {
	m := NewMultipleWindows()
	assert.Zero(t, len(m.Windows))

	m.Add(NewInnerWindow("1", widget.NewLabel("Inside")))
	assert.Equal(t, 1, len(m.Windows))
}

func TestMultipleWindows_Drag(t *testing.T) {
	w := NewInnerWindow("1", widget.NewLabel("Inside"))
	m := NewMultipleWindows(w)
	_ = test.TempWidgetRenderer(t, m) // initialise display
	assert.Equal(t, 1, len(m.Windows))

	assert.True(t, w.Position().IsZero())
	w.OnDragged(&fyne.DragEvent{Dragged: fyne.Delta{DX: 10, DY: 5}})
	assert.Equal(t, float32(10), w.Position().X)
	assert.Equal(t, float32(5), w.Position().Y)
}
