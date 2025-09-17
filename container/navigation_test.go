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
	tr := test.TempWidgetRenderer(t, nav)
	r := tr.(*navigatorRenderer)

	assert.Equal(t, true, r.back.Disabled())
	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, len(nav.titles))

	assert.Nil(t, nav.Back())
	assert.Nil(t, nav.Forward())

	assert.Equal(t, "Title", r.title.Text)

}

func TestNavigation_EmptyPushWithTitle(t *testing.T) {
	a := widget.NewLabel("a")
	b := widget.NewLabel("b")
	c := widget.NewLabel("c")
	nav := NewNavigationWithTitle(nil, "Title")
	tr := test.TempWidgetRenderer(t, nav)
	r := tr.(*navigatorRenderer)

	assert.Equal(t, true, r.back.Disabled())
	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))
	assert.Equal(t, 0, nav.level)

	nav.PushWithTitle(a, "A")
	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, nav.level)
	assert.Equal(t, "A", r.title.Text)
	assert.Equal(t, false, r.back.Disabled())

	nav.PushWithTitle(b, "B")
	assert.Equal(t, 2, len(nav.stack.Objects))
	assert.Equal(t, 2, nav.level)
	assert.Equal(t, "B", r.title.Text)
	assert.Equal(t, false, r.back.Disabled())

	assert.Equal(t, b, nav.Back())
	assert.Equal(t, 1, nav.level)
	assert.Equal(t, b, nav.Forward())
	assert.Equal(t, 2, nav.level)

	nav.PushWithTitle(c, "C")
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, nav.level)
	assert.Equal(t, "C", r.title.Text)
	assert.Equal(t, false, r.back.Disabled())

	assert.Equal(t, 3, len(nav.titles))

	assert.Equal(t, c, nav.Back())
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 2, nav.level)
	assert.Equal(t, "B", r.title.Text)
	assert.Equal(t, false, r.back.Disabled())

	assert.Equal(t, b, nav.Back())
	assert.Equal(t, "A", r.title.Text)
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, len(nav.titles))
	assert.Equal(t, 1, nav.level)
	assert.Equal(t, false, r.back.Disabled())

	assert.Equal(t, a, nav.Back())
	assert.Equal(t, "Title", r.title.Text)
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, len(nav.titles))
	assert.Equal(t, 0, nav.level)
	assert.Equal(t, true, r.back.Disabled())

	assert.Nil(t, nav.Back())
	assert.Equal(t, "Title", r.title.Text)
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, len(nav.titles))
	assert.Equal(t, 0, nav.level)
	assert.Equal(t, true, r.back.Disabled())
}

func TestNavigation_Empty(t *testing.T) {
	nav := NewNavigationWithTitle(nil, "Title")
	tr := test.TempWidgetRenderer(t, nav)
	r := tr.(*navigatorRenderer)

	assert.Equal(t, true, r.back.Disabled())
	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))

	assert.Nil(t, nav.Back())

	assert.Equal(t, "Title", r.title.Text)
}

func TestNavigation_WithoutConstructor(t *testing.T) {
	nav := &Navigation{Title: "Nav Test"}
	_ = test.TempWidgetRenderer(t, nav)

	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))

	assert.Nil(t, nav.Back())

	assert.Equal(t, "Nav Test", nav.Title)
}

func TestNavigation_StructWithRootAndTitle(t *testing.T) {
	nav := &Navigation{
		Title: "Nav Test",
		Root:  widget.NewLabel("Something"),
	}

	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))

	tr := test.TempWidgetRenderer(t, nav)
	r := tr.(*navigatorRenderer)

	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, len(nav.titles))

	assert.Nil(t, nav.Back())

	assert.Equal(t, "Nav Test", nav.Title)
	assert.Equal(t, "Nav Test", r.title.Text)
}
