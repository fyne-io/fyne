package layout_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/layout"
	publicLayout "fyne.io/fyne/layout"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NewRectangle returns a new Rectangle instance
func NewMinSizeRect(min fyne.Size) *canvas.Rectangle {
	rect := &canvas.Rectangle{}
	rect.SetMinSize(min)

	return rect
}

func TestBox_HorizontalSimple(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{Horizontal: true}, obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 50))

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_HorizontalHiddenItem(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj2.Hide()
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{Horizontal: true}, obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(100+(theme.Padding()), 50))

	assert.Equal(t, obj1.Size(), cellSize)
	cell3Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_HorizontalWide(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{Horizontal: true}, obj1, obj2, obj3)
	container.Resize(fyne.NewSize(308, 100))
	assert.Equal(t, fyne.NewSize(150+(theme.Padding()*2), 100), container.MinSize())

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_HorizontalTall(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{Horizontal: true}, obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj2.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(100+theme.Padding()*2, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_HorizontalSpacer(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{Horizontal: true}, publicLayout.NewSpacer(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(300, 100))
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(200-theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(250, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_HorizontalMiddleSpacer(t *testing.T) {
	cellSize := fyne.NewSize(50, 100)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(0, 25)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{Horizontal: true}, obj1, obj2, publicLayout.NewSpacer(), obj3)
	container.Resize(fyne.NewSize(300, 100))
	assert.Equal(t, container.MinSize(), fyne.NewSize(150+(theme.Padding()*2), 100))

	assert.Equal(t, 50, obj1.Size().Width)
	assert.Equal(t, 50, obj2.Size().Width)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(250, 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_HorizontalPadBeforeAndAfter(t *testing.T) {
	require.Greater(t, theme.Padding(), 0)
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{Horizontal: true, PadBeforeAndAfter: true}, obj1, obj2)
	assert.Equal(t, container.MinSize(), fyne.NewSize(100+theme.Padding()+theme.Padding()*2, 50))

	assert.Equal(t, fyne.NewPos(theme.Padding(), 0), obj1.Position())
	assert.Equal(t, cellSize, obj1.Size())
	assert.Equal(t, fyne.NewPos(cellSize.Width+theme.Padding()*2, 0), obj2.Position())
	assert.Equal(t, cellSize, obj2.Size())
}

func TestBox_VerticalSimple(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{}, obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(50, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_VerticalHiddenItem(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)
	obj2.Hide()
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{}, obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(50, 100+(theme.Padding())))

	assert.Equal(t, obj1.Size(), cellSize)
	cell3Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_VerticalWide(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{}, obj1, obj2, obj3)
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj2.Size(), cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_VerticalTall(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{}, obj1, obj2, obj3)
	container.Resize(fyne.NewSize(100, 308))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_VerticalSpacer(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{}, publicLayout.NewSpacer(), obj1, obj2, obj3)
	container.Resize(fyne.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 200-theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_VerticalMiddleSpacer(t *testing.T) {
	cellSize := fyne.NewSize(100, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize.Subtract(fyne.NewSize(25, 0)))
	obj3 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{}, obj1, obj2, publicLayout.NewSpacer(), obj3)
	container.Resize(fyne.NewSize(100, 300))
	assert.Equal(t, container.MinSize(), fyne.NewSize(100, 150+(theme.Padding()*2)))

	assert.Equal(t, 50, obj1.Size().Height)
	assert.Equal(t, 50, obj2.Size().Height)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 250)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestBox_VerticalPadBeforeAndAfter(t *testing.T) {
	require.Greater(t, theme.Padding(), 0)
	cellSize := fyne.NewSize(50, 50)

	obj1 := NewMinSizeRect(cellSize)
	obj2 := NewMinSizeRect(cellSize)

	container := fyne.NewContainerWithLayout(&layout.Box{PadBeforeAndAfter: true}, obj1, obj2)
	assert.Equal(t, container.MinSize(), fyne.NewSize(50, 100+theme.Padding()+theme.Padding()*2))

	assert.Equal(t, fyne.NewPos(0, theme.Padding()), obj1.Position())
	assert.Equal(t, cellSize, obj1.Size())
	assert.Equal(t, fyne.NewPos(0, cellSize.Width+theme.Padding()*2), obj2.Position())
	assert.Equal(t, cellSize, obj2.Size())
}
