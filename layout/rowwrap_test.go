package layout_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestRowWrapLayout_Layout(t *testing.T) {
	t.Run("should arrange objects in a row and wrap overflow objects into next row", func(t *testing.T) {
		// given
		h := float32(10)
		o1 := canvas.NewRectangle(color.Opaque)
		o1.SetMinSize(fyne.NewSize(30, h))
		o2 := canvas.NewRectangle(color.Opaque)
		o2.SetMinSize(fyne.NewSize(80, h))
		o3 := canvas.NewRectangle(color.Opaque)
		o3.SetMinSize(fyne.NewSize(50, h))

		containerSize := fyne.NewSize(125, 125)
		container := &fyne.Container{
			Objects: []fyne.CanvasObject{o1, o2, o3},
		}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayout().Layout(container.Objects, containerSize)

		// then
		p := theme.Padding()
		assert.Equal(t, fyne.NewPos(0, 0), o1.Position())
		assert.Equal(t, fyne.NewPos(o1.Size().Width+p, 0), o2.Position())
		assert.Equal(t, fyne.NewPos(0, o1.Size().Height+p), o3.Position())
	})
	t.Run("should do nothing when container is empty", func(t *testing.T) {
		containerSize := fyne.NewSize(125, 125)
		container := &fyne.Container{}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayout().Layout(container.Objects, containerSize)
	})
	t.Run("should ignore hidden objects", func(t *testing.T) {
		// given
		h := float32(10)
		o1 := canvas.NewRectangle(color.Opaque)
		o1.SetMinSize(fyne.NewSize(30, h))
		o2 := canvas.NewRectangle(color.Opaque)
		o2.SetMinSize(fyne.NewSize(80, h))
		o2.Hide()
		o3 := canvas.NewRectangle(color.Opaque)
		o3.SetMinSize(fyne.NewSize(50, h))

		containerSize := fyne.NewSize(125, 125)
		container := &fyne.Container{
			Objects: []fyne.CanvasObject{o1, o2, o3},
		}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayout().Layout(container.Objects, containerSize)

		// then
		p := theme.Padding()
		assert.Equal(t, fyne.NewPos(0, 0), o1.Position())
		assert.Equal(t, fyne.NewPos(o1.Size().Width+p, 0), o3.Position())
	})
}

func TestRowWrapLayout_MinSize(t *testing.T) {
	t.Run("should return min size of single object when container has only one", func(t *testing.T) {
		// given
		o := canvas.NewRectangle(color.Opaque)
		o.SetMinSize(fyne.NewSize(10, 10))
		container := container.NewWithoutLayout(o)
		layout := layout.NewRowWrapLayout()

		// when/then
		got := layout.MinSize(container.Objects)

		// then
		want := o.MinSize()
		assert.Equal(t, want, got)
	})
	t.Run("should return size 0 when container is empty", func(t *testing.T) {
		// given
		container := container.NewWithoutLayout()
		layout := layout.NewRowWrapLayout()

		// when/then
		got := layout.MinSize(container.Objects)

		// then
		want := fyne.NewSize(0, 0)
		assert.Equal(t, want, got)
	})
	t.Run("should initially return height of first object and width of widest object", func(t *testing.T) {
		// given
		h := float32(10)
		o1 := canvas.NewRectangle(color.Opaque)
		o1.SetMinSize(fyne.NewSize(10, h))
		o2 := canvas.NewRectangle(color.Opaque)
		o2.SetMinSize(fyne.NewSize(20, h))
		container := container.NewWithoutLayout(o1, o2)
		layout := layout.NewRowWrapLayout()

		// when/then
		got := layout.MinSize(container.Objects)

		// then
		want := fyne.NewSize(20, h)
		assert.Equal(t, want, got)
	})
	t.Run("should return actual size of arranged objects after layout was calculated", func(t *testing.T) {
		// given
		h := float32(10)
		o1 := canvas.NewRectangle(color.Opaque)
		o1.SetMinSize(fyne.NewSize(10, h))
		o2 := canvas.NewRectangle(color.Opaque)
		o2.SetMinSize(fyne.NewSize(20, h))
		container := container.New(layout.NewRowWrapLayout(), o1, o2)
		container.Resize(fyne.NewSize(15, 50))

		// when/then
		got := container.MinSize()

		// then
		want := fyne.NewSize(o2.Size().Width, (o1.Size().Height*2)+theme.Padding())
		assert.Equal(t, want, got)
	})
}
