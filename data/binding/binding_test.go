package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type simpleItem struct {
	base
}

func TestBase_AddListener(t *testing.T) {
	data := &simpleItem{}
	assert.Empty(t, data.listeners)

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Len(t, data.listeners, 1)
	assert.True(t, called)
}

func TestBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := &simpleItem{}
	data.listeners = append(data.listeners, fn)

	assert.Len(t, data.listeners, 1)
	data.RemoveListener(fn)
	assert.Empty(t, data.listeners)

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
	require.NoError(t, err)
	val, ok := ival.(int)
	assert.Equal(t, 0, val)
	assert.True(t, ok)

	a.Set(nil)
	ival, err = a.Get()
	require.NoError(t, err)
	assert.Nil(t, ival)
}
