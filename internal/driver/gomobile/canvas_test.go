package gomobile

import (
	"testing"
	"time"

	"fyne.io/fyne"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestCanvas_Tapped(t *testing.T) {
	tapped := false
	altTapped := false
	buttonTap := false
	var tappedObj fyne.Tappable
	button := widget.NewButton("Test", func() {
		buttonTap = true
	})
	c := &mobileCanvas{content: button}
	c.resize(fyne.NewSize(36, 24))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos)
	c.tapUp(tapPos, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		tapped = true
		tappedObj = wid
		wid.Tapped(ev)
	}, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		altTapped = true
		wid.TappedSecondary(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})

	assert.True(t, tapped, "tap primary")
	assert.False(t, altTapped, "don't tap secondary")
	assert.True(t, buttonTap, "button should be tapped")
	assert.Equal(t, button, tappedObj)
}

func TestCanvas_TappedSecondary(t *testing.T) {
	tapped := false
	altTapped := false
	buttonTap := false
	var altTappedObj fyne.Tappable
	button := widget.NewButton("Test", func() {
		buttonTap = false
	})
	c := &mobileCanvas{content: button}
	c.resize(fyne.NewSize(36, 24))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos)
	time.Sleep(310 * time.Millisecond)
	c.tapUp(tapPos, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		tapped = true
		wid.Tapped(ev)
	}, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		altTapped = true
		altTappedObj = wid
		wid.TappedSecondary(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})

	assert.False(t, tapped, "don't tap primary")
	assert.True(t, altTapped, "tap secondary")
	assert.False(t, buttonTap, "button should not be tapped (primary)")
	assert.Equal(t, button, altTappedObj)
}

func TestCanvas_Dragged(t *testing.T) {
	dragged := false
	var draggedObj fyne.Draggable
	scroll := widget.NewScrollContainer(widget.NewLabel("Hi\nHi\nHi"))
	c := &mobileCanvas{content: scroll}
	c.resize(fyne.NewSize(36, 24))
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
