package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindFloatList(t *testing.T) {
	l := []float64{1.0, 5.0, 2.3}
	f := BindFloatList(&l)

	assert.Equal(t, 3, f.Length())
	v, err := f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)

	assert.NotNil(t, f.(*boundFloatList).val)
	assert.Equal(t, 3, len(*(f.(*boundFloatList).val)))

	_, err = f.GetValue(-1)
	assert.NotNil(t, err)
}

func TestNewFloatList(t *testing.T) {
	f := NewFloatList()
	assert.Equal(t, 0, f.Length())

	_, err := f.GetValue(-1)
	assert.NotNil(t, err)
}

func TestFloatList_Append(t *testing.T) {
	f := NewFloatList()
	assert.Equal(t, 0, f.Length())

	f.Append(0.5)
	assert.Equal(t, 1, f.Length())
}

func TestFloatList_Get(t *testing.T) {
	f := NewFloatList()

	err := f.Append(1.3)
	assert.Nil(t, err)
	v, err := f.GetValue(0)
	assert.Nil(t, err)
	assert.Equal(t, 1.3, v)

	err = f.Append(0.2)
	assert.Nil(t, err)
	v, err = f.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, 0.2, v)
}
