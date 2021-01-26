package layout_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestFlexLayout_MainAxisAlignment(t *testing.T) {
	lsize := fyne.NewSize(200, 100)

	obj1 := NewMinSizeRect(fyne.NewSize(20, 40))
	obj2 := NewMinSizeRect(fyne.NewSize(40, 20))
	obj3 := NewMinSizeRect(fyne.NewSize(20, 60))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(lsize)
	alignment := &fyne.AxisAlignment{
		CrossAxisAlignment: fyne.CrossAxisAlignmentStart,
	}

	// -- Start alignment
	alignment.MainAxisAlignment = fyne.MainAxisAlignmentStart
	l := layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	assert.Equal(t, fyne.NewPos(0, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(obj1.Position().X+obj1.MinSize().Width, 0), obj2.Position())
	assert.Equal(t, fyne.NewPos(obj2.Position().X+obj2.MinSize().Width, 0), obj3.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)

	// -- End alignment
	alignment.MainAxisAlignment = fyne.MainAxisAlignmentEnd
	l = layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	assert.Equal(t, fyne.NewPos(lsize.Width-obj3.MinSize().Width, 0), obj3.Position())
	assert.Equal(t, fyne.NewPos(obj3.Position().X-obj2.MinSize().Width, 0), obj2.Position())
	assert.Equal(t, fyne.NewPos(obj2.Position().X-obj1.MinSize().Width, 0), obj1.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)

	// -- Center alignment
	// | x | obj1 | obj2 | obj3 | x |
	alignment.MainAxisAlignment = fyne.MainAxisAlignmentCenter
	l = layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	x := (lsize.Width - obj1.MinSize().Width - obj2.MinSize().Width - obj3.MinSize().Width) / 2
	assert.Equal(t, fyne.NewPos(x, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(obj1.Position().X+obj1.MinSize().Width, 0), obj2.Position())
	assert.Equal(t, fyne.NewPos(obj2.Position().X+obj2.MinSize().Width, 0), obj3.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)

	// -- SpaceBetween alignment
	// Place the free space evenly between the children.
	// | obj1 | x | obj2 | x | obj3 |
	alignment.MainAxisAlignment = fyne.MainAxisAlignmentSpaceBetween
	l = layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	x = (lsize.Width - obj1.MinSize().Width - obj2.MinSize().Width - obj3.MinSize().Width) / 2
	assert.Equal(t, fyne.NewPos(0, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(obj1.MinSize().Width+x, 0), obj2.Position())
	assert.Equal(t, fyne.NewPos(obj2.Position().X+obj2.MinSize().Width+x, 0), obj3.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)

	// -- SpaceAround alignment
	// Place the free space evenly between the children as well as half of that
	// space before and after the first and last child.
	// | x/2 | obj1 | x | obj2 | x | obj3 | x/2 |
	alignment.MainAxisAlignment = fyne.MainAxisAlignmentSpaceAround
	l = layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	x = (lsize.Width - obj1.MinSize().Width - obj2.MinSize().Width - obj3.MinSize().Width) / 3
	assert.Equal(t, fyne.NewPos(x/2, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(obj1.Position().X+obj1.MinSize().Width+x, 0), obj2.Position())
	assert.Equal(t, fyne.NewPos(obj2.Position().X+obj2.MinSize().Width+x, 0), obj3.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)

	// -- SpaceEvenly alignment
	// Place the free space evenly between the children as well as before and
	// after the first and last child.
	// | x | obj1 | x | obj2 | x | obj3 | x |
	alignment.MainAxisAlignment = fyne.MainAxisAlignmentSpaceEvenly
	l = layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	x = (lsize.Width - obj1.MinSize().Width - obj2.MinSize().Width - obj3.MinSize().Width) / 4
	assert.Equal(t, fyne.NewPos(x, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(obj1.Position().X+obj1.MinSize().Width+x, 0), obj2.Position())
	assert.Equal(t, fyne.NewPos(obj2.Position().X+obj2.MinSize().Width+x, 0), obj3.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)
}

func TestFlexLayout_CrossAxisAlignment(t *testing.T) {
	lsize := fyne.NewSize(200, 100)

	obj1 := NewMinSizeRect(fyne.NewSize(20, 40))
	obj2 := NewMinSizeRect(fyne.NewSize(40, 20))
	obj3 := NewMinSizeRect(fyne.NewSize(20, 60))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(lsize)
	crossSize := fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height))
	alignment := &fyne.AxisAlignment{
		MainAxisAlignment: fyne.MainAxisAlignmentStart,
	}

	// -- Start alignment
	alignment.CrossAxisAlignment = fyne.CrossAxisAlignmentStart
	l := layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	assert.Equal(t, fyne.NewPos(0, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(obj1.Position().X+obj1.MinSize().Width, 0), obj2.Position())
	assert.Equal(t, fyne.NewPos(obj2.Position().X+obj2.MinSize().Width, 0), obj3.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)

	// -- End alignment
	alignment.CrossAxisAlignment = fyne.CrossAxisAlignmentEnd
	l = layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	assert.Equal(t, fyne.NewPos(0, crossSize-obj1.MinSize().Height), obj1.Position())
	assert.Equal(t, fyne.NewPos(obj1.MinSize().Width, crossSize-obj2.MinSize().Height),
		obj2.Position())
	assert.Equal(t, fyne.NewPos(obj2.Position().X+obj2.MinSize().Width, crossSize-obj3.MinSize().Height),
		obj3.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)

	// -- Center alignment
	alignment.CrossAxisAlignment = fyne.CrossAxisAlignmentCenter
	l = layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)
	x := (crossSize - obj1.MinSize().Height) / 2
	assert.Equal(t, fyne.NewPos(0, x), obj1.Position())
	x = (crossSize - obj2.MinSize().Height) / 2
	assert.Equal(t, fyne.NewPos(obj1.MinSize().Width, x), obj2.Position())
	x = (crossSize - obj3.MinSize().Height) / 2
	assert.Equal(t, fyne.NewPos(obj2.Position().X+obj2.MinSize().Width, x), obj3.Position())
	assert.Equal(t,
		fyne.NewSize(
			obj1.MinSize().Width+obj2.MinSize().Width+obj3.MinSize().Width,
			fyne.Max(obj1.MinSize().Height, fyne.Max(obj2.MinSize().Height, obj3.MinSize().Height)),
		),
		l.MinSize(container.Objects),
	)

	// -- Baseline alignment
	// TODO (it is necessary to have widgets that implements DistanceToTextBaseline to do this test)
}

func TestFlexLayout_WithFlexibleWidgets(t *testing.T) {
	lsize := fyne.NewSize(200, 300)

	obj1 := widget.NewFlexible(3, NewMinSizeRect(fyne.NewSize(20, 40)))
	obj2 := widget.NewFlexible(1, NewMinSizeRect(fyne.NewSize(40, 20)))
	obj3 := NewMinSizeRect(fyne.NewSize(20, 60))

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(lsize)

	l := layout.NewFlexLayout(fyne.AxisVertical, &fyne.AxisAlignment{
		MainAxisAlignment:  fyne.MainAxisAlignmentStart,
		CrossAxisAlignment: fyne.CrossAxisAlignmentStart,
	})
	l.Layout(container.Objects, lsize)
	allocatedSize := obj3.MinSize().Height
	freeSpace := lsize.Height - allocatedSize
	spacePerFlex := freeSpace / 4 // (3 + 1)
	assert.Equal(t, fyne.NewPos(0, 0), obj1.Position())
	assert.Equal(t, fyne.NewPos(0, obj1.Size().Height), obj2.Position())
	assert.Equal(t, fyne.NewPos(0, obj2.Position().Y+obj2.Size().Height), obj3.Position())
	assert.Equal(t, fyne.NewSize(obj1.MinSize().Width, spacePerFlex*3), obj1.Size())
	assert.Equal(t, fyne.NewSize(obj2.MinSize().Width, spacePerFlex*1), obj2.Size())
	assert.Equal(t, obj3.MinSize(), obj3.Size())

	minSpacePerFlex := float32(0)
	minSpacePerFlex = fyne.Max(minSpacePerFlex, obj1.MinSize().Height/float32(obj1.Flex()))
	minSpacePerFlex = fyne.Max(minSpacePerFlex, obj2.MinSize().Height/float32(obj2.Flex()))
	allocatedFlexSpace := minSpacePerFlex * (float32(obj1.Flex() + obj2.Flex()))
	assert.Equal(t, fyne.NewSize(
		fyne.Max(obj1.MinSize().Width, fyne.Max(obj2.MinSize().Width, obj3.MinSize().Width)),
		allocatedSize+allocatedFlexSpace,
	), l.MinSize(container.Objects))
}

func TestFlexLayout_Hidden(t *testing.T) {
	lsize := fyne.NewSize(200, 100)

	obj1 := NewMinSizeRect(fyne.NewSize(20, 40))
	obj1.Hide()
	obj2 := widget.NewExpanded(NewMinSizeRect(fyne.NewSize(40, 20)))
	obj3 := NewMinSizeRect(fyne.NewSize(20, 60))
	obj3.Hide()

	container := &fyne.Container{
		Objects: []fyne.CanvasObject{obj1, obj2, obj3},
	}
	container.Resize(lsize)
	alignment := &fyne.AxisAlignment{
		MainAxisAlignment:  fyne.MainAxisAlignmentStart,
		CrossAxisAlignment: fyne.CrossAxisAlignmentStart,
	}
	l := layout.NewFlexLayout(fyne.AxisHorizontal, alignment)
	l.Layout(container.Objects, lsize)

	assert.Equal(t, fyne.NewPos(0, 0), obj2.Position())
	assert.Equal(t, fyne.NewSize(lsize.Width, obj2.MinSize().Height), obj2.Size())
	assert.Equal(t,
		obj2.MinSize(),
		l.MinSize(container.Objects),
	)
}
