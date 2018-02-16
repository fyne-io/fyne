package layout

import "testing"

import "reflect"
import "image/color"

import "github.com/fyne-io/fyne/ui"

func TestGridLayout(t *testing.T) {
	gridSize := ui.NewSize(100, 100)
	cellSize := ui.NewSize(50, 50)

	obj1 := ui.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := ui.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := ui.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &ui.Container{
		Size:    gridSize,
		Objects: []ui.CanvasObject{obj1, obj2, obj3},
	}

	NewGridLayout(2).Layout(container, gridSize)

	if !reflect.DeepEqual(obj1.Size, cellSize) {
		t.Fatalf("Expected %s but got %s", cellSize, obj1.Size)
	}
	cell2Pos := ui.NewPos(50, 0)
	if !reflect.DeepEqual(obj2.Position, cell2Pos) {
		t.Fatalf("Expected %s but got %s", cell2Pos, obj2.Position)
	}
}
