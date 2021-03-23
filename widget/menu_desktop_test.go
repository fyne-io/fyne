// +build !mobile

package widget_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenu_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	item1 := fyne.NewMenuItem("A", nil)
	item2 := fyne.NewMenuItem("B (long)", nil)
	sep := fyne.NewMenuItemSeparator()
	item3 := fyne.NewMenuItem("C", nil)
	subItem1 := fyne.NewMenuItem("subitem A", nil)
	subItem2 := fyne.NewMenuItem("subitem B", nil)
	subItem3 := fyne.NewMenuItem("subitem C (long)", nil)
	subsubItem1 := fyne.NewMenuItem("subsubitem A (long)", nil)
	subsubItem2 := fyne.NewMenuItem("subsubitem B", nil)
	subItem3.ChildMenu = fyne.NewMenu("", subsubItem1, subsubItem2)
	item3.ChildMenu = fyne.NewMenu("", subItem1, subItem2, subItem3)
	menu := fyne.NewMenu("", item1, sep, item2, item3)

	for name, tt := range map[string]struct {
		windowSize     fyne.Size
		menuPos        fyne.Position
		mousePositions []fyne.Position
		useTestTheme   bool
		want           string
	}{
		"normal": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			want:       "menu/desktop/layout_normal.xml",
		},
		"normal with submenus": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 100),
				fyne.NewPos(100, 170),
			},
			want: "menu/desktop/layout_normal_with_submenus.xml",
		},
		"background of active submenu parents resets if sibling is focused": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
				fyne.NewPos(300, 170), // focus subsubmenu item
				fyne.NewPos(30, 60),   // focus sibling of submenu parent
			},
			want: "menu/desktop/layout_background_reset.xml",
		},
		"no space on right side for submenu": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(410, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(430, 100), // open submenu
				fyne.NewPos(300, 170), // open subsubmenu
			},
			want: "menu/desktop/layout_no_space_on_right.xml",
		},
		"no space on left & right side for submenu": {
			windowSize: fyne.NewSize(200, 300),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
			},
			want: "menu/desktop/layout_no_space_on_both_sides.xml",
		},
		"window too short for submenu": {
			windowSize: fyne.NewSize(500, 150),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 130), // open subsubmenu
			},
			want: "menu/desktop/layout_window_too_short_for_submenu.xml",
		},
		"theme change": {
			windowSize:   fyne.NewSize(500, 300),
			menuPos:      fyne.NewPos(10, 10),
			useTestTheme: true,
			want:         "menu/desktop/layout_theme_changed.xml",
		},
		"window too short for menu": {
			windowSize: fyne.NewSize(100, 50),
			menuPos:    fyne.NewPos(10, 10),
			want:       "menu/desktop/layout_window_too_short.xml",
		},
	} {
		t.Run(name, func(t *testing.T) {
			w.Resize(tt.windowSize)
			m := widget.NewMenu(menu)
			o := internalWidget.NewOverlayContainer(m, c, nil)
			c.Overlays().Add(o)
			defer c.Overlays().Remove(o)
			m.Move(tt.menuPos)
			m.Resize(m.MinSize())
			for _, pos := range tt.mousePositions {
				test.MoveMouse(c, pos)
			}
			test.AssertRendersToMarkup(t, tt.want, w.Canvas())
			if tt.useTestTheme {
				test.AssertImageMatches(t, "menu/layout_normal.png", c.Capture())
				test.WithTestTheme(t, func() {
					test.AssertImageMatches(t, "menu/layout_theme_changed.png", c.Capture())
				})
			}
		})
	}
}

func TestMenu_Scrolling(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("A", nil),
		fyne.NewMenuItem("B", nil),
		fyne.NewMenuItem("C", nil),
		fyne.NewMenuItem("D", nil),
		fyne.NewMenuItem("E", nil),
		fyne.NewMenuItem("F", nil),
	)

	w.Resize(fyne.NewSize(100, 100))
	m := widget.NewMenu(menu)
	o := internalWidget.NewOverlayContainer(m, c, nil)
	c.Overlays().Add(o)
	defer c.Overlays().Remove(o)
	m.Move(fyne.NewPos(10, 10))
	m.Resize(m.MinSize())
	maxScrollDistance := m.MinSize().Height - 90
	test.AssertRendersToMarkup(t, "menu/desktop/scroll_top.xml", w.Canvas())

	test.Scroll(c, fyne.NewPos(20, 20), 0, -50)
	test.AssertRendersToMarkup(t, "menu/desktop/scroll_middle.xml", w.Canvas())

	test.Scroll(c, fyne.NewPos(20, 20), 0, -maxScrollDistance)
	test.AssertRendersToMarkup(t, "menu/desktop/scroll_bottom.xml", w.Canvas())

	test.Scroll(c, fyne.NewPos(20, 20), 0, maxScrollDistance-50)
	test.AssertRendersToMarkup(t, "menu/desktop/scroll_middle.xml", w.Canvas())

	test.Scroll(c, fyne.NewPos(20, 20), 0, 50)
	test.AssertRendersToMarkup(t, "menu/desktop/scroll_top.xml", w.Canvas())
}

func TestMenu_TraverseMenu(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	w := fyne.CurrentApp().NewWindow("")
	defer w.Close()
	w.SetPadded(false)
	c := w.Canvas()

	itemWithChild := fyne.NewMenuItem("Bar", nil)
	itemWithChild.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("SubA", nil),
		fyne.NewMenuItem("SubB", nil),
	)
	m := widget.NewMenu(fyne.NewMenu("",
		fyne.NewMenuItem("Foo", nil),
		fyne.NewMenuItemSeparator(),
		itemWithChild,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Baz", nil),
	))
	w.SetContent(internalWidget.NewOverlayContainer(m, c, nil))
	w.Resize(m.MinSize())
	m.Resize(m.MinSize())

	// going all the way down …
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_initial.xml", c)

	m.ActivateNext()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_first_active.xml", c)

	m.ActivateNext()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_second_active.xml", c)

	m.ActivateNext()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_third_active.xml", c)

	m.ActivateNext()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_third_active.xml", c, "does not wrap around if last item is already active")

	// … and up again
	m.ActivatePrevious()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_second_active.xml", c)

	m.ActivatePrevious()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_first_active.xml", c)

	m.ActivatePrevious()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_first_active.xml", c, "does not wrap around if on top")

	// activate a submenu (show and activate first item)
	m.ActivateNext()
	assert.True(t, m.ActivateLastSubmenu())
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_submenu_first_active.xml", c)

	assert.False(t, m.ActivateLastSubmenu())
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_submenu_first_active.xml", c, "does nothing if there is no submenu at the last active item")

	// traversing through items of opened submenu
	m.ActivateNext()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_submenu_second_active.xml", c)

	m.ActivateNext()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_submenu_second_active.xml", c, "does not wrap around if last item is already active")

	m.ActivatePrevious()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_submenu_first_active.xml", c)

	m.ActivatePrevious()
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_submenu_first_active.xml", c, "does not wrap around if on top")

	// closing an open submenu
	assert.True(t, m.DeactivateLastSubmenu())
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_second_active.xml", c)

	assert.False(t, m.DeactivateLastSubmenu())
	test.AssertRendersToMarkup(t, "menu/desktop/traverse_second_active.xml", c, "does nothing if there is no submenu opened")
}

func TestMenu_TriggerTraversedMenu(t *testing.T) {
	var triggered string
	var dismissed bool
	setupMenu := func() *widget.Menu {
		triggered = ""
		dismissed = false
		itemWithChild := fyne.NewMenuItem("Bar", func() { triggered = "2nd" })
		itemWithChild.ChildMenu = fyne.NewMenu("",
			fyne.NewMenuItem("SubA", func() { triggered = "1st sub" }),
			fyne.NewMenuItem("SubB", nil),
			fyne.NewMenuItem("SubC", func() { triggered = "3rd sub" }),
		)
		m := widget.NewMenu(fyne.NewMenu("",
			fyne.NewMenuItem("Foo", func() { triggered = "1st" }),
			fyne.NewMenuItemSeparator(),
			itemWithChild,
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Baz", func() { triggered = "3rd" }),
		))
		m.OnDismiss = func() { dismissed = true }
		w := fyne.CurrentApp().NewWindow("")
		w.SetContent(internalWidget.NewOverlayContainer(m, w.Canvas(), nil))
		return m
	}

	t.Run("without active item", func(t *testing.T) {
		m := setupMenu()
		m.TriggerLast()
		assert.Equal(t, "", triggered)
		assert.True(t, dismissed)
	})
	t.Run("first item in submenu", func(t *testing.T) {
		m := setupMenu()
		m.ActivateNext()
		m.ActivateNext()
		require.True(t, m.ActivateLastSubmenu())
		m.TriggerLast()
		assert.Equal(t, "1st sub", triggered)
		assert.True(t, dismissed)
	})
	t.Run("last item in submenu", func(t *testing.T) {
		m := setupMenu()
		m.ActivateNext()
		m.ActivateNext()
		require.True(t, m.ActivateLastSubmenu())
		m.ActivateNext()
		m.ActivateNext()
		m.TriggerLast()
		assert.Equal(t, "3rd sub", triggered)
		assert.True(t, dismissed)
	})
	t.Run("item in menu", func(t *testing.T) {
		m := setupMenu()
		m.ActivateNext()
		m.ActivateNext()
		m.ActivateNext()
		m.TriggerLast()
		assert.Equal(t, "3rd", triggered)
		assert.True(t, dismissed)
	})
	t.Run("item with (closed) submenu", func(t *testing.T) {
		m := setupMenu()
		m.ActivateNext()
		m.ActivateNext()
		m.TriggerLast()
		assert.Equal(t, "2nd", triggered)
		assert.True(t, dismissed)
	})
	t.Run("item without action", func(t *testing.T) {
		m := setupMenu()
		m.ActivateNext()
		m.ActivateNext()
		require.True(t, m.ActivateLastSubmenu())
		m.ActivateNext()
		m.TriggerLast()
		assert.Equal(t, "", triggered)
		assert.True(t, dismissed)
	})
}
