package layout

import "testing"

import "image/color"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/theme"

import "github.com/stretchr/testify/assert"

func TestGridLayout(t *testing.T) {
	gridSize := ui.NewSize(100+theme.Padding(), 100+theme.Padding())
	cellSize := ui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &ui.Container{
		Size:    gridSize,
		Objects: []ui.CanvasObject{obj1, obj2, obj3},
	}

	NewGridLayout(2).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size, cellSize)
	cell2Pos := ui.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, obj2.Position, cell2Pos)
	cell3Pos := ui.NewPos(0, 50+theme.Padding())
	assert.Equal(t, obj3.Position, cell3Pos)
}

func TestGridLayoutRounding(t *testing.T) {
	gridSize := ui.NewSize(100+theme.Padding()*2, 50)

	obj1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &ui.Container{
		Size:    gridSize,
		Objects: []ui.CanvasObject{obj1, obj2, obj3},
	}

	NewGridLayout(3).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Position, ui.NewPos(0, 0))
	assert.Equal(t, obj1.Size, ui.NewSize(33, 50))
	assert.Equal(t, obj2.Position, ui.NewPos(33+theme.Padding(), 0))
	assert.Equal(t, obj2.Size, ui.NewSize(34, 50))
	assert.Equal(t, obj3.Position, ui.NewPos(67+theme.Padding()*2, 0))
	assert.Equal(t, obj3.Size, ui.NewSize(33, 50))
}

func TestGridLayoutMinSize(t *testing.T) {
	text1 := canvas.NewText("Large Text", color.RGBA{0xff, 0, 0, 0})
	text2 := canvas.NewText("small", color.RGBA{0xff, 0, 0, 0})
	minSize := text1.MinSize().Add(ui.NewSize(0, text1.MinSize().Height+theme.Padding()))

	container := ui.NewContainer(text1, text2)
	layoutMin := NewGridLayout(1).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}
