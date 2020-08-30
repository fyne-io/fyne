package widget

import (
	"testing"

	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestGroupSize(t *testing.T) {
	group := NewGroup("Test", NewLabel("Hello"), NewLabel("World"))
	assert.Equal(t, 2, len(group.box.Children))
}

func TestGroupPrepend(t *testing.T) {
	group := NewGroup("Test", NewLabel("World"))
	assert.Equal(t, 1, len(group.box.Children))

	prepend := NewLabel("Hello")
	group.Prepend(prepend)
	assert.Equal(t, 2, len(group.box.Children))
	assert.Equal(t, prepend, group.box.Children[0])
}

func TestGroupAppend(t *testing.T) {
	group := NewGroup("Test", NewLabel("Hello"))
	assert.Equal(t, 1, len(group.box.Children))

	append := NewLabel("World")
	group.Append(append)
	assert.True(t, len(group.box.Children) == 2)
	assert.Equal(t, append, group.box.Children[1])
}

func TestGroup_Text(t *testing.T) {
	group := NewGroup("Test")
	group.Refresh()

	group.Text = "World"
	group.Refresh()

	assert.Equal(t, "World", test.WidgetRenderer(group).(*groupRenderer).label.Text)
}

func TestGroup_Scroll(t *testing.T) {
	group := NewGroupWithScroller("Test",
		NewLabel("A label"), NewLabel("Another item"))
	group.Resize(group.MinSize())

	scroll := group.content.(*ScrollContainer)
	assert.Equal(t, scroll.Direction, ScrollVerticalOnly)
	assert.Equal(t, group.MinSize().Width, scroll.Content.Size().Width)
	totalHeight := test.WidgetRenderer(group).(*groupRenderer).label.MinSize().Height + scroll.Content.Size().Height
	assert.Less(t, group.MinSize().Height, totalHeight)
}
