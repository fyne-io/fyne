//go:build !windows || !ci

package gl

import (
	"image"
	"image/color"
	"testing"

	"github.com/go-text/typesetting/shaping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

// addGlyph rasterises a glyph and files it, which is what the painter does in
// two steps so that it can pick the atlas based on whether the glyph has colour.
func addGlyph(a *glyphGPUAtlas, run shaping.Output, idx, phase, phases int, size, scale float32) (glyphAtlasEntry, image.Rectangle) {
	key := a.cacheKey(run, idx, phase, size, scale)
	if e, ok := a.entries[key]; ok {
		return e, image.Rectangle{}
	}
	img, baseline := paint.RenderGlyphToImage(run, idx, size, scale, float32(phase)/float32(phases))
	return a.add(key, img, baseline)
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
			phase, whole := subpixelPhaseAt(tt.in, subpixelPhases)
			assert.Equal(t, tt.phase, phase)
			assert.Equal(t, tt.wholePart, whole)
		})
	}
}

// TestSubpixelPhaseAtInRange is the property the atlas relies on: the phase is
// always a valid index and the whole part never runs ahead of the position.
func TestSubpixelPhaseAtInRange(t *testing.T) {
	for x := float32(-5); x < 5; x += 0.013 {
		phase, whole := subpixelPhaseAt(x, subpixelPhases)
		assert.GreaterOrEqual(t, phase, 0, "phase out of range at %v", x)
		assert.Less(t, phase, subpixelPhases, "phase out of range at %v", x)
		assert.LessOrEqual(t, whole, x, "whole part should not exceed the position at %v", x)
		assert.Less(t, x-whole, float32(1), "whole part should be within a pixel at %v", x)
	}
}

func TestGlyphAtlasCacheKey(t *testing.T) {
	run, idx := glyphAt(t, "A", 20, 1)
	atlas := newGlyphGPUAtlas(256)

	base := atlas.cacheKey(run, idx, 0, 20, 1)

	assert.Equal(t, base, atlas.cacheKey(run, idx, 0, 20, 1), "same inputs should give the same key")
	assert.NotEqual(t, base, atlas.cacheKey(run, idx, 1, 20, 1), "sub-pixel phase must be part of the key")
	assert.NotEqual(t, base, atlas.cacheKey(run, idx, 0, 21, 1), "font size must be part of the key")
	assert.NotEqual(t, base, atlas.cacheKey(run, idx, 0, 20, 2), "scale must be part of the key")
}

func TestGlyphAtlasGetOrAdd(t *testing.T) {
	run, idx := glyphAt(t, "A", 20, 1)
	atlas := newGlyphGPUAtlas(256)

	entry, dirty := addGlyph(atlas, run, idx, 0, subpixelPhases, 20, 1)
	assert.Positive(t, entry.w, "a newly added glyph should have width")
	assert.Positive(t, entry.h, "a newly added glyph should have height")
	assert.Positive(t, entry.baseline, "a newly added glyph should carry its baseline")
	assert.False(t, dirty.Empty(), "adding a glyph should report the region to upload")
	assert.GreaterOrEqual(t, dirty.Dx(), entry.w, "dirty region should cover the glyph")
	assert.Equal(t, entry.h, dirty.Dy(), "dirty region should cover the glyph")

	// Single byte rows are only read correctly by GL when each starts on its
	// unpack alignment, so slots are placed and sized to that boundary.
	assert.Zero(t, dirty.Min.X%glyphAtlasRowAlign, "upload should start on an aligned column")
	assert.Zero(t, dirty.Dx()%glyphAtlasRowAlign, "upload width should be a multiple of the alignment")

	again, dirtyAgain := addGlyph(atlas, run, idx, 0, subpixelPhases, 20, 1)
	assert.Equal(t, entry, again, "a cached glyph should return the same entry")
	assert.True(t, dirtyAgain.Empty(), "a cached glyph needs no upload")

	// A different sub-pixel phase is a different bitmap and must not collide.
	shifted, dirtyShifted := addGlyph(atlas, run, idx, 2, subpixelPhases, 20, 1)
	assert.False(t, dirtyShifted.Empty(), "a new phase should need uploading")
	assert.NotEqual(t, entry.x, shifted.x, "phases should occupy different atlas slots")
}

// TestGlyphAtlasAsksForResetWhenFull covers the packer running out of room. It
// must not reset itself: a caller part way through a string is holding texture
// coordinates from the current layout, and moving them under it would draw
// those glyphs from whatever now occupies those slots.
func TestGlyphAtlasAsksForResetWhenFull(t *testing.T) {
	run, idx := glyphAt(t, "A", 20, 1)
	atlas := newGlyphGPUAtlas(64)

	start := atlas.generation
	// Vary the size so every call is a fresh entry, filling the small atlas.
	var filled bool
	for i := 0; i < 200 && !filled; i++ {
		addGlyph(atlas, run, idx, 0, subpixelPhases, float32(8+i), 1)
		filled = atlas.resetPending
	}

	require.True(t, filled, "atlas should have run out of room")
	assert.Equal(t, start, atlas.generation, "a full atlas must not reset itself mid-string")
	assert.NotEmpty(t, atlas.entries, "entries packed before the atlas filled should still be addressable")

	// The caller resets between strings, once nothing depends on the layout.
	before := len(atlas.entries)
	atlas.reset()
	assert.Greater(t, atlas.generation, start, "reset should move the generation")
	assert.Empty(t, atlas.entries, "reset should empty the atlas")
	assert.False(t, atlas.resetPending, "reset should clear the request")
	assert.Positive(t, before, "sanity: the atlas held entries before the reset")
}

// TestGlyphAtlasSkipsOversizedGlyph guards a glyph too large for the atlas to
// hold. It cannot be packed at any offset, and writing it anyway would run off
// the end of the texture.
func TestGlyphAtlasSkipsOversizedGlyph(t *testing.T) {
	run, idx := glyphAt(t, "W", 40, 1)
	atlas := newGlyphGPUAtlas(8) // far smaller than any glyph at this size

	assert.NotPanics(t, func() {
		entry, dirty := addGlyph(atlas, run, idx, 0, subpixelPhases, 40, 1)
		assert.Zero(t, entry.w, "an unpackable glyph should report no size")
		assert.True(t, dirty.Empty(), "an unpackable glyph should need no upload")
	})
}

// TestGlyphAtlasColourIndependence is the property that keeps the atlas within
// capacity on a real themed UI. Before this, colour was part of the key, so the
// same text in foreground, disabled and placeholder colours stored three copies
// of every glyph, and a phone at 3x density overflowed on the second colour.
func TestGlyphAtlasColourIndependence(t *testing.T) {
	ascii := " !#$%&()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	face := paint.CachedFontFace(fyne.TextStyle{}, nil, nil)

	// A phone-like pixel density, where the old colour-keyed atlas overflowed.
	const scale = 3
	atlas := newGlyphGPUAtlas(glyphAtlasTexSize)
	paint.WalkStringGlyphs(face.Fonts, ascii, 14, fyne.TextStyle{}, scale,
		func(run shaping.Output, idx int, _, _, _, _ float32) {
			for phase := 0; phase < subpixelPhases; phase++ {
				addGlyph(atlas, run, idx, phase, subpixelPhases, 14, scale)
			}
		})

	require.False(t, atlas.resetPending, "one colour of ASCII should fit comfortably")
	resident := len(atlas.entries)

	// Drawing the same text in any number of further colours must add nothing,
	// because colour is applied when drawing rather than baked into the bitmap.
	for round := 0; round < 4; round++ {
		paint.WalkStringGlyphs(face.Fonts, ascii, 14, fyne.TextStyle{}, scale,
			func(run shaping.Output, idx int, _, _, _, _ float32) {
				for phase := 0; phase < subpixelPhases; phase++ {
					addGlyph(atlas, run, idx, phase, subpixelPhases, 14, scale)
				}
			})
	}

	assert.Equal(t, resident, len(atlas.entries), "colour must not add atlas entries")
	assert.False(t, atlas.resetPending, "repeated colours must not fill the atlas")
}

// TestGlyphCacheKey pins what cached geometry is keyed on. Keying it on the
// object rather than the words looked equivalent and was not: a widget that
// refreshes rebuilds every string it draws, so typing into an Entry rebuilt all
// the visible text on every keystroke instead of the line that changed.
func TestGlyphCacheKey(t *testing.T) {
	base := &canvas.Text{Text: "hello", TextSize: 14}
	key := glyphCacheKey(base, nil)

	same := &canvas.Text{Text: "hello", TextSize: 14}
	assert.Equal(t, key, glyphCacheKey(same, nil),
		"the same words should share geometry across objects")

	// Colour is a uniform applied when drawing, so it must not split the cache.
	coloured := &canvas.Text{Text: "hello", TextSize: 14, Color: color.White}
	assert.Equal(t, key, glyphCacheKey(coloured, nil), "colour must not affect the key")

	for name, differs := range map[string]*canvas.Text{
		"text":  {Text: "hallo", TextSize: 14},
		"size":  {Text: "hello", TextSize: 15},
		"style": {Text: "hello", TextSize: 14, TextStyle: fyne.TextStyle{Bold: true}},
	} {
		t.Run(name, func(t *testing.T) {
			assert.NotEqual(t, key, glyphCacheKey(differs, nil), "%s must be part of the key", name)
		})
	}
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
	points := p.appendGlyphQuad(nil, entry, 10, 20, p.glyphAtlas)

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
	points := p.appendGlyphQuad(nil, entry, 7, 3, p.glyphAtlas)

	for i := 0; i < len(points); i += coordinateSize2DWithTexture {
		assert.Contains(t, []float32{7, 11}, points[i], "x should span the requested position")
	}
}

func TestAppendGlyphQuadAccumulates(t *testing.T) {
	p := &painter{pixScale: 1}
	p.glyphAtlas = newGlyphGPUAtlas(64)

	entry := glyphAtlasEntry{x: 0, y: 0, w: 4, h: 4}
	points := p.appendGlyphQuad(nil, entry, 0, 0, p.glyphAtlas)
	points = p.appendGlyphQuad(points, entry, 4, 0, p.glyphAtlas)

	assert.Len(t, points, 2*floatsPerGlyph, "each glyph should add to the batch rather than replace it")
}

// TestIsColour separates glyphs that reduce to coverage from those carrying
// their own colours. Getting it wrong drew emoji as flat silhouettes, because
// only the alpha channel of a colour bitmap was kept.
func TestIsColour(t *testing.T) {
	face := paint.CachedFontFace(fyne.TextStyle{}, nil, nil)

	check := func(s string) bool {
		var colour bool
		var seen bool
		paint.WalkStringGlyphs(face.Fonts, s, 32, fyne.TextStyle{}, 1,
			func(run shaping.Output, idx int, _, _, _, _ float32) {
				img, _ := paint.RenderGlyphToImage(run, idx, 32, 1, 0)
				colour, seen = isColour(img), true
			})
		require.True(t, seen, "expected a glyph for %q", s)
		return colour
	}

	assert.False(t, check("A"), "an outline glyph is coverage only")
	assert.True(t, check("\U0001F600"), "an emoji carries its own colours")
}
