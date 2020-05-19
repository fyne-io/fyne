package widget_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestPopUpMenu_Move(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	m.Show()
	test.AssertImageMatches(t, "popup_menu_shown.png", c.Capture())

	m.Move(fyne.NewPos(20, 20))
	test.AssertImageMatches(t, "popup_menu_moved.png", c.Capture())

	m.Move(fyne.NewPos(190, 10))
	test.AssertImageMatches(t, "popup_menu_no_space_right.png", c.Capture())

	m.Move(fyne.NewPos(10, 190))
	test.AssertImageMatches(t, "popup_menu_no_space_down.png", c.Capture())
}

func TestPopUpMenu_Resize(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "popup_menu_shown_at_pos.png", c.Capture())

	m.Resize(m.Size().Add(fyne.NewSize(10, 10)))
	test.AssertImageMatches(t, "popup_menu_grown.png", c.Capture())

	largeSize := c.Size().Add(fyne.NewSize(10, 10))
	m.Resize(largeSize)
	test.AssertImageMatches(t, "popup_menu_canvas_too_small.png", c.Capture())
	assert.Equal(t, largeSize, m.Size())
}

func TestPopUpMenu_Show(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	test.AssertImageMatches(t, "popup_menu_hidden.png", c.Capture())

	m.Show()
	test.AssertImageMatches(t, "popup_menu_shown.png", c.Capture())
}

func TestPopUpMenu_ShowAtPosition(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	test.AssertImageMatches(t, "popup_menu_hidden.png", c.Capture())

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "popup_menu_shown_at_pos.png", c.Capture())

	m.Hide()
	test.AssertImageMatches(t, "popup_menu_hidden.png", c.Capture())

	m.ShowAtPosition(fyne.NewPos(190, 10))
	test.AssertImageMatches(t, "popup_menu_no_space_right.png", c.Capture())

	m.Hide()
	test.AssertImageMatches(t, "popup_menu_hidden.png", c.Capture())

	m.ShowAtPosition(fyne.NewPos(10, 190))
	test.AssertImageMatches(t, "popup_menu_no_space_down.png", c.Capture())

	m.Hide()
	test.AssertImageMatches(t, "popup_menu_hidden.png", c.Capture())
	menuSize := c.Size().Add(fyne.NewSize(10, 10))
	m.Resize(menuSize)

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "popup_menu_canvas_too_small.png", c.Capture())
	assert.Equal(t, menuSize, m.Size())
}

func setupPopUpMenuTest() (*widget.PopUpMenu, fyne.Window) {
	app := test.NewApp()
	app.Settings().SetTheme(theme.DarkTheme())

	w := test.NewWindow(canvas.NewRectangle(color.NRGBA{G: 150, B: 150, A: 255}))
	w.Resize(fyne.NewSize(200, 200))
	m := widget.NewPopUpMenu(fyne.NewMenu(
		"",
		fyne.NewMenuItem("Option A", nil),
		fyne.NewMenuItem("Option B", nil),
	), w.Canvas())
	return m, w
}

func tearDownPopUpMenuTest(w fyne.Window) {
	w.Close()
	test.NewApp()
}
