package widget_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"
)

func TestMenu_ItemHovered(t *testing.T) {
	m := widget.NewMenu(
		fyne.NewMenu("",
			fyne.NewMenuItem("Foo", nil),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Bar", nil),
		),
	)
	size := m.MinSize()
	m.Resize(size)

	mi := m.Items[0].(*widget.MenuItem)
	r1 := cache.Renderer(mi)
	assert.Equal(t, color.Transparent, r1.BackgroundColor())
	mi.MouseIn(nil)
	assert.Equal(t, theme.HoverColor(), r1.BackgroundColor())
	mi.MouseOut()
	assert.Equal(t, color.Transparent, r1.BackgroundColor())
}

func TestMenu_ItemTapped(t *testing.T) {
	tapped := false
	item1 := fyne.NewMenuItem("Foo", nil)
	item2 := fyne.NewMenuItem("Bar", func() { tapped = true })
	m := widget.NewMenu(fyne.NewMenu("", item1, item2))
	size := m.MinSize()
	m.Resize(size)
	dismissed := false
	m.DismissAction = func() { dismissed = true }

	mi1 := m.Items[0].(*widget.MenuItem)
	mi2 := m.Items[1].(*widget.MenuItem)
	assert.Equal(t, item1, mi1.Item)
	assert.Equal(t, item2, mi2.Item)

	// tap on item without action does not panic
	test.Tap(mi1)
	assert.False(t, tapped)
	assert.False(t, dismissed, "tap on item w/o action does not dismiss the menu")
	assert.True(t, m.Visible(), "tap on item w/o action does not hide the menu")

	test.Tap(mi2)
	assert.True(t, tapped)
	assert.True(t, dismissed, "tap on item dismisses the menu")
	assert.True(t, m.Visible(), "tap on item does not hide the menu … the DismissAction is reponsible for that")
}

func TestMenu_Layout(t *testing.T) {
	item1 := fyne.NewMenuItem("A", nil)
	item2 := fyne.NewMenuItem("B (long)", nil)
	sep := fyne.NewMenuItemSeparator()
	m := widget.NewMenu(fyne.NewMenu("", item1, sep, item2))
	itemSize := canvas.NewText("B (long)", color.Black).MinSize().Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
	size := fyne.NewSize(itemSize.Width, 2*itemSize.Height+2+4*theme.Padding()) // 2 for the separator; padding between items (2) & at start/end of menu (2)
	objects := test.LaidOutObjects(m)
	require.GreaterOrEqual(t, len(objects), 1)

	assert.Equal(t, objects[0], m)
	assert.Equal(t, size, m.MinSize())

	containerFound := false
	for _, object := range objects {
		if c, ok := object.(*fyne.Container); ok {
			containerFound = true
			assert.Equal(t, layout.NewVBoxLayout(), c.Layout)
			assert.Len(t, c.Objects, 3)

			assert.Equal(t, item1, c.Objects[0].(*widget.MenuItem).Item)
			assert.IsType(t, (*canvas.Rectangle)(nil), c.Objects[1])
			assert.Equal(t, item2, c.Objects[2].(*widget.MenuItem).Item)

			assert.Equal(t, size.Subtract(fyne.NewSize(0, theme.Padding()*2)), c.Size(), "container size does not include leading & trailing padding")
			assert.Equal(t, itemSize, c.Objects[0].Size())
			assert.Equal(t, fyne.NewSize(size.Width, 2), c.Objects[1].Size())
			assert.Equal(t, itemSize, c.Objects[2].Size())

			assert.Equal(t, fyne.NewPos(0, theme.Padding()), c.Position())
			y1 := 0
			assert.Equal(t, fyne.NewPos(0, y1), c.Objects[0].Position())
			y2 := y1 + itemSize.Height + theme.Padding()
			assert.Equal(t, fyne.NewPos(0, y2), c.Objects[1].Position())
			y3 := y2 + 2 + theme.Padding()
			assert.Equal(t, fyne.NewPos(0, y3), c.Objects[2].Position())
		}
	}
	assert.True(t, containerFound, "expect container")
}
