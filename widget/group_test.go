package widget

import (
	"testing"

	_ "github.com/fyne-io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestGroupSize(t *testing.T) {
	group := NewGroup("Test", NewLabel("Hello"), NewLabel("World"))
	assert.Equal(t, 2, len(group.box.Children))
}

func TestGroupPrepend(t *testing.T) {
	group := NewGroup("Test", NewLabel("World"))
	assert.Equal(t, 1, len(group.box.Children))

	prepend := NewLabel("Hello")
	group.Prepend(prepend)
	assert.Equal(t, 2, len(group.box.Children))
	assert.Equal(t, prepend, group.box.Children[0])
}

func TestGroupAppend(t *testing.T) {
	group := NewGroup("Test", NewLabel("Hello"))
	assert.Equal(t, 1, len(group.box.Children))

	append := NewLabel("World")
	group.Append(append)
	assert.True(t, len(group.box.Children) == 2)
	assert.Equal(t, append, group.box.Children[1])
}
