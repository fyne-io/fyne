// +build mobile

package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"

	"github.com/stretchr/testify/assert"
)

func TestMenu_ItemWithChildTapped(t *testing.T) {
	sub1 := fyne.NewMenuItem("sub1", nil)
	sub1.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub1 A", nil),
		fyne.NewMenuItem("sub1 B", nil),
	)
	sub2sub := fyne.NewMenuItem("sub2sub", nil)
	sub2sub.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub2sub A", nil),
		fyne.NewMenuItem("sub2sub B", nil),
	)
	sub2 := fyne.NewMenuItem("sub2", nil)
	sub2.ChildMenu = fyne.NewMenu("",
		fyne.NewMenuItem("sub2 A", nil),
		fyne.NewMenuItem("sub2 B", nil),
		sub2sub,
	)
	m := NewMenu(
		fyne.NewMenu("",
			fyne.NewMenuItem("Foo", nil),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Bar", nil),
			sub1,
			sub2,
		),
	)
	size := m.MinSize()
	m.Resize(size)

	sub1Widget := m.Items[3].(*menuItem)
	assert.Equal(t, sub1, sub1Widget.Item)
	sub2Widget := m.Items[4].(*menuItem)
	assert.Equal(t, sub2, sub2Widget.Item)
	assert.False(t, sub1Widget.child.Visible(), "submenu is invisible initially")
	assert.False(t, sub2Widget.child.Visible(), "submenu is invisible initially")
	test.Tap(sub1Widget)
	assert.True(t, sub1Widget.child.Visible(), "tapping item shows submenu")
	assert.False(t, sub2Widget.child.Visible(), "other child menu stays hidden")
	test.Tap(sub2Widget)
	assert.False(t, sub1Widget.child.Visible(), "tapping other item hides current submenu")
	assert.True(t, sub2Widget.child.Visible(), "other child menu shows up")

	sub2subWidget := sub2Widget.child.Items[2].(*menuItem)
	assert.Equal(t, sub2sub, sub2subWidget.Item)
	assert.False(t, sub2subWidget.child.Visible(), "2nd level submenu is invisible initially")
	test.Tap(sub2subWidget)
	assert.True(t, sub2Widget.child.Visible(), "1st level submenu stays visible")
	assert.True(t, sub2subWidget.child.Visible(), "2nd level submenu shows up")

	test.Tap(sub1Widget)
	assert.False(t, sub2Widget.child.Visible(), "1st level submenu is hidden by other submenu")
	test.Tap(sub2Widget)
	assert.False(t, sub2subWidget.child.Visible(), "2nd level submenu is hidden when re-entering its parent")
}
