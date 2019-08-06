package driver

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	_ "fyne.io/fyne/test"
	"fyne.io/fyne/widget"
)

func TestWalkObjectTree(t *testing.T) {
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(100, 100))
	child := canvas.NewRectangle(color.Black)
	base := fyne.NewContainerWithLayout(layout.NewGridLayout(1), rect, widget.NewScrollContainer(child))

	clipPos := fyne.NewPos(0, 0)
	clipSize := rect.MinSize()

	WalkObjectTree(base, func(object fyne.CanvasObject, position fyne.Position, clippingPos fyne.Position, clippingSize fyne.Size) bool {
		if _, ok := object.(*widget.ScrollContainer); ok {
			clipPos = clippingPos
			clipSize = clippingSize
		}
		return false
	}, nil)

	assert.Equal(t, fyne.NewPos(0, 104), clipPos)
	assert.Equal(t, fyne.NewSize(100, 100), clipSize)
}
