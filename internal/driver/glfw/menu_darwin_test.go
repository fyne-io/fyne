//go:build !ci && !no_native_menus && !mobile
// +build !ci,!no_native_menus,!mobile

package glfw

import (
	"testing"
	"unsafe"

	"fyne.io/fyne/v2/driver/desktop"
	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestDarwinMenu(t *testing.T) {
	setExceptionCallback(func(msg string) { t.Error("Obj-C exception:", msg) })
	defer setExceptionCallback(nil)

	resetMainMenu()

	w := createWindow("Test").(*window)

	var lastAction string
	assertLastAction := func(wantAction string) {
		w.WaitForEvents()
		assert.Equal(t, wantAction, lastAction)
	}

	assertNSMenuItemSeparator := func(m unsafe.Pointer, i int) {
		item := testNSMenuItemAtIndex(m, i)
		assert.True(t, testNSMenuItemIsSeparatorItem(item), "item is expected to be a separator")
	}

	itemNew := fyne.NewMenuItem("New", func() { lastAction = "new" })
	itemNew.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierShortcutDefault}
	itemOpen := fyne.NewMenuItem("Open", func() { lastAction = "open" })
	itemOpen.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierAlt}
	itemRecent := fyne.NewMenuItem("Recent", nil)
	itemFoo := fyne.NewMenuItem("Foo", func() { lastAction = "foo" })
	itemRecent.ChildMenu = fyne.NewMenu("", itemFoo)
	menuEdit := fyne.NewMenu("File", itemNew, itemOpen, fyne.NewMenuItemSeparator(), itemRecent)

	itemHelp := fyne.NewMenuItem("Help", func() { lastAction = "Help!!!" })
	itemHelp.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyH, Modifier: fyne.KeyModifierControl}
	itemHelpMe := fyne.NewMenuItem("Help Me", func() { lastAction = "Help me!!!" })
	itemHelpMe.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyH, Modifier: fyne.KeyModifierShift}
	menuHelp := fyne.NewMenu("Help", itemHelp, itemHelpMe)

	itemHelloWorld := fyne.NewMenuItem("Hello World", func() { lastAction = "Hello World!" })
	itemHelloWorld.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyH, Modifier: fyne.KeyModifierControl | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierSuper}
	itemPrefs := fyne.NewMenuItem("Preferences", func() { lastAction = "prefs" })
	itemMore := fyne.NewMenuItem("More", func() { lastAction = "more" })
	itemMorePrefs := fyne.NewMenuItem("Preferences…", func() { lastAction = "more prefs" })
	menuMore := fyne.NewMenu("More Stuff", itemHelloWorld, itemPrefs, itemMore, itemMorePrefs)

	itemSettings := fyne.NewMenuItem("Settings", func() { lastAction = "settings" })
	itemMoreSetings := fyne.NewMenuItem("Settings…", func() { lastAction = "more settings" })
	menuSettings := fyne.NewMenu("Settings", itemSettings, fyne.NewMenuItemSeparator(), itemMoreSetings)

	mainMenu := fyne.NewMainMenu(menuEdit, menuHelp, menuMore, menuSettings)
	setupNativeMenu(w, mainMenu)

	mm := testDarwinMainMenu()
	// The custom “Preferences” menu should be moved to the system app menu completely.
	// -> only three custom menus
	assert.Equal(t, 5, testNSMenuNumberOfItems(mm), "two built-in + three custom")

	m := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 0))
	assert.Equal(t, "", testNSMenuTitle(m), "app menu doesn’t have a title")
	assertNSMenuItemSeparator(m, 1)
	assertNSMenuItem(t, "Preferences", "", 0, m, 2)
	assertLastAction("prefs")
	assertNSMenuItem(t, "Preferences…", "", 0, m, 3)
	assertLastAction("more prefs")
	assertNSMenuItemSeparator(m, 4)
	assertNSMenuItem(t, "Settings", "", 0, m, 5)
	assertLastAction("settings")
	assertNSMenuItem(t, "Settings…", "", 0, m, 6)
	assertLastAction("more settings")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
	assert.Equal(t, "File", testNSMenuTitle(m))
	assert.Equal(t, 4, testNSMenuNumberOfItems(m))
	// NSEventModifierFlagCommand = 1 << 20
	assertNSMenuItem(t, "New", "n", 0b100000000000000000000, m, 0)
	assertLastAction("new")
	// NSEventModifierFlagOption = 1 << 19
	assertNSMenuItem(t, "Open", "o", 0b10000000000000000000, m, 1)
	assertLastAction("open")
	assertNSMenuItemSeparator(m, 2)
	i := testNSMenuItemAtIndex(m, 3)
	assert.Equal(t, "Recent", testNSMenuItemTitle(i))
	sm := testNSMenuItemSubmenu(i)
	assert.NotNil(t, sm, "item has submenu")
	assert.Equal(t, 1, testNSMenuNumberOfItems(sm))
	assertNSMenuItem(t, "Foo", "", 0, sm, 0)
	assertLastAction("foo")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 2))
	assert.Equal(t, "More Stuff", testNSMenuTitle(m))
	assert.Equal(t, 2, testNSMenuNumberOfItems(m))
	assertNSMenuItem(t, "Hello World", "h", 0b111100000000000000000, m, 0)
	assertLastAction("Hello World!")
	assertNSMenuItem(t, "More", "", 0, m, 1)
	assertLastAction("more")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 3))
	assert.Equal(t, "Window", testNSMenuTitle(m))

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 4))
	assert.Equal(t, "Help", testNSMenuTitle(m))
	assert.Equal(t, 2, testNSMenuNumberOfItems(m))
	// NSEventModifierFlagControl = 1 << 18
	assertNSMenuItem(t, "Help", "h", 0b1000000000000000000, m, 0)
	assertLastAction("Help!!!")
	// NSEventModifierFlagShift = 1 << 17
	assertNSMenuItem(t, "Help Me", "h", 0b100000000000000000, m, 1)
	assertLastAction("Help me!!!")

	// change action works
	itemOpen.Action = func() { lastAction = "new open" }
	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
	assertNSMenuItem(t, "Open", "", 0, m, 1)
	assertLastAction("new open")
}

var initialAppMenuItems []string
var initialMenus []string

func assertNSMenuItem(t *testing.T, wantTitle, wantKey string, wantModifier uint64, m unsafe.Pointer, i int) {
	item := testNSMenuItemAtIndex(m, i)
	assert.Equal(t, wantTitle, testNSMenuItemTitle(item))
	if wantKey != "" {
		assert.Equal(t, wantKey, testNSMenuItemKeyEquivalent(item))
		assert.Equal(t, wantModifier, testNSMenuItemKeyEquivalentModifierMask(item))
	}
	testNSMenuPerformActionForItemAtIndex(m, i)
}

func initMainMenu() {
	createWindow("Test").Close() // ensure GLFW has performed [NSApp run]
	mainMenu := testDarwinMainMenu()
	for i := 0; i < testNSMenuNumberOfItems(mainMenu); i++ {
		menu := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mainMenu, i))
		initialMenus = append(initialMenus, testNSMenuTitle(menu))
	}
	appMenu := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mainMenu, 0))
	for i := 0; i < testNSMenuNumberOfItems(appMenu); i++ {
		item := testNSMenuItemAtIndex(appMenu, i)
		initialAppMenuItems = append(initialAppMenuItems, testNSMenuItemTitle(item))
	}
}

func resetMainMenu() {
	mainMenu := testDarwinMainMenu()
	j := 0
	for i := 0; i < testNSMenuNumberOfItems(mainMenu); i++ {
		menu := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mainMenu, i))
		if j < len(initialMenus) && testNSMenuTitle(menu) == initialMenus[j] {
			j++
			continue
		}
		testNSMenuRemoveItemAtIndex(mainMenu, i)
		i--
	}
	j = 0
	appMenu := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mainMenu, 0))
	for i := 0; i < testNSMenuNumberOfItems(appMenu); i++ {
		item := testNSMenuItemAtIndex(appMenu, i)
		if testNSMenuItemTitle(item) == initialAppMenuItems[j] {
			j++
			continue
		}
		testNSMenuRemoveItemAtIndex(appMenu, i)
		i--
	}
}
