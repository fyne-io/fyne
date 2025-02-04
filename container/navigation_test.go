package container

import (
	"testing"

	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestNavigation_Basics(t *testing.T) {
	l := widget.NewLabel("something")
	nav := NewNavigation("Title", l)

	assert.Equal(t, true, nav.button.Disabled())
	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, len(nav.titles))

	assert.Equal(t, l, nav.Pop())

	assert.Equal(t, "Title", nav.label.Text)
}

func TestNavigation_NewWithMultiple(t *testing.T) {
	a := widget.NewLabel("a")
	b := widget.NewLabel("b")
	c := widget.NewLabel("c")
	nav := NewNavigation("Title", a, b, c)

	assert.Equal(t, false, nav.button.Disabled())
	assert.Equal(t, 3, len(nav.stack.Objects))
	assert.Equal(t, 3, len(nav.titles))

	assert.Equal(t, c, nav.Pop())
	assert.Equal(t, "Title", nav.label.Text)
	assert.Equal(t, 2, len(nav.stack.Objects))
	assert.Equal(t, 2, len(nav.titles))
	assert.Equal(t, false, nav.button.Disabled())

	assert.Equal(t, b, nav.Pop())
	assert.Equal(t, "Title", nav.label.Text)
	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, len(nav.titles))
	assert.Equal(t, true, nav.button.Disabled())
}

func TestNavigation_PushWithTitle(t *testing.T) {
	a := widget.NewLabel("a")
	b := widget.NewLabel("b")
	c := widget.NewLabel("c")
	nav := NewNavigation("Title")

	assert.Equal(t, true, nav.button.Disabled())
	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))

	nav.PushWithTitle(a, "A")
	assert.Equal(t, "A", nav.label.Text)
	assert.Equal(t, true, nav.button.Disabled())

	nav.PushWithTitle(b, "B")
	assert.Equal(t, "B", nav.label.Text)
	assert.Equal(t, false, nav.button.Disabled())

	nav.PushWithTitle(c, "C")
	assert.Equal(t, "C", nav.label.Text)
	assert.Equal(t, false, nav.button.Disabled())

	assert.Equal(t, c, nav.Pop())
	assert.Equal(t, "B", nav.label.Text)
	assert.Equal(t, 2, len(nav.stack.Objects))
	assert.Equal(t, 2, len(nav.titles))
	assert.Equal(t, false, nav.button.Disabled())

	assert.Equal(t, b, nav.Pop())
	assert.Equal(t, "A", nav.label.Text)
	assert.Equal(t, 1, len(nav.stack.Objects))
	assert.Equal(t, 1, len(nav.titles))
	assert.Equal(t, true, nav.button.Disabled())

	assert.Equal(t, a, nav.Pop())
	assert.Equal(t, "Title", nav.label.Text)
	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))
	assert.Equal(t, true, nav.button.Disabled())

	assert.Nil(t, nav.Pop())
	assert.Equal(t, "Title", nav.label.Text)
	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))
	assert.Equal(t, true, nav.button.Disabled())
}

func TestNavigation_Empty(t *testing.T) {
	nav := NewNavigation("Title")

	assert.Equal(t, true, nav.button.Disabled())
	assert.Equal(t, 0, len(nav.stack.Objects))
	assert.Equal(t, 0, len(nav.titles))

	assert.Nil(t, nav.Pop())

	assert.Equal(t, "Title", nav.label.Text)
}
