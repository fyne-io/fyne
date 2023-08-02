package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeBase_AddListener(t *testing.T) {
	data := newSimpleTree()
	assert.Equal(t, 0, syncMapLen(&data.listeners))

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Equal(t, 1, syncMapLen(&data.listeners))

	data.trigger()
	waitForItems()
	assert.True(t, called)
}

func TestTreeBase_GetItem(t *testing.T) {
	data := newSimpleTree()
	f := 0.5
	data.appendItem(BindFloat(&f), "f", "")
	assert.Equal(t, 1, len(data.items))

	item, err := data.GetItem("f")
	assert.Nil(t, err)
	val, err := item.(Float).Get()
	assert.Nil(t, err)
	assert.Equal(t, f, val)

	_, err = data.GetItem("g")
	assert.NotNil(t, err)
}

func TestListBase_IDs(t *testing.T) {
	data := newSimpleTree()
	assert.Equal(t, 0, len(data.ChildIDs("")))

	data.appendItem(NewFloat(), "1", "")
	assert.Equal(t, 1, len(data.ChildIDs("")))
	assert.Equal(t, "1", data.ChildIDs("")[0])
}

func TestTreeBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := newSimpleTree()
	data.listeners.Store(fn, true)

	assert.Equal(t, 1, syncMapLen(&data.listeners))
	data.RemoveListener(fn)
	assert.Equal(t, 0, syncMapLen(&data.listeners))

	data.trigger()
	waitForItems()
	assert.False(t, called)
}

type simpleTree struct {
	treeBase
}

func newSimpleTree() *simpleTree {
	t := &simpleTree{}
	t.ids = map[string][]string{}
	t.items = map[string]DataItem{}

	return t
}
