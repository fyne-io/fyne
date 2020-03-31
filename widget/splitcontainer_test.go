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
	split := (100 - dividerThickness()) / 2

	objA := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	objB := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	t.Run("Horizontal", func(t *testing.T) {
		NewHSplitContainer(objA, objB).Resize(size)

		sizeA := objA.Size()
		sizeB := objB.Size()

		assert.Equal(t, split, sizeA.Width)
		assert.Equal(t, 100, sizeA.Height)
		assert.Equal(t, split, sizeB.Width)
		assert.Equal(t, 100, sizeB.Height)
	})
	t.Run("Vertical", func(t *testing.T) {
		NewVSplitContainer(objA, objB).Resize(size)

		sizeA := objA.Size()
		sizeB := objB.Size()

		assert.Equal(t, 100, sizeA.Width)
		assert.Equal(t, split, sizeA.Height)
		assert.Equal(t, 100, sizeB.Width)
		assert.Equal(t, split, sizeB.Height)
	})
}

func TestSplitContainer_MinSize(t *testing.T) {
	textA := canvas.NewText("TEXTA", color.RGBA{0, 0xff, 0, 0})
	textB := canvas.NewText("TEXTB", color.RGBA{0, 0xff, 0, 0})
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

func TestSplitContainer_updateOffset(t *testing.T) {
	size := fyne.NewSize(100, 100)
	split := (100 - dividerThickness()) / 2

	objA := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	objB := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})

	t.Run("Horizontal", func(t *testing.T) {
		sc := NewHSplitContainer(objA, objB)
		sc.Resize(size)
		t.Run("Positive", func(t *testing.T) {
			sc.updateOffset(8)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, split+8, sizeA.Width)
			assert.Equal(t, 100, sizeA.Height)
			assert.Equal(t, split-8, sizeB.Width)
			assert.Equal(t, 100, sizeB.Height)
		})
		t.Run("Negative", func(t *testing.T) {
			sc.updateOffset(-12)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, split-12, sizeA.Width)
			assert.Equal(t, 100, sizeA.Height)
			assert.Equal(t, split+12, sizeB.Width)
			assert.Equal(t, 100, sizeB.Height)
		})
	})
	t.Run("Vertical", func(t *testing.T) {
		sc := NewVSplitContainer(objA, objB)
		sc.Resize(size)
		t.Run("Positive", func(t *testing.T) {
			sc.updateOffset(8)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 100, sizeA.Width)
			assert.Equal(t, split+8, sizeA.Height)
			assert.Equal(t, 100, sizeB.Width)
			assert.Equal(t, split-8, sizeB.Height)
		})
		t.Run("Negative", func(t *testing.T) {
			sc.updateOffset(-12)
			sizeA := objA.Size()
			sizeB := objB.Size()
			assert.Equal(t, 100, sizeA.Width)
			assert.Equal(t, split-12, sizeA.Height)
			assert.Equal(t, 100, sizeB.Width)
			assert.Equal(t, split+12, sizeB.Height)
		})
	})
}

func TestSplitContainer_updateOffset_limits(t *testing.T) {
	size := fyne.NewSize(50, 50)
	objA := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	objA.SetMinSize(size)
	objB := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	objB.SetMinSize(size)
	t.Run("Horizontal", func(t *testing.T) {
		sc := NewHSplitContainer(objA, objB)
		min := sc.MinSize()
		sc.Resize(fyne.NewSize(min.Width+6, min.Height))
		t.Run("Positive", func(t *testing.T) {
			sc.updateOffset(8)
			assert.Equal(t, 3, sc.offset)
		})
		t.Run("Negative", func(t *testing.T) {
			sc.updateOffset(-12)
			assert.Equal(t, -3, sc.offset)
		})
	})
	t.Run("Vertical", func(t *testing.T) {
		sc := NewVSplitContainer(objA, objB)
		min := sc.MinSize()
		sc.Resize(fyne.NewSize(min.Width, min.Height+6))
		t.Run("Positive", func(t *testing.T) {
			sc.updateOffset(8)
			assert.Equal(t, 3, sc.offset)
		})
		t.Run("Negative", func(t *testing.T) {
			sc.updateOffset(-12)
			assert.Equal(t, -3, sc.offset)
		})
	})
}

func TestSplitContainer_divider_drag(t *testing.T) {
	size := fyne.NewSize(10, 10)
	objA := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	objA.SetMinSize(size)
	objB := canvas.NewRectangle(color.RGBA{0, 0, 0, 0})
	objB.SetMinSize(size)
	t.Run("Horizontal", func(t *testing.T) {
		split := NewHSplitContainer(objA, objB)
		split.Resize(fyne.NewSize(100, 100))
		divider := newDivider(split)
		assert.Equal(t, 0, split.offset)

		divider.Dragged(&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(20, 9)},
			DraggedX:   10,
			DraggedY:   -1,
		})
		assert.Equal(t, 10, split.offset)

		divider.DragEnd()
		assert.Equal(t, 10, split.offset)
	})
	t.Run("Vertical", func(t *testing.T) {
		split := NewVSplitContainer(objA, objB)
		split.Resize(fyne.NewSize(100, 100))
		divider := newDivider(split)
		assert.Equal(t, 0, split.offset)

		divider.Dragged(&fyne.DragEvent{
			PointEvent: fyne.PointEvent{Position: fyne.NewPos(9, 20)},
			DraggedX:   -1,
			DraggedY:   10,
		})
		assert.Equal(t, 10, split.offset)

		divider.DragEnd()
		assert.Equal(t, 10, split.offset)
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
