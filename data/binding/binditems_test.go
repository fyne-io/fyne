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

	err = f.Set(0.3)
	assert.Nil(t, err)
	assert.Equal(t, 0.3, val)
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
