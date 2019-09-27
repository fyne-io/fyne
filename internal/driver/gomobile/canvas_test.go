package gomobile

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	canv "fyne.io/fyne/canvas"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestCanvas_Resize(t *testing.T) {
	c := &canvas{content: canv.NewRectangle(color.White), padded: true, scale: 2}
	screenSize := fyne.NewSize(480, 640)
	c.Resize(screenSize)

	theme := fyne.CurrentApp().Settings().Theme()
	assert.Equal(t, screenSize.Width-theme.Padding()*2, c.Content().Size().Width)
	assert.Greater(t, screenSize.Height-theme.Padding()*2, c.Content().Size().Height) // a status bar...

	assert.Equal(t, theme.Padding(), c.Content().Position().X)
	assert.Less(t, theme.Padding(), c.Content().Position().Y)
}

func TestCanvas_Tapped(t *testing.T) {
	tapped := false
	buttonTap := false
	var tappedObj fyne.Tappable
	button := widget.NewButton("Test", func() {
		buttonTap = true
	})
	c := &canvas{content: button}
	c.Resize(fyne.NewSize(36, 24))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos)
	c.tapUp(tapPos, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		tapped = true
		tappedObj = wid
		wid.Tapped(ev)
	}, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		wid.TappedSecondary(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})

	assert.True(t, tapped)
	assert.True(t, buttonTap)
	assert.Equal(t, button, tappedObj)
}

func TestCanvas_TappedSecondary(t *testing.T) {
	tapped := false
	buttonTap := false
	var tappedObj fyne.Tappable
	button := widget.NewButton("Test", func() {
		buttonTap = false
	})
	c := &canvas{content: button}
	c.Resize(fyne.NewSize(36, 24))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos)
	c.tapUp(tapPos, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		tapped = true
		tappedObj = wid
		wid.Tapped(ev)
	}, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		tapped = true
		tappedObj = wid
		wid.TappedSecondary(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})

	assert.True(t, tapped)
	assert.False(t, buttonTap)
	assert.Equal(t, button, tappedObj)
}

func TestCanvas_Dragged(t *testing.T) {
	dragged := false
	var draggedObj fyne.Draggable
	scroll := widget.NewScrollContainer(widget.NewLabel("Hi\nHi\nHi"))
	c := &canvas{content: scroll}
	c.Resize(fyne.NewSize(36, 24))
	assert.Equal(t, 0, scroll.Offset.Y)

	c.tapDown(fyne.NewPos(35, 3))
	c.tapMove(fyne.NewPos(35, 10), func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged = true
		draggedObj = wid
	})

	offset := scroll.Offset.Y
	assert.True(t, dragged)
	assert.NotNil(t, draggedObj)
	assert.Greater(t, offset, 0)

	c.tapMove(fyne.NewPos(35, 5), func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
	})
	assert.Less(t, scroll.Offset.Y, offset)
}
