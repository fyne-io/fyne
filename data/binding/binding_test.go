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
	assert.Equal(t, 0, len(data.listeners))

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Equal(t, 1, len(data.listeners))
	assert.True(t, called)
}

func TestBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := &simpleItem{}
	data.listeners = append(data.listeners, fn)

	assert.Equal(t, 1, len(data.listeners))
	data.RemoveListener(fn)
	assert.Equal(t, 0, len(data.listeners))

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

	a.Set(0)

	ival, err := a.Get()
	assert.NoError(t, err)
	val, ok := ival.(int)
	assert.Equal(t, 0, val)
	assert.True(t, ok)

	a.Set(nil)
	ival, err = a.Get()
	assert.NoError(t, err)
	assert.Equal(t, nil, ival)
}
