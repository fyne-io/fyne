package app_test

import (
	"testing"

	"fyne.io/fyne/internal/app"
	"fyne.io/fyne/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFocusManager_Focus(t *testing.T) {
	t.Run("focusing and unfocusing", func(t *testing.T) {
		manager, entry1, _, entry2, _, entry3 := setupFocusManager(t)

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
	})

	itBehavesLikeUnfocus := func(manager *app.FocusManager, notFocusableObj *focusable, focusableObj *focusable) {
		manager.Focus(notFocusableObj)
		assert.Nil(t, manager.Focused())
		assert.False(t, notFocusableObj.focused)

		manager.Focus(focusableObj)
		require.True(t, focusableObj.focused)
		manager.Focus(notFocusableObj)
		assert.Nil(t, manager.Focused())
		assert.False(t, notFocusableObj.focused)
		assert.False(t, focusableObj.focused)
	}

	t.Run("focus disabled", func(t *testing.T) {
		manager, entry1, _, _, disabled, _ := setupFocusManager(t)
		itBehavesLikeUnfocus(manager, disabled, entry1)
	})

	t.Run("focus hidden", func(t *testing.T) {
		manager, entry1, hidden, _, _, _ := setupFocusManager(t)
		itBehavesLikeUnfocus(manager, hidden, entry1)
	})

	t.Run("focus foreign", func(t *testing.T) {
		manager, entry1, _, _, _, _ := setupFocusManager(t)
		foreigner := &focusable{}

		manager.Focus(foreigner)
		assert.Nil(t, manager.Focused())
		assert.False(t, foreigner.focused)

		manager.Focus(entry1)
		require.Equal(t, entry1, manager.Focused())
		require.True(t, entry1.focused)
		manager.Focus(foreigner)
		assert.False(t, foreigner.focused)
		assert.Equal(t, entry1, manager.Focused())
		assert.True(t, entry1.focused)
	})
}

func TestFocusManager_FocusLost_FocusGained(t *testing.T) {
	manager, entry1, _, entry2, _, _ := setupFocusManager(t)

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
	manager, entry1, _, entry2, _, entry3 := setupFocusManager(t)

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
	manager, entry1, _, entry2, _, entry3 := setupFocusManager(t)

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

func setupFocusManager(t *testing.T) (m *app.FocusManager, entry1, hidden, entry2, disabled, entry3 *focusable) {
	entry1 = &focusable{}
	hidden = &focusable{}
	hidden.Hide()
	entry2 = &focusable{}
	disabled = &focusable{}
	disabled.Disable()
	entry3 = &focusable{}
	m = app.NewFocusManager(widget.NewVBox(entry1, hidden, entry2, disabled, entry3))
	assert.Nil(t, m.Focused())
	return
}
