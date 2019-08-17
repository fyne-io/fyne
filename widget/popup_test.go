package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewPopUp(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	assert.Equal(t, pop, test.Canvas().Overlay())
}

func TestPopUp_Hide(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	pop.Hide()
	assert.False(t, pop.Visible())
	assert.Nil(t, test.Canvas().Overlay())
}

func TestPopUp_MinSize(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	inner := pop.Content.MinSize()
	assert.Equal(t, label.MinSize().Width, inner.Width)
	assert.Equal(t, label.MinSize().Height, inner.Height)

	min := pop.MinSize()
	assert.Equal(t, label.MinSize().Width+theme.Padding()*2, min.Width)
	assert.Equal(t, label.MinSize().Height+theme.Padding()*2, min.Height)
}

func TestPopUp_Move(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	win.Resize(fyne.NewSize(50, 50))
	pop := NewPopUp(label, win.Canvas())

	pos := fyne.NewPos(10, 10)
	pop.Move(pos)

	innerPos := pop.Content.Position()
	assert.Equal(t, pos.X+theme.Padding(), innerPos.X)
	assert.Equal(t, pos.Y+theme.Padding(), innerPos.Y)

	popPos := pop.Position()
	assert.Equal(t, 0, popPos.X) // these are 0 as the popUp must fill our overlay
	assert.Equal(t, 0, popPos.Y)
}

func TestPopUp_Move_Constrained(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	win.Resize(fyne.NewSize(50, 30))
	pop := NewPopUp(label, win.Canvas())

	pos := fyne.NewPos(10, 10)
	pop.Move(pos)

	innerPos := pop.Content.Position()
	assert.Equal(t, pos.X+theme.Padding(), innerPos.X)
	assert.Equal(t, win.Canvas().Size().Height-pop.Content.Size().Height-theme.Padding(), innerPos.Y)
}

func TestPopUp_Tapped(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	test.Tap(pop)
	assert.False(t, pop.Visible())
	assert.Nil(t, test.Canvas().Overlay())
}

func TestPopUp_TappedSecondary(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	test.TapSecondary(pop)
	assert.False(t, pop.Visible())
	assert.Nil(t, test.Canvas().Overlay())
}

func TestModalPopUp_Tapped(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewModalPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	test.Tap(pop)
	assert.True(t, pop.Visible())
	assert.Equal(t, pop, test.Canvas().Overlay())
}

func TestModalPopUp_TappedSecondary(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewModalPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	test.TapSecondary(pop)
	assert.True(t, pop.Visible())
	assert.Equal(t, pop, test.Canvas().Overlay())
}
