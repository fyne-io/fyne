package widget

import (
	"testing"

	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestNewPopOver(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopOver(label, test.Canvas())

	assert.True(t, pop.Visible())
	assert.Equal(t, pop, test.Canvas().Overlay())
}

func TestPopOver_Hide(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopOver(label, test.Canvas())

	assert.True(t, pop.Visible())
	pop.Hide()
	assert.False(t, pop.Visible())
	assert.Nil(t, test.Canvas().Overlay())
}

func TestPopOver_MinSize(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopOver(label, test.Canvas())
	min := pop.MinSize()

	assert.True(t, min.Width > label.MinSize().Width)
	assert.True(t, min.Height > label.MinSize().Height)
}

func TestPopOver_Tapped(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewPopOver(label, test.Canvas())

	assert.True(t, pop.Visible())
	test.Tap(pop)
	assert.False(t, pop.Visible())
	assert.Nil(t, test.Canvas().Overlay())
}

func TestModalPopOver_Tapped(t *testing.T) {
	label := NewLabel("Hi")
	pop := NewModalPopOver(label, test.Canvas())

	assert.True(t, pop.Visible())
	test.Tap(pop)
	assert.True(t, pop.Visible())
	assert.Equal(t, pop, test.Canvas().Overlay())
}
