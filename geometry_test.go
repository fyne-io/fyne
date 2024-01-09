package fyne_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewSquareOffsetPosition(t *testing.T) {
	pos := fyne.NewSquareOffsetPos(10)
	assert.Equal(t, pos.X, pos.Y)
}

func TestPosition_Add(t *testing.T) {
	pos1 := fyne.NewPos(10, 10)
	pos2 := fyne.NewPos(25, 25)

	pos3 := pos1.Add(pos2)

	assert.Equal(t, float32(35), pos3.X)
	assert.Equal(t, float32(35), pos3.Y)
}

func TestPosition_Add_Size(t *testing.T) {
	pos1 := fyne.NewPos(10, 10)
	s := fyne.NewSize(25, 25)

	pos2 := pos1.Add(s)

	assert.Equal(t, float32(35), pos2.X)
	assert.Equal(t, float32(35), pos2.Y)
}

func TestPosition_Add_Vector(t *testing.T) {
	pos1 := fyne.NewPos(10, 10)
	v := fyne.NewDelta(25, 25)

	pos2 := pos1.Add(v)

	assert.Equal(t, float32(35), pos2.X)
	assert.Equal(t, float32(35), pos2.Y)
}

func TestPosition_AddXY(t *testing.T) {
	pos1 := fyne.NewPos(10, 10)

	pos2 := pos1.AddXY(25, 25)

	assert.Equal(t, float32(35), pos2.X)
	assert.Equal(t, float32(35), pos2.Y)
}

func TestPosition_IsZero(t *testing.T) {
	for name, tt := range map[string]struct {
		p    fyne.Position
		want bool
	}{
		"zero value":       {fyne.Position{}, true},
		"0,0":              {fyne.NewPos(0, 0), true},
		"zero X":           {fyne.NewPos(0, 42), false},
		"zero Y":           {fyne.NewPos(17, 0), false},
		"non-zero X and Y": {fyne.NewPos(6, 9), false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.IsZero())
		})
	}
}

func TestPosition_Subtract(t *testing.T) {
	pos1 := fyne.NewPos(25, 25)
	pos2 := fyne.NewPos(10, 10)

	pos3 := pos1.Subtract(pos2)

	assert.Equal(t, float32(15), pos3.X)
	assert.Equal(t, float32(15), pos3.Y)
}

func TestPosition_SubtractXY(t *testing.T) {
	pos1 := fyne.NewPos(25, 25)

	pos2 := pos1.SubtractXY(10, 10)

	assert.Equal(t, float32(15), pos2.X)
	assert.Equal(t, float32(15), pos2.Y)
}

func TestSize_Add(t *testing.T) {
	size1 := fyne.NewSize(10, 10)
	size2 := fyne.NewSize(25, 25)

	size3 := size1.Add(size2)

	assert.Equal(t, float32(35), size3.Width)
	assert.Equal(t, float32(35), size3.Height)
}

func TestSize_Add_Position(t *testing.T) {
	size1 := fyne.NewSize(10, 10)
	p := fyne.NewSize(25, 25)

	size2 := size1.Add(p)

	assert.Equal(t, float32(35), size2.Width)
	assert.Equal(t, float32(35), size2.Height)
}

func TestSize_Add_Vector(t *testing.T) {
	size1 := fyne.NewSize(10, 10)
	v := fyne.NewDelta(25, 25)

	size2 := size1.Add(v)

	assert.Equal(t, float32(35), size2.Width)
	assert.Equal(t, float32(35), size2.Height)
}

func TestSize_AddWidthHeight(t *testing.T) {
	size1 := fyne.NewSize(10, 10)

	size2 := size1.AddWidthHeight(25, 25)

	assert.Equal(t, float32(35), size2.Width)
	assert.Equal(t, float32(35), size2.Height)
}

func TestSize_IsZero(t *testing.T) {
	for name, tt := range map[string]struct {
		s    fyne.Size
		want bool
	}{
		"zero value":    {fyne.Size{}, true},
		"0x0":           {fyne.NewSize(0, 0), true},
		"zero width":    {fyne.NewSize(0, 42), false},
		"zero height":   {fyne.NewSize(17, 0), false},
		"non-zero area": {fyne.NewSize(6, 9), false},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.IsZero())
		})
	}
}

func TestSize_Max(t *testing.T) {
	size1 := fyne.NewSize(10, 100)
	size2 := fyne.NewSize(100, 10)

	size3 := size1.Max(size2)

	assert.Equal(t, float32(100), size3.Width)
	assert.Equal(t, float32(100), size3.Height)
}

func TestSize_Min(t *testing.T) {
	size1 := fyne.NewSize(10, 100)
	size2 := fyne.NewSize(100, 10)

	size3 := size1.Min(size2)

	assert.Equal(t, float32(10), size3.Width)
	assert.Equal(t, float32(10), size3.Height)
}

func TestSize_Subtract(t *testing.T) {
	size1 := fyne.NewSize(25, 25)
	size2 := fyne.NewSize(10, 10)

	size3 := size1.Subtract(size2)

	assert.Equal(t, float32(15), size3.Width)
	assert.Equal(t, float32(15), size3.Height)
}

func TestSize_SubtractWidthHeight(t *testing.T) {
	size1 := fyne.NewSize(25, 25)

	size2 := size1.SubtractWidthHeight(10, 10)

	assert.Equal(t, float32(15), size2.Width)
	assert.Equal(t, float32(15), size2.Height)
}

func TestSquareSize(t *testing.T) {
	size := fyne.NewSquareSize(10)
	assert.Equal(t, size.Height, size.Width)
}

func TestVector_IsZero(t *testing.T) {
	v := fyne.NewDelta(0, 0)

	assert.True(t, v.IsZero())

	v.DX = 1
	assert.False(t, v.IsZero())
}
