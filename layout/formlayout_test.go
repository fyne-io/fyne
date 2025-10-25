package layout_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

var minSize fyne.Size

func BenchmarkFormLayout(b *testing.B) {

	min := fyne.Size{}
	form := layout.NewFormLayout()
	label1 := canvas.NewRectangle(color.Black)
	content1 := canvas.NewRectangle(color.Black)
	label2 := canvas.NewRectangle(color.Black)
	content2 := canvas.NewRectangle(color.Black)

	objects := []fyne.CanvasObject{label1, content1, label2, content2}

	b.ReportAllocs()
	for b.Loop() {
		min = form.MinSize(objects)
	}

	minSize = min
}

func TestFormLayout(t *testing.T) {
	gridSize := fyne.NewSize(125, 125)

	label1 := canvas.NewRectangle(color.Black)
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.Black)
	content1.SetMinSize(fyne.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.Black)
	label2.SetMinSize(fyne.NewSize(70, 30))
	content2 := canvas.NewRectangle(color.Black)
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

func TestFormLayout_Text(t *testing.T) {
	size := fyne.NewSize(120, 50)
	label := canvas.NewText("Label", color.Black)
	content := canvas.NewText("Content", color.Black)

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label, content},
	}
	container.Resize(size)
	layout.NewFormLayout().Layout(container.Objects, size)

	assert.Equal(t, label.Size().Height, content.Size().Height)
}

func TestFormLayout_Hidden(t *testing.T) {
	gridSize := fyne.NewSize(190+theme.Padding(), 125)

	label1 := canvas.NewRectangle(color.Black)
	label1.SetMinSize(fyne.NewSize(70, 50))
	label1.Hide()
	content1 := canvas.NewRectangle(color.Black)
	content1.SetMinSize(fyne.NewSize(120, 100))
	content1.Hide()

	label2 := canvas.NewRectangle(color.Black)
	label2.SetMinSize(fyne.NewSize(50, 30))
	content2 := canvas.NewRectangle(color.Black)
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

	label1 := canvas.NewRectangle(color.Black)
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.Black)
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
	label1 := canvas.NewRectangle(color.Black)
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.Black)
	content1.SetMinSize(fyne.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.Black)
	label2.SetMinSize(fyne.NewSize(70, 30))
	content2 := canvas.NewRectangle(color.Black)
	content2.SetMinSize(fyne.NewSize(120, 80))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{label1, content1, label2, content2},
	}

	l := layout.NewFormLayout()
	layoutMin := l.MinSize(container.Objects)
	expectedRowWidth := 70 + 120 + theme.Padding()
	expectedRowHeight := 100 + 80 + theme.Padding()
	assert.Equal(t, fyne.NewSize(expectedRowWidth, expectedRowHeight), layoutMin)

	text := canvas.NewText("Text", color.Black)
	value := widget.NewLabel("Text")
	l = layout.NewFormLayout()
	layoutMin = l.MinSize([]fyne.CanvasObject{text, value})
	// check that the text minsize is padded to match a label
	assert.Equal(t, value.MinSize().Width*2+theme.Padding(), layoutMin.Width)
}

func TestFormLayout_MinSize_Hidden(t *testing.T) {
	label1 := canvas.NewRectangle(color.Black)
	label1.SetMinSize(fyne.NewSize(50, 50))
	content1 := canvas.NewRectangle(color.Black)
	content1.SetMinSize(fyne.NewSize(100, 100))

	label2 := canvas.NewRectangle(color.Black)
	label2.SetMinSize(fyne.NewSize(70, 30))
	label2.Hide()
	content2 := canvas.NewRectangle(color.Black)
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

func TestFormLayout_MinSize_CanvasText_SingleRow(t *testing.T) {
	text1 := canvas.NewText("First text", color.Black)
	text2 := canvas.NewText("Second text", color.Black)

	l := layout.NewFormLayout()
	layoutMin := l.MinSize([]fyne.CanvasObject{text1, text2})

	inner := theme.InnerPadding()
	min1 := text1.MinSize()
	min2 := text2.MinSize()
	expectedWidth := (min1.Width + inner*2) + (min2.Width + inner*2) + theme.Padding()

	expectedHeight := fyne.Max(min1.Height, min2.Height) + inner*2

	assert.Equal(t, fyne.NewSize(expectedWidth, expectedHeight), layoutMin)
}

func TestFormLayout_MinSize_CanvasText_TwoRows(t *testing.T) {
	label1 := canvas.NewText("First Text", color.Black)
	value1 := canvas.NewText("First Value", color.Black)
	label2 := canvas.NewText("Second Text", color.Black)
	value2 := canvas.NewText("Second Value", color.Black)

	l := layout.NewFormLayout()
	layoutMin := l.MinSize([]fyne.CanvasObject{label1, value1, label2, value2})

	inner := theme.InnerPadding()
	l1 := label1.MinSize()
	l2 := label2.MinSize()
	v1 := value1.MinSize()
	v2 := value2.MinSize()
	labelCol := fyne.Max(l1.Width+inner*2, l2.Width+inner*2)
	valueCol := fyne.Max(v1.Width+inner*2, v2.Width+inner*2)
	expectedWidth := labelCol + valueCol + theme.Padding()
	row1 := fyne.Max(l1.Height+inner*2, v1.Height+inner*2)
	row2 := fyne.Max(l2.Height+inner*2, v2.Height+inner*2)
	expectedHeight := row1 + row2 + theme.Padding()

	assert.Equal(t, fyne.NewSize(expectedWidth, expectedHeight), layoutMin)
}
