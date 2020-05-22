// +build !mobile
// +build !ci !windows

package glfw_test

import (
	"strconv"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/glfw"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestMenuBar(t *testing.T) {
	app := test.NewApp()
	defer test.NewApp()
	app.Settings().SetTheme(theme.LightTheme())

	w := test.NewWindow(nil)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(fyne.NewSize(300, 300))
	c := w.Canvas()

	var lastAction string

	m1i3i3i1 := fyne.NewMenuItem("Old 1", func() { lastAction = "old 1" })
	m1i3i3i2 := fyne.NewMenuItem("Old 2", func() { lastAction = "old 2" })
	m1i3i1 := fyne.NewMenuItem("File 1", func() { lastAction = "file 1" })
	m1i3i2 := fyne.NewMenuItem("File 2", func() { lastAction = "file 2" })
	m1i3i3 := fyne.NewMenuItem("Older", nil)
	m1i3i3.ChildMenu = fyne.NewMenu("", m1i3i3i1, m1i3i3i2)
	m1i1 := fyne.NewMenuItem("New", func() { lastAction = "new" })
	m1i2 := fyne.NewMenuItem("Open", func() { lastAction = "open" })
	m1i3 := fyne.NewMenuItem("Recent", nil)
	m1i3.ChildMenu = fyne.NewMenu("", m1i3i1, m1i3i2, m1i3i3)
	// TODO: remove useless separators: trailing, leading & double
	// m1 := fyne.NewMenu("File", m1i1, m1i2, fyne.newMenuItemSeparator(), m1i3)
	m1 := fyne.NewMenu("File", m1i1, m1i2, m1i3)

	m2i1 := fyne.NewMenuItem("Copy", func() { lastAction = "copy" })
	m2i2 := fyne.NewMenuItem("Paste", func() { lastAction = "paste" })
	m2 := fyne.NewMenu("Edit", m2i1, m2i2)

	m3i1 := fyne.NewMenuItem("Help!", func() { lastAction = "help" })
	m3 := fyne.NewMenu("Help", m3i1)

	menu := fyne.NewMainMenu(m1, m2, m3)
	menuBar := glfw.NewMenuBar(menu, c)

	themeCounter := 0
	button := widget.NewButton("Button", func() {
		switch themeCounter % 3 {
		case 0:
			test.ApplyTheme(t, theme.DarkTheme())
		case 1:
			test.ApplyTheme(t, test.NewTheme())
		default:
			test.ApplyTheme(t, theme.LightTheme())
		}
		themeCounter++
	})

	container := fyne.NewContainer(button, menuBar)

	w.SetContent(container)
	w.Resize(fyne.NewSize(300, 300))
	button.Resize(button.MinSize())
	button.Move(fyne.NewPos(100, 50))
	menuBar.Resize(fyne.NewSize(300, 0).Max(menuBar.MinSize()))

	_ = lastAction
	buttonPos := fyne.NewPos(110, 60)
	fileMenuPos := fyne.NewPos(20, 10)
	fileNewPos := fyne.NewPos(20, 50)
	fileOpenPos := fyne.NewPos(20, 70)
	fileRecentPos := fyne.NewPos(20, 100)
	fileRecentOlderPos := fyne.NewPos(120, 170)
	fileRecentOlderOld1Pos := fyne.NewPos(200, 170)
	editMenuPos := fyne.NewPos(70, 10)
	helpMenuPos := fyne.NewPos(120, 10)
	type action struct {
		typ string
		pos fyne.Position
	}
	type step struct {
		actions    []action
		wantImage  string
		wantAction string
	}
	for name, tt := range map[string]struct {
		steps []step
	}{
		"switch theme": {
			steps: []step{
				{
					actions:   []action{{"move", buttonPos}},
					wantImage: "menu_bar_hovered_content.png",
				},
				{
					actions:   []action{{"tap", buttonPos}},
					wantImage: "menu_bar_hovered_content_dark.png",
				},
				{
					actions:   []action{{"tap", buttonPos}},
					wantImage: "menu_bar_hovered_content_test_theme.png",
				},
				{
					actions:   []action{{"tap", buttonPos}},
					wantImage: "menu_bar_hovered_content.png",
				},
			},
		},
		"activate and deactivate menu": {
			[]step{
				{
					actions:   []action{{"move", fileMenuPos}},
					wantImage: "menu_bar_inactive_file.png",
				},
				{
					actions:   []action{{"tap", fileMenuPos}},
					wantImage: "menu_bar_active_file.png",
				},
				{
					actions:   []action{{"tap", fileMenuPos}},
					wantImage: "menu_bar_inactive_file.png",
				},
				{
					actions:   []action{{"move", buttonPos}},
					wantImage: "menu_bar_hovered_content.png",
				},
			},
		},
		"active menu deactivates content": {
			[]step{
				{
					actions:   []action{{"move", fileMenuPos}},
					wantImage: "menu_bar_inactive_file.png",
				},
				{
					actions:   []action{{"tap", fileMenuPos}},
					wantImage: "menu_bar_active_file.png",
				},
				{
					actions:   []action{{"move", buttonPos}},
					wantImage: "menu_bar_content_not_hoverable_with_active_menu.png",
				},
				{
					// TODO: should be hovered over content here
					// -> is this canvas logic?
					// -> it would be for overlays (probably the same issue present)
					// -> it would not be for current menu implementation (menu is not an overlay, canvas does not know about activation state)
					actions:   []action{{"tap", buttonPos}},
					wantImage: "menu_bar_tap_content_with_active_menu_does_not_trigger_action_but_dismisses_menu.png",
				},
				{
					// menu bar is inactive again (menu not shown at hover)
					actions:   []action{{"move", fileMenuPos}},
					wantImage: "menu_bar_inactive_file.png",
				},
				{
					actions:   []action{{"move", buttonPos}},
					wantImage: "menu_bar_hovered_content.png",
				},
			},
		},
		"menu action File->New": {
			[]step{
				{
					actions:   []action{{"tap", fileMenuPos}},
					wantImage: "menu_bar_active_file.png",
				},
				{
					actions:   []action{{"move", fileNewPos}},
					wantImage: "menu_bar_hovered_file_new.png",
				},
				{
					actions:    []action{{"tap", fileNewPos}},
					wantAction: "new",
					wantImage:  "menu_bar_initial.png",
				},
			},
		},
		"menu action File->Open": {
			[]step{
				{
					actions: []action{
						{"tap", fileMenuPos},
						{"move", fileOpenPos},
					},
					wantImage: "menu_bar_hovered_file_open.png",
				},
				{
					actions:    []action{{"tap", fileOpenPos}},
					wantAction: "open",
					wantImage:  "menu_bar_initial.png",
				},
			},
		},
		"menu action File->Recent->Older->Old 1": {
			[]step{
				{
					actions: []action{
						{"tap", fileMenuPos},
						{"move", fileRecentPos},
					},
					wantImage: "menu_bar_hovered_file_recent.png",
				},
				{
					actions:   []action{{"move", fileRecentOlderPos}},
					wantImage: "menu_bar_hovered_file_recent_older.png",
				},
				{
					actions:   []action{{"move", fileRecentOlderOld1Pos}},
					wantImage: "menu_bar_hovered_file_recent_older_old1.png",
				},
				{
					actions:    []action{{"tap", fileRecentOlderOld1Pos}},
					wantAction: "old 1",
					wantImage:  "menu_bar_initial.png",
				},
			},
		},
		"move mouse outside does not hide menu": {
			[]step{
				{
					actions: []action{
						{"tap", fileMenuPos},
						{"move", fileRecentPos},
						{"move", fileRecentOlderPos},
						{"move", fileRecentOlderOld1Pos},
					},
					wantImage: "menu_bar_hovered_file_recent_older_old1.png",
				},
				{
					actions:   []action{{"move", buttonPos}},
					wantImage: "menu_bar_hovered_file_recent_older.png",
				},
			},
		},
		"hover other menu item hides previous submenus": {
			[]step{
				{
					actions: []action{
						{"tap", fileMenuPos},
						{"move", fileRecentPos},
						{"move", fileRecentOlderPos},
						{"move", fileRecentOlderOld1Pos},
					},
					wantImage: "menu_bar_hovered_file_recent_older_old1.png",
				},
				{
					actions:   []action{{"move", fileNewPos}},
					wantImage: "menu_bar_hovered_file_new.png",
				},
				{
					actions:   []action{{"move", fileRecentPos}},
					wantImage: "menu_bar_hovered_file_recent.png",
				},
			},
		},
		"hover other menu bar item changes active menu": {
			[]step{
				{
					actions:   []action{{"tap", fileMenuPos}},
					wantImage: "menu_bar_active_file.png",
				},
				{
					actions:   []action{{"move", editMenuPos}},
					wantImage: "menu_bar_active_edit.png",
				},
				{
					actions:   []action{{"move", helpMenuPos}},
					wantImage: "menu_bar_active_help.png",
				},
			},
		},
		"hover other menu bar item hides previous submenus": {
			[]step{
				{
					actions: []action{
						{"tap", fileMenuPos},
						{"move", fileRecentPos},
						{"move", fileRecentOlderPos},
						{"move", fileRecentOlderOld1Pos},
					},
					wantImage: "menu_bar_hovered_file_recent_older_old1.png",
				},
				{
					actions:   []action{{"move", helpMenuPos}},
					wantImage: "menu_bar_active_help.png",
				},
				{
					actions:   []action{{"move", fileMenuPos}},
					wantImage: "menu_bar_active_file.png",
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			test.MoveMouse(c, fyne.NewPos(0, 0))
			test.TapCanvas(c, fyne.NewPos(0, 0))
			test.AssertImageMatches(t, "menu_bar_initial.png", c.Capture())
			for i, s := range tt.steps {
				t.Run("step "+strconv.Itoa(i+1), func(t *testing.T) {
					lastAction = ""
					for _, a := range s.actions {
						switch a.typ {
						case "move":
							test.MoveMouse(c, a.pos)
						case "tap":
							test.MoveMouse(c, a.pos)
							test.TapCanvas(c, a.pos)
						}
					}
					test.AssertImageMatches(t, s.wantImage, c.Capture())
					assert.Equal(t, s.wantAction, lastAction, "last action should match expected")
				})
			}
		})
	}
}
