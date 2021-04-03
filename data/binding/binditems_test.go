package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindFloat(t *testing.T) {
	val := 0.5
	f := BindFloat(&val)
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.5, v)

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	f.AddListener(fn)
	waitForItems()
	assert.True(t, called)

	called = false
	err = f.Set(0.3)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 0.3, val)
	assert.True(t, called)

	called = false
	val = 1.2
	f.Reload()
	waitForItems()
	assert.True(t, called)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 1.2, v)
}

func TestNewFloat(t *testing.T) {
	f := NewFloat()
	v, err := f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, v)

	err = f.Set(0.3)
	assert.Nil(t, err)
	v, err = f.Get()
	assert.Nil(t, err)
	assert.Equal(t, 0.3, v)
}

func TestNewFloat_TriggerOnlyWhenChange(t *testing.T) {
	f := NewFloat()
	called := 0
	f.AddListener(NewDataListener(func() {
		called++
	}))
	waitForItems()
	assert.Equal(t, 1, called)

	for i := 0; i < 5; i++ {
		err := f.Set(9.0)
		assert.Nil(t, err)
		waitForItems()
		assert.Equal(t, 2, called)
	}

	err := f.Set(9.1)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 3, called)

	f.Set(8.0)
	assert.Nil(t, err)
	waitForItems()
	assert.Equal(t, 4, called)
}
