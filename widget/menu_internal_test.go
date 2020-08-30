package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
)

func TestMenu_ItemTapped(t *testing.T) {
	tapped := false
	item1 := fyne.NewMenuItem("Foo", nil)
	item2 := fyne.NewMenuItem("Bar", func() { tapped = true })
	item3 := fyne.NewMenuItem("Sub", nil)
	subItem := fyne.NewMenuItem("Foo", func() {})
	item3.ChildMenu = fyne.NewMenu("", subItem)
	m := NewMenu(fyne.NewMenu("", item1, item2, item3))
	size := m.MinSize()
	m.Resize(size)
	dismissed := false
	m.OnDismiss = func() { dismissed = true }

	mi1 := m.Items[0].(*menuItem)
	mi2 := m.Items[1].(*menuItem)
	mi3 := m.Items[2].(*menuItem)
	assert.Equal(t, item1, mi1.Item)
	assert.Equal(t, item2, mi2.Item)
	assert.Equal(t, item3, mi3.Item)

	// tap on item without action does not panic
	test.Tap(mi1)
	assert.False(t, tapped)
	assert.False(t, dismissed, "tap on item w/o action does not dismiss the menu")
	assert.True(t, m.Visible(), "tap on item w/o action does not hide the menu")

	test.Tap(mi2)
	assert.True(t, tapped)
	assert.True(t, dismissed, "tap on item dismisses the menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … OnDismiss is responsible for that")

	dismissed = false // reset
	mi3.MouseIn(nil)
	sm := mi3.child
	smi := sm.Items[0].(*menuItem)
	assert.Equal(t, subItem, smi.Item)
	assert.True(t, sm.Visible(), "sub menu is visible")

	test.Tap(smi)
	assert.True(t, dismissed, "tap on sub item dismisses the root menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … OnDismiss is responsible for that")
	assert.False(t, sm.Visible(), "tap on sub item hides the sub menu")

	newActionTapped := false
	item2.Action = func() { newActionTapped = true }
	test.Tap(mi2)
	assert.True(t, newActionTapped, "tap on item performs its current action")
}
