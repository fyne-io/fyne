package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolToString(t *testing.T) {
	b := NewBool()
	s := BoolToString(b)
	assert.Equal(t, "false", s.Get())

	b.Set(true)
	assert.Equal(t, "true", s.Get())

	s.Set("false")
	assert.Equal(t, false, b.Get())
}

func TestFloatToString(t *testing.T) {
	f := NewFloat()
	s := FloatToString(f)
	assert.Equal(t, "0.000000", s.Get())

	f.Set(0.3)
	assert.Equal(t, "0.300000", s.Get())

	s.Set("5.00")
	assert.Equal(t, 5.0, f.Get())
}

func TestIntToString(t *testing.T) {
	i := NewInt()
	s := IntToString(i)
	assert.Equal(t, "0", s.Get())

	i.Set(3)
	assert.Equal(t, "3", s.Get())

	s.Set("5")
	assert.Equal(t, 5, i.Get())
}
