package widget_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestPopUpMenu_Move(t *testing.T) {
	m, w := setupPopUpMenuTest(t)
	c := w.Canvas()

	m.Show()
	test.AssertRendersToMarkup(t, "popup_menu/shown.xml", c)

	m.Move(fyne.NewPos(20, 20))
	test.AssertRendersToMarkup(t, "popup_menu/moved.xml", c)

	m.Move(fyne.NewPos(190, 10))
	test.AssertRendersToMarkup(t, "popup_menu/no_space_right.xml", c)

	m.Move(fyne.NewPos(10, 190))
	test.AssertRendersToMarkup(t, "popup_menu/no_space_down.xml", c)
}

func TestPopUpMenu_Resize(t *testing.T) {
	m, w := setupPopUpMenuTest(t)
	c := w.Canvas()

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "popup_menu/shown_at_pos.xml", c)

	m.Resize(m.Size().Add(fyne.NewSize(10, 10)))
	test.AssertRendersToMarkup(t, "popup_menu/grown.xml", c)

	largeSize := c.Size().Add(fyne.NewSize(10, 10))
	m.Resize(largeSize)
	test.AssertRendersToMarkup(t, "popup_menu/canvas_too_small.xml", c)
	_, areaSize := c.InteractiveArea()
	assert.Equal(t, fyne.NewSize(largeSize.Width, areaSize.Height), m.Size(), "width is larger than canvas; height is limited by canvas (menu scrolls)")
}

func TestPopUpMenu_Show(t *testing.T) {
	m, w := setupPopUpMenuTest(t)
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.Show()
	test.AssertRendersToMarkup(t, "popup_menu/shown.xml", c)
}

func TestPopUpMenu_ShowAtPosition(t *testing.T) {
	m, w := setupPopUpMenuTest(t)
	c := w.Canvas()

	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertRendersToMarkup(t, "popup_menu/shown_at_pos.xml", c)

	m.Hide()
	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.ShowAtPosition(fyne.NewPos(190, 10))
	test.AssertRendersToMarkup(t, "popup_menu/no_space_right.xml", c)

	m.Hide()
	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.ShowAtPosition(fyne.NewPos(10, 190))
	test.AssertRendersToMarkup(t, "popup_menu/no_space_down.xml", c)

	m.Hide()
	test.AssertRendersToMarkup(t, "popup_menu/hidden.xml", c)

	m.ShowAtPosition(fyne.NewPos(10, 10))
	menuSize := c.Size().Add(fyne.NewSize(10, 10))
	m.Resize(menuSize)
	test.AssertRendersToMarkup(t, "popup_menu/canvas_too_small.xml", c)
	_, areaSize := c.InteractiveArea()
	assert.Equal(t, fyne.NewSize(menuSize.Width, areaSize.Height), m.Size(), "width is larger than canvas; height is limited by canvas (menu scrolls)")
}

func setupPopUpMenuTest(t *testing.T) (*widget.PopUpMenu, fyne.Window) {
	test.NewApp()

	w := test.NewTempWindow(t, canvas.NewRectangle(color.NRGBA{G: 150, B: 150, A: 255}))
	t.Cleanup(func() { test.NewApp() })

	w.Resize(fyne.NewSize(200, 200))
	m := widget.NewPopUpMenu(fyne.NewMenu(
		"",
		fyne.NewMenuItem("Option A", nil),
		fyne.NewMenuItem("Option B", nil),
	), w.Canvas())
	return m, w
}

func setupPopUpMenuWithSubmenusTest(t *testing.T, callback func(string)) (*widget.PopUpMenu, fyne.Window) {
	test.NewApp()

	w := test.NewTempWindow(t, canvas.NewRectangle(color.NRGBA{G: 150, B: 150, A: 255}))
	t.Cleanup(func() { test.NewApp() })

	w.Resize(fyne.NewSize(800, 600))
	itemA := fyne.NewMenuItem("Option A", func() { callback("Option A") })
	itemB := fyne.NewMenuItem("Option B", func() { callback("Option B") })
	itemBA := fyne.NewMenuItem("Sub Option A", func() { callback("Sub Option A") })
	itemBB := fyne.NewMenuItem("Sub Option B", func() { callback("Sub Option B") })
	itemBBA := fyne.NewMenuItem("Sub Sub Option A", func() { callback("Sub Sub Option A") })
	itemBBB := fyne.NewMenuItem("Sub Sub Option B", func() { callback("Sub Sub Option B") })
	itemB.ChildMenu = fyne.NewMenu("", itemBA, itemBB)
	itemBB.ChildMenu = fyne.NewMenu("", itemBBA, itemBBB)
	m := widget.NewPopUpMenu(fyne.NewMenu("", itemA, itemB), w.Canvas())
	return m, w
}
