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

func TestRowWrapLayout_MinSize(t *testing.T) {
	p := theme.Padding()
	t.Run("should return min size of single object when container has only one", func(t *testing.T) {
		// given
		a := makeObject(10, 10)
		container := container.NewWithoutLayout(a)
		layout := layout.NewRowWrapLayout()

		// when/then
		got := layout.MinSize(container.Objects)

		// then
		want := a.MinSize()
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
	t.Run("should estimate min size when layout not yet known", func(t *testing.T) {
		// given
		a := makeObject(10, 10)
		b := makeObject(20, 10)
		container := container.NewWithoutLayout(a, b)
		layout := layout.NewRowWrapLayout()

		// when/then
		got := layout.MinSize(container.Objects)

		// then
		want := fyne.NewSize(20, 10+p+10)
		assert.Equal(t, want, got)
	})
	t.Run("should use custom padding when estimating min size", func(t *testing.T) {
		// given
		a := makeObject(10, 10)
		b := makeObject(20, 10)
		container := container.NewWithoutLayout(a, b)
		layout := layout.NewRowWrapLayoutWithCustomPadding(5, 7)

		// when/then
		got := layout.MinSize(container.Objects)

		// then
		want := fyne.NewSize(20, 10+7+10)
		assert.Equal(t, want, got)
	})
	t.Run("should ignore invisible objects when estimating min size", func(t *testing.T) {
		// given
		a := makeObject(10, 10)
		b := makeObject(20, 10)
		b.Hide()
		container := container.NewWithoutLayout(a, b)
		layout := layout.NewRowWrapLayout()

		// when/then
		got := layout.MinSize(container.Objects)

		// then
		want := fyne.NewSize(10, 10)
		assert.Equal(t, want, got)
	})

	t.Run("should return actual size of arranged objects after layout was calculated", func(t *testing.T) {
		// given
		a := makeObject(10, 10)
		b := makeObject(20, 10)
		c := makeObject(20, 10)
		container := container.New(layout.NewRowWrapLayout(), a, b, c)
		container.Resize(fyne.NewSize(55, 50))

		// when/then
		got := container.MinSize()

		// then
		p := theme.Padding()
		want := fyne.NewSize(10+p+20, 10+p+10)
		assert.Equal(t, want, got)
	})
}

func TestRowWrapLayout_Layout(t *testing.T) {
	p := theme.Padding()
	t.Run("should arrange single object", func(t *testing.T) {
		// given
		a := makeObject(30, 10)
		containerSize := fyne.NewSize(120, 30)
		container := &fyne.Container{
			Objects: []fyne.CanvasObject{a},
		}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayout().Layout(container.Objects, containerSize)

		// then
		assert.Equal(t, fyne.NewPos(0, 0), a.Position())
	})

	t.Run("should arrange objects in single row when they fit", func(t *testing.T) {
		// given
		a := makeObject(30, 10)
		b := makeObject(80, 10)
		containerSize := fyne.NewSize(120, 30)
		container := &fyne.Container{
			Objects: []fyne.CanvasObject{a, b},
		}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayout().Layout(container.Objects, containerSize)

		// then
		assert.Equal(t, fyne.NewPos(0, 0), a.Position())
		assert.Equal(t, fyne.NewPos(30+p, 0), b.Position())
	})
	t.Run("should wrap overflowing object into new row with multiple objects in a row", func(t *testing.T) {
		// given
		a := makeObject(30, 10)
		b := makeObject(80, 10)
		c := makeObject(50, 10)
		containerSize := fyne.NewSize(125, 125)
		container := &fyne.Container{
			Objects: []fyne.CanvasObject{a, b, c},
		}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayout().Layout(container.Objects, containerSize)

		// then
		assert.Equal(t, fyne.NewPos(0, 0), a.Position())
		assert.Equal(t, fyne.NewPos(30+p, 0), b.Position())
		assert.Equal(t, fyne.NewPos(0, 10+p), c.Position())
	})
	t.Run("should wrap overflowing object into new row with one object on a row", func(t *testing.T) {
		// given
		a := makeObject(80, 10)
		b := makeObject(30, 10)
		containerSize := fyne.NewSize(40, 30)
		container := &fyne.Container{
			Objects: []fyne.CanvasObject{a, b},
		}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayout().Layout(container.Objects, containerSize)

		// then
		assert.Equal(t, fyne.NewPos(0, 0), a.Position())
		assert.Equal(t, fyne.NewPos(0, 10+p), b.Position())
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
		a := makeObject(30, 10)
		b := makeObject(80, 10)
		b.Hide()
		c := makeObject(50, 10)

		containerSize := fyne.NewSize(125, 125)
		container := &fyne.Container{
			Objects: []fyne.CanvasObject{a, b, c},
		}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayout().Layout(container.Objects, containerSize)

		// then
		assert.Equal(t, fyne.NewPos(0, 0), a.Position())
		assert.Equal(t, fyne.NewPos(30+p, 0), c.Position())
	})

	t.Run("should arrange objects with custom padding", func(t *testing.T) {
		// given
		a := makeObject(30, 10)
		b := makeObject(80, 10)
		c := makeObject(50, 10)
		containerSize := fyne.NewSize(125, 125)
		container := &fyne.Container{
			Objects: []fyne.CanvasObject{a, b, c},
		}
		container.Resize(containerSize)

		// when
		layout.NewRowWrapLayoutWithCustomPadding(5, 7).Layout(container.Objects, containerSize)

		// then
		assert.Equal(t, fyne.NewPos(0, 0), a.Position())
		assert.Equal(t, fyne.NewPos(30+5, 0), b.Position())
		assert.Equal(t, fyne.NewPos(0, 10+7), c.Position())
	})
}

func makeObject(w, h float32) fyne.CanvasObject {
	a := canvas.NewRectangle(color.Opaque)
	a.SetMinSize(fyne.NewSize(w, h))
	return a
}
