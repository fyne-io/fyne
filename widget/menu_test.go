package widget

import (
	"image/color"
	"testing"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/test"
	"fyne.io/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewPopUpMenu(t *testing.T) {
	c := test.Canvas()
	menu := fyne.NewMenu("Foo", fyne.NewMenuItem("Bar", func() {}))

	pop := NewPopUpMenu(menu, c)
	assert.Equal(t, 1, len(c.Overlays().List()))
	assert.Equal(t, pop, c.Overlays().List()[0])

	pop.Hide()
	assert.Equal(t, 0, len(c.Overlays().List()))
}

func TestPopUpMenu_Size(t *testing.T) {
	win := test.NewWindow(NewLabel("OK"))
	defer win.Close()
	win.Resize(fyne.NewSize(100, 100))
	menu := fyne.NewMenu("Foo",
		fyne.NewMenuItem("A", func() {}),
		fyne.NewMenuItem("A", func() {}),
	)
	menuItemSize := canvas.NewText("A", color.Black).MinSize().Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2))
	expectedSize := menuItemSize.Add(fyne.NewSize(0, menuItemSize.Height)).Add(fyne.NewSize(0, theme.Padding()))
	c := win.Canvas()

	pop := NewPopUpMenu(menu, c)
	defer pop.Hide()
	assert.Equal(t, expectedSize, pop.Content.Size())

	for _, o := range test.LaidOutObjects(pop) {
		if s, ok := o.(*widget.Shadow); ok {
			assert.Equal(t, expectedSize.Add(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)), s.Size())
		}
	}
}
