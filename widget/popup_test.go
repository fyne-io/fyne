package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/software"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPopUp(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())
	defer pop.Hide()

	assert.Empty(t, test.Canvas().Overlays().List())
	pop.Show()

	assert.True(t, pop.Visible())
	assert.Len(t, test.Canvas().Overlays().List(), 1)
	assert.Equal(t, pop, test.Canvas().Overlays().List()[0].(*widget.OverlayContainer).Content)
}

func TestShowPopUp(t *testing.T) {
	test.NewTempApp(t)

	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(200, 200))
	require.Nil(t, w.Canvas().Overlays().Top())

	label := NewLabel("Hi")
	ShowPopUp(label, w.Canvas())
	pop := w.Canvas().Overlays().Top().(*widget.OverlayContainer).Content
	if assert.NotNil(t, pop) {
		defer pop.Hide()

		assert.True(t, pop.Visible())
		assert.Len(t, w.Canvas().Overlays().List(), 1)
	}

	test.AssertRendersToMarkup(t, "popup/normal.xml", w.Canvas())
}

func TestShowPopUpAtPosition(t *testing.T) {
	c := test.NewCanvas()
	c.Resize(fyne.NewSize(100, 100))
	pos := fyne.NewPos(6, 9)
	label := NewLabel("Hi")
	ShowPopUpAtPosition(label, c, pos)
	pop := c.Overlays().Top().(*widget.OverlayContainer).Content
	if assert.NotNil(t, pop) {
		assert.True(t, pop.Visible())
		assert.Len(t, c.Overlays().List(), 1)
		assert.Equal(t, pos, pop.(*PopUp).Position())
	}
}

func TestShowPopUpAtRelativePosition(t *testing.T) {
	pos := fyne.NewPos(6, 9)
	label := NewLabel("Hi")
	parent1 := NewLabel("Parent1")
	parent2 := NewLabel("Parent2")
	w := test.NewTempWindow(
		t, &fyne.Container{Layout: layout.NewVBoxLayout(), Objects: []fyne.CanvasObject{parent1, parent2}},
	)
	w.Resize(fyne.NewSize(100, 200))

	ShowPopUpAtRelativePosition(label, w.Canvas(), pos, parent2)
	pop := w.Canvas().Overlays().Top().(*widget.OverlayContainer).Content
	if assert.NotNil(t, pop) {
		assert.True(t, pop.Visible())
		assert.Len(t, w.Canvas().Overlays().List(), 1)
		areaPos, _ := w.Canvas().InteractiveArea()
		assert.Equal(t, pos.Add(parent2.Position()).Add(fyne.NewPos(theme.Padding(), theme.Padding())).Subtract(areaPos), pop.(*PopUp).Position())
	}
}

func TestShowModalPopUp(t *testing.T) {
	test.NewTempApp(t)

	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(200, 199))
	require.Nil(t, w.Canvas().Overlays().Top())

	label := NewLabel("Hi")
	ShowModalPopUp(label, w.Canvas())
	pop := w.Canvas().Overlays().Top().(*widget.OverlayContainer).Content
	if assert.NotNil(t, pop) {
		defer pop.Hide()

		assert.True(t, pop.Visible())
		assert.Len(t, w.Canvas().Overlays().List(), 1)
	}

	test.AssertRendersToMarkup(t, "popup/modal.xml", w.Canvas())
}

func TestPopUp_Show(t *testing.T) {
	c := test.NewCanvas()
	cSize := fyne.NewSize(100, 100)
	c.Resize(cSize)
	label := NewLabel("Hi")
	pop := newPopUp(label, c)
	require.Nil(t, c.Overlays().Top())

	pop.Show()
	assert.Equal(t, pop, c.Overlays().Top().(*widget.OverlayContainer).Content)
	assert.Len(t, c.Overlays().List(), 1)
	assert.Equal(t, c.Overlays().Top().(*widget.OverlayContainer).Content.Size(), pop.Size())
	assert.Equal(t, label.MinSize(), pop.Content.Size())
}

func TestPopUp_ShowAtPosition(t *testing.T) {
	c := test.NewCanvas()
	cSize := fyne.NewSize(100, 100)
	c.Resize(cSize)
	label := NewLabel("Hi")
	pop := newPopUp(label, c)
	pos := fyne.NewPos(6, 9)
	require.Nil(t, c.Overlays().Top())

	pop.ShowAtPosition(pos)
	assert.Equal(t, pop, c.Overlays().Top().(*widget.OverlayContainer).Content)
	assert.Len(t, c.Overlays().List(), 1)
	assert.Equal(t, c.Overlays().Top().(*widget.OverlayContainer).Content.Size(), pop.Size())
	assert.Equal(t, label.MinSize(), pop.Content.Size())
	assert.Equal(t, pos, pop.Position())
}

func TestPopUp_Hide(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())
	pop.Show()

	assert.True(t, pop.Visible())
	pop.Hide()
	assert.False(t, pop.Visible())
	assert.Empty(t, test.Canvas().Overlays().List())
}

func TestPopUp_MinSize(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopUp(label, test.Canvas())

	inner := pop.Content.MinSize()
	assert.Equal(t, label.MinSize().Width, inner.Width)
	assert.Equal(t, label.MinSize().Height, inner.Height)

	min := pop.MinSize()
	assert.Equal(t, label.MinSize().Width, min.Width)
	assert.Equal(t, label.MinSize().Height, min.Height)
}

func TestPopUp_Move(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(70, 70))
	pop := newPopUp(label, win.Canvas())
	defer pop.Hide()

	pos := fyne.NewPos(10, 10)
	pop.Move(pos)
	pop.Show()

	assert.Equal(t, pos, pop.Position())

	popPos := pop.Position()
	fullPos, _ := win.Canvas().InteractiveArea()
	assert.Equal(t, fullPos, win.Canvas().Overlays().Top().Position()) // these are edge of safe area as the popUp must fill our overlay
	assert.Equal(t, float32(10), popPos.X)
	assert.Equal(t, float32(10), popPos.Y)
}

func TestPopUp_Move_Constrained(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(60, 48))
	pop := NewPopUp(label, win.Canvas())
	pop.Show()
	defer pop.Hide()

	pos := fyne.NewPos(30, 20)
	pop.Move(pos)

	innerPos := pop.Position()
	assert.Less(t, innerPos.X-theme.Padding(), pos.X,
		"content X position is adjusted to keep the content inside the window")
	assert.Less(t, innerPos.Y-theme.Padding(), pos.Y,
		"content Y position is adjusted to keep the content inside the window")
	// TODO constrain after a move
	// assert.Equal(t, win.Canvas().Size().Width-pop.Size().Width, innerPos.X,
	//	"content X position is adjusted to keep the content inside the window")
	// assert.Equal(t, win.Canvas().Size().Height-pop.Size().Height-theme.Padding(), innerPos.Y,
	//	"content Y position is adjusted to keep the content inside the window")
}

func TestPopUp_Move_ConstrainedWindowToSmall(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(10, 5))
	pop := NewPopUp(label, win.Canvas())
	pop.Show()
	defer pop.Hide()

	pos := fyne.NewPos(20, 10)
	pop.Move(pos)

	// innerPos := pop.Position()
	// TODO this constrain too
	// assert.Equal(t, theme.Padding(), innerPos.X, "content X position is adjusted but the window is too small")
	// assert.Equal(t, theme.Padding(), innerPos.Y, "content Y position is adjusted but the window is too small")
}

func TestPopUp_Resize(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(80, 80))

	pop := newPopUp(label, win.Canvas())
	pop.Show()
	defer pop.Hide()

	size := fyne.NewSize(60, 50)
	pop.Resize(size)
	assert.Equal(t, size, pop.Content.Size())

	popSize := pop.Size()
	assert.Equal(t, float32(60), popSize.Width)
	assert.Equal(t, float32(50), popSize.Height)
}

func TestPopUp_Tapped(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewTempWindow(t, NewLabel(""))
	c := win.Canvas()
	win.Resize(fyne.NewSize(120, 30))
	pop := NewPopUp(label, c)
	pop.Show()

	assert.True(t, pop.Visible())
	test.Tap(pop)
	assert.True(t, pop.Visible())
	test.TapCanvas(c, fyne.NewPos(100, 20))
	assert.False(t, pop.Visible())
	assert.Empty(t, test.Canvas().Overlays().List())
}

func TestPopUp_TappedSecondary(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewTempWindow(t, NewLabel(""))
	c := win.Canvas()
	win.Resize(fyne.NewSize(120, 30))
	pop := NewPopUp(label, c)
	pop.Show()

	assert.True(t, pop.Visible())
	test.TapSecondary(pop)
	assert.True(t, pop.Visible())
	test.TapCanvas(c, fyne.NewPos(100, 20))
	assert.False(t, pop.Visible())
	assert.Empty(t, test.Canvas().Overlays().List())
}

func TestPopUp_Stacked(t *testing.T) {
	assert.Nil(t, test.Canvas().Overlays().Top())
	assert.Empty(t, test.Canvas().Overlays().List())

	pop1 := NewPopUp(NewLabel("Hi"), test.Canvas())
	pop1.Show()
	assert.True(t, pop1.Visible())
	assert.Equal(t, pop1, test.Canvas().Overlays().Top().(*widget.OverlayContainer).Content)

	pop2 := NewPopUp(NewLabel("Hi"), test.Canvas())
	pop2.Show()
	assert.True(t, pop1.Visible())
	assert.True(t, pop2.Visible())
	assert.Equal(t, pop2, test.Canvas().Overlays().Top().(*widget.OverlayContainer).Content)

	pop3 := NewPopUp(NewLabel("Hi"), test.Canvas())
	pop3.Show()
	assert.True(t, pop1.Visible())
	assert.True(t, pop2.Visible())
	assert.True(t, pop3.Visible())
	assert.Equal(t, pop3, test.Canvas().Overlays().Top().(*widget.OverlayContainer).Content)

	pop3.Hide()
	assert.True(t, pop1.Visible())
	assert.True(t, pop2.Visible())
	assert.False(t, pop3.Visible())
	assert.Equal(t, pop2, test.Canvas().Overlays().Top().(*widget.OverlayContainer).Content)

	// hiding a pop-up cuts stack
	pop1.Hide()
	assert.False(t, pop1.Visible())
	assert.Nil(t, test.Canvas().Overlays().Top())
	assert.Empty(t, test.Canvas().Overlays().List())
}

func TestPopUp_Layout(t *testing.T) {
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(80, 80))

	content := NewLabel("Hi")
	pop := newPopUp(content, win.Canvas())
	pos := fyne.NewPos(6, 9)
	pop.ShowAtPosition(pos)
	defer pop.Hide()

	size := fyne.NewSize(60, 50)
	pop.Resize(size)
	r := cache.Renderer(pop)
	require.GreaterOrEqual(t, len(r.Objects()), 3)

	if s, ok := r.Objects()[0].(*widget.Shadow); assert.True(t, ok, "first rendered object is a shadow") {
		assert.Equal(t, size, s.Size())
	}
	if bg, ok := r.Objects()[1].(*canvas.Rectangle); assert.True(t, ok, "a background rectangle is rendered before the content") {
		assert.Equal(t, size, bg.Size())
		assert.Equal(t, theme.Color(theme.ColorNameOverlayBackground), bg.FillColor)
	}
	assert.Equal(t, r.Objects()[2], content)
}

func TestPopUp_ApplyThemeOnShow(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(200, 300))

	pop := NewPopUp(NewLabel("Label"), w.Canvas())

	test.ApplyTheme(t, test.Theme())
	pop.Show()
	test.AssertImageMatches(t, "popup/normal-onshow-theme-default.png", w.Canvas().Capture())
	pop.Hide()

	test.ApplyTheme(t, test.NewTheme())
	pop.Show()
	test.AssertImageMatches(t, "popup/normal-onshow-theme-changed.png", w.Canvas().Capture())
	pop.Hide()
}

func TestPopUp_ResizeOnShow(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	size := fyne.NewSize(200, 300)
	w.Resize(size)

	pop := NewPopUp(NewLabel("Label"), w.Canvas())

	pop.Show()
	_, fullSize := w.Canvas().InteractiveArea()
	assert.Equal(t, fullSize, w.Canvas().Overlays().Top().Size())
	pop.Hide()

	size = fyne.NewSize(500, 500)
	w.Resize(size)
	_, fullSize = w.Canvas().InteractiveArea()
	pop.Show()
	assert.Equal(t, fullSize, w.Canvas().Overlays().Top().Size())
	pop.Hide()
}

func TestPopUp_ResizeBeforeShow_CanvasSizeZero(t *testing.T) {
	test.NewTempApp(t)

	// Simulate canvas size {0,0}
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(0, 0))
	w := test.NewTempWindow(t, rect)
	w.SetPadded(false)
	w.Resize(fyne.NewSize(0, 0))
	assert.Zero(t, w.Canvas().Size())

	pop := NewPopUp(NewLabel("Label"), w.Canvas())
	popBgSize := fyne.NewSize(200, 200)
	pop.Resize(popBgSize)

	winSize := fyne.NewSize(300, 300)
	w.Resize(winSize)
	pop.Show()

	// get content padding dynamically
	popContentPadding := pop.MinSize().Subtract(pop.Content.MinSize())

	_, fullSize := w.Canvas().InteractiveArea()
	assert.Equal(t, popBgSize.Subtract(popContentPadding), pop.Content.Size())
	assert.Equal(t, fullSize, w.Canvas().Overlays().Top().Size())
}

func TestModalPopUp_Tapped(t *testing.T) {
	label := NewLabel("Hi")
	c := test.Canvas().(software.WindowlessCanvas)
	c.Resize(fyne.NewSquareSize(200))
	pop := NewModalPopUp(label, c)
	pop.Show()
	defer pop.Hide()

	assert.True(t, pop.Visible())
	test.TapCanvas(c, fyne.NewSquareOffsetPos(195))
	assert.True(t, pop.Visible())
	assert.Len(t, test.Canvas().Overlays().List(), 1)
	assert.Equal(t, pop, test.Canvas().Overlays().List()[0].(*widget.OverlayContainer).Content)
}

func TestModalPopUp_TappedSecondary(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewModalPopUp(label, test.Canvas())
	pop.Show()
	defer pop.Hide()

	assert.True(t, pop.Visible())
	test.TapSecondary(pop)
	assert.True(t, pop.Visible())
	assert.Len(t, test.Canvas().Overlays().List(), 1)
	assert.Equal(t, pop, test.Canvas().Overlays().List()[0].(*widget.OverlayContainer).Content)
}

func TestModalPopUp_Resize(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(80, 80))

	pop := newModalPopUp(label, win.Canvas())

	size := fyne.NewSize(50, 48)
	pop.Resize(size)
	assert.Equal(t, size, pop.Content.Size())
	pop.Show()
	defer pop.Hide()

	popSize := pop.Size()
	topSize := win.Canvas().Overlays().Top().Size()
	_, fullSize := win.Canvas().InteractiveArea()
	assert.Equal(t, fullSize, topSize) // these are full as the background must fill our overlay
	assert.Equal(t, float32(50), popSize.Width)
	assert.Equal(t, float32(48), popSize.Height)
}

func TestModalPopUp_TappedInside(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(80, 80))

	pop := newPopUp(label, win.Canvas())
	pop.Show()
	defer pop.Hide()

	size := fyne.NewSize(50, 48)
	pop.Resize(size)
	pop.Move(fyne.NewPos(10, 10))
	assert.Equal(t, size, pop.Content.Size())

	pop.Tapped(&fyne.PointEvent{Position: fyne.NewPos(30, 30)})
	assert.False(t, pop.Hidden)
	test.TapCanvas(win.Canvas(), fyne.NewPos(5, 5))
	assert.True(t, pop.Hidden)
}

func TestModalPopUp_Resize_Constrained(t *testing.T) {
	label := NewLabel("Hi")
	win := test.NewTempWindow(t, NewLabel("OK"))
	win.Resize(fyne.NewSize(80, 80))
	pop := NewModalPopUp(label, win.Canvas())

	pop.Resize(fyne.NewSize(90, 100))
	pop.Show()
	_, safe := win.Canvas().InteractiveArea()

	assert.Equal(t, safe.Width, pop.Content.Size().Width)
	assert.Equal(t, safe.Height, pop.Content.Size().Height)
	assert.Equal(t, safe.Width, pop.Size().Width)
	assert.Equal(t, safe.Height, pop.Size().Height)
}

func TestModalPopUp_ApplyThemeOnShow(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(200, 300))

	pop := NewModalPopUp(NewLabel("Label"), w.Canvas())

	test.ApplyTheme(t, test.Theme())
	pop.Show()
	test.AssertImageMatches(t, "popup/modal-onshow-theme-default.png", w.Canvas().Capture())
	pop.Hide()

	pop.Show()
	test.ApplyTheme(t, test.NewTheme())
	test.AssertImageMatches(t, "popup/modal-onshow-theme-changed.png", w.Canvas().Capture())
	pop.Hide()
}

func TestModalPopUp_ResizeOnShow(t *testing.T) {
	test.NewTempApp(t)
	w := test.NewTempWindow(t, canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(200, 300))

	pop := NewModalPopUp(NewLabel("Label"), w.Canvas())
	size := pop.MinSize()

	pop.Show()
	assert.Equal(t, size, pop.Size())
	pop.Hide()

	w.Resize(fyne.NewSize(500, 500))
	pop.Show()
	fullPos, fullSize := w.Canvas().InteractiveArea()
	assert.Equal(t, fullPos, w.Canvas().Overlays().Top().Position())
	assert.Equal(t, fullSize, w.Canvas().Overlays().Top().Size())
	assert.Equal(t, size, pop.Size())
	pop.Hide()
}

// Repro for fyne-io/fyne#5965.
//
// Root cause: internal/widget.overlayRenderer.MinSize dereferences
// OverlayContainer.canvas unconditionally. When the canvas reference is nil
// (e.g. the owning widget was detached from its window between the tap event
// and the popup creation) the call panics with a nil pointer dereference.
//
// Higher-level trigger seen in the wild: widget.Select.Tapped ->
// showPopUp -> NewPopUpMenu(menu, c) where Driver.CanvasForObject(s.super())
// returns nil because the Select was already removed from its parent tree.
// That path is driver-specific (real GLFW driver consults the canvas cache;
// the test driver fakes CanvasForObject), so the underlying contract is
// exercised here directly.
func TestNewPopUpMenu_NilCanvas_NoPanic(t *testing.T) {
	test.NewTempApp(t)

	menu := fyne.NewMenu(
		"",
		fyne.NewMenuItem("a", func() {}),
		fyne.NewMenuItem("b", func() {}),
		fyne.NewMenuItem("c", func() {}),
	)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("NewPopUpMenu with nil canvas panicked: %v", r)
		}
	}()

	pop := NewPopUpMenu(menu, nil)
	if pop == nil {
		return
	}
	pop.ShowAtPosition(fyne.NewPos(0, 0))
}

func TestShowPopUpMenuAtPosition_NilCanvas_NoPanic(t *testing.T) {
	test.NewTempApp(t)

	menu := fyne.NewMenu(
		"",
		fyne.NewMenuItem("only", func() {}),
	)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ShowPopUpMenuAtPosition with nil canvas panicked: %v", r)
		}
	}()

	ShowPopUpMenuAtPosition(menu, nil, fyne.NewPos(0, 0))
}

func TestModelPopUp_ResizeBeforeShow_CanvasSizeZero(t *testing.T) {
	test.NewTempApp(t)

	// Simulate canvas size {0,0}
	rect := canvas.NewRectangle(color.Black)
	rect.SetMinSize(fyne.NewSize(0, 0))
	w := test.NewTempWindow(t, rect)
	w.SetPadded(false)
	w.Resize(fyne.NewSize(0, 0))
	assert.Zero(t, w.Canvas().Size())

	pop := NewModalPopUp(NewLabel("Label"), w.Canvas())
	popBgSize := fyne.NewSize(200, 200)
	pop.Resize(popBgSize)

	winSize := fyne.NewSize(300, 300)
	w.Resize(winSize)
	pop.Show()

	// get content padding dynamically
	popContentPadding := pop.MinSize().Subtract(pop.Content.MinSize())
	assert.Equal(t, popBgSize.Subtract(popContentPadding), pop.Content.Size())
	assert.Equal(t, popBgSize, pop.Size())
}
