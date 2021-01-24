// +build !ci
// +build !mobile
// +build !darwin no_native_menus

package glfw

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/stretchr/testify/assert"
)

func TestGlCanvas_FocusHandlingWhenActivatingOrDeactivatingTheMenu(t *testing.T) {
	w := createWindow("Test")
	w.SetPadded(false)
	w.SetMainMenu(
		fyne.NewMainMenu(
			fyne.NewMenu("test", fyne.NewMenuItem("item", func() {})),
			fyne.NewMenu("other", fyne.NewMenuItem("item", func() {})),
		),
	)
	c := w.Canvas().(*glCanvas)

	ce1 := &focusable{id: "ce1"}
	ce2 := &focusable{id: "ce2"}
	content := container.NewVBox(ce1, ce2)
	w.SetContent(content)

	assert.Nil(t, c.Focused())
	m := c.menu.(*MenuBar)
	assert.False(t, m.IsActive())

	c.FocusPrevious()
	assert.Equal(t, ce2, c.Focused())
	assert.True(t, ce2.focused)

	m.Items[0].(*menuBarItem).Tapped(&fyne.PointEvent{})
	assert.True(t, m.IsActive())
	ctxt := "activating the menu changes focus handler and focuses the menu bar item but does not remove focus from content"
	assert.True(t, ce2.focused, ctxt)
	assert.Equal(t, m.Items[0], c.Focused(), ctxt)

	c.FocusNext()
	ctxt = "changing focus with active menu does not affect content focus"
	assert.True(t, ce2.focused, ctxt)
	assert.Equal(t, m.Items[1], c.Focused(), ctxt)

	m.Items[0].(*menuBarItem).Tapped(&fyne.PointEvent{})
	assert.False(t, m.IsActive())
	ctxt = "deactivating the menu restores focus handler from content"
	assert.True(t, ce2.focused, ctxt)
	assert.Equal(t, ce2, c.Focused(), ctxt)

	c.FocusPrevious()
	assert.Equal(t, ce1, c.Focused(), ctxt)
	assert.True(t, ce1.focused, ctxt)
	assert.False(t, ce2.focused, ctxt)
}
