// +build !mobile
// +build !ci !windows

package glfw_test

import (
	"strconv"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/glfw"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenuBar(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

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

	t.Run("mouse control and basic behaviour", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(fyne.NewSize(300, 300))
		c := w.Canvas()

		menuBar := glfw.NewMenuBar(menu, c)
		themeCounter := 0
		button := widget.NewButton("Button", func() {
			switch themeCounter % 2 {
			case 0:
				test.ApplyTheme(t, test.NewTheme())
			case 1:
				test.ApplyTheme(t, test.Theme())
			}
			themeCounter++
		})
		container := fyne.NewContainerWithoutLayout(button, menuBar)
		w.SetContent(container)
		w.Resize(fyne.NewSize(300, 300))
		button.Resize(button.MinSize())
		button.Move(fyne.NewPos(100, 50))
		menuBar.Resize(fyne.NewSize(300, 0).Max(menuBar.MinSize()))

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
				if test.AssertImageMatches(t, "menu_bar_initial.png", c.Capture()) {
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
				}
			})
		}
	})

	t.Run("keyboard control", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(fyne.NewSize(300, 300))
		c := w.Canvas()

		menuBar := glfw.NewMenuBar(menu, c)
		themeCounter := 0
		button := widget.NewButton("Button", func() {
			switch themeCounter % 2 {
			case 0:
				test.ApplyTheme(t, test.NewTheme())
			case 1:
				test.ApplyTheme(t, test.Theme())
			}
			themeCounter++
		})
		container := fyne.NewContainerWithoutLayout(button, menuBar)
		w.SetContent(container)
		w.Resize(fyne.NewSize(300, 300))
		button.Resize(button.MinSize())
		button.Move(fyne.NewPos(100, 50))
		menuBar.Resize(fyne.NewSize(300, 0).Max(menuBar.MinSize()))

		fileMenuPos := fyne.NewPos(20, 10)
		for name, tt := range map[string]struct {
			keys       []fyne.KeyName
			wantAction string
		}{
			"traverse_menu_bar_items_right_1": {
				keys: []fyne.KeyName{fyne.KeyRight},
			},
			"traverse_menu_bar_items_right_2": {
				keys: []fyne.KeyName{fyne.KeyRight, fyne.KeyRight},
			},
			"traverse_menu_bar_items_right_3": {
				keys: []fyne.KeyName{fyne.KeyRight, fyne.KeyRight, fyne.KeyRight},
			},
			"traverse_menu_bar_items_left_1": {
				keys: []fyne.KeyName{fyne.KeyLeft},
			},
			"traverse_menu_bar_items_left_2": {
				keys: []fyne.KeyName{fyne.KeyLeft, fyne.KeyLeft},
			},
			"traverse_menu_bar_items_left_3": {
				keys: []fyne.KeyName{fyne.KeyLeft, fyne.KeyLeft, fyne.KeyLeft},
			},
			"traverse_menu_down_1": {
				keys: []fyne.KeyName{fyne.KeyDown},
			},
			"traverse_menu_down_2": {
				keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown},
			},
			"traverse_menu_down_3": {
				keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown},
			},
			"traverse_menu_up_1": {
				keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyUp},
			},
			"traverse_menu_up_2": {
				keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyUp, fyne.KeyUp},
			},
			"open_submenu_1": {
				keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight},
			},
			"open_submenu_2": {
				keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight},
			},
			"close_submenu_1": {
				keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyLeft},
			},
			"close_submenu_2": {
				keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyLeft, fyne.KeyLeft},
			},
			"trigger_with_enter": {
				keys:       []fyne.KeyName{fyne.KeyDown, fyne.KeyEnter},
				wantAction: "new",
			},
			"trigger_with_return": {
				keys:       []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyReturn},
				wantAction: "open",
			},
			"trigger_with_space": {
				keys:       []fyne.KeyName{fyne.KeyRight, fyne.KeyDown, fyne.KeySpace},
				wantAction: "copy",
			},
			"trigger_submenu_item": {
				keys:       []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyReturn},
				wantAction: "old 1",
			},
			"trigger_without_active_item": {
				keys:       []fyne.KeyName{fyne.KeyEnter},
				wantAction: "",
			},
		} {
			t.Run(name, func(t *testing.T) {
				test.MoveMouse(c, fyne.NewPos(0, 0))
				test.TapCanvas(c, fyne.NewPos(0, 0))
				test.TapCanvas(c, fileMenuPos) // activate menu
				require.Equal(t, menuBar.Items[0], c.Focused())
				if test.AssertImageMatches(t, "menu_bar_active_file.png", c.Capture()) {
					lastAction = ""
					for _, key := range tt.keys {
						c.Focused().TypedKey(&fyne.KeyEvent{
							Name: key,
						})
					}
					test.AssertRendersToMarkup(t, "menu_bar_kbdctrl_"+name+".xml", c)
					assert.Equal(t, tt.wantAction, lastAction, "last action should match expected")
				}
			})
		}

		t.Run("moving mouse over unfocused item moves focus", func(t *testing.T) {
			test.MoveMouse(c, fyne.NewPos(0, 0))
			test.TapCanvas(c, fyne.NewPos(0, 0))
			test.MoveMouse(c, fileMenuPos)
			test.TapCanvas(c, fileMenuPos) // activate menu
			require.Equal(t, menuBar.Items[0], c.Focused())
			c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
			c.Focused().TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
			require.Equal(t, menuBar.Items[2], c.Focused())

			test.MoveMouse(c, fileMenuPos.Add(fyne.NewPos(1, 0)))
			assert.Equal(t, menuBar.Items[0], c.Focused())
			test.AssertImageMatches(t, "menu_bar_active_file.png", c.Capture())
		})
	})
}

func TestMenuBar_Toggle(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	m1i1 := fyne.NewMenuItem("New", nil)
	m1i2 := fyne.NewMenuItem("Open", nil)
	m1 := fyne.NewMenu("File", m1i1, m1i2)

	m2i1 := fyne.NewMenuItem("Copy", nil)
	m2i2 := fyne.NewMenuItem("Paste", nil)
	m2 := fyne.NewMenu("Edit", m2i1, m2i2)

	menu := fyne.NewMainMenu(m1, m2)

	t.Run("when menu bar is inactive", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(fyne.NewSize(300, 300))
		c := w.Canvas()
		menuBar := glfw.NewMenuBar(menu, c)
		w.SetContent(fyne.NewContainerWithoutLayout(menuBar))
		w.Resize(fyne.NewSize(300, 300))
		menuBar.Resize(fyne.NewSize(300, 0).Max(menuBar.MinSize()))

		require.False(t, menuBar.IsActive())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_deactivated.xml", c)

		menuBar.Toggle()
		assert.True(t, menuBar.IsActive())
		assert.Equal(t, c.Focused(), menuBar.Items[0])
		test.AssertRendersToMarkup(t, "menu_bar_toggle_first_item_active.xml", c)
	})

	t.Run("when menu bar is active (first menu item active)", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(fyne.NewSize(300, 300))
		c := w.Canvas()
		menuBar := glfw.NewMenuBar(menu, c)
		w.SetContent(fyne.NewContainerWithoutLayout(menuBar))
		w.Resize(fyne.NewSize(300, 300))
		menuBar.Resize(fyne.NewSize(300, 0).Max(menuBar.MinSize()))

		menuBar.Toggle()
		require.True(t, menuBar.IsActive())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_first_item_active.xml", c)

		menuBar.Toggle()
		assert.False(t, menuBar.IsActive())
		assert.Nil(t, c.Focused())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_deactivated.xml", c)
	})

	t.Run("when menu bar is active (second menu item active)", func(t *testing.T) {
		w := test.NewWindow(nil)
		defer w.Close()
		w.SetPadded(false)
		w.Resize(fyne.NewSize(300, 300))
		c := w.Canvas()
		menuBar := glfw.NewMenuBar(menu, c)
		w.SetContent(fyne.NewContainerWithoutLayout(menuBar))
		w.Resize(fyne.NewSize(300, 300))
		menuBar.Resize(fyne.NewSize(300, 0).Max(menuBar.MinSize()))

		menuBar.Toggle()
		c.(test.WindowlessCanvas).FocusNext()
		require.True(t, menuBar.IsActive())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_second_item_active.xml", c)

		menuBar.Toggle()
		assert.False(t, menuBar.IsActive())
		assert.Nil(t, c.Focused())
		test.AssertRendersToMarkup(t, "menu_bar_toggle_deactivated.xml", c)
	})
}
