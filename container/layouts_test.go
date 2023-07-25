package container

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2/canvas"
)

func TestNewBorder_Nil(t *testing.T) {
	b := NewBorder(nil, nil, nil, nil)
	assert.Equal(t, 0, len(b.Objects))
	b = NewBorder(canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black))
	assert.Equal(t, 4, len(b.Objects))

	b = NewBorder(canvas.NewRectangle(color.Black), nil, nil, nil, nil) // a common error - adding nil to component list
	assert.Equal(t, 1, len(b.Objects))
	b = NewBorder(canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), canvas.NewRectangle(color.Black), nil)
	assert.Equal(t, 4, len(b.Objects))
}
