package layout_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestGridLayout(t *testing.T) {
	gridSize := fyne.NewSize(100+theme.Padding(), 100+theme.Padding())
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	layout.NewGridLayout(2).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestGridLayoutRounding(t *testing.T) {
	gridSize := fyne.NewSize(100+theme.Padding()*2, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	layout.NewGridLayout(3).Layout(container.Objects, gridSize)

	assert.Equal(t, fyne.NewPos(0, 0), obj1.Position())
	assert.Equal(t, fyne.NewSize(33, 50), obj1.Size())
	assert.Equal(t, fyne.NewPos(33+theme.Padding(), 0), obj2.Position())
	assert.Equal(t, fyne.NewSize(34, 50), obj2.Size())
	assert.Equal(t, fyne.NewPos(67+theme.Padding()*2, 0), obj3.Position())
	assert.Equal(t, fyne.NewSize(33, 50), obj3.Size())
}

func TestGridLayout_Vertical(t *testing.T) {
	gridSize := fyne.NewSize(100+theme.Padding(), 100+theme.Padding())
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	layout.NewGridLayoutWithRows(2).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, cell2Pos, obj2.Position())
	cell3Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, cell3Pos, obj3.Position())
}

func TestGridLayout_MinSize(t *testing.T) {
	text1 := canvas.NewText("Large Text", color.NRGBA{0xff, 0, 0, 0})
	text2 := canvas.NewText("small", color.NRGBA{0xff, 0, 0, 0})
	minSize := text1.MinSize().Add(fyne.NewSize(0, text2.MinSize().Height+theme.Padding()))

	container := fyne.NewContainer(text1, text2)
	layoutMin := layout.NewGridLayout(1).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestGridLayout_MinSize_Vertical(t *testing.T) {
	text1 := canvas.NewText("Text", color.NRGBA{0xff, 0, 0, 0})
	text2 := canvas.NewText("Text", color.NRGBA{0xff, 0, 0, 0})
	minSize := text1.MinSize().Add(fyne.NewSize(text2.MinSize().Width+theme.Padding(), 0))

	container := fyne.NewContainer(text1, text2)
	layoutMin := layout.NewGridLayoutWithRows(1).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestGridLayout_MinSize_HiddenItem(t *testing.T) {
	text1 := canvas.NewText("Large Text", color.NRGBA{0xff, 0, 0, 0})
	text2 := canvas.NewText("hidden", color.NRGBA{0xff, 0, 0, 0})
	text2.Hide()
	text3 := canvas.NewText("small", color.NRGBA{0xff, 0, 0, 0})
	minSize := text1.MinSize().Add(fyne.NewSize(0, text3.MinSize().Height+theme.Padding()))

	container := fyne.NewContainer(text1, text2, text3)
	layoutMin := layout.NewGridLayout(1).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}
