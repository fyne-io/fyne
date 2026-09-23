package widget_test

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenu_RefreshOptions(t *testing.T) {
	test.NewTempApp(t)

	w := fyne.CurrentApp().NewWindow("")
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	itemFoo := fyne.NewMenuItem("Foo", nil)
	itemBar := fyne.NewMenuItem("Bar", nil)
	itemBar.ChildMenu = fyne.NewMenu("", fyne.NewMenuItem("Sub", nil))
	itemBar.Icon = theme.AccountIcon()
	itemBaz := fyne.NewMenuItem("Baz", nil)

	m := widget.NewMenu(fyne.NewMenu(
		"",
		itemFoo,
		fyne.NewMenuItemSeparator(),
		itemBar,
		fyne.NewMenuItemSeparator(),
		itemBaz,
	))
	w.SetContent(internalWidget.NewOverlayContainer(m, c, nil))
	// + 4,5 for canvas’ safe area
	w.Resize(m.MinSize().AddWidthHeight(4, 5))
	m.Resize(m.MinSize())
	test.AssertRendersToMarkup(t, "menu/refresh_initial.xml", c)

	itemBar.Disabled = true
	m.Refresh()

	test.AssertRendersToMarkup(t, "menu/refresh_disabled.xml", c)

	itemBaz.Checked = true
	m.Refresh()

	test.AssertRendersToMarkup(t, "menu/refresh_checkmark.xml", c)

	itemBar.Checked = true
	m.Refresh()

	test.AssertRendersToMarkup(t, "menu/refresh_2nd_checkmark.xml", c)

	itemBar.Checked = false
	itemBar.Disabled = false
	m.Refresh()

	itemBaz.Checked = false
	m.Refresh()

	test.AssertRendersToMarkup(t, "menu/refresh_initial.xml", c)
}

func TestMenu_TappedPaddingOrSeparator(t *testing.T) {
	test.NewTempApp(t)

	w := fyne.CurrentApp().NewWindow("")
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	var item1Hit, item2Hit, overlayContainerHit bool
	m := widget.NewMenu(fyne.NewMenu(
		"",
		fyne.NewMenuItem("Foo", func() { item1Hit = true }),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Bar", func() { item2Hit = true }),
	))
	size := m.MinSize()
	w.Resize(size.Add(fyne.NewSize(4, 4)))
	o := internalWidget.NewOverlayContainer(m, c, func() { overlayContainerHit = true })
	w.SetContent(o)

	// tap on separator
	p := fyne.NewPos(5, size.Height/2)
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
	test.TapCanvas(c, fyne.NewPos(size.Width+2, size.Height+2))
	assert.True(t, overlayContainerHit, "hit the overlay container")
}

func TestMenu_SubmenuTallerThanCanvasScrolls(t *testing.T) {
	w := test.NewTempWindow(t, widget.NewLabel(""))
	w.SetPadded(false)
	w.Resize(fyne.NewSize(200, 300))
	c := w.Canvas()
	d := fyne.CurrentApp().Driver()

	items := make([]*fyne.MenuItem, 50)
	for i := range items {
		items[i] = fyne.NewMenuItem(fmt.Sprintf("Item %d", i), nil)
	}
	parent := fyne.NewMenuItem("Alpha", nil)
	parent.ChildMenu = fyne.NewMenu("", items...)
	m := widget.NewPopUpMenu(fyne.NewMenu("", parent), c)
	m.ShowAtPosition(fyne.NewPos(0, 150))
	test.MoveMouse(c, fyne.NewPos(20, 160)) // hover "Alpha" to open its submenu on desktop ...
	test.TapCanvas(c, fyne.NewPos(20, 160)) // ... or tap it on mobile

	var child *widget.Menu
	for _, o := range test.WidgetRenderer(m.Menu).Objects() {
		if sub, ok := o.(*widget.Menu); ok {
			child = sub
		}
	}
	require.NotNil(t, child)
	_, areaSize := c.InteractiveArea()
	require.Greater(t, child.MinSize().Height, areaSize.Height, "the submenu must not fit the canvas")

	// the submenu is pinned to the top of the canvas and clamped to its height, like an over-tall menu
	childPos := d.AbsolutePositionForObject(child)
	assert.Equal(t, float32(0), childPos.Y)
	assert.Equal(t, areaSize.Height, child.Size().Height)

	// the items below the canvas can be scrolled into view
	last := child.Items[len(child.Items)-1]
	require.Greater(t, d.AbsolutePositionForObject(last).Y, areaSize.Height)
	test.Scroll(c, fyne.NewPos(childPos.X+10, 10), 0, -child.MinSize().Height)
	assert.LessOrEqual(t, d.AbsolutePositionForObject(last).Y+last.Size().Height, areaSize.Height)
}

func TestMenu_ResizeBelowMinSizeGrowsBack(t *testing.T) {
	test.NewTempApp(t)

	m := widget.NewMenu(fyne.NewMenu("", fyne.NewMenuItem("Foo", nil), fyne.NewMenuItem("Bar", nil)))
	m.Resize(fyne.NewSize(100, 10)) // e.g. a stale size after a theme change
	assert.Equal(t, m.MinSize().Height, m.Size().Height)
}
