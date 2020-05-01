package app_test

import (
	"testing"

	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestFocusManager_FocusNext(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	hidden := widget.NewCheck("test", func(bool) {})
	hidden.Hide()
	entry2 := widget.NewEntry()
	disabled := widget.NewCheck("test", func(bool) {})
	disabled.Disable()
	entry3 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, hidden, entry2, disabled, entry3))

	manager := app.NewFocusManager(c)
	assert.Nil(t, c.Focused())

	manager.FocusNext(nil)
	assert.Equal(t, entry1, c.Focused())

	manager.FocusNext(entry1)
	assert.Equal(t, entry2, c.Focused())

	manager.FocusNext(entry2)
	assert.Equal(t, entry3, c.Focused())

	manager.FocusNext(entry3)
	assert.Equal(t, entry1, c.Focused())
}

func TestFocusManager_FocusPrevious(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	hidden := widget.NewCheck("test", func(bool) {})
	hidden.Hide()
	entry2 := widget.NewEntry()
	disabled := widget.NewCheck("test", func(bool) {})
	disabled.Disable()
	entry3 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, hidden, entry2, disabled, entry3))

	manager := app.NewFocusManager(c)
	assert.Nil(t, c.Focused())

	manager.FocusPrevious(nil)
	assert.Equal(t, entry3, c.Focused())

	manager.FocusPrevious(entry3)
	assert.Equal(t, entry2, c.Focused())

	manager.FocusPrevious(entry2)
	assert.Equal(t, entry1, c.Focused())

	manager.FocusPrevious(entry1)
	assert.Equal(t, entry3, c.Focused())
}
