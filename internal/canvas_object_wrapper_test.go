package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal"
)

var _ fyne.CanvasObject = (*simpleCanvasObject)(nil)

func TestSimpleCanvasObjectWrapper_Hide(t *testing.T) {
	o := &simpleCanvasObject{}
	w := &internal.SimpleCanvasObjectWrapper{O: o}

	assert.True(t, o.Visible())
	w.Hide()
	assert.False(t, o.Visible())
}

func TestSimpleCanvasObjectWrapper_MinSize(t *testing.T) {
	o := &simpleCanvasObject{minSize: fyne.NewSize(12, 34)}
	w := &internal.SimpleCanvasObjectWrapper{O: o}
	assert.Equal(t, fyne.NewSize(12, 34), w.MinSize())
	o = &simpleCanvasObject{minSize: fyne.NewSize(13, 42)}
	w = &internal.SimpleCanvasObjectWrapper{O: o}
	assert.Equal(t, fyne.NewSize(13, 42), w.MinSize())
}

func TestSimpleCanvasObjectWrapper_MoveAndPosition(t *testing.T) {
	o := &simpleCanvasObject{}
	w := &internal.SimpleCanvasObjectWrapper{O: o}

	assert.Equal(t, fyne.NewPos(0, 0), w.Position())
	assert.Equal(t, fyne.NewPos(0, 0), o.Position())
	w.Move(fyne.NewPos(17, 389))
	assert.Equal(t, fyne.NewPos(17, 389), w.Position())
	assert.Equal(t, fyne.NewPos(0, 0), o.Position())
}

func TestSimpleCanvasObjectWrapper_Refresh(t *testing.T) {
	o := &simpleCanvasObject{}
	w := &internal.SimpleCanvasObjectWrapper{O: o}

	assert.Equal(t, 0, o.refreshCount)
	w.Refresh()
	assert.Equal(t, 1, o.refreshCount)
	w.Refresh()
	w.Refresh()
	assert.Equal(t, 3, o.refreshCount)
}

func TestSimpleCanvasObjectWrapper_Resize(t *testing.T) {
	o := &simpleCanvasObject{}
	w := &internal.SimpleCanvasObjectWrapper{O: o}

	assert.Equal(t, fyne.NewSize(0, 0), o.Size())
	w.Resize(fyne.NewSize(17, 389))
	assert.Equal(t, fyne.NewSize(17, 389), o.Size())
}

func TestSimpleCanvasObjectWrapper_Show(t *testing.T) {
	o := &simpleCanvasObject{hidden: true}
	w := &internal.SimpleCanvasObjectWrapper{O: o}

	assert.False(t, o.Visible())
	w.Show()
	assert.True(t, o.Visible())
}

func TestSimpleCanvasObjectWrapper_Size(t *testing.T) {
	o := &simpleCanvasObject{size: fyne.NewSize(12, 34)}
	w := &internal.SimpleCanvasObjectWrapper{O: o}
	assert.Equal(t, fyne.NewSize(12, 34), w.Size())
	o = &simpleCanvasObject{size: fyne.NewSize(13, 42)}
	w = &internal.SimpleCanvasObjectWrapper{O: o}
	assert.Equal(t, fyne.NewSize(13, 42), w.Size())
}

func TestSimpleCanvasObjectWrapper_Visible(t *testing.T) {
	o := &simpleCanvasObject{hidden: true}
	w := &internal.SimpleCanvasObjectWrapper{O: o}
	assert.False(t, w.Visible())
	o = &simpleCanvasObject{hidden: false}
	w = &internal.SimpleCanvasObjectWrapper{O: o}
	assert.True(t, w.Visible())
}

func TestSimpleCanvasObjectWrapper_WrappedObject(t *testing.T) {
	o := &simpleCanvasObject{}
	w := &internal.SimpleCanvasObjectWrapper{O: o}
	assert.Same(t, o, w.WrappedObject())
}

type simpleCanvasObject struct {
	hidden       bool
	minSize      fyne.Size
	pos          fyne.Position
	refreshCount int
	size         fyne.Size
}

func (o *simpleCanvasObject) MinSize() fyne.Size {
	return o.minSize
}

func (o *simpleCanvasObject) Move(position fyne.Position) {
	o.pos = position
}

func (o *simpleCanvasObject) Position() fyne.Position {
	return o.pos
}

func (o *simpleCanvasObject) Resize(size fyne.Size) {
	o.size = size
}

func (o *simpleCanvasObject) Size() fyne.Size {
	return o.size
}

func (o *simpleCanvasObject) Hide() {
	o.hidden = true
}

func (o *simpleCanvasObject) Visible() bool {
	return !o.hidden
}

func (o *simpleCanvasObject) Show() {
	o.hidden = false
}

func (o *simpleCanvasObject) Refresh() {
	o.refreshCount++
}
