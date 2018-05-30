package widget

import "testing"

import "github.com/stretchr/testify/assert"

import _ "github.com/fyne-io/fyne/test"

func TestListSize(t *testing.T) {
	list := NewList(NewLabel("World"))
	assert.True(t, list.MinSize().Height > 0)

	min := list.MinSize()
	list.Prepend(NewLabel("Hello"))
	assert.True(t, list.MinSize().Height > min.Height)
}
