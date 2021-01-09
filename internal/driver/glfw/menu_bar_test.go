// +build !mobile
// +build !ci !windows

package glfw_test

import (
	"strconv"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/glfw"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"

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

		markupMenuBarPrefix := `
			<canvas size="300x300">
				<content>
					<container size="300x300">
						<widget pos="100,50" size="77x37" type="*widget.Button">
							<widget pos="2,2" size="73x33" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-2,-2" size="2x2" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-2" size="73x2"/>
								<radialGradient centerOffset="-0.5,0.5" pos="73,-2" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" pos="73,0" size="2x33" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="73,33" size="2x2" startColor="shadow"/>
								<linearGradient pos="0,33" size="73x2" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-2,33" size="2x2" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-2,0" size="2x33"/>
							</widget>
							<rectangle fillColor="button" pos="2,2" size="73x33"/>
							<text bold pos="12,8" size="53x21">Button</text>
						</widget>
						<widget backgroundColor="button" size="300x29" type="*glfw.MenuBar">
							<widget size="300x29" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-4" size="300x4"/>
								<radialGradient centerOffset="-0.5,0.5" pos="300,-4" size="4x4" startColor="shadow"/>
								<linearGradient angle="270" pos="300,0" size="4x29" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="300,29" size="4x4" startColor="shadow"/>
								<linearGradient pos="0,29" size="300x4" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-4,29" size="4x4" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x29"/>
							</widget>`
		markupMenuBarActivated := `
							<widget size="300x300" type="*glfw.menuBarBackground">
							</widget>`
		markupMenuBarDeactivated := `
							<widget size="0x0" type="*glfw.menuBarBackground">
							</widget>`
		markupMenuBarEditActive := `
							<container pos="4,0" size="292x29">
								<widget backgroundColor="focus" size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>`
		markupFileMenuPrefix := `
							<widget pos="4,29" size="85x103" type="*widget.Menu">
								<widget size="85x103" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="85x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="85,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="85,0" size="4x103" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="85,103" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,103" size="85x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
								</widget>
								<widget size="85x103" type="*widget.ScrollContainer">
									<widget size="85x103" type="*widget.menuBox">
										<container pos="0,4" size="85x111">`
		markupFileMenuItems := `
											<widget size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>`
		markupFileMenuSuffix := `
										</container>
									</widget>
								</widget>
							</widget>`
		markupFileMenu := markupFileMenuPrefix + markupFileMenuItems + markupFileMenuSuffix
		markupSubmenu1Prefix := `
								<widget pos="85,66" size="76x103" type="*widget.Menu">
									<widget size="76x103" type="*widget.Shadow">
										<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
										<linearGradient endColor="shadow" pos="0,-4" size="76x4"/>
										<radialGradient centerOffset="-0.5,0.5" pos="76,-4" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" pos="76,0" size="4x103" startColor="shadow"/>
										<radialGradient centerOffset="-0.5,-0.5" pos="76,103" size="4x4" startColor="shadow"/>
										<linearGradient pos="0,103" size="76x4" startColor="shadow"/>
										<radialGradient centerOffset="0.5,-0.5" pos="-4,103" size="4x4" startColor="shadow"/>
										<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x103"/>
									</widget>
									<widget size="76x103" type="*widget.ScrollContainer">
										<widget size="76x103" type="*widget.menuBox">
											<container pos="0,4" size="76x111">`
		markupSubmenu1Suffix := `
											</container>
										</widget>
									</widget>
								</widget>`
		markupSubmenu2Prefix := `
									<widget pos="76,66" size="54x70" type="*widget.Menu">
										<widget size="54x70" type="*widget.Shadow">
											<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
											<linearGradient endColor="shadow" pos="0,-4" size="54x4"/>
											<radialGradient centerOffset="-0.5,0.5" pos="54,-4" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" pos="54,0" size="4x70" startColor="shadow"/>
											<radialGradient centerOffset="-0.5,-0.5" pos="54,70" size="4x4" startColor="shadow"/>
											<linearGradient pos="0,70" size="54x4" startColor="shadow"/>
											<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
											<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
										</widget>
										<widget size="54x70" type="*widget.ScrollContainer">
											<widget size="54x70" type="*widget.menuBox">
												<container pos="0,4" size="54x78">`
		markupSubmenu2Suffix := `
												</container>
											</widget>
										</widget>
									</widget>`
		markupEditMenu := `
							<widget pos="49,29" size="55x70" type="*widget.Menu">
								<widget size="55x70" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="55x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="55,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="55,0" size="4x70" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="55,70" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,70" size="55x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
								</widget>
								<widget size="55x70" type="*widget.ScrollContainer">
									<widget size="55x70" type="*widget.menuBox">
										<container pos="0,4" size="55x78">
											<widget size="55x29" type="*widget.menuItem">
												<text pos="8,4" size="36x21">Copy</text>
											</widget>
											<widget pos="0,33" size="55x29" type="*widget.menuItem">
												<text pos="8,4" size="39x21">Paste</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>`
		markupHelpMenu := `
							<widget pos="97,29" size="54x37" type="*widget.Menu">
								<widget size="54x37" type="*widget.Shadow">
									<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
									<linearGradient endColor="shadow" pos="0,-4" size="54x4"/>
									<radialGradient centerOffset="-0.5,0.5" pos="54,-4" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" pos="54,0" size="4x37" startColor="shadow"/>
									<radialGradient centerOffset="-0.5,-0.5" pos="54,37" size="4x4" startColor="shadow"/>
									<linearGradient pos="0,37" size="54x4" startColor="shadow"/>
									<radialGradient centerOffset="0.5,-0.5" pos="-4,37" size="4x4" startColor="shadow"/>
									<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x37"/>
								</widget>
								<widget size="54x37" type="*widget.ScrollContainer">
									<widget size="54x37" type="*widget.menuBox">
										<container pos="0,4" size="54x45">
											<widget size="54x29" type="*widget.menuItem">
												<text pos="8,4" size="38x21">Help!</text>
											</widget>
										</container>
									</widget>
								</widget>
							</widget>`
		markupMenuBarSuffix := `
						</widget>
					</container>
				</content>
			</canvas>
		`
		fileMenuPos := fyne.NewPos(20, 10)
		type step struct {
			keys       []fyne.KeyName
			wantMarkup string
			wantAction string
		}
		for name, tt := range map[string]struct {
			steps []step
		}{
			"traverse menu bar items right #1": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyRight},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget backgroundColor="focus" pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupEditMenu +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu bar items right #2": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyRight, fyne.KeyRight},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget backgroundColor="focus" pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupHelpMenu +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu bar items right #3": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyRight, fyne.KeyRight, fyne.KeyRight},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated + `
							<container pos="4,0" size="292x29">
								<widget backgroundColor="focus" size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupFileMenu +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu bar items left #1": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyLeft},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget backgroundColor="focus" pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupHelpMenu +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu bar items left #2": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyLeft, fyne.KeyLeft},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget backgroundColor="focus" pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupEditMenu +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu bar items left #3": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyLeft, fyne.KeyLeft, fyne.KeyLeft},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated + `
							<container pos="4,0" size="292x29">
								<widget backgroundColor="focus" size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupFileMenu +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu down #1": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget backgroundColor="focus" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>` +
							markupFileMenuSuffix +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu down #2": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget backgroundColor="focus" pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>` +
							markupFileMenuSuffix +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu down #3": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget backgroundColor="focus" pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>` +
							markupFileMenuSuffix +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu up #1": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyUp},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget backgroundColor="focus" pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>` +
							markupFileMenuSuffix +
							markupMenuBarSuffix,
					},
				},
			},
			"traverse menu up #2": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyUp, fyne.KeyUp},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget backgroundColor="focus" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>` +
							markupFileMenuSuffix +
							markupMenuBarSuffix,
					},
				},
			},
			"open submenu #1": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget backgroundColor="focus" pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>
										</container>
									</widget>
								</widget>` +
							markupSubmenu1Prefix + `
												<widget backgroundColor="focus" size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="37x21">File 1</text>
												</widget>
												<widget pos="0,33" size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="37x21">File 2</text>
												</widget>
												<widget pos="0,66" size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="40x21">Older</text>
													<image pos="56,4" rsc="menuExpandIcon" size="iconInlineSize"/>
												</widget>` +
							markupSubmenu1Suffix + `
							</widget>` +
							markupMenuBarSuffix,
					},
				},
			},
			"open submenu #2": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget backgroundColor="focus" pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>
										</container>
									</widget>
								</widget>` +
							markupSubmenu1Prefix + `
												<widget size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="37x21">File 1</text>
												</widget>
												<widget pos="0,33" size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="37x21">File 2</text>
												</widget>
												<widget backgroundColor="focus" pos="0,66" size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="40x21">Older</text>
													<image pos="56,4" rsc="menuExpandIcon" size="iconInlineSize"/>
												</widget>
											</container>
										</widget>
									</widget>` +
							markupSubmenu2Prefix + `
													<widget backgroundColor="focus" size="54x29" type="*widget.menuItem">
														<text pos="8,4" size="38x21">Old 1</text>
													</widget>
													<widget pos="0,33" size="54x29" type="*widget.menuItem">
														<text pos="8,4" size="38x21">Old 2</text>
													</widget>` + markupSubmenu2Suffix + `
								</widget>
							</widget>` +
							markupMenuBarSuffix,
					},
				},
			},
			"close submenu #1": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyLeft},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget backgroundColor="focus" pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>
										</container>
									</widget>
								</widget>` +
							markupSubmenu1Prefix + `
												<widget size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="37x21">File 1</text>
												</widget>
												<widget pos="0,33" size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="37x21">File 2</text>
												</widget>
												<widget backgroundColor="focus" pos="0,66" size="76x29" type="*widget.menuItem">
													<text pos="8,4" size="40x21">Older</text>
													<image pos="56,4" rsc="menuExpandIcon" size="iconInlineSize"/>
												</widget>` +
							markupSubmenu1Suffix + `
							</widget>` +
							markupMenuBarSuffix,
					},
				},
			},
			"close submenu #2": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyLeft, fyne.KeyLeft},
						wantMarkup: markupMenuBarPrefix + markupMenuBarActivated +
							markupMenuBarEditActive +
							markupFileMenuPrefix + `
											<widget size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="33x21">New</text>
											</widget>
											<widget pos="0,33" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="40x21">Open</text>
											</widget>
											<widget backgroundColor="focus" pos="0,66" size="85x29" type="*widget.menuItem">
												<text pos="8,4" size="49x21">Recent</text>
												<image pos="65,4" rsc="menuExpandIcon" size="iconInlineSize"/>
											</widget>` +
							markupFileMenuSuffix +
							markupMenuBarSuffix,
					},
				},
			},
			"trigger with Enter": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyEnter},
						wantMarkup: markupMenuBarPrefix + markupMenuBarDeactivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupMenuBarSuffix,
						wantAction: "new",
					},
				},
			},
			"trigger with Return": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyReturn},
						wantMarkup: markupMenuBarPrefix + markupMenuBarDeactivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupMenuBarSuffix,
						wantAction: "open",
					},
				},
			},
			"trigger with Space": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyRight, fyne.KeyDown, fyne.KeySpace},
						wantMarkup: markupMenuBarPrefix + markupMenuBarDeactivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupMenuBarSuffix,
						wantAction: "copy",
					},
				},
			},
			"trigger submenu item": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyDown, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyDown, fyne.KeyDown, fyne.KeyRight, fyne.KeyReturn},
						wantMarkup: markupMenuBarPrefix + markupMenuBarDeactivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupMenuBarSuffix,
						wantAction: "old 1",
					},
				},
			},
			"trigger without active item": {
				[]step{
					{
						keys: []fyne.KeyName{fyne.KeyEnter},
						wantMarkup: markupMenuBarPrefix + markupMenuBarDeactivated + `
							<container pos="4,0" size="292x29">
								<widget size="41x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="25x21">File</text>
								</widget>
								<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="28x21">Edit</text>
								</widget>
								<widget pos="93,0" size="50x29" type="*glfw.menuBarItem">
									<text pos="8,4" size="34x21">Help</text>
								</widget>
							</container>` +
							markupMenuBarSuffix,
						wantAction: "",
					},
				},
			},
		} {
			t.Run(name, func(t *testing.T) {
				test.MoveMouse(c, fyne.NewPos(0, 0))
				test.TapCanvas(c, fyne.NewPos(0, 0))
				test.TapCanvas(c, fileMenuPos) // activate menu
				require.Equal(t, menuBar.Items[0], c.Focused())
				if test.AssertImageMatches(t, "menu_bar_active_file.png", c.Capture()) {
					for i, s := range tt.steps {
						t.Run("step "+strconv.Itoa(i+1), func(t *testing.T) {
							lastAction = ""
							for _, key := range s.keys {
								c.Focused().TypedKey(&fyne.KeyEvent{
									Name: key,
								})
							}
							test.AssertRendersToMarkup(t, s.wantMarkup, c)
							assert.Equal(t, s.wantAction, lastAction, "last action should match expected")
						})
					}
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

	markupMenuBarPrefix := `
		<canvas size="300x300">
			<content>
				<container size="300x300">
					<widget backgroundColor="button" size="300x29" type="*glfw.MenuBar">
						<widget size="300x29" type="*widget.Shadow">
							<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
							<linearGradient endColor="shadow" pos="0,-4" size="300x4"/>
							<radialGradient centerOffset="-0.5,0.5" pos="300,-4" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" pos="300,0" size="4x29" startColor="shadow"/>
							<radialGradient centerOffset="-0.5,-0.5" pos="300,29" size="4x4" startColor="shadow"/>
							<linearGradient pos="0,29" size="300x4" startColor="shadow"/>
							<radialGradient centerOffset="0.5,-0.5" pos="-4,29" size="4x4" startColor="shadow"/>
							<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x29"/>
						</widget>`
	markupMenuBarActivated := `
						<widget size="300x300" type="*glfw.menuBarBackground">
						</widget>`
	markupMenuBarDeactivated := `
						<widget size="0x0" type="*glfw.menuBarBackground">
						</widget>`
	markupMenuBar := `
						<container pos="4,0" size="292x29">
							<widget size="41x29" type="*glfw.menuBarItem">
								<text pos="8,4" size="25x21">File</text>
							</widget>
							<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
								<text pos="8,4" size="28x21">Edit</text>
							</widget>
						</container>`
	markupMenuBarFileActive := `
						<container pos="4,0" size="292x29">
							<widget backgroundColor="focus" size="41x29" type="*glfw.menuBarItem">
								<text pos="8,4" size="25x21">File</text>
							</widget>
							<widget pos="45,0" size="44x29" type="*glfw.menuBarItem">
								<text pos="8,4" size="28x21">Edit</text>
							</widget>
						</container>`
	markupMenuBarEditActive := `
						<container pos="4,0" size="292x29">
							<widget size="41x29" type="*glfw.menuBarItem">
								<text pos="8,4" size="25x21">File</text>
							</widget>
							<widget backgroundColor="focus" pos="45,0" size="44x29" type="*glfw.menuBarItem">
								<text pos="8,4" size="28x21">Edit</text>
							</widget>
						</container>`
	markupFileMenu := `
						<widget pos="4,29" size="56x70" type="*widget.Menu">
							<widget size="56x70" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-4" size="56x4"/>
								<radialGradient centerOffset="-0.5,0.5" pos="56,-4" size="4x4" startColor="shadow"/>
								<linearGradient angle="270" pos="56,0" size="4x70" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="56,70" size="4x4" startColor="shadow"/>
								<linearGradient pos="0,70" size="56x4" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
							</widget>
							<widget size="56x70" type="*widget.ScrollContainer">
								<widget size="56x70" type="*widget.menuBox">
									<container pos="0,4" size="56x78">
										<widget size="56x29" type="*widget.menuItem">
											<text pos="8,4" size="33x21">New</text>
										</widget>
										<widget pos="0,33" size="56x29" type="*widget.menuItem">
											<text pos="8,4" size="40x21">Open</text>
										</widget>
									</container>
								</widget>
							</widget>
						</widget>`
	markupEditMenu := `
						<widget pos="49,29" size="55x70" type="*widget.Menu">
							<widget size="55x70" type="*widget.Shadow">
								<radialGradient centerOffset="0.5,0.5" pos="-4,-4" size="4x4" startColor="shadow"/>
								<linearGradient endColor="shadow" pos="0,-4" size="55x4"/>
								<radialGradient centerOffset="-0.5,0.5" pos="55,-4" size="4x4" startColor="shadow"/>
								<linearGradient angle="270" pos="55,0" size="4x70" startColor="shadow"/>
								<radialGradient centerOffset="-0.5,-0.5" pos="55,70" size="4x4" startColor="shadow"/>
								<linearGradient pos="0,70" size="55x4" startColor="shadow"/>
								<radialGradient centerOffset="0.5,-0.5" pos="-4,70" size="4x4" startColor="shadow"/>
								<linearGradient angle="270" endColor="shadow" pos="-4,0" size="4x70"/>
							</widget>
							<widget size="55x70" type="*widget.ScrollContainer">
								<widget size="55x70" type="*widget.menuBox">
									<container pos="0,4" size="55x78">
										<widget size="55x29" type="*widget.menuItem">
											<text pos="8,4" size="36x21">Copy</text>
										</widget>
										<widget pos="0,33" size="55x29" type="*widget.menuItem">
											<text pos="8,4" size="39x21">Paste</text>
										</widget>
									</container>
								</widget>
							</widget>
						</widget>`
	markupMenuBarSuffix := `
					</widget>
				</container>
			</content>
		</canvas>
	`

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

		menuBar.Toggle()
		assert.True(t, menuBar.IsActive())
		assert.Equal(t, c.Focused(), menuBar.Items[0])
		test.AssertRendersToMarkup(t, markupMenuBarPrefix+markupMenuBarActivated+markupMenuBarFileActive+markupFileMenu+markupMenuBarSuffix, c)
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

		menuBar.Toggle()
		assert.False(t, menuBar.IsActive())
		assert.Nil(t, c.Focused())
		test.AssertRendersToMarkup(t, markupMenuBarPrefix+markupMenuBarDeactivated+markupMenuBar+markupMenuBarSuffix, c)
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
		test.AssertRendersToMarkup(t, markupMenuBarPrefix+markupMenuBarActivated+markupMenuBarEditActive+markupEditMenu+markupMenuBarSuffix, c)

		menuBar.Toggle()
		assert.False(t, menuBar.IsActive())
		assert.Nil(t, c.Focused())
		test.AssertRendersToMarkup(t, markupMenuBarPrefix+markupMenuBarDeactivated+markupMenuBar+markupMenuBarSuffix, c)
	})
}
