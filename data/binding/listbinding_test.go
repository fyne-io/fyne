package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleList struct {
	listBase
}

func TestListBase_AddListener(t *testing.T) {
	data := &simpleList{}
	assert.Equal(t, 0, len(data.listeners))

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Equal(t, 1, len(data.listeners))

	data.trigger()
	waitForItems()
	assert.True(t, called)
}

func TestListBase_GetItem(t *testing.T) {
	data := &simpleList{}
	f := 0.5
	data.appendItem(BindFloat(&f))
	assert.Equal(t, 1, len(data.items))

	item, err := data.GetItem(0)
	assert.Nil(t, err)
	val, err := item.(Float).Get()
	assert.Nil(t, err)
	assert.Equal(t, f, val)

	_, err = data.GetItem(5)
	assert.NotNil(t, err)
}

func TestListBase_Length(t *testing.T) {
	data := &simpleList{}
	assert.Equal(t, 0, data.Length())

	data.appendItem(NewFloat())
	assert.Equal(t, 1, data.Length())
}

func TestListBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := &simpleList{}
	data.listeners = []DataListener{fn}

	assert.Equal(t, 1, len(data.listeners))
	data.RemoveListener(fn)
	assert.Equal(t, 0, len(data.listeners))

	data.trigger()
	waitForItems()
	assert.False(t, called)
}

func TestNewDataListListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})

	fn.DataChanged()
	assert.True(t, called)
}
