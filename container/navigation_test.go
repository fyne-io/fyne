package container

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestNavigation_RootWithTitle(t *testing.T) {
	l := widget.NewLabel("something")
	nav := NewNavigationWithTitle(l, "Title")

	assert.Equal(t, true, nav.Back.Disabled())
	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, len(nav.titles))

	assert.Nil(t, nav.Pop())

	assert.Equal(t, "Title", nav.Label.Text)
}

func TestNavigation_EmptyPushWithTitle(t *testing.T) {
	a := widget.NewLabel("a")
	b := widget.NewLabel("b")
	c := widget.NewLabel("c")
	nav := NewNavigationWithTitle(nil, "Title")

	assert.Equal(t, true, nav.Back.Disabled())
	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))
	assert.Equal(t, 0, nav.level)

	nav.PushWithTitle(a, "A")
	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, nav.level)
	assert.Equal(t, "A", nav.Label.Text)
	assert.Equal(t, false, nav.Back.Disabled())

	nav.PushWithTitle(b, "B")
	assert.Equal(t, 2, len(nav.stack.Objects))
	assert.Equal(t, 2, nav.level)
	assert.Equal(t, "B", nav.Label.Text)
	assert.Equal(t, false, nav.Back.Disabled())

	nav.PushWithTitle(c, "C")
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, nav.level)
	assert.Equal(t, "C", nav.Label.Text)
	assert.Equal(t, false, nav.Back.Disabled())

	assert.Equal(t, 3, len(nav.titles))

	assert.Equal(t, c, nav.Pop())
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 2, nav.level)
	assert.Equal(t, "B", nav.Label.Text)
	assert.Equal(t, false, nav.Back.Disabled())

	assert.Equal(t, b, nav.Pop())
	assert.Equal(t, "A", nav.Label.Text)
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, len(nav.titles))
	assert.Equal(t, 1, nav.level)
	assert.Equal(t, false, nav.Back.Disabled())

	assert.Equal(t, a, nav.Pop())
	assert.Equal(t, "Title", nav.Label.Text)
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, len(nav.titles))
	assert.Equal(t, 0, nav.level)
	assert.Equal(t, true, nav.Back.Disabled())

	assert.Nil(t, nav.Pop())
	assert.Equal(t, "Title", nav.Label.Text)
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, len(nav.titles))
	assert.Equal(t, 0, nav.level)
	assert.Equal(t, true, nav.Back.Disabled())
}

func TestNavigation_Empty(t *testing.T) {
	nav := NewNavigationWithTitle(nil, "Title")

	assert.Equal(t, true, nav.Back.Disabled())
	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))

	assert.Nil(t, nav.Pop())

	assert.Equal(t, "Title", nav.Label.Text)
}

func TestNavigation_WithoutConstructor(t *testing.T) {
	nav := &Navigation{Title: "Nav Test"}

	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))

	assert.Nil(t, nav.Pop())

	assert.Nil(t, nav.Back)
	assert.Nil(t, nav.Next)
	assert.Nil(t, nav.Label)
	assert.Equal(t, "Nav Test", nav.Title)
}

func TestNavigation_StructWithRootAndTitle(t *testing.T) {
	nav := &Navigation{
		Title: "Nav Test",
		Root:  widget.NewLabel("Something"),
	}

	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))

	_ = test.TempWidgetRenderer(t, nav)

	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, len(nav.titles))

	assert.Nil(t, nav.Pop())

	assert.Nil(t, nav.Back)
	assert.Nil(t, nav.Next)
	assert.Nil(t, nav.Label)
	assert.Equal(t, "Nav Test", nav.Title)
}
