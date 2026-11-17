package gl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	paint "fyne.io/fyne/v2/internal/painter"
)

// TestPushQuad pins the batch vertex layout the glyph shader reads: position,
// then texture coordinate, then colour, six vertices covering the quad.
func TestPushQuad(t *testing.T) {
	p := &painter{}
	q := paint.GlyphQuad{X1: 0, Y1: 0, X2: 50, Y2: 25, U1: 0.1, V1: 0.2, U2: 0.3, V2: 0.4}
	col := [4]float32{1, 0.5, 0.25, 0.8}
	p.pushQuad(q, col, 100, 50)
	require.Len(t, p.glyphPending, 6*glyphVertexFloats)

	// The top-left pixel corner is clip space (-1, 1); the quad spans half of
	// each axis, so its far corner lands on the origin.
	corners := [][4]float32{
		{-1, 1, 0.1, 0.2}, {0, 1, 0.3, 0.2}, {-1, 0, 0.1, 0.4}, // first triangle
		{-1, 0, 0.1, 0.4}, {0, 1, 0.3, 0.2}, {0, 0, 0.3, 0.4}, // second
	}
	for i, want := range corners {
		v := p.glyphPending[i*glyphVertexFloats : (i+1)*glyphVertexFloats]
		assert.Equal(t, want[:], v[:4], "vertex %d position and texture coordinate", i)
		assert.Equal(t, col[:], v[4:], "vertex %d colour", i)
	}

	p.glyphPending = p.glyphPending[:0]
	p.FlushGlyphs() // an empty queue draws nothing, and needs no context to know it
}
