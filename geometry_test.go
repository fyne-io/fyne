package fyne

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizeAdd(t *testing.T) {
	size1 := NewSize(10, 10)
	size2 := NewSize(25, 25)

	size3 := size1.Add(size2)

	assert.Equal(t, size3.Width, size1.Width+size2.Width)
	assert.Equal(t, size3.Height, size1.Height+size2.Height)
}

func TestSizeSubtract(t *testing.T) {
	size1 := NewSize(25, 25)
	size2 := NewSize(10, 10)

	size3 := size1.Subtract(size2)

	assert.Equal(t, size3.Width, size1.Width-size2.Width)
	assert.Equal(t, size3.Height, size1.Height-size2.Height)
}

func TestSizeUnion(t *testing.T) {
	size1 := NewSize(10, 100)
	size2 := NewSize(100, 10)

	size3 := size1.Union(size2)

	maxW := Max(size1.Width, size2.Width)
	maxH := Max(size1.Height, size2.Height)

	assert.Equal(t, size3.Width, maxW)
	assert.Equal(t, size3.Height, maxH)
}

func TestSize_FitsInto(t *testing.T) {
	tests := map[string]struct {
		width       int
		height      int
		otherWidth  int
		otherHeight int
		wantFit     bool
	}{
		"fits loose": {
			width:       6,
			height:      9,
			otherWidth:  42,
			otherHeight: 13,
			wantFit:     true,
		},
		"fits perfectly": {
			width:       17,
			height:      71,
			otherWidth:  17,
			otherHeight: 71,
			wantFit:     true,
		},
		"too wide": {
			width:       42,
			height:      9,
			otherWidth:  6,
			otherHeight: 13,
			wantFit:     false,
		},
		"too tall": {
			width:       6,
			height:      13,
			otherWidth:  42,
			otherHeight: 9,
			wantFit:     false,
		},
		"too big": {
			width:       42,
			height:      13,
			otherWidth:  6,
			otherHeight: 9,
			wantFit:     false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s1 := NewSize(tt.width, tt.height)
			s2 := NewSize(tt.otherWidth, tt.otherHeight)
			assert.EqualValues(t, tt.wantFit, s1.FitsInto(s2))
		})
	}
}

func TestPosition_Add(t *testing.T) {
	pos1 := NewPos(10, 10)
	pos2 := NewPos(25, 25)

	pos3 := pos1.Add(pos2)

	assert.Equal(t, pos3.X, pos1.X+pos2.X)
	assert.Equal(t, pos3.Y, pos1.Y+pos2.Y)
}

func TestPosition_Subtract(t *testing.T) {
	pos1 := NewPos(25, 25)
	pos2 := NewPos(10, 10)

	pos3 := pos1.Subtract(pos2)

	assert.Equal(t, pos3.X, pos1.X-pos2.X)
	assert.Equal(t, pos3.Y, pos1.Y-pos2.Y)
}
