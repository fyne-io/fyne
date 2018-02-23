package layout

import "testing"

import "reflect"
import "image/color"

import "github.com/fyne-io/fyne/ui"

func TestMaxLayout(t *testing.T) {
	size := ui.NewSize(100, 100)

	obj := ui.NewRectangle(color.RGBA{0, 0, 0, 0})
	container := &ui.Container{
		Size:    size,
		Objects: []ui.CanvasObject{obj},
	}

	NewMaxLayout().Layout(container.Objects, size)

	if !reflect.DeepEqual(obj.Size, size) {
		t.Fatal("Expected", size, "but got", obj.Size)
	}
}

func TestMaxLayoutMinSize(t *testing.T) {
	text := ui.NewText("Padding")
	minSize := text.MinSize()

	container := ui.NewContainer(text)
	layoutMin := NewMaxLayout().MinSize(container.Objects)

	if !reflect.DeepEqual(minSize, layoutMin) {
		t.Fatal("Expected", minSize, "but got", layoutMin)
	}
}
