package layout

import "testing"

import "reflect"
import "image/color"

import "github.com/fyne-io/fyne/ui"
import "github.com/fyne-io/fyne/ui/canvas"
import "github.com/fyne-io/fyne/ui/theme"

func TestFixedGridLayout(t *testing.T) {
	gridSize := ui.NewSize(125, 125)
	cellSize := ui.NewSize(50, 50)

	obj1 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj2 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	obj3 := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	container := &ui.Container{
		Size:    gridSize,
		Objects: []ui.CanvasObject{obj1, obj2, obj3},
	}

	NewFixedGridLayout(cellSize).Layout(container.Objects, gridSize)

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

func TestFixedGridLayoutMinSize(t *testing.T) {
	cellSize := ui.NewSize(50, 50)
	minSize := cellSize

	container := ui.NewContainer(canvas.NewRectangle(color.RGBA{0, 0, 0, 0}))
	layout := NewFixedGridLayout(cellSize)

	layoutMin := layout.MinSize(container.Objects)
	if !reflect.DeepEqual(minSize, layoutMin) {
		t.Fatal("Expected", minSize, "but got", layoutMin)
	}
}
