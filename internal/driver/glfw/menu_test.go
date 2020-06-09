// +build !ci
// +build !mobile

package glfw

import (
	"testing"

	"fyne.io/fyne"

	"github.com/stretchr/testify/assert"
)

func Test_Menu_Empty(t *testing.T) {
	w := createWindow("Menu Test").(*window)
	bar := buildMenuOverlay(fyne.NewMainMenu(), w.canvas)
	assert.Nil(t, bar) // no bar but does not crash
}

func Test_Menu_AddsQuit(t *testing.T) {
	w := createWindow("Menu Test").(*window)
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File"))
	bar := buildMenuOverlay(mainMenu, w.canvas)
	assert.NotNil(t, bar)
	assert.Equal(t, 1, len(mainMenu.Items))
	assert.Equal(t, 2, len(mainMenu.Items[0].Items)) // separator+quit inserted
}

func Test_Menu_LeaveQuit(t *testing.T) {
	w := createWindow("Menu Test").(*window)
	mainMenu := fyne.NewMainMenu(fyne.NewMenu("File", fyne.NewMenuItem("Quit", nil)))
	bar := buildMenuOverlay(mainMenu, w.canvas)
	assert.NotNil(t, bar)
	assert.Equal(t, 1, len(mainMenu.Items[0].Items)) // no separator added
	assert.Nil(t, mainMenu.Items[0].Items[0].Action) // no quit handler added
}
