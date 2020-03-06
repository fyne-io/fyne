package widget

import (
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestNewPopUpMenu(t *testing.T) {
	c := test.Canvas()
	menu := fyne.NewMenu("Foo", fyne.NewMenuItem("Bar", func() {}))

	pop := NewPopUpMenu(menu, c)
	assert.Equal(t, 1, len(c.Overlays().All()))
	assert.Equal(t, pop, c.Overlays().All()[0])

	pop.Hide()
	assert.Equal(t, 0, len(c.Overlays().All()))
}
