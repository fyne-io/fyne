package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type simpleItem struct {
	base
}

func TestBase_AddListener(t *testing.T) {
	data := &simpleItem{}
	assert.Equal(t, 0, data.listeners.Len())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Equal(t, 1, data.listeners.Len())
	assert.True(t, called)
}

func TestBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := &simpleItem{}
	data.listeners.Store(fn, true)

	assert.Equal(t, 1, data.listeners.Len())
	data.RemoveListener(fn)
	assert.Equal(t, 0, data.listeners.Len())

	data.trigger()
	assert.False(t, called)
}

func TestNewDataItemListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})

	fn.DataChanged()
	assert.True(t, called)
}

func TestBindAnyWithNil(t *testing.T) {
	a := NewUntyped()
	a.Set(nil)
	b := 1
	a.Set(b)
	var tr any = nil
	a.Set(tr)
}
