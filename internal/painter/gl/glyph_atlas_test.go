//go:build !windows || !ci

package gl

import (
	"image/color"
	"testing"

	"github.com/go-text/typesetting/shaping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	paint "fyne.io/fyne/v2/internal/painter"
)

// glyphAt shapes s and returns the run and index of its first glyph, giving
// tests a real shaped glyph to rasterise rather than a synthetic one.
func glyphAt(t *testing.T, s string, size, scale float32) (shaping.Output, int) {
	t.Helper()

	face := paint.CachedFontFace(fyne.TextStyle{}, nil, nil)
	require.NotNil(t, face)

	var run shaping.Output
	var idx int
	found := false
	paint.WalkStringGlyphs(face.Fonts, s, size, fyne.TextStyle{}, scale,
		func(r shaping.Output, i int, _, _, _, _ float32) {
			if !found {
				run, idx, found = r, i, true
			}
		})
	require.True(t, found, "expected at least one glyph for %q", s)
	return run, idx
}

func TestSubpixelPhaseAt(t *testing.T) {
	for name, tt := range map[string]struct {
		in        float32
		phase     int
		wholePart float32
	}{
		"whole pixel":      {10, 0, 10},
		"first quarter":    {10.1, 0, 10},
		"second quarter":   {10.3, 1, 10},
		"third quarter":    {10.5, 2, 10},
		"fourth quarter":   {10.8, 3, 10},
		"just under whole": {10.999, 3, 10},
		"zero":             {0, 0, 0},
		"negative whole":   {-3, 0, -3},
		"negative frac":    {-2.5, 2, -3},
	} {
		t.Run(name, func(t *testing.T) {
			phase, whole := subpixelPhaseAt(tt.in)
			assert.Equal(t, tt.phase, phase)
			assert.Equal(t, tt.wholePart, whole)
		})
	}
}

// TestSubpixelPhaseAtInRange is the property the atlas relies on: the phase is
// always a valid index and the whole part never runs ahead of the position.
func TestSubpixelPhaseAtInRange(t *testing.T) {
	for x := float32(-5); x < 5; x += 0.013 {
		phase, whole := subpixelPhaseAt(x)
		assert.GreaterOrEqual(t, phase, 0, "phase out of range at %v", x)
		assert.Less(t, phase, subpixelPhases, "phase out of range at %v", x)
		assert.LessOrEqual(t, whole, x, "whole part should not exceed the position at %v", x)
		assert.Less(t, x-whole, float32(1), "whole part should be within a pixel at %v", x)
	}
}

func TestGlyphAtlasCacheKey(t *testing.T) {
	run, idx := glyphAt(t, "A", 20, 1)
	atlas := newGlyphGPUAtlas(256)

	base := atlas.cacheKey(run, idx, 0, 20, 1, color.White)

	assert.Equal(t, base, atlas.cacheKey(run, idx, 0, 20, 1, color.White), "same inputs should give the same key")
	assert.NotEqual(t, base, atlas.cacheKey(run, idx, 1, 20, 1, color.White), "sub-pixel phase must be part of the key")
	assert.NotEqual(t, base, atlas.cacheKey(run, idx, 0, 21, 1, color.White), "font size must be part of the key")
	assert.NotEqual(t, base, atlas.cacheKey(run, idx, 0, 20, 2, color.White), "scale must be part of the key")
	assert.NotEqual(t, base, atlas.cacheKey(run, idx, 0, 20, 1, color.Black), "colour must be part of the key")
}

func TestGlyphAtlasGetOrAdd(t *testing.T) {
	run, idx := glyphAt(t, "A", 20, 1)
	atlas := newGlyphGPUAtlas(256)

	entry, dirty := atlas.getOrAdd(run, idx, 0, 20, 1, color.White)
	assert.Positive(t, entry.w, "a newly added glyph should have width")
	assert.Positive(t, entry.h, "a newly added glyph should have height")
	assert.Positive(t, entry.baseline, "a newly added glyph should carry its baseline")
	assert.False(t, dirty.Empty(), "adding a glyph should report the region to upload")
	assert.Equal(t, entry.w, dirty.Dx(), "dirty region should cover the glyph")
	assert.Equal(t, entry.h, dirty.Dy(), "dirty region should cover the glyph")

	again, dirtyAgain := atlas.getOrAdd(run, idx, 0, 20, 1, color.White)
	assert.Equal(t, entry, again, "a cached glyph should return the same entry")
	assert.True(t, dirtyAgain.Empty(), "a cached glyph needs no upload")

	// A different sub-pixel phase is a different bitmap and must not collide.
	shifted, dirtyShifted := atlas.getOrAdd(run, idx, 2, 20, 1, color.White)
	assert.False(t, dirtyShifted.Empty(), "a new phase should need uploading")
	assert.NotEqual(t, entry.x, shifted.x, "phases should occupy different atlas slots")
}

// TestGlyphAtlasResetsWhenFull covers the packer running out of room: entries
// are dropped and the generation moves, which is how callers know the texture
// coordinates they are holding have gone stale.
func TestGlyphAtlasResetsWhenFull(t *testing.T) {
	run, idx := glyphAt(t, "A", 20, 1)
	atlas := newGlyphGPUAtlas(64)

	start := atlas.generation
	// Vary the colour so every call is a fresh entry, filling the small atlas.
	for i := 0; i < 200 && atlas.generation == start; i++ {
		col := color.NRGBA{R: uint8(i), G: uint8(i * 3), B: uint8(i * 7), A: 0xff}
		atlas.getOrAdd(run, idx, 0, 20, 1, col)
	}

	require.Greater(t, atlas.generation, start, "atlas should have reset once full")
	assert.NotEmpty(t, atlas.entries, "the glyph that triggered the reset should be in the fresh atlas")
	assert.Less(t, len(atlas.entries), 200, "a reset should have dropped earlier entries")
}

// TestGlyphAtlasSkipsOversizedGlyph guards a glyph too large for the atlas to
// hold. It cannot be packed at any offset, and writing it anyway would run off
// the end of the texture.
func TestGlyphAtlasSkipsOversizedGlyph(t *testing.T) {
	run, idx := glyphAt(t, "W", 40, 1)
	atlas := newGlyphGPUAtlas(8) // far smaller than any glyph at this size

	assert.NotPanics(t, func() {
		entry, dirty := atlas.getOrAdd(run, idx, 0, 40, 1, color.White)
		assert.Zero(t, entry.w, "an unpackable glyph should report no size")
		assert.True(t, dirty.Empty(), "an unpackable glyph should need no upload")
	})
}

// TestTextVerticesUsable covers when cached glyph geometry may be reused. Both
// negative cases would draw the wrong thing rather than fail loudly: stale
// texture coordinates point at whatever now occupies that part of the atlas,
// and stale scale draws glyphs rasterised for a different pixel density.
func TestTextVerticesUsable(t *testing.T) {
	cached := &textVertices{generation: 3, pixScale: 2}

	assert.True(t, cached.usable(3, 2), "unchanged atlas and scale should be reusable")
	assert.False(t, cached.usable(4, 2), "an atlas reset should invalidate the geometry")
	assert.False(t, cached.usable(3, 1), "a scale change should invalidate the geometry")

	var missing *textVertices
	assert.False(t, missing.usable(3, 2), "absent geometry is never usable")
}

func TestAppendGlyphQuad(t *testing.T) {
	p := &painter{pixScale: 1}
	p.glyphAtlas = newGlyphGPUAtlas(64)

	entry := glyphAtlasEntry{x: 2, y: 4, w: 6, h: 8}
	points := p.appendGlyphQuad(nil, entry, 10, 20)

	require.Len(t, points, floatsPerGlyph, "a glyph should emit two triangles")

	// Position spans offX..offX+w and offY..offY+h, in device pixels relative
	// to the string origin.
	xs, ys := map[float32]bool{}, map[float32]bool{}
	us, vs := map[float32]bool{}, map[float32]bool{}
	for i := 0; i < len(points); i += coordinateSize2DWithTexture {
		xs[points[i]] = true
		ys[points[i+1]] = true
		us[points[i+2]] = true
		vs[points[i+3]] = true
	}
	assert.Equal(t, map[float32]bool{10: true, 16: true}, xs)
	assert.Equal(t, map[float32]bool{20: true, 28: true}, ys)

	// Texture coordinates address the glyph's slot within the atlas.
	assert.Equal(t, map[float32]bool{2.0 / 64: true, 8.0 / 64: true}, us)
	assert.Equal(t, map[float32]bool{4.0 / 64: true, 12.0 / 64: true}, vs)
}

// TestAppendGlyphQuadKeepsOrigin checks that the quad is placed exactly where
// asked, without re-rounding: the caller has already split the position into the
// whole pixel drawn on and the fraction rasterised into the bitmap.
func TestAppendGlyphQuadKeepsOrigin(t *testing.T) {
	p := &painter{pixScale: 1}
	p.glyphAtlas = newGlyphGPUAtlas(64)

	entry := glyphAtlasEntry{x: 0, y: 0, w: 4, h: 4}
	points := p.appendGlyphQuad(nil, entry, 7, 3)

	for i := 0; i < len(points); i += coordinateSize2DWithTexture {
		assert.Contains(t, []float32{7, 11}, points[i], "x should span the requested position")
	}
}

func TestAppendGlyphQuadAccumulates(t *testing.T) {
	p := &painter{pixScale: 1}
	p.glyphAtlas = newGlyphGPUAtlas(64)

	entry := glyphAtlasEntry{x: 0, y: 0, w: 4, h: 4}
	points := p.appendGlyphQuad(nil, entry, 0, 0)
	points = p.appendGlyphQuad(points, entry, 4, 0)

	assert.Len(t, points, 2*floatsPerGlyph, "each glyph should add to the batch rather than replace it")
}
