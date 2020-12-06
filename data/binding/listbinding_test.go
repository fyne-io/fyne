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
	assert.Equal(t, 1, len(data.val))

	assert.Equal(t, f, data.GetItem(0).(Float).Get())
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
