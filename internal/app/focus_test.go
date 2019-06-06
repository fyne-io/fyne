package app

import (
	"testing"

	"fyne.io/fyne/test"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestFocusManager_nextInChain(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	entry2 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, entry2))

	manager := NewFocusManager(c)
	next := manager.nextInChain(entry1)
	assert.Equal(t, entry2, next)
}

func TestFocusManager_nextInChain_Nil(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	entry2 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, entry2))

	manager := NewFocusManager(c)
	next := manager.nextInChain(nil)
	assert.Equal(t, entry1, next)
}

func TestFocusManager_nextInChain_Disableable(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	check := widget.NewCheck("test", func(bool) {})
	check.Disable()
	entry2 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, check, entry2))

	manager := NewFocusManager(c)
	next := manager.nextInChain(entry1)
	assert.Equal(t, entry2, next)
}

func TestFocusManager_FocusNext(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	entry2 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, entry2))

	manager := NewFocusManager(c)
	c.Focus(entry1)
	assert.Equal(t, entry1, c.Focused())

	manager.FocusNext(entry1)
	assert.Equal(t, entry2, c.Focused())
}

func TestFocusManager_previousInChain(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	entry2 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, entry2))

	manager := NewFocusManager(c)
	previous := manager.previousInChain(entry2)
	assert.Equal(t, entry1, previous)
}

func TestFocusManager_previousInChain_Nil(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	entry2 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, entry2))

	manager := NewFocusManager(c)
	previous := manager.previousInChain(nil)
	assert.Equal(t, entry2, previous)
}

func TestFocusManager_previousInChain_Disableable(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	check := widget.NewCheck("test", func(bool) {})
	check.Disable()
	entry2 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, check, entry2))

	manager := NewFocusManager(c)
	next := manager.previousInChain(entry2)
	assert.Equal(t, entry1, next)
}

func TestFocusManager_FocusPrevious(t *testing.T) {
	c := test.NewCanvas()
	entry1 := widget.NewEntry()
	entry2 := widget.NewEntry()
	c.SetContent(widget.NewVBox(entry1, entry2))

	manager := NewFocusManager(c)
	c.Focus(entry2)
	assert.Equal(t, entry2, c.Focused())

	manager.FocusPrevious(entry2)
	assert.Equal(t, entry1, c.Focused())
}
