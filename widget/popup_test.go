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
	defer test.Canvas().Overlays().Remove(pop)

	assert.True(t, pop.Visible())
	assert.Equal(t, 1, len(test.Canvas().Overlays().List()))
	assert.Equal(t, pop, test.Canvas().Overlays().List()[0])
}

func TestPopUp_Hide(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	pop.Hide()
	assert.False(t, pop.Visible())
	assert.Equal(t, 0, len(test.Canvas().Overlays().List()))
}

func TestPopUp_MinSize(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())
	defer test.Canvas().Overlays().Remove(pop)

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
	defer win.Close()
	win.Resize(fyne.NewSize(50, 50))
	pop := NewPopUp(label, win.Canvas())
	defer test.Canvas().Overlays().Remove(pop)

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
	defer win.Close()
	win.Resize(fyne.NewSize(60, 40))
	pop := NewPopUp(label, win.Canvas())
	defer test.Canvas().Overlays().Remove(pop)

	pos := fyne.NewPos(30, 20)
	pop.Move(pos)

	innerPos := pop.Content.Position()
	assert.Less(t, innerPos.X-theme.Padding(), pos.X,
		"content X position is adjusted to keep the content inside the window")
	assert.Less(t, innerPos.Y-theme.Padding(), pos.Y,
		"content Y position is adjusted to keep the content inside the window")
	assert.Equal(t, win.Canvas().Size().Width-pop.Content.Size().Width-theme.Padding(), innerPos.X,
		"content X position is adjusted to keep the content inside the window")
	assert.Equal(t, win.Canvas().Size().Height-pop.Content.Size().Height-theme.Padding(), innerPos.Y,
		"content Y position is adjusted to keep the content inside the window")
}

func TestPopUp_Move_ConstrainedWindowToSmall(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(10, 5))
	pop := NewPopUp(label, win.Canvas())
	defer test.Canvas().Overlays().Remove(pop)

	pos := fyne.NewPos(20, 10)
	pop.Move(pos)

	innerPos := pop.Content.Position()
	assert.Equal(t, theme.Padding(), innerPos.X, "content X position is adjusted but the window is too small")
	assert.Equal(t, theme.Padding(), innerPos.Y, "content Y position is adjusted but the window is too small")
}

func TestPopUp_Resize(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(80, 80))
	pop := NewPopUp(label, win.Canvas())
	defer test.Canvas().Overlays().Remove(pop)

	size := fyne.NewSize(50, 40)
	pop.Resize(size)

	innerSize := pop.Content.Size()
	assert.Equal(t, size.Width-theme.Padding()*2, innerSize.Width)
	assert.Equal(t, size.Height-theme.Padding()*2, innerSize.Height)

	popSize := pop.Size()
	assert.Equal(t, 80, popSize.Width) // these are 80 as the popUp must fill our overlay
	assert.Equal(t, 80, popSize.Height)
}

func TestPopUp_Tapped(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	test.Tap(pop)
	assert.False(t, pop.Visible())
	assert.Equal(t, 0, len(test.Canvas().Overlays().List()))
}

func TestPopUp_TappedSecondary(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	assert.True(t, pop.Visible())
	test.TapSecondary(pop)
	assert.False(t, pop.Visible())
	assert.Equal(t, 0, len(test.Canvas().Overlays().List()))
}

func TestPopUp_Stacked(t *testing.T) {
	assert.Nil(t, test.Canvas().Overlays().Top())
	assert.Empty(t, test.Canvas().Overlays().List())

	pop1 := NewPopUp(NewLabel("Hi"), test.Canvas())
	assert.True(t, pop1.Visible())
	assert.Equal(t, pop1, test.Canvas().Overlays().Top())
	assert.Equal(t, []fyne.CanvasObject{pop1}, test.Canvas().Overlays().List())

	pop2 := NewPopUp(NewLabel("Hi"), test.Canvas())
	assert.True(t, pop1.Visible())
	assert.True(t, pop2.Visible())
	assert.Equal(t, pop2, test.Canvas().Overlays().Top())
	assert.Equal(t, []fyne.CanvasObject{pop1, pop2}, test.Canvas().Overlays().List())

	pop3 := NewPopUp(NewLabel("Hi"), test.Canvas())
	assert.True(t, pop1.Visible())
	assert.True(t, pop2.Visible())
	assert.True(t, pop3.Visible())
	assert.Equal(t, pop3, test.Canvas().Overlays().Top())
	assert.Equal(t, []fyne.CanvasObject{pop1, pop2, pop3}, test.Canvas().Overlays().List())

	pop3.Hide()
	assert.True(t, pop1.Visible())
	assert.True(t, pop2.Visible())
	assert.False(t, pop3.Visible())
	assert.Equal(t, pop2, test.Canvas().Overlays().Top())
	assert.Equal(t, []fyne.CanvasObject{pop1, pop2}, test.Canvas().Overlays().List())

	// hiding a pop-up cuts stack
	pop1.Hide()
	assert.False(t, pop1.Visible())
	assert.Nil(t, test.Canvas().Overlays().Top())
	assert.Empty(t, test.Canvas().Overlays().List())
}

func TestModalPopUp_Tapped(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewModalPopUp(label, test.Canvas())
	defer test.Canvas().Overlays().Remove(pop)

	assert.True(t, pop.Visible())
	test.Tap(pop)
	assert.True(t, pop.Visible())
	assert.Equal(t, 1, len(test.Canvas().Overlays().List()))
	assert.Equal(t, pop, test.Canvas().Overlays().List()[0])
}

func TestModalPopUp_TappedSecondary(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewModalPopUp(label, test.Canvas())
	defer test.Canvas().Overlays().Remove(pop)

	assert.True(t, pop.Visible())
	test.TapSecondary(pop)
	assert.True(t, pop.Visible())
	assert.Equal(t, 1, len(test.Canvas().Overlays().List()))
	assert.Equal(t, pop, test.Canvas().Overlays().List()[0])
}

func TestModalPopUp_Resize(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	win.Resize(fyne.NewSize(80, 80))
	pop := NewModalPopUp(label, win.Canvas())
	defer win.Canvas().Overlays().Remove(pop)

	assert.Less(t, pop.Content.Size().Width, 70-theme.Padding()*2)
	assert.Less(t, pop.Content.Size().Height, 50-theme.Padding()*2)

	pop.Resize(fyne.NewSize(70, 50))
	assert.Equal(t, 70-theme.Padding()*2, pop.Content.Size().Width)
	assert.Equal(t, 50-theme.Padding()*2, pop.Content.Size().Height)
	assert.Equal(t, 80, pop.Size().Width) // these are 80 as the popUp must fill our overlay
	assert.Equal(t, 80, pop.Size().Height)
}

func TestModalPopUp_Resize_Constrained(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	win.Resize(fyne.NewSize(80, 80))
	pop := NewModalPopUp(label, win.Canvas())
	defer win.Canvas().Overlays().Remove(pop)

	pop.Resize(fyne.NewSize(90, 100))
	assert.Equal(t, 80, pop.Content.Size().Width)
	assert.Equal(t, 80, pop.Content.Size().Height)
	assert.Equal(t, 80, pop.Size().Width)
	assert.Equal(t, 80, pop.Size().Height)
}

func TestModalPopUp_Resize_ConstrainedMin(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	win.Resize(fyne.NewSize(80, 80))
	pop := NewModalPopUp(label, win.Canvas())
	defer win.Canvas().Overlays().Remove(pop)

	win.Resize(fyne.NewSize(30, 20))
	assert.Greater(t, pop.Content.Size().Width, 30)
	assert.Greater(t, pop.Content.Size().Height, 20)
	assert.Greater(t, pop.Size().Width, 30)
	assert.Greater(t, pop.Size().Height, 20)
}
