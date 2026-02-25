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
	assert.False(t, pixCloseEnough([]byte{wr, wg, wb, wa}, []byte{wr + 1, wg - 1, wb, wa}))
}

func uint8color(c color.Color) (r, g, b, a uint8) {
	rr, gg, bb, aa := c.RGBA()
	return uint8(rr >> 8), uint8(gg >> 8), uint8(bb >> 8), uint8(aa >> 8)
}
