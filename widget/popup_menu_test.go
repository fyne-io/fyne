package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestPopUpMenu_Move(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	m.Show()
	test.AssertImageMatches(t, "popup_menu/shown.png", c.Capture())

	m.Move(fyne.NewPos(20, 20))
	test.AssertImageMatches(t, "popup_menu/moved.png", c.Capture())

	m.Move(fyne.NewPos(190, 10))
	test.AssertImageMatches(t, "popup_menu/no_space_right.png", c.Capture())

	m.Move(fyne.NewPos(10, 190))
	test.AssertImageMatches(t, "popup_menu/no_space_down.png", c.Capture())
}

func TestPopUpMenu_Resize(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "popup_menu/shown_at_pos.png", c.Capture())

	m.Resize(m.Size().Add(fyne.NewSize(10, 10)))
	test.AssertImageMatches(t, "popup_menu/grown.png", c.Capture())

	largeSize := c.Size().Add(fyne.NewSize(10, 10))
	m.Resize(largeSize)
	test.AssertImageMatches(t, "popup_menu/canvas_too_small.png", c.Capture())
	assert.Equal(t, fyne.NewSize(largeSize.Width, c.Size().Height), m.Size(), "width is larger than canvas; height is limited by canvas (menu scrolls)")
}

func TestPopUpMenu_Show(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	test.AssertImageMatches(t, "popup_menu/hidden.png", c.Capture())

	m.Show()
	test.AssertImageMatches(t, "popup_menu/shown.png", c.Capture())
}

func TestPopUpMenu_ShowAtPosition(t *testing.T) {
	m, w := setupPopUpMenuTest()
	defer tearDownPopUpMenuTest(w)
	c := w.Canvas()

	test.AssertImageMatches(t, "popup_menu/hidden.png", c.Capture())

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "popup_menu/shown_at_pos.png", c.Capture())

	m.Hide()
	test.AssertImageMatches(t, "popup_menu/hidden.png", c.Capture())

	m.ShowAtPosition(fyne.NewPos(190, 10))
	test.AssertImageMatches(t, "popup_menu/no_space_right.png", c.Capture())

	m.Hide()
	test.AssertImageMatches(t, "popup_menu/hidden.png", c.Capture())

	m.ShowAtPosition(fyne.NewPos(10, 190))
	test.AssertImageMatches(t, "popup_menu/no_space_down.png", c.Capture())

	m.Hide()
	test.AssertImageMatches(t, "popup_menu/hidden.png", c.Capture())
	menuSize := c.Size().Add(fyne.NewSize(10, 10))
	m.Resize(menuSize)

	m.ShowAtPosition(fyne.NewPos(10, 10))
	test.AssertImageMatches(t, "popup_menu/canvas_too_small.png", c.Capture())
	assert.Equal(t, fyne.NewSize(menuSize.Width, c.Size().Height), m.Size(), "width is larger than canvas; height is limited by canvas (menu scrolls)")
}

func setupPopUpMenuTest() (*PopUpMenu, fyne.Window) {
	app := test.NewApp()
	app.Settings().SetTheme(theme.DarkTheme())

	w := test.NewWindow(canvas.NewRectangle(color.NRGBA{G: 150, B: 150, A: 255}))
	w.Resize(fyne.NewSize(200, 200))
	m := newPopUpMenu(fyne.NewMenu(
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

//
// Old pop-up menu tests
//

func TestNewPopUpMenu(t *testing.T) {
	c := test.Canvas()
	menu := fyne.NewMenu("Foo", fyne.NewMenuItem("Bar", func() {}))

	pop := NewPopUpMenu(menu, c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	assert.Equal(t, pop, c.Overlays().List()[0])

	pop.Hide()
	assert.Equal(t, 0, len(c.Overlays().List()))
}

func TestPopUpMenu_Size(t *testing.T) {
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(100, 100))
	menu := fyne.NewMenu("Foo",
		fyne.NewMenuItem("A", func() {}),
		fyne.NewMenuItem("A", func() {}),
	)
	menuItemSize := canvas.NewText("A", color.Black).MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	expectedSize := menuItemSize.Add(fyne.NewSize(0, menuItemSize.Height)).Add(fyne.NewSize(0, theme.Padding()))
	c := win.Canvas()

	pop := NewPopUpMenu(menu, c)
	defer pop.Hide()
	assert.Equal(t, expectedSize, pop.Content.Size())

	for _, o := range test.LaidOutObjects(pop) {
		if s, ok := o.(*widget.Shadow); ok {
			assert.Equal(t, expectedSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)), s.Size())
		}
	}
}
