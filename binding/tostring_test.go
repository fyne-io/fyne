package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatToString(t *testing.T) {
	f := NewFloat()
	s := FloatToString(f)
	assert.Equal(t, "0.00", s.Get())

	f.Set(0.3)
	assert.Equal(t, "0.30", s.Get())

	s.Set("5.00")
	assert.Equal(t, 5.0, f.Get())
}
