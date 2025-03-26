//go:build !no_glfw && !mobile

package glfw

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"

	"github.com/stretchr/testify/assert"
)

func Test_Menu_Empty(t *testing.T) {
	w := createWindow("Menu Test")
	bar := buildMenuOverlay(fyne.NewMainMenu(), w.window)
	assert.Nil(t, bar) // no bar but does not crash
}

func Test_Menu_AddsQuit(t *testing.T) {
	w := createWindow("Menu Test")
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File"))
	bar := buildMenuOverlay(mainMenu, w.window)
	assert.NotNil(t, bar)
	assert.Len(t, mainMenu.Items, 1)
	assert.Len(t, mainMenu.Items[0].Items, 2) // separator+quit inserted
	assert.True(t, mainMenu.Items[0].Items[1].IsQuit)
}

func Test_Menu_LeaveQuit(t *testing.T) {
	w := createWindow("Menu Test")
	quitFunc := func() {}
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File", fyne.NewMenuItem(lang.L("Quit"), quitFunc)))
	bar := buildMenuOverlay(mainMenu, w.window)
	assert.NotNil(t, bar)
	assert.Len(t, mainMenu.Items[0].Items, 1) // no separator added
	assert.Equal(t, reflect.ValueOf(quitFunc).Pointer(), reflect.ValueOf(mainMenu.Items[0].Items[0].Action).Pointer())
}
func Test_Menu_LeaveQuit_AddAction(t *testing.T) {
	w := createWindow("Menu Test")
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File", fyne.NewMenuItem(lang.L("Quit"), nil)))
	bar := buildMenuOverlay(mainMenu, w.window)
	assert.NotNil(t, bar)
	assert.Len(t, mainMenu.Items[0].Items, 1)           // no separator added
	assert.NotNil(t, mainMenu.Items[0].Items[0].Action) // quit action was added
}

func Test_Menu_CustomQuit(t *testing.T) {
	w := createWindow("Menu Test")

	quitFunc := func() {}
	quitItem := fyne.NewMenuItem("Beenden", quitFunc)
	quitItem.IsQuit = true

	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File", quitItem))
	bar := buildMenuOverlay(mainMenu, w.window)

	assert.NotNil(t, bar)
	assert.Len(t, mainMenu.Items[0].Items, 1) // no separator added
	assert.Equal(t, reflect.ValueOf(quitFunc).Pointer(), reflect.ValueOf(mainMenu.Items[0].Items[0].Action).Pointer())
}

func Test_Menu_CustomQuit_NoAction(t *testing.T) {
	w := createWindow("Menu Test")
	quitItem := fyne.NewMenuItem("Beenden", nil)
	quitItem.IsQuit = true
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File", quitItem))
	bar := buildMenuOverlay(mainMenu, w.window)

	assert.NotNil(t, bar)
	assert.Len(t, mainMenu.Items[0].Items, 1)           // no separator added
	assert.NotNil(t, mainMenu.Items[0].Items[0].Action) // quit action was added
}
