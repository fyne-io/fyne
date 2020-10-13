// +build !windows !ci

package gomobile

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/mobile"
	"fyne.io/fyne/layout"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestCanvas_ChildMinSizeChangeAffectsAncestorsUpToRoot(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)
	leftObj1 := canvas.NewRectangle(color.Black)
	leftObj1.SetMinSize(fyne.NewSize(100, 50))
	leftObj2 := canvas.NewRectangle(color.Black)
	leftObj2.SetMinSize(fyne.NewSize(100, 50))
	leftCol := widget.NewVBox(leftObj1, leftObj2)
	rightObj1 := canvas.NewRectangle(color.Black)
	rightObj1.SetMinSize(fyne.NewSize(100, 50))
	rightObj2 := canvas.NewRectangle(color.Black)
	rightObj2.SetMinSize(fyne.NewSize(100, 50))
	rightCol := widget.NewVBox(rightObj1, rightObj2)
	content := widget.NewHBox(leftCol, rightCol)
	c.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), content))
	c.resize(fyne.NewSize(300, 300))

	oldContentSize := fyne.NewSize(200+theme.Padding(), 100+theme.Padding())
	assert.Equal(t, oldContentSize, content.Size())

	leftObj1.SetMinSize(fyne.NewSize(110, 60))
	c.ensureMinSize()

	expectedContentSize := oldContentSize.Add(fyne.NewSize(10, 10))
	assert.Equal(t, expectedContentSize, content.Size())
}

func TestCanvas_PixelCoordinateAtPosition(t *testing.T) {
	c := NewCanvas().(*mobileCanvas)

	pos := fyne.NewPos(4, 4)
	c.scale = 2.5
	x, y := c.PixelCoordinateForPosition(pos)
	assert.Equal(t, 10, x)
	assert.Equal(t, 10, y)
}

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
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
		altTapped = true
		wid.TappedSecondary(ev)
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
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
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})

	assert.False(t, buttonTap, "button should not be tapped")
}

func TestCanvas_TappedSecondary(t *testing.T) {
	var pointEvent *fyne.PointEvent
	var altTappedObj fyne.SecondaryTappable
	obj := &tappableLabel{}
	obj.ExtendBaseWidget(obj)
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(obj)
	c.resize(fyne.NewSize(36, 24))
	obj.Move(fyne.NewPos(3, 3))

	tapPos := fyne.NewPos(6, 6)
	c.tapDown(tapPos, 0)
	time.Sleep(310 * time.Millisecond)
	c.tapUp(tapPos, 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		obj.tap = true
		wid.Tapped(ev)
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
		obj.altTap = true
		altTappedObj = wid
		pointEvent = ev
		wid.TappedSecondary(ev)
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})

	assert.False(t, obj.tap, "don't tap primary")
	assert.True(t, obj.altTap, "tap secondary")
	assert.Equal(t, obj, altTappedObj)
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

func TestCanvas_Tappable(t *testing.T) {
	content := &touchableLabel{Label: widget.NewLabel("Hi\nHi\nHi")}
	content.ExtendBaseWidget(content)
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(content)
	c.resize(fyne.NewSize(36, 24))
	content.Resize(fyne.NewSize(24, 24))

	c.tapDown(fyne.NewPos(15, 15), 0)
	assert.True(t, content.down)

	c.tapUp(fyne.NewPos(15, 15), 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})
	assert.True(t, content.up)

	c.tapDown(fyne.NewPos(15, 15), 0)
	c.tapMove(fyne.NewPos(35, 15), 0, func(wid fyne.Draggable, ev *fyne.DragEvent) {
		wid.Dragged(ev)
	})
	assert.True(t, content.cancel)
}

func TestWindow_TappedAndDoubleTapped(t *testing.T) {
	tapped := 0
	but := newDoubleTappableButton()
	but.OnTapped = func() {
		tapped = 1
	}
	but.onDoubleTap = func() {
		tapped = 2
	}

	c := NewCanvas().(*mobileCanvas)
	c.SetContent(fyne.NewContainerWithLayout(layout.NewMaxLayout(), but))
	c.resize(fyne.NewSize(36, 24))

	c.tapDown(fyne.NewPos(15, 15), 0)
	c.tapUp(fyne.NewPos(15, 15), 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		wid.Tapped(ev)
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})
	time.Sleep(700 * time.Millisecond)
	assert.Equal(t, tapped, 1)

	c.tapDown(fyne.NewPos(15, 15), 0)
	c.tapUp(fyne.NewPos(15, 15), 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		wid.Tapped(ev)
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})
	c.tapDown(fyne.NewPos(15, 15), 0)
	c.tapUp(fyne.NewPos(15, 15), 0, func(wid fyne.Tappable, ev *fyne.PointEvent) {
		wid.Tapped(ev)
	}, func(wid fyne.SecondaryTappable, ev *fyne.PointEvent) {
	}, func(wid fyne.DoubleTappable, ev *fyne.PointEvent) {
		wid.DoubleTapped(ev)
	}, func(wid fyne.Draggable, ev *fyne.DragEvent) {
	})
	time.Sleep(700 * time.Millisecond)
	assert.Equal(t, tapped, 1)
}

func TestCanvas_Focusable(t *testing.T) {
	content := newFocusableEntry()
	c := NewCanvas().(*mobileCanvas)
	c.SetContent(content)

	c.tapDown(fyne.NewPos(10, 10), 0)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	c.tapDown(fyne.NewPos(10, 10), 1)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	c.Focus(content)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 0, content.unfocusedTimes)

	c.Unfocus()
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)

	content.Disable()
	c.Focus(content)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)

	c.tapDown(fyne.NewPos(10, 10), 2)
	assert.Equal(t, 1, content.focusedTimes)
	assert.Equal(t, 1, content.unfocusedTimes)
}

type touchableLabel struct {
	*widget.Label
	down, up, cancel bool
}

func (t *touchableLabel) TouchDown(event *mobile.TouchEvent) {
	t.down = true
}

func (t *touchableLabel) TouchUp(event *mobile.TouchEvent) {
	t.up = true
}

func (t *touchableLabel) TouchCancel(event *mobile.TouchEvent) {
	t.cancel = true
}

type tappableLabel struct {
	widget.Label
	tap, altTap bool
}

func (t *tappableLabel) Tapped(_ *fyne.PointEvent) {
	t.tap = true
}

func (t *tappableLabel) TappedSecondary(_ *fyne.PointEvent) {
	t.altTap = true
}

type focusableEntry struct {
	widget.Entry
	focusedTimes   int
	unfocusedTimes int
}

func newFocusableEntry() *focusableEntry {
	entry := &focusableEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (f *focusableEntry) FocusGained() {
	f.focusedTimes++
	f.Entry.FocusGained()
}

func (f *focusableEntry) FocusLost() {
	f.unfocusedTimes++
	f.Entry.FocusLost()
}

type doubleTappableButton struct {
	widget.Button

	onDoubleTap func()
}

func (t *doubleTappableButton) DoubleTapped(_ *fyne.PointEvent) {
	t.onDoubleTap()
}

func newDoubleTappableButton() *doubleTappableButton {
	but := &doubleTappableButton{}
	but.ExtendBaseWidget(but)

	return but
}
