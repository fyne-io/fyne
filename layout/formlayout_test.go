package layout

import (
	"image/color"
	"testing"

	"fyne.io/fyne/theme"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/stretchr/testify/assert"
)

func TestFormLayout(t *testing.T) {
	gridSize := fyne.NewSize(125, 125)

	label1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	content1.SetMinSize(fyne.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	label2.SetMinSize(fyne.NewSize(70, 30))
	content2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	content2.SetMinSize(fyne.NewSize(120, 80))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label1, content1, label2, content2},
	}
	container.Resize(gridSize)

	NewFormLayout().Layout(container.Objects, gridSize)

	assert.Equal(t, fyne.NewSize(70, 100), label1.Size())
	assert.Equal(t, fyne.NewSize(120, 100), content1.Size())
	assert.Equal(t, fyne.NewSize(70, 80), label2.Size())
	assert.Equal(t, fyne.NewSize(120, 80), content2.Size())

}

func TestFormLayoutMinSize(t *testing.T) {

	label1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	content1.SetMinSize(fyne.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	label2.SetMinSize(fyne.NewSize(70, 30))
	content2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	content2.SetMinSize(fyne.NewSize(120, 80))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label1, content1, label2, content2},
	}

	layout := NewFormLayout()
	layoutMin := layout.MinSize(container.Objects)
	expectedRowWidth := 70 + 120 + theme.Padding()
	expectedRowHeight := 100 + 80 + theme.Padding()
	assert.Equal(t, fyne.NewSize(expectedRowWidth, expectedRowHeight), layoutMin)
}
