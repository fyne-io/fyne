package layout

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/theme"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MinSizeRect struct {
	Size     fyne.Size
	Position fyne.Position

	minSize fyne.Size
}

// CurrentSize returns the current size of this rectangle object
func (r *MinSizeRect) CurrentSize() fyne.Size {
	return r.Size
}

// Resize sets a new size for the rectangle object
func (r *MinSizeRect) Resize(size fyne.Size) {
	r.Size = size
}

// CurrentPosition gets the current position of this rectangle object, relative to it's parent / canvas
func (r *MinSizeRect) CurrentPosition() fyne.Position {
	return r.Position
}

// Move the rectangle object to a new position, relative to it's parent / canvas
func (r *MinSizeRect) Move(pos fyne.Position) {
	r.Position = pos
}

// MinSize for a Rectangle simply returns Size{1, 1} as there is no
// explicit content
func (r *MinSizeRect) MinSize() fyne.Size {
	return r.minSize
}

// NewRectangle returns a new Rectangle instance
func NewMinSizeRect(min fyne.Size) *MinSizeRect {
	return &MinSizeRect{
		minSize: min,
	}
}

func TestSimpleListLayout(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewListLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(50, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size, cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position)
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position)
}

func TestWideListLayout(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewListLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size, cellSize)
	assert.Equal(t, obj2.Size, cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position)
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position)
}

func TestTallListLayout(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewListLayout(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(100, 308))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size.Height, 50)
	assert.Equal(t, obj2.Size.Height, 50)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position)
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position)
}

func TestSpacerListLayout(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewListLayout(), fyne.NewSpacer(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size.Height, 50)
	assert.Equal(t, obj2.Size.Height, 50)
	cell2Pos := fyne.NewPos(0, 200-theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position)
	cell3Pos := fyne.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position)
}

func TestMiddleSpacerListLayout(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewListLayout(), obj1, obj2, fyne.NewSpacer(), obj3)
	container.Resize(fyne.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size.Height, 50)
	assert.Equal(t, obj2.Size.Height, 50)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position)
	cell3Pos := fyne.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position)
}
