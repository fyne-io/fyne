//go:build windows && directx

package dx

// The atlas bookkeeping - packing, glyph rasterisation, shaping and laying a
// string out as quads - is shared with the GL painter in internal/painter's
// atlas.go. This file holds what is particular to the Direct3D batch; the GPU
// half of the atlas lives in atlas_gpu.go.

// glyphBatchMax is the number of glyph quads one DrawInstanced can carry, and
// must match the gGlyphs array length in common.hlsl. A longer run simply
// flushes and starts a new batch, so this is a memory choice rather than a
// limit worth tuning: 256 instances is 12KB of constant buffer.
const glyphBatchMax = 256

// glyphInst mirrors the `GlyphInst` struct in common.hlsl, one per glyph quad.
// The two declarations must stay in the same order; TestGlyphInstMatchesShader
// pins them, because a mirror that has drifted is exactly the mistake no
// compiler will catch.
type glyphInst struct {
	NDC   [4]float32
	UV    [4]float32
	Color [4]float32
}
