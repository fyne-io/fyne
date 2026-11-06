//go:build (android || ios || mobile) && (!wasm || !test_web_driver)

package gl

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestToLEByteOrder checks that float32s encode little-endian and that every
// upload gets its own immutable backing array. On the gomobile backend
// BufferData/BufferSubData only enqueue the payload pointer; the GL worker
// reads it later from the work queue (internal/driver/mobile/gl), so reusing
// one scratch buffer would let a later upload overwrite an earlier queued
// payload and corrupt vertex/UV uploads outside debug builds.
func TestToLEByteOrder(t *testing.T) {
	values := []float32{1, -2.5, math.MaxFloat32}
	got := toLEByteOrder(values...)
	assert.Len(t, got, 4*len(values))
	for i, v := range values {
		assert.Equal(t, math.Float32bits(v), binary.LittleEndian.Uint32(got[i*4:]), "value %d", i)
	}

	first := toLEByteOrder(1, 2, 3, 4)
	second := toLEByteOrder(9, 9, 9, 9)
	assert.NotSame(t, &first[0], &second[0], "each queued upload needs its own backing array")
	assert.Equal(t, math.Float32bits(1), binary.LittleEndian.Uint32(first[0:]), "first payload must stay intact after a later upload")
	assert.Equal(t, math.Float32bits(9), binary.LittleEndian.Uint32(second[0:]))
}
