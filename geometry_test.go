package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition_Add(t *testing.T) {
	pos1 := NewPos(10, 10)
	pos2 := NewPos(25, 25)

	pos3 := pos1.Add(pos2)

	assert.Equal(t, float32(35), pos3.X)
	assert.Equal(t, float32(35), pos3.Y)
}

func TestPosition_Add_Size(t *testing.T) {
	pos1 := NewPos(10, 10)
	s := NewSize(25, 25)

	pos2 := pos1.Add(s)

	assert.Equal(t, float32(35), pos2.X)
	assert.Equal(t, float32(35), pos2.Y)
}

func TestPosition_Add_Vector(t *testing.T) {
	pos1 := NewPos(10, 10)
	v := NewDelta(25, 25)

	pos2 := pos1.Add(v)

	assert.Equal(t, float32(35), pos2.X)
	assert.Equal(t, float32(35), pos2.Y)
}

func TestPosition_IsZero(t *testing.T) {
	for name, tt := range map[string]struct {
		p    Position
		want bool
	}{
		"zero value":       {Position{}, true},
		"0,0":              {NewPos(0, 0), true},
		"zero X":           {NewPos(0, 42), false},
		"zero Y":           {NewPos(17, 0), false},
		"non-zero X and Y": {NewPos(6, 9), false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.IsZero())
		})
	}
}

func TestPosition_Subtract(t *testing.T) {
	pos1 := NewPos(25, 25)
	pos2 := NewPos(10, 10)

	pos3 := pos1.Subtract(pos2)

	assert.Equal(t, float32(15), pos3.X)
	assert.Equal(t, float32(15), pos3.Y)
}

func TestSize_Add(t *testing.T) {
	size1 := NewSize(10, 10)
	size2 := NewSize(25, 25)

	size3 := size1.Add(size2)

	assert.Equal(t, float32(35), size3.Width)
	assert.Equal(t, float32(35), size3.Height)
}

func TestSize_Add_Position(t *testing.T) {
	size1 := NewSize(10, 10)
	p := NewSize(25, 25)

	size2 := size1.Add(p)

	assert.Equal(t, float32(35), size2.Width)
	assert.Equal(t, float32(35), size2.Height)
}

func TestSize_Add_Vector(t *testing.T) {
	size1 := NewSize(10, 10)
	v := NewDelta(25, 25)

	size2 := size1.Add(v)

	assert.Equal(t, float32(35), size2.Width)
	assert.Equal(t, float32(35), size2.Height)
}

func TestSize_IsZero(t *testing.T) {
	for name, tt := range map[string]struct {
		s    Size
		want bool
	}{
		"zero value":    {Size{}, true},
		"0x0":           {NewSize(0, 0), true},
		"zero width":    {NewSize(0, 42), false},
		"zero height":   {NewSize(17, 0), false},
		"non-zero area": {NewSize(6, 9), false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.IsZero())
		})
	}
}

func TestSize_Max(t *testing.T) {
	size1 := NewSize(10, 100)
	size2 := NewSize(100, 10)

	size3 := size1.Max(size2)

	assert.Equal(t, float32(100), size3.Width)
	assert.Equal(t, float32(100), size3.Height)
}

func TestSize_Min(t *testing.T) {
	size1 := NewSize(10, 100)
	size2 := NewSize(100, 10)

	size3 := size1.Min(size2)

	assert.Equal(t, float32(10), size3.Width)
	assert.Equal(t, float32(10), size3.Height)
}

func TestSize_Subtract(t *testing.T) {
	size1 := NewSize(25, 25)
	size2 := NewSize(10, 10)

	size3 := size1.Subtract(size2)

	assert.Equal(t, float32(15), size3.Width)
	assert.Equal(t, float32(15), size3.Height)
}

func TestVector_IsZero(t *testing.T) {
	v := NewDelta(0, 0)

	assert.True(t, v.IsZero())

	v.DX = 1
	assert.False(t, v.IsZero())
}
