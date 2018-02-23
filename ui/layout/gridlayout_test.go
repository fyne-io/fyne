package layout

import "testing"

import "reflect"
import "image/color"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/theme"

func TestGridLayout(t *testing.T) {
	gridSize := ui.NewSize(100+theme.Padding(), 100+theme.Padding())
	cellSize := ui.NewSize(50, 50)

	obj1 := ui.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := ui.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := ui.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &ui.Container{
		Size:    gridSize,
		Objects: []ui.CanvasObject{obj1, obj2, obj3},
	}

	NewGridLayout(2).Layout(container.Objects, gridSize)

	if !reflect.DeepEqual(obj1.Size, cellSize) {
		t.Fatal("Expected", cellSize, "but got", obj1.Size)
	}
	cell2Pos := ui.NewPos(50+theme.Padding(), 0)
	if !reflect.DeepEqual(obj2.Position, cell2Pos) {
		t.Fatal("Expected", cell2Pos, "but got", obj2.Position)
	}
	cell3Pos := ui.NewPos(0, 50+theme.Padding())
	if !reflect.DeepEqual(obj3.Position, cell3Pos) {
		t.Fatal("Expected", cell3Pos, "but got", obj3.Position)
	}
}

func TestGridLayoutMinSize(t *testing.T) {
	text1 := ui.NewText("Large Text")
	text2 := ui.NewText("small")
	minSize := text1.MinSize().Add(ui.NewSize(0, text1.MinSize().Height+theme.Padding()))

	container := ui.NewContainer(text1, text2)
	layoutMin := NewGridLayout(1).MinSize(container.Objects)

	if !reflect.DeepEqual(minSize, layoutMin) {
		t.Fatal("Expected", minSize, "but got", layoutMin)
	}
}
