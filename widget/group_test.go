package widget

import "testing"

import "github.com/stretchr/testify/assert"

import _ "github.com/fyne-io/fyne/test"

func TestGroupSize(t *testing.T) {
	group := NewGroup("Test", NewLabel("Hello"), NewLabel("World"))
	assert.Equal(t, 2, len(group.list.Children))
}

func TestGroupPrepend(t *testing.T) {
	group := NewGroup("Test", NewLabel("World"))
	assert.Equal(t, 1, len(group.list.Children))

	prepend := NewLabel("Hello")
	group.Prepend(prepend)
	assert.Equal(t, 2, len(group.list.Children))
	assert.Equal(t, prepend, group.list.Children[0])
}

func TestGroupAppend(t *testing.T) {
	group := NewGroup("Test", NewLabel("Hello"))
	assert.Equal(t, 1, len(group.list.Children))

	append := NewLabel("World")
	group.Append(append)
	assert.True(t, len(group.list.Children) == 2)
	assert.Equal(t, append, group.list.Children[1])
}
