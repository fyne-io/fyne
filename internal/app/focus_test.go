package app_test

import (
	"testing"

	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestFocusManager_Focus(t *testing.T) {
	entry1 := widget.NewEntry()
	hidden := widget.NewCheck("test", func(bool) {})
	hidden.Hide()
	entry2 := widget.NewEntry()
	disabled := widget.NewCheck("test", func(bool) {})
	disabled.Disable()
	entry3 := widget.NewEntry()
	c := widget.NewVBox(entry1, hidden, entry2, disabled, entry3)

	manager := app.NewFocusManager(c)
	assert.Nil(t, manager.Focused())

	manager.Focus(entry2)
	assert.Equal(t, entry2, manager.Focused())

	manager.Focus(entry1)
	assert.Equal(t, entry1, manager.Focused())

	manager.Focus(entry3)
	assert.Equal(t, entry3, manager.Focused())

	manager.Focus(nil)
	assert.Nil(t, manager.Focused())
}

func TestFocusManager_FocusNext(t *testing.T) {
	entry1 := widget.NewEntry()
	hidden := widget.NewCheck("test", func(bool) {})
	hidden.Hide()
	entry2 := widget.NewEntry()
	disabled := widget.NewCheck("test", func(bool) {})
	disabled.Disable()
	entry3 := widget.NewEntry()
	c := widget.NewVBox(entry1, hidden, entry2, disabled, entry3)

	manager := app.NewFocusManager(c)
	assert.Nil(t, manager.Focused())

	manager.FocusNext()
	assert.Equal(t, entry1, manager.Focused())

	manager.FocusNext()
	assert.Equal(t, entry2, manager.Focused())

	manager.FocusNext()
	assert.Equal(t, entry3, manager.Focused())

	manager.FocusNext()
	assert.Equal(t, entry1, manager.Focused())
}

func TestFocusManager_FocusPrevious(t *testing.T) {
	entry1 := widget.NewEntry()
	hidden := widget.NewCheck("test", func(bool) {})
	hidden.Hide()
	entry2 := widget.NewEntry()
	disabled := widget.NewCheck("test", func(bool) {})
	disabled.Disable()
	entry3 := widget.NewEntry()
	c := widget.NewVBox(entry1, hidden, entry2, disabled, entry3)

	manager := app.NewFocusManager(c)
	assert.Nil(t, manager.Focused())

	manager.FocusPrevious()
	assert.Equal(t, entry3, manager.Focused())

	manager.FocusPrevious()
	assert.Equal(t, entry2, manager.Focused())

	manager.FocusPrevious()
	assert.Equal(t, entry1, manager.Focused())

	manager.FocusPrevious()
	assert.Equal(t, entry3, manager.Focused())
}
