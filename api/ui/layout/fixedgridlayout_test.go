package layout

import "testing"

import "image/color"

import "github.com/fyne-io/fyne/api/ui"
import "github.com/fyne-io/fyne/api/ui/canvas"
import "github.com/fyne-io/fyne/api/ui/theme"

import "github.com/stretchr/testify/assert"

func TestFixedGridLayout(t *testing.T) {
	gridSize := ui.NewSize(125, 125)
	cellSize := ui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &ui.Container{
		Size:    gridSize,
		Objects: []ui.CanvasObject{obj1, obj2, obj3},
	}

	NewFixedGridLayout(cellSize).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size, cellSize)
	cell2Pos := ui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, obj2.Position, cell2Pos)
	cell3Pos := ui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, obj3.Position, cell3Pos)
}

func TestFixedGridLayoutMinSize(t *testing.T) {
	cellSize := ui.NewSize(50, 50)
	minSize := cellSize

	container := ui.NewContainer(canvas.NewRectangle(color.RGBA{0, 0, 0, 0}))
	layout := NewFixedGridLayout(cellSize)

	layoutMin := layout.MinSize(container.Objects)
	assert.Equal(t, minSize, layoutMin)
}
