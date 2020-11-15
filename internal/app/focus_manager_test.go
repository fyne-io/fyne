package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
)

func TestFocusManager_Focus(t *testing.T) {
	entry1 := &focusable{}
	hidden := widget.NewCheck("test", func(bool) {})
	hidden.Hide()
	entry2 := &focusable{}
	disabled := widget.NewCheck("test", func(bool) {})
	disabled.Disable()
	entry3 := &focusable{}
	c := widget.NewVBox(entry1, hidden, entry2, disabled, entry3)

	manager := app.NewFocusManager(c)
	assert.Nil(t, manager.Focused())

	manager.Focus(entry2)
	assert.Equal(t, entry2, manager.Focused())
	assert.True(t, entry2.focused)

	manager.Focus(entry1)
	assert.Equal(t, entry1, manager.Focused())
	assert.True(t, entry1.focused)
	assert.False(t, entry2.focused)

	manager.Focus(entry3)
	assert.Equal(t, entry3, manager.Focused())
	assert.True(t, entry3.focused)
	assert.False(t, entry1.focused)

	manager.Focus(nil)
	assert.Nil(t, manager.Focused())
	assert.False(t, entry3.focused)
}

func TestFocusManager_FocusLost_FocusGained(t *testing.T) {
	entry1 := &focusable{}
	entry2 := &focusable{}
	entry3 := &focusable{}
	c := widget.NewVBox(entry1, entry2, entry3)

	manager := app.NewFocusManager(c)
	manager.Focus(entry2)
	require.Equal(t, entry2, manager.Focused())
	require.False(t, entry1.focused)
	require.True(t, entry2.focused)

	manager.FocusLost()
	assert.Equal(t, entry2, manager.Focused(), "losing focus does not mean that manager loses track of focused element")
	assert.False(t, entry1.focused)
	assert.False(t, entry2.focused, "focused entry loses focus if manager loses it")

	manager.FocusGained()
	assert.Equal(t, entry2, manager.Focused())
	assert.False(t, entry1.focused)
	assert.True(t, entry2.focused, "focused entry gains focus if manager gains it")
}

func TestFocusManager_FocusNext(t *testing.T) {
	entry1 := &focusable{}
	hidden := widget.NewCheck("test", func(bool) {})
	hidden.Hide()
	entry2 := &focusable{}
	disabled := widget.NewCheck("test", func(bool) {})
	disabled.Disable()
	entry3 := &focusable{}
	c := widget.NewVBox(entry1, hidden, entry2, disabled, entry3)

	manager := app.NewFocusManager(c)
	assert.Nil(t, manager.Focused())

	manager.FocusNext()
	assert.Equal(t, entry1, manager.Focused())
	assert.True(t, entry1.focused)

	manager.FocusNext()
	assert.Equal(t, entry2, manager.Focused())
	assert.True(t, entry2.focused)
	assert.False(t, entry1.focused)

	manager.FocusNext()
	assert.Equal(t, entry3, manager.Focused())
	assert.True(t, entry3.focused)
	assert.False(t, entry2.focused)

	manager.FocusNext()
	assert.Equal(t, entry1, manager.Focused())
	assert.True(t, entry1.focused)
	assert.False(t, entry3.focused)
}

func TestFocusManager_FocusPrevious(t *testing.T) {
	entry1 := &focusable{}
	hidden := widget.NewCheck("test", func(bool) {})
	hidden.Hide()
	entry2 := &focusable{}
	disabled := widget.NewCheck("test", func(bool) {})
	disabled.Disable()
	entry3 := &focusable{}
	c := widget.NewVBox(entry1, hidden, entry2, disabled, entry3)

	manager := app.NewFocusManager(c)
	assert.Nil(t, manager.Focused())

	manager.FocusPrevious()
	assert.Equal(t, entry3, manager.Focused())
	assert.True(t, entry3.focused)

	manager.FocusPrevious()
	assert.Equal(t, entry2, manager.Focused())
	assert.True(t, entry2.focused)
	assert.False(t, entry3.focused)

	manager.FocusPrevious()
	assert.Equal(t, entry1, manager.Focused())
	assert.True(t, entry1.focused)
	assert.False(t, entry2.focused)

	manager.FocusPrevious()
	assert.Equal(t, entry3, manager.Focused())
	assert.True(t, entry3.focused)
	assert.False(t, entry1.focused)
}

type focusable struct {
	widget.Entry
	focused bool
}

func (f *focusable) FocusGained() {
	if f.Disabled() {
		return
	}
	f.focused = true
}

func (f *focusable) FocusLost() {
	f.focused = false
}
