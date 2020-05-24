// +build mobile

package widget_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	internalWidget "fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func TestMenu_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.DarkTheme())

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
		wantImage    string
	}{
		"normal": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			wantImage:  "menu/layout_normal.png",
		},
		"normal with submenus": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),
				fyne.NewPos(100, 170),
			},
			wantImage: "menu/mobile/layout_normal_with_submenus.png",
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
			wantImage: "menu/mobile/layout_background_reset.png",
		},
		"no space on right side for submenu": {
			windowSize: fyne.NewSize(500, 300),
			menuPos:    fyne.NewPos(410, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(430, 100), // open submenu
				fyne.NewPos(300, 170), // open subsubmenu
			},
			wantImage: "menu/mobile/layout_no_space_on_right.png",
		},
		"no space on left & right side for submenu": {
			windowSize: fyne.NewSize(200, 300),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 170), // open subsubmenu
			},
			wantImage: "menu/mobile/layout_no_space_on_both_sides.png",
		},
		"window too short for submenu": {
			windowSize: fyne.NewSize(500, 150),
			menuPos:    fyne.NewPos(10, 10),
			tapPositions: []fyne.Position{
				fyne.NewPos(30, 100),  // open submenu
				fyne.NewPos(100, 130), // open subsubmenu
			},
			wantImage: "menu/mobile/layout_window_too_short_for_submenu.png",
		},
		"theme change": {
			windowSize:   fyne.NewSize(500, 300),
			menuPos:      fyne.NewPos(10, 10),
			useTestTheme: true,
			wantImage:    "menu/layout_theme_changed.png",
		},
		"window too short for menu": {
			windowSize: fyne.NewSize(100, 50),
			menuPos:    fyne.NewPos(10, 10),
			wantImage:  "menu/layout_window_too_short.png",
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

func TestMenu_Dragging(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.DarkTheme())

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
	test.AssertImageMatches(t, "menu/mobile/drag_top.png", c.Capture())

	test.Drag(c, fyne.NewPos(20, 20), 0, -50)
	test.AssertImageMatches(t, "menu/mobile/drag_middle.png", c.Capture())

	test.Drag(c, fyne.NewPos(20, 20), 0, -maxDragDistance)
	test.AssertImageMatches(t, "menu/mobile/drag_bottom.png", c.Capture())

	test.Drag(c, fyne.NewPos(20, 20), 0, maxDragDistance-50)
	test.AssertImageMatches(t, "menu/mobile/drag_middle.png", c.Capture())

	test.Drag(c, fyne.NewPos(20, 20), 0, 50)
	test.AssertImageMatches(t, "menu/mobile/drag_top.png", c.Capture())
}
