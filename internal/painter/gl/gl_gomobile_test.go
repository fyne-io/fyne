//go:build (android || ios || mobile) && (!wasm || !test_web_driver)

package gl

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMobileContext_toLEByteOrder checks that the shared scratch buffer encodes
// float32s in little endian byte order, grows for a wider upload and is handed
// back to a later, narrower one without a fresh allocation.
func TestMobileContext_toLEByteOrder(t *testing.T) {
	c := &mobileContext{}

	values := []float32{1, -2.5, math.MaxFloat32}
	got := c.toLEByteOrder(values)
	assert.Len(t, got, 4*len(values))
	for i, v := range values {
		assert.Equal(t, math.Float32bits(v), binary.LittleEndian.Uint32(got[i*4:]), "value %d", i)
	}

	// A wider upload has to grow the buffer...
	grown := c.toLEByteOrder(make([]float32, 64))
	assert.Len(t, grown, 256)

	// ...and a narrower one after it reuses the same backing array.
	reused := c.toLEByteOrder([]float32{3})
	assert.Len(t, reused, 4)
	assert.Equal(t, math.Float32bits(3), binary.LittleEndian.Uint32(reused))
	assert.True(t, &grown[0] == &reused[0], "the buffer should be reused between uploads")
}
