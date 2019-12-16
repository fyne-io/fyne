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
	var pointEvent *fyne.PointEvent
	var tappedObj fyne.Tappable
	button := widget.NewButton("Test", func() {
		buttonTap = true
	})
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(button)
	c.resize(fyne.NewSize(36, 24))
	button.Move(fyne.NewPos(3, 3))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	c.tapUp(tapPos, 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		tapped = true
		tappedObj = wid
		pointEvent = ev
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
	if assert.NotNil(t, pointEvent) {
		assert.Equal(t, fyne.NewPos(6, 6), pointEvent.AbsolutePosition)
		assert.Equal(t, fyne.NewPos(3, 3), pointEvent.Position)
	}
}

func TestCanvas_Tapped_Multi(t *testing.T) {
	buttonTap := false
	button := widget.NewButton("Test", func() {
		buttonTap = true
	})
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(button)
	c.resize(fyne.NewSize(36, 24))
	button.Move(fyne.NewPos(3, 3))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	c.tapUp(tapPos, 1, func(wid fyne.Tappable, ev *fyne.PointEvent) { // different tapID
		wid.Tapped(ev)
	}, func(wid fyne.Tappable, ev *fyne.PointEvent) {
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})

	assert.False(t, buttonTap, "button should not be tapped")
}

func TestCanvas_TappedSecondary(t *testing.T) {
	tapped := false
	altTapped := false
	buttonTap := false
	var pointEvent *fyne.PointEvent
	var altTappedObj fyne.Tappable
	button := widget.NewButton("Test", func() {
		buttonTap = false
	})
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(button)
	c.resize(fyne.NewSize(36, 24))
	button.Move(fyne.NewPos(3, 3))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	time.Sleep(310 * time.Millisecond)
	c.tapUp(tapPos, 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		tapped = true
		wid.Tapped(ev)
	}, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		altTapped = true
		altTappedObj = wid
		pointEvent = ev
		wid.TappedSecondary(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})

	assert.False(t, tapped, "don't tap primary")
	assert.True(t, altTapped, "tap secondary")
	assert.False(t, buttonTap, "button should not be tapped (primary)")
	assert.Equal(t, button, altTappedObj)
	if assert.NotNil(t, pointEvent) {
		assert.Equal(t, fyne.NewPos(6, 6), pointEvent.AbsolutePosition)
		assert.Equal(t, fyne.NewPos(3, 3), pointEvent.Position)
	}
}

func TestCanvas_Dragged(t *testing.T) {
	dragged := false
	var draggedObj fyne.Draggable
	scroll := widget.NewScrollContainer(widget.NewLabel("Hi\nHi\nHi"))
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(scroll)
	c.resize(fyne.NewSize(40, 24))
	assert.Equal(t, 0, scroll.Offset.Y)

	c.tapDown(fyne.NewPos(32, 3), 0)
	c.tapMove(fyne.NewPos(32, 10), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged = true
		draggedObj = wid
	})

	assert.True(t, dragged)
	assert.Equal(t, scroll, draggedObj)
	// TODO find a way to get the test driver to report as mobile
	dragged = false
	c.tapMove(fyne.NewPos(32, 5), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
		dragged = true
	})
}
