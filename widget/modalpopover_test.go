package widget

import (
	"testing"

	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestNewModalPopOver(t *testing.T) {
	label := NewLabel("Modal")
	pop := NewModalPopOver(label, test.Canvas())

	assert.True(t, pop.Visible())
	assert.Equal(t, pop, test.Canvas().Overlay())
}

func TestModalPopOver_Hide(t *testing.T) {
	label := NewLabel("Modal")
	pop := NewModalPopOver(label, test.Canvas())

	assert.True(t, pop.Visible())
	pop.Hide()
	assert.False(t, pop.Visible())
	assert.Nil(t, test.Canvas().Overlay())
}

func TestModalPopOver_MinSize(t *testing.T) {
	label := NewLabel("Modal")
	pop := NewModalPopOver(label, test.Canvas())
	min := pop.MinSize()

	assert.True(t, min.Width > label.MinSize().Width)
	assert.True(t, min.Height > label.MinSize().Height)
}
