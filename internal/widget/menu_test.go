package widget_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestMenu_ItemHovered(t *testing.T) {
	sub1 := fyne.NewMenuItem("sub1", nil)
	sub1.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub1 A", nil),
		fyne.NewMenuItem("sub1 B", nil),
	)
	sub2sub := fyne.NewMenuItem("sub2sub", nil)
	sub2sub.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub2sub A", nil),
		fyne.NewMenuItem("sub2sub B", nil),
	)
	sub2 := fyne.NewMenuItem("sub2", nil)
	sub2.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub2 A", nil),
		fyne.NewMenuItem("sub2 B", nil),
		sub2sub,
	)
	m := widget.NewMenu(
		fyne.NewMenu("",
			fyne.NewMenuItem("Foo", nil),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Bar", nil),
			sub1,
			sub2,
		),
	)
	size := m.MinSize()
	m.Resize(size)

	mi := m.Items[0].(*widget.MenuItem)
	r1 := cache.Renderer(mi)
	assert.Equal(t, color.Transparent, r1.BackgroundColor())
	mi.MouseIn(nil)
	assert.Equal(t, theme.HoverColor(), r1.BackgroundColor())
	mi.MouseOut()
	assert.Equal(t, color.Transparent, r1.BackgroundColor())

	sub1Widget := m.Items[3].(*widget.MenuItem)
	assert.Equal(t, sub1, sub1Widget.Item)
	sub2Widget := m.Items[4].(*widget.MenuItem)
	assert.Equal(t, sub2, sub2Widget.Item)
	assert.False(t, sub1Widget.Child.Visible(), "submenu is invisible initially")
	assert.False(t, sub2Widget.Child.Visible(), "submenu is invisible initially")
	sub1Widget.MouseIn(nil)
	assert.True(t, sub1Widget.Child.Visible(), "hovering item shows submenu")
	assert.False(t, sub2Widget.Child.Visible(), "other Child menu stays hidden")
	sub1Widget.MouseOut()
	assert.True(t, sub1Widget.Child.Visible(), "hover out does not hide submenu")
	assert.False(t, sub2Widget.Child.Visible(), "other Child menu still hidden")
	sub2Widget.MouseIn(nil)
	assert.False(t, sub1Widget.Child.Visible(), "hovering other item hides current submenu")
	assert.True(t, sub2Widget.Child.Visible(), "other Child menu shows up")

	sub2subWidget := sub2Widget.Child.Items[2].(*widget.MenuItem)
	assert.Equal(t, sub2sub, sub2subWidget.Item)
	assert.False(t, sub2subWidget.Child.Visible(), "2nd level submenu is invisible initially")
	sub2Widget.MouseOut()
	sub2subWidget.MouseIn(nil)
	assert.True(t, sub2Widget.Child.Visible(), "1st level submenu stays visible")
	assert.True(t, sub2subWidget.Child.Visible(), "2nd level submenu shows up")
	sub2subWidget.MouseOut()
	assert.True(t, sub2Widget.Child.Visible(), "1st level submenu still visible")
	assert.True(t, sub2subWidget.Child.Visible(), "2nd level submenu still visible")

	sub1Widget.MouseIn(nil)
	assert.False(t, sub2Widget.Child.Visible(), "1st level submenu is hidden by other submenu")
	sub2Widget.MouseIn(nil)
	assert.False(t, sub2subWidget.Child.Visible(), "2nd level submenu is hidden when re-entering its parent")
}

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
	o := widget.NewOverlayContainer(m, c, func() { overlayContainerHit = true })
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

func TestMenu_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.DarkTheme())

	w := test.NewWindow(nil)
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
		wantImage      string
	}{
		"normal": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			wantImage:  "menu_layout_normal.png",
		},
		"normal with submenus": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 100),
				fyne.NewPos(100, 170),
			},
			wantImage: "menu_layout_normal_with_submenus.png",
		},
		"background of active submenu parents resets if sibling is hovered": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
				fyne.NewPos(300, 170), // hover subsubmenu item
				fyne.NewPos(30, 60),   // hover sibling of submenu parent
			},
			wantImage: "menu_layout_background_reset.png",
		},
		"no space on right side for submenu": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(410, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(430, 100), // open submenu
				fyne.NewPos(300, 170), // open subsubmenu
			},
			wantImage: "menu_layout_no_space_on_right.png",
		},
		"no space on left & right side for submenu": {
			windowSize: fyne.NewSize(200, 300),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
			},
			wantImage: "menu_layout_no_space_on_both_sides.png",
		},
		"window too short for submenu": {
			windowSize: fyne.NewSize(500, 150),
			menuPos:    fyne.NewPos(10, 10),
			mousePositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 130), // open subsubmenu
			},
			wantImage: "menu_layout_window_too_short.png",
		},
		"theme change": {
			windowSize:   fyne.NewSize(500, 300),
			menuPos:      fyne.NewPos(10, 10),
			useTestTheme: true,
			wantImage:    "menu_layout_theme_changed.png",
		},
	} {
		t.Run(name, func(t *testing.T) {
			w.Resize(tt.windowSize)
			m := widget.NewMenu(menu)
			w.SetContent(m)
			w.Resize(tt.windowSize) // SetContent changes window’s size
			m.Resize(m.MinSize())
			m.Move(tt.menuPos)
			for _, pos := range tt.mousePositions {
				test.MoveMouse(c, pos)
			}
			if tt.useTestTheme {
				test.WithTestTheme(t, func() {
					test.AssertImageMatches(t, tt.wantImage, c.Capture())
				})
			} else {
				test.AssertImageMatches(t, tt.wantImage, c.Capture())
			}
		})
	}
}
