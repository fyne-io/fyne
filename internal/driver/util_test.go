package driver_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/driver"
	"fyne.io/fyne/layout"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestWalkVisibleObjectTree(t *testing.T) {
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	child := canvas.NewRectangle(color.Black)
	child.Hide()
	base := widget.NewHBox(rect, child)

	walked := 0
	driver.WalkVisibleObjectTree(base, func(object fyne.CanvasObject, position fyne.Position, clippingPos fyne.Position, clippingSize fyne.Size) bool {
		walked++
		return false
	}, nil)

	assert.Equal(t, 2, walked)
}

func TestWalkWholeObjectTree(t *testing.T) {
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	child := canvas.NewRectangle(color.Black)
	child.Hide()
	base := widget.NewHBox(rect, child)

	walked := 0
	driver.WalkCompleteObjectTree(base, func(object fyne.CanvasObject, position fyne.Position, clippingPos fyne.Position, clippingSize fyne.Size) bool {
		walked++
		return false
	}, nil)

	assert.Equal(t, 3, walked)
}

func TestWalkVisibleObjectTree_Clip(t *testing.T) {
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	child := canvas.NewRectangle(color.Black)
	base := fyne.NewContainerWithLayout(layout.NewGridLayout(1), rect, widget.NewScrollContainer(child))

	clipPos := fyne.NewPos(0, 0)
	clipSize := rect.MinSize()

	driver.WalkVisibleObjectTree(base, func(object fyne.CanvasObject, position fyne.Position, clippingPos fyne.Position, clippingSize fyne.Size) bool {
		if _, ok := object.(*widget.ScrollContainer); ok {
			clipPos = clippingPos
			clipSize = clippingSize
		}
		return false
	}, nil)

	assert.Equal(t, fyne.NewPos(0, 104), clipPos)
	assert.Equal(t, fyne.NewSize(100, 100), clipSize)
}

func TestAbsolutePositionForObject(t *testing.T) {
	t1r1c1 := widget.NewLabel("row 1 col 1")
	t1r1c2 := widget.NewLabel("row 1 col 2")
	t1r2c1 := widget.NewLabel("row 2 col 1")
	t1r2c2 := widget.NewLabel("row 2 col 2")
	t1r2c2.Hide()
	t1r1 := fyne.NewContainer(t1r1c1, t1r1c2)
	t1r2 := fyne.NewContainer(t1r2c1, t1r2c2)
	tree1 := fyne.NewContainer(t1r1, t1r2)

	t1r1c1.Move(fyne.NewPos(111, 111))
	t1r1c2.Move(fyne.NewPos(112, 112))
	t1r2c1.Move(fyne.NewPos(121, 121))
	t1r2c2.Move(fyne.NewPos(122, 122))
	t1r1.Move(fyne.NewPos(11, 11))
	t1r2.Move(fyne.NewPos(12, 12))
	tree1.Move(fyne.NewPos(1, 1))

	t2r1c1 := widget.NewLabel("row 1 col 1")
	t2r1c2 := widget.NewLabel("row 1 col 2")
	t2r2c1 := widget.NewLabel("row 2 col 1")
	t2r2c2 := widget.NewLabel("row 2 col 2")
	t2r1 := fyne.NewContainer(t2r1c1, t2r1c2)
	t2r2 := fyne.NewContainer(t2r2c1, t2r2c2)
	tree2 := fyne.NewContainer(t2r1, t2r2)

	t2r1c1.Move(fyne.NewPos(211, 211))
	t2r1c2.Move(fyne.NewPos(212, 212))
	t2r2c1.Move(fyne.NewPos(221, 221))
	t2r2c2.Move(fyne.NewPos(222, 222))
	t2r1.Move(fyne.NewPos(21, 21))
	t2r2.Move(fyne.NewPos(22, 22))
	tree2.Move(fyne.NewPos(2, 2))

	t3r1 := widget.NewLabel("row 1")
	t3r2 := widget.NewLabel("row 2")
	tree3 := fyne.NewContainer(t3r1, t3r2)

	t3r1.Move(fyne.NewPos(31, 31))
	t3r2.Move(fyne.NewPos(32, 32))
	tree3.Move(fyne.NewPos(3, 3))

	trees := []fyne.CanvasObject{tree1, tree2, tree3}

	outsideTrees := widget.NewLabel("outside trees")
	outsideTrees.Move(fyne.NewPos(10, 10))

	tests := map[string]struct {
		object fyne.CanvasObject
		want   fyne.Position
	}{
		"tree 1: a cell": {
			object: t1r1c2,
			want:   fyne.NewPos(124, 124), // 1 (root) + 11 (row 1) + 112 (cell 2)
		},
		"tree 1: a row": {
			object: t1r2,
			want:   fyne.NewPos(13, 13), // 1 (root) + 12 (row 2)
		},
		"tree 1: root": {
			object: tree1,
			want:   fyne.NewPos(1, 1),
		},
		"tree 1: a hidden element": {
			object: t1r2c2,
			want:   fyne.NewPos(0, 0),
		},

		"tree 2: a row": {
			object: t2r2,
			want:   fyne.NewPos(24, 24), // 2 (root) + 22 (row 2)
		},
		"tree 2: root": {
			object: tree2,
			want:   fyne.NewPos(2, 2),
		},

		"tree 3: a row": {
			object: t3r2,
			want:   fyne.NewPos(35, 35), // 3 (root) + 32 (row 2)
		},
		"tree 3: root": {
			object: tree3,
			want:   fyne.NewPos(3, 3),
		},

		"an object not inside any tree": {
			object: outsideTrees,
			want:   fyne.NewPos(0, 0),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, driver.AbsolutePositionForObject(tt.object, trees))
		})
	}
}
