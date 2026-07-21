package test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPixCloseEnough(t *testing.T) {
	wr, wg, wb, wa := uint8color(color.White)
	assert.True(t, pixCloseEnough([]byte{wr, wg, wb, wa}, []byte{wr, wg, wb, wa}))
	br, bg, bb, ba := uint8color(color.Black)
	assert.False(t, pixCloseEnough([]byte{wr, wg, wb, wa}, []byte{br, bg, bb, ba}))

	// Overflow case (255 vs 0)
	assert.False(t, pixCloseEnough([]byte{wr, wg, wb, wa}, []byte{wr + 1, wg - 1, wb, wa}))

	// Test small differences
	assert.True(t, pixCloseEnough([]byte{100, 100, 100, 100}, []byte{101, 100, 100, 100}))
	assert.False(t, pixCloseEnough([]byte{100, 100, 100, 100}, []byte{101, 99, 100, 100}))

	// Test a single large difference
	assert.False(t, pixCloseEnough([]byte{100, 100, 100, 100}, []byte{121, 100, 100, 100}))
	assert.False(t, pixCloseEnough([]byte{100, 100, 100, 100}, []byte{120, 100, 100, 100}))

	// Test larger array
	bigA := make([]byte, 1000)
	bigB := make([]byte, 1000)
	for i := 0; i < 100; i++ {
		bigB[i] = 1 // 10% mismatch by 1
	}
	assert.True(t, pixCloseEnough(bigA, bigB)) // diff = 100, avg = 100/1000 = 0.1

	for i := 0; i < 100; i++ {
		bigB[i] = 10 // 10% mismatch by 10
	}
	assert.False(t, pixCloseEnough(bigA, bigB)) // diff = 1000, avg = 1000/1000 = 1.0
}

func uint8color(c color.Color) (r, g, b, a uint8) {
	rgba, _ := color.RGBAModel.Convert(c).(color.RGBA)
	return rgba.R, rgba.G, rgba.B, rgba.A
}
