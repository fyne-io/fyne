package layout_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestGridLWrapLayout_Layout(t *testing.T) {
	gridSize := fyne.NewSize(125, 125)
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	c := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	c.Resize(gridSize)

	layout.NewGridWrapLayout(cellSize).Layout(c.Objects, gridSize)

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

	c := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}

	layout.NewGridWrapLayout(cellSize).Layout(c.Objects, c.MinSize())

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

	c := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3, obj4},
	}
	c.Resize(gridSize)

	layout.NewGridWrapLayout(cellSize).Layout(c.Objects, gridSize)

	assert.Equal(t, obj1.Size(), cellSize)
	assert.Equal(t, obj3.Position(), fyne.NewPos(50+theme.Padding(), 0))
	assert.Equal(t, obj4.Position(), fyne.NewPos(0, 50+theme.Padding()))
}

func TestGridLWrapLayout_MinSize(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)
	minSize := cellSize

	c := container.NewWithoutLayout(canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}))
	l := layout.NewGridWrapLayout(cellSize)

	layoutMin := l.MinSize(c.Objects)
	assert.Equal(t, minSize, layoutMin)

	// This has a dynamic minSize so we need to check again after layout!
	l.Layout(c.Objects, minSize)
	layoutMin = l.MinSize(c.Objects)
	assert.Equal(t, minSize, layoutMin)

	// vertical 2 items
	c.Add(canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}))
	l.Layout(c.Objects, fyne.NewSize(minSize.Width, minSize.Height*2.5))
	layoutMin = l.MinSize(c.Objects)
	assert.Equal(t, minSize.Height*2+theme.Padding(), layoutMin.Height)

	// horizontal 2 items
	l.Layout(c.Objects, fyne.NewSize(minSize.Width*2.5, minSize.Height))
	layoutMin = l.MinSize(c.Objects)
	assert.Equal(t, minSize.Height, layoutMin.Height)
}

func TestGridLWrapLayout_MinSize_Hidden(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	obj2.Hide()
	obj3 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	c := container.NewWithoutLayout(obj1, obj2, obj3)
	l := layout.NewGridWrapLayout(cellSize)

	layoutMin := l.MinSize(c.Objects)
	assert.Equal(t, fyne.NewSize(50, 50), layoutMin)

	// This has a dynamic minSize so we need to check again after layout!
	l.Layout(c.Objects, fyne.NewSize(50, 75))
	layoutMin = l.MinSize(c.Objects)
	assert.Equal(t, fyne.NewSize(50, 100+theme.Padding()), layoutMin)
}

func TestGridLWrapLayout_Resize_LessThanMinSize(t *testing.T) {
	cellSize := fyne.NewSize(50, 50)
	minSize := cellSize

	c := container.NewWithoutLayout(canvas.NewRectangle(color.NRGBA{0, 0, 0, 0}))
	l := layout.NewGridWrapLayout(cellSize)
	c.Resize(fyne.NewSize(25, 25))

	layoutMin := l.MinSize(c.Objects)
	assert.Equal(t, minSize, layoutMin)
}
