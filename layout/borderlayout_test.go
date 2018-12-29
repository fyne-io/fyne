package layout

import (
	"image/color"
	"testing"

	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestBorderLayoutEmpty(t *testing.T) {
	size := fyne.NewSize(100, 100)

	obj := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj},
	}
	container.Resize(size)

	NewBorderLayout(nil, nil, nil, nil).Layout(container.Objects, size)

	assert.Equal(t, obj.Size(), size)
}

func TestBorderLayoutTopBottom(t *testing.T) {
	size := fyne.NewSize(100, 100)

	obj1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(size)

	NewBorderLayout(obj1, obj2, nil, nil).Layout(container.Objects, size)

	innerSize := fyne.NewSize(size.Width, size.Height-obj1.Size().Height-obj2.Size().Height-theme.Padding()*2)
	assert.Equal(t, innerSize, obj3.Size())
	assert.Equal(t, fyne.NewPos(0, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(0, size.Height-obj2.Size().Height), obj2.Position())
	assert.Equal(t, fyne.NewPos(0, obj1.Size().Height+theme.Padding()), obj3.Position())
}

func TestBorderLayoutLeftRight(t *testing.T) {
	size := fyne.NewSize(100, 100)

	obj1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(size)

	NewBorderLayout(nil, nil, obj1, obj2).Layout(container.Objects, size)

	innerSize := fyne.NewSize(size.Width-obj1.Size().Width-obj2.Size().Width-theme.Padding()*2, size.Height)
	assert.Equal(t, innerSize, obj3.Size())
	assert.Equal(t, fyne.NewPos(0, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(size.Width-obj2.Size().Width, 0), obj2.Position())
	assert.Equal(t, fyne.NewPos(obj1.Size().Width+theme.Padding(), 0), obj3.Position())
}

func TestBorderCenterLayoutMinSize(t *testing.T) {
	text := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := text.MinSize()

	container := fyne.NewContainer(text)
	layoutMin := NewBorderLayout(nil, nil, nil, nil).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestBorderTopBottomLayoutMinSize(t *testing.T) {
	text1 := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	text2 := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	text3 := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := fyne.NewSize(text3.MinSize().Width, text1.MinSize().Height+text2.MinSize().Height+text3.MinSize().Height+theme.Padding()*2)

	container := fyne.NewContainer(text1, text2, text3)
	layoutMin := NewBorderLayout(text1, text2, nil, nil).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestBorderTopOnlyLayoutMinSize(t *testing.T) {
	text1 := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := fyne.NewSize(text1.MinSize().Width, text1.MinSize().Height+theme.Padding())

	container := fyne.NewContainer(text1)
	layoutMin := NewBorderLayout(text1, nil, nil, nil).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestBorderLeftRightLayoutMinSize(t *testing.T) {
	text1 := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	text2 := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	text3 := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := fyne.NewSize(text1.MinSize().Width+text2.MinSize().Width+text3.MinSize().Width+theme.Padding()*2, text3.MinSize().Height)

	container := fyne.NewContainer(text1, text2, text3)
	layoutMin := NewBorderLayout(nil, nil, text1, text2).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}

func TestBorderLeftOnlyLayoutMinSize(t *testing.T) {
	text1 := canvas.NewText("Padding", color.RGBA{0, 0xff, 0, 0})
	minSize := fyne.NewSize(text1.MinSize().Width+theme.Padding(), text1.MinSize().Height)

	container := fyne.NewContainer(text1)
	layoutMin := NewBorderLayout(nil, nil, text1, nil).MinSize(container.Objects)

	assert.Equal(t, minSize, layoutMin)
}
