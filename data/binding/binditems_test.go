package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindFloat(t *testing.T) {
	val := 0.5
	f := BindFloat(&val)
	assert.Equal(t, 0.5, f.Get())

	f.Set(0.3)
	assert.Equal(t, 0.3, val)
}

func TestBindFloatUpdate(t *testing.T) {
	val := 0.5
	called := false
	f := BindFloat(&val)
	f.AddListener(NewDataListener(func() {
		called = true
	}))
	waitForItems()
	assert.Equal(t, 0.5, f.Get())
	assert.True(t, called)

	called = false
	val = 1.2
	f.Reload()
	waitForItems()
	assert.Equal(t, 1.2, val)
	assert.True(t, called)
}

func TestNewFloat(t *testing.T) {
	f := NewFloat()
	assert.Equal(t, 0.0, f.Get())

	f.Set(0.3)
	assert.Equal(t, 0.3, f.Get())
}
