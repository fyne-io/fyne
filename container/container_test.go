package container

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/stretchr/testify/assert"
)

func TestContainer_NoResize(t *testing.T) {
	rect := &canvas.Rectangle{}
	rect.SetMinSize(fyne.NewSize(100, 100))

	container := NewHBox(rect)
	assert.Equal(t, fyne.NewSize(0, 0), container.Size())

	container.Resize(container.MinSize())
	assert.Equal(t, rect.MinSize(), container.Size())
}
