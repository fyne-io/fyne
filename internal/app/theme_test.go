package app_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestApplySettings_BeforeContentSet(*testing.T) {
	a := test.NewApp()
	w := a.NewWindow("NoContent")
	defer w.Close()

	app.ApplySettings(a.Settings(), a)
}

type countingObject struct {
	fyne.CanvasObject
	refreshes int
}

func (c *countingObject) Refresh() {
	c.refreshes++
}

func TestApplyThemeTo_RefreshesEachObjectOnce(t *testing.T) {
	a := test.NewApp()
	defer a.Quit()

	leaf := &countingObject{CanvasObject: canvas.NewRectangle(color.Transparent)}
	inner := container.NewVBox(leaf)
	middle := container.NewVBox(inner)
	root := container.NewScroll(container.NewVBox(middle))

	app.ApplyThemeTo(root, nil)

	assert.Equal(t, 1, leaf.refreshes, "a leaf must be refreshed once per theme change, not once per ancestor")

	leaf.refreshes = 0
	root.Refresh()
	assert.Equal(t, 1, leaf.refreshes, "a normal refresh still cascades to children")
}
