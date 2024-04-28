package binding

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func syncMapLen(m *sync.Map) (n int) {
	m.Range(func(_, _ any) bool {
		n++
		return true
	})
	return
}

type simpleItem struct {
	base
}

func TestBase_AddListener(t *testing.T) {
	data := &simpleItem{}
	assert.Equal(t, 0, syncMapLen(&data.listeners))

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Equal(t, 1, syncMapLen(&data.listeners))

	waitForItems()
	assert.True(t, called)
}

func TestBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := &simpleItem{}
	data.listeners.Store(fn, true)

	assert.Equal(t, 1, syncMapLen(&data.listeners))
	data.RemoveListener(fn)
	assert.Equal(t, 0, syncMapLen(&data.listeners))

	waitForItems()
	data.trigger()
	assert.False(t, called)
}

func TestNewDataItemListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})

	waitForItems()
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
