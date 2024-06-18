package widget_test

import (
	"testing"

	"fyne.io/fyne/v2"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
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

	m := widget.NewMenu(fyne.NewMenu("",
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
	m := widget.NewMenu(fyne.NewMenu("",
		fyne.NewMenuItem("Foo", func() { item1Hit = true }),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Bar", func() { item2Hit = true }),
	))
	size := m.MinSize()
	w.Resize(size.Add(fyne.NewSize(10, 10)))
	m.Resize(size)
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
	test.TapCanvas(c, fyne.NewPos(size.Width+1, size.Height+1))
	assert.True(t, overlayContainerHit, "hit the overlay container")
}
