package driver_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/driver"
	internal_widget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

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

func TestFindObjectAtPositionMatching(t *testing.T) {
	col1cell1 := &objectTree{
		pos:  fyne.NewPos(10, 10),
		size: fyne.NewSize(15, 15),
	}
	col1cell2 := &objectTree{
		pos:  fyne.NewPos(10, 35),
		size: fyne.NewSize(15, 15),
	}
	col1cell3 := &objectTree{
		pos:  fyne.NewPos(10, 60),
		size: fyne.NewSize(15, 15),
	}
	col1 := &objectTree{
		children: []fyne.CanvasObject{col1cell1, col1cell2, col1cell3},
		pos:      fyne.NewPos(10, 10),
		size:     fyne.NewSize(35, 80),
	}
	col2cell1 := &objectTree{
		pos:  fyne.NewPos(10, 10),
		size: fyne.NewSize(15, 15),
	}
	col2cell2 := &objectTree{
		pos:  fyne.NewPos(10, 35),
		size: fyne.NewSize(15, 15),
	}
	col2cell3 := &objectTree{
		pos:  fyne.NewPos(10, 60),
		size: fyne.NewSize(15, 15),
	}
	col2 := &objectTree{
		children: []fyne.CanvasObject{col2cell1, col2cell2, col2cell3},
		pos:      fyne.NewPos(55, 10),
		size:     fyne.NewSize(35, 80),
	}
	colTree := &objectTree{
		children: []fyne.CanvasObject{col1, col2},
		pos:      fyne.NewPos(10, 10),
		size:     fyne.NewSize(100, 100),
	}
	row1cell1 := &objectTree{
		pos:  fyne.NewPos(10, 10),
		size: fyne.NewSize(15, 15),
	}
	row1cell2 := &objectTree{
		pos:  fyne.NewPos(35, 10),
		size: fyne.NewSize(15, 15),
	}
	row1cell3 := &objectTree{
		pos:  fyne.NewPos(60, 10),
		size: fyne.NewSize(15, 15),
	}
	row1 := &objectTree{
		children: []fyne.CanvasObject{row1cell1, row1cell2, row1cell3},
		pos:      fyne.NewPos(10, 10),
		size:     fyne.NewSize(80, 35),
	}
	row2cell1 := &objectTree{
		pos:  fyne.NewPos(10, 10),
		size: fyne.NewSize(15, 15),
	}
	row2cell2 := &objectTree{
		pos:  fyne.NewPos(35, 10),
		size: fyne.NewSize(15, 15),
	}
	row2cell3 := &objectTree{
		pos:  fyne.NewPos(60, 10),
		size: fyne.NewSize(15, 15),
	}
	row2 := &objectTree{
		children: []fyne.CanvasObject{row2cell1, row2cell2, row2cell3},
		pos:      fyne.NewPos(10, 55),
		size:     fyne.NewSize(80, 35),
	}
	rowTree := &objectTree{
		children: []fyne.CanvasObject{row1, row2},
		pos:      fyne.NewPos(10, 10),
		size:     fyne.NewSize(100, 100),
	}
	tree1 := &objectTree{
		pos:  fyne.NewPos(100, 100),
		size: fyne.NewSize(5, 5),
	}
	tree2 := &objectTree{
		pos:  fyne.NewPos(0, 0),
		size: fyne.NewSize(5, 5),
	}
	tree3 := &objectTree{
		pos:  fyne.NewPos(50, 50),
		size: fyne.NewSize(5, 5),
	}
	for name, tt := range map[string]struct {
		matcher    func(object fyne.CanvasObject) bool
		overlay    fyne.CanvasObject
		pos        fyne.Position
		roots      []fyne.CanvasObject
		wantObject fyne.CanvasObject
		wantPos    fyne.Position
		wantLayer  int
	}{
		"match in overlay and roots": {
			matcher:    func(o fyne.CanvasObject) bool { return o.Size().Width == 15 },
			overlay:    colTree,
			pos:        fyne.NewPos(35, 60),
			roots:      []fyne.CanvasObject{rowTree},
			wantObject: col1cell2,
			wantPos:    fyne.NewPos(5, 5),
			wantLayer:  0,
		},
		"match in root but overlay without match present": {
			matcher:    func(o fyne.CanvasObject) bool { return o.Size().Width == 15 },
			overlay:    tree1,
			pos:        fyne.NewPos(35, 60),
			roots:      []fyne.CanvasObject{colTree, rowTree},
			wantObject: nil,
			wantPos:    fyne.Position{},
			wantLayer:  0,
		},
		"match in multiple roots without overlay": {
			matcher:    func(o fyne.CanvasObject) bool { return o.Size().Width == 15 },
			overlay:    nil,
			pos:        fyne.NewPos(83, 83),
			roots:      []fyne.CanvasObject{tree1, rowTree, tree2, colTree},
			wantObject: row2cell3,
			wantPos:    fyne.NewPos(3, 8),
			wantLayer:  2,
		},
		"no match in roots without overlay": {
			matcher:    func(o fyne.CanvasObject) bool { return true },
			overlay:    nil,
			pos:        fyne.NewPos(66, 66),
			roots:      []fyne.CanvasObject{tree1, tree2, tree3},
			wantObject: nil,
			wantPos:    fyne.Position{},
			wantLayer:  3,
		},
		"no overlay and no roots": {
			matcher:    func(o fyne.CanvasObject) bool { return true },
			overlay:    nil,
			pos:        fyne.NewPos(66, 66),
			roots:      nil,
			wantObject: nil,
			wantPos:    fyne.Position{},
			wantLayer:  0,
		},
	} {
		t.Run(name, func(t *testing.T) {
			o, p, l := driver.FindObjectAtPositionMatching(tt.pos, tt.matcher, tt.overlay, tt.roots...)
			assert.Equal(t, tt.wantObject, o, "found object")
			assert.Equal(t, tt.wantPos, p, "position of found object")
			assert.Equal(t, tt.wantLayer, l, "layer of found object (0 - overlay, 1, 2, 3… - roots")
		})
	}
}

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

var _ fyne.Widget = (*objectTree)(nil)

type objectTree struct {
	children []fyne.CanvasObject
	hidden   bool
	pos      fyne.Position
	size     fyne.Size
}

func (o objectTree) Size() fyne.Size {
	return o.size
}

func (o objectTree) Resize(size fyne.Size) {
	o.size = size
}

func (o objectTree) Position() fyne.Position {
	return o.pos
}

func (o objectTree) Move(position fyne.Position) {
	o.pos = position
}

func (o objectTree) MinSize() fyne.Size {
	return o.size
}

func (o objectTree) Visible() bool {
	return !o.hidden
}

func (o objectTree) Show() {
	o.hidden = false
}

func (o objectTree) Hide() {
	o.hidden = true
}

func (o objectTree) Refresh() {
}

func (o objectTree) CreateRenderer() fyne.WidgetRenderer {
	r := &objectTreeRenderer{}
	r.SetObjects(o.children)
	return r
}

var _ fyne.WidgetRenderer = (*objectTreeRenderer)(nil)

type objectTreeRenderer struct {
	internal_widget.BaseRenderer
}

func (o objectTreeRenderer) Layout(_ fyne.Size) {
}

func (o objectTreeRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (o objectTreeRenderer) Refresh() {
}
