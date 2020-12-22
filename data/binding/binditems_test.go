package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindFloat(t *testing.T) {
	val := 0.5
	f := BindFloat(&val)
	assert.Equal(t, float64(0.5), f.Get())

	f.Set(0.3)
	assert.Equal(t, float64(0.3), val)
}

func TestNewFloat(t *testing.T) {
	f := NewFloat()
	assert.Equal(t, float64(0.0), f.Get())

	f.Set(0.3)
	assert.Equal(t, float64(0.3), f.Get())
}
