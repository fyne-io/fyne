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

func TestFormLayout(t *testing.T) {
	gridSize := fyne.NewSize(125, 125)

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(fyne.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label2.SetMinSize(fyne.NewSize(70, 30))
	content2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content2.SetMinSize(fyne.NewSize(120, 80))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label1, content1, label2, content2},
	}
	container.Resize(gridSize)

	layout.NewFormLayout().Layout(container.Objects, gridSize)

	assert.Equal(t, fyne.NewSize(70, 100), label1.Size())
	assert.Equal(t, fyne.NewSize(120, 100), content1.Size())
	assert.Equal(t, fyne.NewSize(70, 80), label2.Size())
	assert.Equal(t, fyne.NewSize(120, 80), content2.Size())
}

func TestFormLayout_Hidden(t *testing.T) {
	gridSize := fyne.NewSize(190+theme.Padding(), 125)

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(fyne.NewSize(70, 50))
	label1.Hide()
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(fyne.NewSize(120, 100))
	content1.Hide()

	label2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label2.SetMinSize(fyne.NewSize(50, 30))
	content2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content2.SetMinSize(fyne.NewSize(100, 80))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label1, content1, label2, content2},
	}
	container.Resize(gridSize)

	layout.NewFormLayout().Layout(container.Objects, gridSize)

	assert.Equal(t, fyne.NewSize(50, 80), label2.Size())
	assert.Equal(t, fyne.NewSize(140, 80), content2.Size())
	assert.Equal(t, fyne.NewPos(0, 0), label2.Position())
	assert.Equal(t, fyne.NewPos(50+theme.Padding(), 0), content2.Position())
}

func TestFormLayout_StretchX(t *testing.T) {
	wideSize := fyne.NewSize(150, 50)

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(fyne.NewSize(50, 50))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label1, content1},
	}
	container.Resize(wideSize)

	layout.NewFormLayout().Layout(container.Objects, wideSize)

	assert.Equal(t, fyne.NewSize(50, 50), label1.Size())
	assert.Equal(t, fyne.NewSize(wideSize.Width-50-theme.Padding(), 50), content1.Size())
}

func TestFormLayout_MinSize(t *testing.T) {

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(fyne.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label2.SetMinSize(fyne.NewSize(70, 30))
	content2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content2.SetMinSize(fyne.NewSize(120, 80))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label1, content1, label2, content2},
	}

	l := layout.NewFormLayout()
	layoutMin := l.MinSize(container.Objects)
	expectedRowWidth := 70 + 120 + theme.Padding()
	expectedRowHeight := 100 + 80 + theme.Padding()
	assert.Equal(t, fyne.NewSize(expectedRowWidth, expectedRowHeight), layoutMin)
}

func TestFormLayout_MinSize_Hidden(t *testing.T) {

	label1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content1.SetMinSize(fyne.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	label2.SetMinSize(fyne.NewSize(70, 30))
	label2.Hide()
	content2 := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	content2.SetMinSize(fyne.NewSize(120, 80))
	content2.Hide()

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label1, content1, label2, content2},
	}

	l := layout.NewFormLayout()
	layoutMin := l.MinSize(container.Objects)
	expectedRowWidth := 50 + 100 + theme.Padding()
	expectedRowHeight := float32(100)
	assert.Equal(t, fyne.NewSize(expectedRowWidth, expectedRowHeight), layoutMin)
}
