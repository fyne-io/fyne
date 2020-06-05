package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"github.com/stretchr/testify/assert"
)

func TestSplitContainer(t *testing.T) {
	size := fyne.NewSize(100, 100)

	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	t.Run("Horizontal", func(t *testing.T) {
		NewHSplitContainer(objA, objB).Resize(size)

		sizeA := objA.Size()
		sizeB := objB.Size()

		assert.Equal(t, 50-halfDividerThickness(), sizeA.Width)
		assert.Equal(t, 100, sizeA.Height)
		assert.Equal(t, 50-halfDividerThickness(), sizeB.Width)
		assert.Equal(t, 100, sizeB.Height)
	})
	t.Run("Vertical", func(t *testing.T) {
		NewVSplitContainer(objA, objB).Resize(size)

		sizeA := objA.Size()
		sizeB := objB.Size()

		assert.Equal(t, 100, sizeA.Width)
		assert.Equal(t, 50-halfDividerThickness(), sizeA.Height)
		assert.Equal(t, 100, sizeB.Width)
		assert.Equal(t, 50-halfDividerThickness(), sizeB.Height)
	})
}

func TestSplitContainer_MinSize(t *testing.T) {
	textA := canvas.NewText("TEXTA", color.NRGBA{0, 0xff, 0, 0})
	textB := canvas.NewText("TEXTB", color.NRGBA{0, 0xff, 0, 0})
	t.Run("Horizontal", func(t *testing.T) {
		min := NewHSplitContainer(textA, textB).MinSize()
		assert.Equal(t, textA.MinSize().Width+textB.MinSize().Width+dividerThickness(), min.Width)
		assert.Equal(t, fyne.Max(textA.MinSize().Height, fyne.Max(textB.MinSize().Height, dividerLength())), min.Height)
	})
	t.Run("Vertical", func(t *testing.T) {
		min := NewVSplitContainer(textA, textB).MinSize()
		assert.Equal(t, fyne.Max(textA.MinSize().Width, fyne.Max(textB.MinSize().Width, dividerLength())), min.Width)
		assert.Equal(t, textA.MinSize().Height+textB.MinSize().Height+dividerThickness(), min.Height)
	})
}

func TestSplitContainer_SetRatio(t *testing.T) {
	size := fyne.NewSize(100, 100)

	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})

	t.Run("Horizontal", func(t *testing.T) {
		sc := NewHSplitContainer(objA, objB)
		sc.Resize(size)
		t.Run("Leading", func(t *testing.T) {
			sc.SetOffset(0.75)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 75-halfDividerThickness(), sizeA.Width)
			assert.Equal(t, 100, sizeA.Height)
			assert.Equal(t, 25-halfDividerThickness(), sizeB.Width)
			assert.Equal(t, 100, sizeB.Height)
		})
		t.Run("Trailing", func(t *testing.T) {
			sc.SetOffset(0.25)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 25-halfDividerThickness(), sizeA.Width)
			assert.Equal(t, 100, sizeA.Height)
			assert.Equal(t, 75-halfDividerThickness(), sizeB.Width)
			assert.Equal(t, 100, sizeB.Height)
		})
	})
	t.Run("Vertical", func(t *testing.T) {
		sc := NewVSplitContainer(objA, objB)
		sc.Resize(size)
		t.Run("Leading", func(t *testing.T) {
			sc.SetOffset(0.75)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 100, sizeA.Width)
			assert.Equal(t, 75-halfDividerThickness(), sizeA.Height)
			assert.Equal(t, 100, sizeB.Width)
			assert.Equal(t, 25-halfDividerThickness(), sizeB.Height)
		})
		t.Run("Trailing", func(t *testing.T) {
			sc.SetOffset(0.25)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 100, sizeA.Width)
			assert.Equal(t, 25-halfDividerThickness(), sizeA.Height)
			assert.Equal(t, 100, sizeB.Width)
			assert.Equal(t, 75-halfDividerThickness(), sizeB.Height)
		})
	})
}

func TestSplitContainer_SetRatio_limits(t *testing.T) {
	size := fyne.NewSize(50, 50)
	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objA.SetMinSize(size)
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB.SetMinSize(size)
	t.Run("Horizontal", func(t *testing.T) {
		sc := NewHSplitContainer(objA, objB)
		t.Run("Leading", func(t *testing.T) {
			sc.SetOffset(1.0)
			sc.Resize(fyne.NewSize(200, 50))
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 150-dividerThickness(), sizeA.Width)
			assert.Equal(t, 50, sizeA.Height)
			assert.Equal(t, 50, sizeB.Width)
			assert.Equal(t, 50, sizeB.Height)
		})
		t.Run("Trailing", func(t *testing.T) {
			sc.SetOffset(0.0)
			sc.Resize(fyne.NewSize(200, 50))
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 50, sizeA.Width)
			assert.Equal(t, 50, sizeA.Height)
			assert.Equal(t, 150-dividerThickness(), sizeB.Width)
			assert.Equal(t, 50, sizeB.Height)
		})
	})
	t.Run("Vertical", func(t *testing.T) {
		sc := NewVSplitContainer(objA, objB)
		t.Run("Leading", func(t *testing.T) {
			sc.SetOffset(1.0)
			sc.Resize(fyne.NewSize(50, 200))
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 50, sizeA.Width)
			assert.Equal(t, 150-dividerThickness(), sizeA.Height)
			assert.Equal(t, 50, sizeB.Width)
			assert.Equal(t, 50, sizeB.Height)
		})
		t.Run("Trailing", func(t *testing.T) {
			sc.SetOffset(0.0)
			sc.Resize(fyne.NewSize(50, 200))
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 50, sizeA.Width)
			assert.Equal(t, 50, sizeA.Height)
			assert.Equal(t, 50, sizeB.Width)
			assert.Equal(t, 150-dividerThickness(), sizeB.Height)
		})
	})
}

func TestSplitContainer_divider_cursor(t *testing.T) {
	t.Run("Horizontal", func(t *testing.T) {
		divider := newDivider(&SplitContainer{Horizontal: true})
		assert.Equal(t, desktop.HResizeCursor, divider.Cursor())
	})
	t.Run("Vertical", func(t *testing.T) {
		divider := newDivider(&SplitContainer{Horizontal: false})
		assert.Equal(t, desktop.VResizeCursor, divider.Cursor())
	})
}

func TestSplitContainer_divider_drag(t *testing.T) {
	size := fyne.NewSize(10, 10)
	objA := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objA.SetMinSize(size)
	objB := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	objB.SetMinSize(size)
	t.Run("Horizontal", func(t *testing.T) {
		split := NewHSplitContainer(objA, objB)
		split.Resize(fyne.NewSize(100, 100))
		divider := newDivider(split)
		assert.Equal(t, 0.5, split.Offset)

		divider.Dragged(&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(20, 9)},
			DraggedX:   10,
			DraggedY:   -1,
		})
		assert.Equal(t, 0.6, split.Offset)

		divider.DragEnd()
		assert.Equal(t, 0.6, split.Offset)
	})
	t.Run("Vertical", func(t *testing.T) {
		split := NewVSplitContainer(objA, objB)
		split.Resize(fyne.NewSize(100, 100))
		divider := newDivider(split)
		assert.Equal(t, 0.5, split.Offset)

		divider.Dragged(&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(9, 20)},
			DraggedX:   -1,
			DraggedY:   10,
		})
		assert.Equal(t, 0.6, split.Offset)

		divider.DragEnd()
		assert.Equal(t, 0.6, split.Offset)
	})
}

func TestSplitContainer_divider_hover(t *testing.T) {
	t.Run("Horizontal", func(t *testing.T) {
		divider := newDivider(&SplitContainer{Horizontal: true})
		assert.False(t, divider.hovered)

		divider.MouseIn(&desktop.MouseEvent{})
		assert.True(t, divider.hovered)

		divider.MouseOut()
		assert.False(t, divider.hovered)
	})
	t.Run("Vertical", func(t *testing.T) {
		divider := newDivider(&SplitContainer{Horizontal: false})
		assert.False(t, divider.hovered)

		divider.MouseIn(&desktop.MouseEvent{})
		assert.True(t, divider.hovered)

		divider.MouseOut()
		assert.False(t, divider.hovered)
	})
}

func TestSplitContainer_divider_MinSize(t *testing.T) {
	t.Run("Horizontal", func(t *testing.T) {
		divider := newDivider(&SplitContainer{Horizontal: true})
		min := divider.MinSize()
		assert.Equal(t, dividerThickness(), min.Width)
		assert.Equal(t, dividerLength(), min.Height)
	})
	t.Run("Vertical", func(t *testing.T) {
		divider := newDivider(&SplitContainer{Horizontal: false})
		min := divider.MinSize()
		assert.Equal(t, dividerLength(), min.Width)
		assert.Equal(t, dividerThickness(), min.Height)
	})
}
