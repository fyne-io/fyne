package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFloatList(t *testing.T) {
	f := NewFloatList()
	assert.Equal(t, 0, f.Length())
}

func TestFloatList_Append(t *testing.T) {
	f := NewFloatList()
	assert.Equal(t, 0, f.Length())

	f.Append(0.5)
	assert.Equal(t, 1, f.Length())
}

func TestFloatList_Get(t *testing.T) {
	f := NewFloatList()

	f.Append(1.3)
	assert.Equal(t, 1.3, f.Get(0))

	f.Append(0.2)
	assert.Equal(t, 0.2, f.Get(1))
}
