// +build mobile

package widget_test

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestMenu_ItemWithChildTapped(t *testing.T) {
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

	sub1Widget := m.Items[3].(*widget.MenuItem)
	assert.Equal(t, sub1, sub1Widget.Item)
	sub2Widget := m.Items[4].(*widget.MenuItem)
	assert.Equal(t, sub2, sub2Widget.Item)
	assert.False(t, sub1Widget.Child.Visible(), "submenu is invisible initially")
	assert.False(t, sub2Widget.Child.Visible(), "submenu is invisible initially")
	test.Tap(sub1Widget)
	assert.True(t, sub1Widget.Child.Visible(), "tapping item shows submenu")
	assert.False(t, sub2Widget.Child.Visible(), "other Child menu stays hidden")
	test.Tap(sub2Widget)
	assert.False(t, sub1Widget.Child.Visible(), "tapping other item hides current submenu")
	assert.True(t, sub2Widget.Child.Visible(), "other Child menu shows up")

	sub2subWidget := sub2Widget.Child.Items[2].(*widget.MenuItem)
	assert.Equal(t, sub2sub, sub2subWidget.Item)
	assert.False(t, sub2subWidget.Child.Visible(), "2nd level submenu is invisible initially")
	test.Tap(sub2subWidget)
	assert.True(t, sub2Widget.Child.Visible(), "1st level submenu stays visible")
	assert.True(t, sub2subWidget.Child.Visible(), "2nd level submenu shows up")

	test.Tap(sub1Widget)
	assert.False(t, sub2Widget.Child.Visible(), "1st level submenu is hidden by other submenu")
	test.Tap(sub2Widget)
	assert.False(t, sub2subWidget.Child.Visible(), "2nd level submenu is hidden when re-entering its parent")
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
		windowSize   fyne.Size
		menuPos      fyne.Position
		tapPositions []fyne.Position
		useTestTheme bool
		wantImage    string
	}{
		"normal": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			wantImage:  "menu_layout_normal.png",
		},
		"normal with submenus": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),
				fyne.NewPos(100, 170),
			},
			wantImage: "menu_layout_mobile_normal_with_submenus.png",
		},
		"background of active submenu parents resets if sibling is hovered": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
				fyne.NewPos(300, 170), // hover subsubmenu item
				fyne.NewPos(30, 60),   // hover sibling of submenu parent
			},
			wantImage: "menu_layout_mobile_background_reset.png",
		},
		"no space on right side for submenu": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(410, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(430, 100), // open submenu
				fyne.NewPos(300, 170), // open subsubmenu
			},
			wantImage: "menu_layout_mobile_no_space_on_right.png",
		},
		"no space on left & right side for submenu": {
			windowSize: fyne.NewSize(200, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
			},
			wantImage: "menu_layout_mobile_no_space_on_both_sides.png",
		},
		"window too short for submenu": {
			windowSize: fyne.NewSize(500, 150),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 130), // open subsubmenu
			},
			wantImage: "menu_layout_mobile_window_too_short.png",
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
			for _, pos := range tt.tapPositions {
				test.TapCanvas(c, pos)
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
