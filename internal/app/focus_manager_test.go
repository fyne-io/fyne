package app_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/app"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFocusManager_Focus(t *testing.T) {
	t.Run("focusing and unfocusing", func(t *testing.T) {
		manager, entry1, _, _, entry2, _, entry3 := setupFocusManager(t)

		assert.True(t, manager.Focus(entry2))
		assert.Equal(t, entry2, manager.Focused())
		assert.True(t, entry2.focused)

		assert.True(t, manager.Focus(entry1))
		assert.Equal(t, entry1, manager.Focused())
		assert.True(t, entry1.focused)
		assert.False(t, entry2.focused)

		assert.True(t, manager.Focus(entry3))
		assert.Equal(t, entry3, manager.Focused())
		assert.True(t, entry3.focused)
		assert.False(t, entry1.focused)

		assert.True(t, manager.Focus(nil))
		assert.Nil(t, manager.Focused())
		assert.False(t, entry3.focused)
	})

	itBehavesLikeDoingNothing := func(t *testing.T, manager *app.FocusManager, notFocusableObj *focusable, focusableObj *focusable, objectIsPartOfContent bool) {
		if objectIsPartOfContent {
			assert.True(t, manager.Focus(notFocusableObj))
		} else {
			assert.False(t, manager.Focus(notFocusableObj))
		}
		assert.Nil(t, manager.Focused())
		assert.False(t, notFocusableObj.focused)

		require.True(t, manager.Focus(focusableObj))
		require.Equal(t, focusableObj, manager.Focused())
		require.True(t, focusableObj.focused)
		if objectIsPartOfContent {
			assert.True(t, manager.Focus(notFocusableObj))
		} else {
			assert.False(t, manager.Focus(notFocusableObj))
		}
		assert.False(t, notFocusableObj.focused)
		assert.Equal(t, focusableObj, manager.Focused())
		assert.True(t, focusableObj.focused)
	}

	t.Run("focus disabled", func(t *testing.T) {
		manager, entry1, _, _, _, disabled, _ := setupFocusManager(t)
		itBehavesLikeDoingNothing(t, manager, disabled, entry1, true)
	})

	t.Run("focus hidden", func(t *testing.T) {
		manager, entry1, hidden, _, _, _, _ := setupFocusManager(t)
		itBehavesLikeDoingNothing(t, manager, hidden, entry1, true)
	})

	t.Run("focus visible inside hidden", func(t *testing.T) {
		manager, entry1, _, visibleInsideHidden, _, _, _ := setupFocusManager(t)
		itBehavesLikeDoingNothing(t, manager, visibleInsideHidden, entry1, true)
	})

	t.Run("focus foreign", func(t *testing.T) {
		manager, entry1, _, _, _, _, _ := setupFocusManager(t)
		foreigner := &focusable{}
		itBehavesLikeDoingNothing(t, manager, foreigner, entry1, false)
	})
}

func TestFocusManager_FocusLost_FocusGained(t *testing.T) {
	manager, entry1, _, _, entry2, _, _ := setupFocusManager(t)

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
	manager, entry1, _, _, entry2, _, entry3 := setupFocusManager(t)

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
	manager, entry1, _, _, entry2, _, entry3 := setupFocusManager(t)

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

var _ fyne.Widget = (*focusable)(nil)
var _ fyne.Focusable = (*focusable)(nil)
var _ fyne.Disableable = (*focusable)(nil)

type focusable struct {
	widget.DisableableWidget
	children []fyne.CanvasObject
	focused  bool
}

func (f *focusable) CreateRenderer() fyne.WidgetRenderer {
	return &focusableRenderer{BaseRenderer: internalWidget.NewBaseRenderer(f.children)}
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

func (f *focusable) TypedRune(_ rune) {
}

func (f *focusable) TypedKey(_ *fyne.KeyEvent) {
}

var _ fyne.WidgetRenderer = (*focusableRenderer)(nil)

type focusableRenderer struct {
	internalWidget.BaseRenderer
}

func (f focusableRenderer) Layout(_ fyne.Size) {
}

func (f focusableRenderer) MinSize() fyne.Size {
	return fyne.NewSize(0, 0)
}

func (f focusableRenderer) Refresh() {
}

func setupFocusManager(t *testing.T) (m *app.FocusManager, entry1, hidden, visibleInsideHidden, entry2, disabled, entry3 *focusable) {
	entry1 = &focusable{}
	visibleInsideHidden = &focusable{}
	hidden = &focusable{
		children: []fyne.CanvasObject{visibleInsideHidden},
	}
	hidden.Hide()
	entry2 = &focusable{}
	disabled = &focusable{}
	disabled.Disable()
	entry3 = &focusable{}
	m = app.NewFocusManager(container.NewVBox(entry1, hidden, entry2, disabled, entry3))
	require.Nil(t, m.Focused())
	require.False(t, hidden.Visible())
	require.True(t, visibleInsideHidden.Visible())
	return
}
