package layout

import (
	"testing"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

// NewRectangle returns a new Rectangle instance
func NewMinSizeRect(min fyne.Size) *canvas.Rectangle {
	rect := &canvas.Rectangle{}
	rect.SetMinSize(min)

	return rect
}

func TestSimpleHBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewHBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 50))

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestWideHBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewHBoxLayout(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(308, 100))
	assert.Equal(t, fyne.NewSize(150+(theme.Padding()*2), 100), container.MinSize())

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestTallHBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewHBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj2.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestSpacerHBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewHBoxLayout(), NewSpacer(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(300, 100))
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(200-theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(250, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestMiddleSpacerHBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewHBoxLayout(), obj1, obj2, NewSpacer(), obj3)
	container.Resize(fyne.NewSize(300, 100))
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(250, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestNewHBoxLayout(t *testing.T) {
	lay := NewHBoxLayout()

	assert.Equal(t, true, lay.(*boxLayout).horizontal)
}

func TestSimpleVBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewVBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(50, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestWideVBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewVBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj2.Size(), cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestTallVBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewVBoxLayout(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(100, 308))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestSpacerVBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewVBoxLayout(), NewSpacer(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 200-theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestMiddleSpacerVBoxLayout(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(NewVBoxLayout(), obj1, obj2, NewSpacer(), obj3)
	container.Resize(fyne.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestNewVBoxLayout(t *testing.T) {
	lay := NewVBoxLayout()

	assert.Equal(t, false, lay.(*boxLayout).horizontal)
}
