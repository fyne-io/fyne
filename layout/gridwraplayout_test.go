package layout_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestGridLWrapLayout_Layout(t *testing.T) {
	gridSize := fyne.NewSize(125, 125)
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(gridSize)

	layout.NewGridWrapLayout(cellSize).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(50+theme.Padding(), 0)
	assert.Equal(t, obj2.Position(), cell2Pos)
	cell3Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, obj3.Position(), cell3Pos)
}

func TestGridLWrapLayout_Layout_Min(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}

	layout.NewGridWrapLayout(cellSize).Layout(container.Objects, container.MinSize())

	assert.Equal(t, obj1.Size(), cellSize)
	cell2Pos := fyne.NewPos(0, 50+theme.Padding())
	assert.Equal(t, obj2.Position(), cell2Pos)
	cell3Pos := fyne.NewPos(0, 100+theme.Padding()*2)
	assert.Equal(t, obj3.Position(), cell3Pos)
}

func TestGridLWrapLayout_Layout_HiddenItem(t *testing.T) {
	gridSize := fyne.NewSize(125, 125)
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2.Hide()
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj4 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3, obj4},
	}
	container.Resize(gridSize)

	layout.NewGridWrapLayout(cellSize).Layout(container.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj3.Position(), fyne.NewPos(50+theme.Padding(), 0))
	assert.Equal(t, obj4.Position(), fyne.NewPos(0, 50+theme.Padding()))
}

func TestGridLWrapLayout_MinSize(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)
	minSize := cellSize

	container := fyne.NewContainer(canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}))
	layout := layout.NewGridWrapLayout(cellSize)

	layoutMin := layout.MinSize(container.Objects)
	assert.Equal(t, minSize, layoutMin)

	// This has a dynamic minSize so we need to check again after layout!
	layout.Layout(container.Objects, minSize)
	layoutMin = layout.MinSize(container.Objects)
	assert.Equal(t, minSize, layoutMin)
}

func TestGridLWrapLayout_MinSize_Hidden(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2.Hide()
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	container := fyne.NewContainer(obj1, obj2, obj3)
	layout := layout.NewGridWrapLayout(cellSize)

	layoutMin := layout.MinSize(container.Objects)
	assert.Equal(t, fyne.NewSize(50, 50), layoutMin)

	// This has a dynamic minSize so we need to check again after layout!
	layout.Layout(container.Objects, fyne.NewSize(50, 75))
	layoutMin = layout.MinSize(container.Objects)
	assert.Equal(t, fyne.NewSize(50, 100+theme.Padding()), layoutMin)
}

func TestGridLWrapLayout_Resize_LessThanMinSize(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)
	minSize := cellSize

	container := fyne.NewContainer(canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}))
	l := layout.NewGridWrapLayout(cellSize)
	container.Resize(fyne.NewSize(25, 25))

	layoutMin := l.MinSize(container.Objects)
	assert.Equal(t, minSize, layoutMin)
}
