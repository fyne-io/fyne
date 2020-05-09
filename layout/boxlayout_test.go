package layout_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

// NewRectangle returns a new Rectangle instance
func NewMinSizeRect(min fyne.Size) *canvas.Rectangle {
	rect := &canvas.Rectangle{}
	rect.SetMinSize(min)

	return rect
}

func TestHBoxLayout_Simple(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 50))

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_HiddenItem(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj2.Hide()
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(100+(theme.Padding()), 50))

	assert.Equal(t, obj1.Size(), cellSize)
	cell3Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_Wide(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(308, 100))
	assert.Equal(t, fyne.NewSize(150+(theme.Padding()*2), 100), container.MinSize())

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_Tall(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj2.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_Spacer(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(300, 100))
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(200-theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(250, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestHBoxLayout_MiddleSpacer(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), obj1, obj2, layout.NewSpacer(), obj3)
	container.Resize(fyne.NewSize(300, 100))
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(250, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_Simple(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(50, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_HiddenItem(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj2.Hide()
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(50, 100+(theme.Padding())))

	assert.Equal(t, obj1.Size(), cellSize)
	cell3Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_Wide(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj2.Size(), cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_Tall(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(100, 308))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_Spacer(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), layout.NewSpacer(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 200-theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestVBoxLayout_MiddleSpacer(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), obj1, obj2, layout.NewSpacer(), obj3)
	container.Resize(fyne.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position())
}
