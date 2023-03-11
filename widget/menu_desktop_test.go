//go:build !mobile
// +build !mobile

package widget_test

import (
	"image/color"
	"runtime"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
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
	subItem1.Checked = true
	subItem2 := fyne.NewMenuItem("subitem B", nil)
	subItem2.Checked = true
	subItem2.Icon = theme.InfoIcon()
	subItem3 := fyne.NewMenuItem("subitem C (long)", nil)
	subItem3.Icon = theme.MenuIcon()
	subsubItem1 := fyne.NewMenuItem("subsubitem A (long)", nil)
	subsubItem1.Icon = theme.FileIcon()
	subsubItem2 := fyne.NewMenuItem("subsubitem B", nil)
	subItem3.ChildMenu = fyne.NewMenu("", subsubItem1, subsubItem2)
	item3.ChildMenu = fyne.NewMenu("", subItem1, subItem2, subItem3)
	item4 := fyne.NewMenuItem("D", nil)
	subItem4a := fyne.NewMenuItem("a", nil)
	subItem4a.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyA, Modifier: fyne.KeyModifierControl}
	subItem4b := fyne.NewMenuItem("b", nil)
	subItem4b.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyBackspace, Modifier: fyne.KeyModifierAlt}
	subItem4c := fyne.NewMenuItem("c", nil)
	subItem4c.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyDelete, Modifier: fyne.KeyModifierSuper}
	subItem4d := fyne.NewMenuItem("d", nil)
	subItem4d.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyDown, Modifier: fyne.KeyModifierShift}
	subItem4e := fyne.NewMenuItem("e", nil)
	subItem4e.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyEnd, Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt}
	subItem4f := fyne.NewMenuItem("f", nil)
	subItem4f.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyEnter, Modifier: fyne.KeyModifierControl | fyne.KeyModifierSuper}
	subItem4g := fyne.NewMenuItem("g", nil)
	subItem4g.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyEscape, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift}
	subItem4h := fyne.NewMenuItem("h", nil)
	subItem4h.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyHome, Modifier: fyne.KeyModifierAlt | fyne.KeyModifierSuper}
	subItem4i := fyne.NewMenuItem("i", nil)
	subItem4i.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyLeft, Modifier: fyne.KeyModifierAlt | fyne.KeyModifierShift}
	subItem4j := fyne.NewMenuItem("j", nil)
	subItem4j.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyPageDown, Modifier: fyne.KeyModifierSuper | fyne.KeyModifierShift}
	subItem4k := fyne.NewMenuItem("k", nil)
	subItem4k.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyPageUp, Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt | fyne.KeyModifierSuper}
	subItem4l := fyne.NewMenuItem("l", nil)
	subItem4l.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt | fyne.KeyModifierShift}
	subItem4m := fyne.NewMenuItem("m", nil)
	subItem4m.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyRight, Modifier: fyne.KeyModifierControl | fyne.KeyModifierSuper | fyne.KeyModifierShift}
	subItem4n := fyne.NewMenuItem("n", nil)
	subItem4n.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeySpace, Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt | fyne.KeyModifierSuper | fyne.KeyModifierShift}
	subItem4o := fyne.NewMenuItem("o", nil)
	subItem4o.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyTab, Modifier: fyne.KeyModifierControl}
	subItem4p := fyne.NewMenuItem("p", nil)
	subItem4p.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyUp, Modifier: fyne.KeyModifierControl}
	subItem4q := fyne.NewMenuItem("q", nil)
	subItem4q.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF6}
	item4.ChildMenu = fyne.NewMenu("", subItem4a, subItem4b, subItem4c, subItem4d, subItem4e, subItem4f, subItem4g, subItem4h, subItem4i, subItem4j, subItem4k, subItem4l, subItem4m, subItem4n, subItem4o, subItem4p, subItem4q)

	menu := fyne.NewMenu("", item1, sep, item2, item3, sep, item4)

	var shortcutsMasterPrefixPath string
	if runtime.GOOS == "darwin" {
		shortcutsMasterPrefixPath = "menu/desktop/layout_shortcuts_darwin"
	} else {
		shortcutsMasterPrefixPath = "menu/desktop/layout_shortcuts_other"
	}

	for name, tt := range map[string]struct {
		windowSize         fyne.Size
		menuPos            fyne.Position
		mousePositions     []fyne.Position
		want               string
		wantImage          string
		wantTestThemeImage string
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
			windowSize:         fyne.NewSize(500, 300),
			menuPos:            fyne.NewPos(10, 10),
			want:               "menu/desktop/layout_theme_changed.xml",
			wantImage:          "menu/desktop/layout_normal.png",
			wantTestThemeImage: "menu/desktop/layout_theme_changed.png",
		},
		"window too short for menu": {
			windowSize: fyne.NewSize(100, 50),
			menuPos:    fyne.NewPos(10, 10),
			want:       "menu/desktop/layout_window_too_short.xml",
		},
		"menu with shortcuts": {
			windowSize: fyne.NewSize(300, 800),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 140), // open submenu
			},
			want:               shortcutsMasterPrefixPath + ".xml",
			wantImage:          shortcutsMasterPrefixPath + ".png",
			wantTestThemeImage: shortcutsMasterPrefixPath + "_theme_changed.png",
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
			if tt.wantImage != "" {
				test.AssertImageMatches(t, tt.wantImage, c.Capture())
			}
			if tt.wantTestThemeImage != "" {
				test.WithTestTheme(t, func() {
					test.AssertImageMatches(t, tt.wantTestThemeImage, c.Capture())
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
