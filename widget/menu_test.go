package widget_test

import (
	"testing"

	"fyne.io/fyne"
	internalWidget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestMenu_ItemTapped(t *testing.T) {
	tapped := false
	item1 := fyne.NewMenuItem("Foo", nil)
	item2 := fyne.NewMenuItem("Bar", func() { tapped = true })
	item3 := fyne.NewMenuItem("Sub", nil)
	subItem := fyne.NewMenuItem("Foo", func() {})
	item3.ChildMenu = fyne.NewMenu("", subItem)
	m := widget.NewMenu(fyne.NewMenu("", item1, item2, item3))
	size := m.MinSize()
	m.Resize(size)
	dismissed := false
	m.DismissAction = func() { dismissed = true }

	mi1 := m.Items[0].(*widget.MenuItem)
	mi2 := m.Items[1].(*widget.MenuItem)
	mi3 := m.Items[2].(*widget.MenuItem)
	assert.Equal(t, item1, mi1.Item)
	assert.Equal(t, item2, mi2.Item)
	assert.Equal(t, item3, mi3.Item)

	// tap on item without action does not panic
	test.Tap(mi1)
	assert.False(t, tapped)
	assert.False(t, dismissed, "tap on item w/o action does not dismiss the menu")
	assert.True(t, m.Visible(), "tap on item w/o action does not hide the menu")

	test.Tap(mi2)
	assert.True(t, tapped)
	assert.True(t, dismissed, "tap on item dismisses the menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … the DismissAction is reponsible for that")

	dismissed = false // reset
	mi3.MouseIn(nil)
	sm := mi3.Child
	smi := sm.Items[0].(*widget.MenuItem)
	assert.Equal(t, subItem, smi.Item)
	assert.True(t, sm.Visible(), "sub menu is visible")

	test.Tap(smi)
	assert.True(t, dismissed, "tap on sub item dismisses the root menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … the DismissAction is reponsible for that")
	assert.False(t, sm.Visible(), "tap on sub item hides the sub menu")

	newActionTapped := false
	item2.Action = func() { newActionTapped = true }
	test.Tap(mi2)
	assert.True(t, newActionTapped, "tap on item performs its current action")
}

func TestMenu_TappedPaddingOrSeparator(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.DarkTheme())

	w := app.NewWindow("")
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	var item1Hit, item2Hit, overlayContainerHit bool
	m := widget.NewMenu(fyne.NewMenu("",
		fyne.NewMenuItem("Foo", func() { item1Hit = true }),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Bar", func() { item2Hit = true }),
	))
	size := m.MinSize()
	m.Resize(size)
	o := internalWidget.NewOverlayContainer(m, c, func() { overlayContainerHit = true })
	w.SetContent(o)
	w.Resize(size.Add(fyne.NewSize(10, 10)))

	// tap on top padding
	p := fyne.NewPos(5, 1)
	if test.AssertCanvasTappableAt(t, c, p) {
		test.TapCanvas(c, p)
		assert.False(t, item1Hit, "item 1 should not be hit")
		assert.False(t, item2Hit, "item 2 should not be hit")
		assert.False(t, overlayContainerHit, "the overlay container should not be hit")
	}

	// tap on separator
	fyne.NewPos(5, size.Height/2)
	if test.AssertCanvasTappableAt(t, c, p) {
		test.TapCanvas(c, p)
		assert.False(t, item1Hit, "item 1 should not be hit")
		assert.False(t, item2Hit, "item 2 should not be hit")
		assert.False(t, overlayContainerHit, "the overlay container should not be hit")
	}

	// tap bottom padding
	p = fyne.NewPos(5, size.Height-1)
	if test.AssertCanvasTappableAt(t, c, p) {
		test.TapCanvas(c, p)
		assert.False(t, item1Hit, "item 1 should not be hit")
		assert.False(t, item2Hit, "item 2 should not be hit")
		assert.False(t, overlayContainerHit, "the overlay container should not be hit")
	}

	// verify test setup: we can hit the items and the container
	test.TapCanvas(c, fyne.NewPos(5, size.Height/4))
	assert.True(t, item1Hit, "hit item 1")
	assert.False(t, item2Hit, "item 2 should not be hit")
	assert.False(t, overlayContainerHit, "the overlay container should not be hit")
	test.TapCanvas(c, fyne.NewPos(5, 3*size.Height/4))
	assert.True(t, item2Hit, "hit item 2")
	assert.False(t, overlayContainerHit, "the overlay container should not be hit")
	test.TapCanvas(c, fyne.NewPos(size.Width+1, size.Height+1))
	assert.True(t, overlayContainerHit, "hit the overlay container")
}
