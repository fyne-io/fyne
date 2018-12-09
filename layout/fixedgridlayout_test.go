package layout

import "testing"

import "image/color"

import "github.com/fyne-io/fyne"
import "github.com/fyne-io/fyne/canvas"
import "github.com/fyne-io/fyne/theme"

import "github.com/stretchr/testify/assert"

func TestFixedGridLayout(t *testing.T) {
	gridSize := fyne.NewSize(125, 125)
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	NewFixedGridLayout(cellSize).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, obj2.Position(), cell2Pos)
	cell3Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, obj3.Position(), cell3Pos)
}

func TestFixedGridLayoutMinSize(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)
	minSize := cellSize

	container := fyne.NewContainer(canvas.NewRectangle(color.RGBA{0, 0, 0, 0}))
	layout := NewFixedGridLayout(cellSize)

	layoutMin := layout.MinSize(container.Objects)
	assert.Equal(t, minSize, layoutMin)
}
