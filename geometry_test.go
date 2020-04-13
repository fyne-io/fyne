package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosition_Add(t *testing.T) {
	pos1 := NewPos(10, 10)
	pos2 := NewPos(25, 25)

	pos3 := pos1.Add(pos2)

	assert.Equal(t, 35, pos3.X)
	assert.Equal(t, 35, pos3.Y)
}

func TestPosition_Subtract(t *testing.T) {
	pos1 := NewPos(25, 25)
	pos2 := NewPos(10, 10)

	pos3 := pos1.Subtract(pos2)

	assert.Equal(t, 15, pos3.X)
	assert.Equal(t, 15, pos3.Y)
}

func TestSizeAdd(t *testing.T) {
	size1 := NewSize(10, 10)
	size2 := NewSize(25, 25)

	size3 := size1.Add(size2)

	assert.Equal(t, 35, size3.Width)
	assert.Equal(t, 35, size3.Height)
}

func TestSizeMax(t *testing.T) {
	size1 := NewSize(10, 100)
	size2 := NewSize(100, 10)

	size3 := size1.Max(size2)

	assert.Equal(t, 100, size3.Width)
	assert.Equal(t, 100, size3.Height)
}

func TestSizeMin(t *testing.T) {
	size1 := NewSize(10, 100)
	size2 := NewSize(100, 10)

	size3 := size1.Min(size2)

	assert.Equal(t, 10, size3.Width)
	assert.Equal(t, 10, size3.Height)
}

func TestSizeSubtract(t *testing.T) {
	size1 := NewSize(25, 25)
	size2 := NewSize(10, 10)

	size3 := size1.Subtract(size2)

	assert.Equal(t, 15, size3.Width)
	assert.Equal(t, 15, size3.Height)
}

func TestSizeUnion(t *testing.T) {
	size1 := NewSize(10, 100)
	size2 := NewSize(100, 10)

	size3 := size1.Union(size2)

	assert.Equal(t, 100, size3.Width)
	assert.Equal(t, 100, size3.Height)
}
