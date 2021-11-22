//go:build !ci && !no_native_menus && !mobile
// +build !ci,!no_native_menus,!mobile

package glfw

import (
	"fmt"
	"testing"
	"unsafe"

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
	itemOpen := fyne.NewMenuItem("Open", func() { lastAction = "open" })
	itemRecent := fyne.NewMenuItem("Recent", nil)
	itemFoo := fyne.NewMenuItem("Foo", func() { lastAction = "foo" })
	itemRecent.ChildMenu = fyne.NewMenu("", itemFoo)
	menuEdit := fyne.NewMenu("File", itemNew, itemOpen, fyne.NewMenuItemSeparator(), itemRecent)

	itemHelp := fyne.NewMenuItem("Help", func() { lastAction = "Help!!!" })
	itemHelpMe := fyne.NewMenuItem("Help Me", func() { lastAction = "Help me!!!" })
	menuHelp := fyne.NewMenu("Help", itemHelp, itemHelpMe)

	itemHelloWorld := fyne.NewMenuItem("Hello World", func() { lastAction = "Hello World!" })
	itemPrefs := fyne.NewMenuItem("Preferences", func() { lastAction = "prefs" })
	itemMore := fyne.NewMenuItem("More", func() { lastAction = "more" })
	itemMorePrefs := fyne.NewMenuItem("Preferences…", func() { lastAction = "more prefs" })
	menuMore := fyne.NewMenu("More Stuff", itemHelloWorld, itemPrefs, itemMore, itemMorePrefs)

	itemSettings := fyne.NewMenuItem("Settings", func() { lastAction = "settings" })
	itemMoreSetings := fyne.NewMenuItem("Settings…", func() { lastAction = "more settings" })
	menuSettings := fyne.NewMenu("Settings", itemSettings, fyne.NewMenuItemSeparator(), itemMoreSetings)

	mainMenu := fyne.NewMainMenu(menuEdit, menuHelp, menuMore, menuSettings)
	setupNativeMenu(w, mainMenu)
	fmt.Println(lastAction)

	mm := testDarwinMainMenu()
	// The custom “Preferences” menu should be moved to the system app menu completely.
	// -> only three custom menus
	assert.Equal(t, 5, testNSMenuNumberOfItems(mm), "two built-in + three custom")

	m := testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 0))
	assert.Equal(t, "", testNSMenuTitle(m), "app menu doesn’t have a title")
	assertNSMenuItemSeparator(m, 1)
	assertNSMenuItem(t, "Preferences", m, 2)
	assertLastAction("prefs")
	assertNSMenuItem(t, "Preferences…", m, 3)
	assertLastAction("more prefs")
	assertNSMenuItemSeparator(m, 4)
	assertNSMenuItem(t, "Settings", m, 5)
	assertLastAction("settings")
	assertNSMenuItem(t, "Settings…", m, 6)
	assertLastAction("more settings")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
	assert.Equal(t, "File", testNSMenuTitle(m))
	assert.Equal(t, 4, testNSMenuNumberOfItems(m))
	assertNSMenuItem(t, "New", m, 0)
	assertLastAction("new")
	assertNSMenuItem(t, "Open", m, 1)
	assertLastAction("open")
	assertNSMenuItemSeparator(m, 2)
	i := testNSMenuItemAtIndex(m, 3)
	assert.Equal(t, "Recent", testNSMenuItemTitle(i))
	sm := testNSMenuItemSubmenu(i)
	assert.NotNil(t, sm, "item has submenu")
	assert.Equal(t, 1, testNSMenuNumberOfItems(sm))
	assertNSMenuItem(t, "Foo", sm, 0)
	assertLastAction("foo")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 2))
	assert.Equal(t, "More Stuff", testNSMenuTitle(m))
	assert.Equal(t, 2, testNSMenuNumberOfItems(m))
	assertNSMenuItem(t, "Hello World", m, 0)
	assertLastAction("Hello World!")
	assertNSMenuItem(t, "More", m, 1)
	assertLastAction("more")

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 3))
	assert.Equal(t, "Window", testNSMenuTitle(m))

	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 4))
	assert.Equal(t, "Help", testNSMenuTitle(m))
	assert.Equal(t, 2, testNSMenuNumberOfItems(m))
	assertNSMenuItem(t, "Help", m, 0)
	assertLastAction("Help!!!")
	assertNSMenuItem(t, "Help Me", m, 1)
	assertLastAction("Help me!!!")

	// change action works
	itemOpen.Action = func() { lastAction = "new open" }
	m = testNSMenuItemSubmenu(testNSMenuItemAtIndex(mm, 1))
	assertNSMenuItem(t, "Open", m, 1)
	assertLastAction("new open")
}

var initialAppMenuItems []string
var initialMenus []string

func assertNSMenuItem(t *testing.T, wantTitle string, m unsafe.Pointer, i int) {
	item := testNSMenuItemAtIndex(m, i)
	assert.Equal(t, wantTitle, testNSMenuItemTitle(item))
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
