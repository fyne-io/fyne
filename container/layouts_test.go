package container

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/canvas"
)

func TestNewBorder_Nil(t *testing.T) {
	b := NewBorder(nil, nil, nil, nil)
	assert.Empty(t, b.Objects)
	b = NewBorder(canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black))
	assert.Len(t, b.Objects, 4)

	b = NewBorder(canvas.NewRectangle(color.Black), nil, nil, nil, nil) // a common error - adding nil to component list
	assert.Len(t, b.Objects, 1)
	b = NewBorder(canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), nil)
	assert.Len(t, b.Objects, 4)
}
