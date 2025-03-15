package binding

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNot(t *testing.T) {
	var b bool = true

	bb := BindBool(&b)
	notbb := Not(bb)

	assert.NotNil(t, bb)
	assert.NotNil(t, notbb)

	notb, err := notbb.Get()
	require.NoError(t, err)
	assert.Equal(t, !b, notb)
	assert.False(t, notb)

	err = notbb.Set(true)
	require.NoError(t, err)
	assert.False(t, b)
}

func TestAnd(t *testing.T) {
	b := []bool{false, false, false, false, false}

	var bb []Bool
	for idx := range b {
		bb = append(bb, BindBool(&b[idx]))
	}

	andbb := And(bb...)
	assert.NotNil(t, andbb)

	setAtOffset := func(offset, value int) {
		if value&(1<<offset) != 0 {
			bb[offset].Set(true)
		} else {
			bb[offset].Set(false)
		}
	}

	for i := 0; i < 32; i++ {
		for idx := range bb {
			setAtOffset(idx, i)
		}
		log.Println(b)

		a := true
		for _, v := range b {
			if v == false {
				a = false
			}
		}

		andb, err := andbb.Get()
		require.NoError(t, err)
		assert.Equal(t, a, andb)
	}
	for _, v := range b {
		assert.True(t, v)
	}
}

func TestOr(t *testing.T) {
	b := []bool{false, false, false, false, false}

	var bb []Bool
	for idx := range b {
		bb = append(bb, BindBool(&b[idx]))
	}

	andbb := Or(bb...)
	assert.NotNil(t, andbb)

	setAtOffset := func(offset, value int) {
		if value&(1<<offset) != 0 {
			bb[offset].Set(true)
		} else {
			bb[offset].Set(false)
		}
	}

	for i := 0; i < 32; i++ {
		for idx := range bb {
			setAtOffset(idx, i)
		}
		log.Println(b)

		a := false
		for _, v := range b {
			if v == true {
				a = true
			}
		}

		andb, err := andbb.Get()
		require.NoError(t, err)
		assert.Equal(t, a, andb)
	}
	for _, v := range b {
		assert.True(t, v)
	}
}
