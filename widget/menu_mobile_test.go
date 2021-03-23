// +build mobile

package widget_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
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
		windowSize   fyne.Size
		menuPos      fyne.Position
		tapPositions []fyne.Position
		useTestTheme bool
		want         string
	}{
		"normal": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			want:       "menu/mobile/layout_normal.xml",
		},
		"normal with submenus": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),
				fyne.NewPos(100, 170),
			},
			want: "menu/mobile/layout_normal_with_submenus.xml",
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
			want: "menu/mobile/layout_background_reset.xml",
		},
		"no space on right side for submenu": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(410, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(430, 100), // open submenu
				fyne.NewPos(300, 170), // open subsubmenu
			},
			want: "menu/mobile/layout_no_space_on_right.xml",
		},
		"no space on left & right side for submenu": {
			windowSize: fyne.NewSize(200, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
			},
			want: "menu/mobile/layout_no_space_on_both_sides.xml",
		},
		"window too short for submenu": {
			windowSize: fyne.NewSize(500, 150),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 130), // open subsubmenu
			},
			want: "menu/mobile/layout_window_too_short_for_submenu.xml",
		},
		"theme change": {
			windowSize:   fyne.NewSize(500, 300),
			menuPos:      fyne.NewPos(10, 10),
			useTestTheme: true,
			want:         "menu/mobile/layout_theme_changed.xml",
		},
		"window too short for menu": {
			windowSize: fyne.NewSize(100, 50),
			menuPos:    fyne.NewPos(10, 10),
			want:       "menu/mobile/layout_window_too_short.xml",
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
			for _, pos := range tt.tapPositions {
				test.TapCanvas(c, pos)
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

func TestMenu_Dragging(t *testing.T) {
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
	maxDragDistance := m.MinSize().Height - 90
	test.AssertRendersToMarkup(t, "menu/mobile/drag_top.xml", w.Canvas())

	test.Drag(c, fyne.NewPos(20, 20), 0, -50)
	test.AssertRendersToMarkup(t, "menu/mobile/drag_middle.xml", w.Canvas())

	test.Drag(c, fyne.NewPos(20, 20), 0, -maxDragDistance)
	test.AssertRendersToMarkup(t, "menu/mobile/drag_bottom.xml", w.Canvas())

	test.Drag(c, fyne.NewPos(20, 20), 0, maxDragDistance-50)
	test.AssertRendersToMarkup(t, "menu/mobile/drag_middle.xml", w.Canvas())

	test.Drag(c, fyne.NewPos(20, 20), 0, 50)
	test.AssertRendersToMarkup(t, "menu/mobile/drag_top.xml", w.Canvas())
}
