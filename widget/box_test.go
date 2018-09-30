package widget

import "testing"

import "github.com/stretchr/testify/assert"

import _ "github.com/fyne-io/fyne/test"

func TestBoxSize(t *testing.T) {
	list := NewVBox(NewLabel("Hello"), NewLabel("World"))
	assert.Equal(t, 2, len(list.Children))
}

func TestBoxPrepend(t *testing.T) {
	list := NewVBox(NewLabel("World"))
	assert.Equal(t, 1, len(list.Children))

	prepend := NewLabel("Hello")
	list.Prepend(prepend)
	assert.Equal(t, 2, len(list.Children))
	assert.Equal(t, prepend, list.Children[0])
}

func TestBoxAppend(t *testing.T) {
	list := NewVBox(NewLabel("Hello"))
	assert.Equal(t, 1, len(list.Children))

	append := NewLabel("World")
	list.Append(append)
	assert.True(t, len(list.Children) == 2)
	assert.Equal(t, append, list.Children[1])
}
