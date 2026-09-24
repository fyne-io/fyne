//go:build !windows || !ci

package gl

import (
	"testing"

	"github.com/go-text/typesetting/shaping"

	"fyne.io/fyne/v2"
	paint "fyne.io/fyne/v2/internal/painter"
)

const benchText = "The quick brown fox jumps over the lazy dog 0123456789"

// benchGlyphs shapes a line once so the benchmarks below measure the atlas
// rather than the shaper.
func benchGlyphs(b *testing.B, scale float32) []struct {
	run shaping.Output
	idx int
} {
	b.Helper()

	face := paint.CachedFontFace(fyne.TextStyle{}, nil, nil)
	var out []struct {
		run shaping.Output
		idx int
	}
	paint.WalkStringGlyphs(face.Fonts, benchText, 14, fyne.TextStyle{}, scale,
		func(run shaping.Output, idx int, _, _, _, _ float32) {
			out = append(out, struct {
				run shaping.Output
				idx int
			}{run, idx})
		})
	return out
}

// BenchmarkGlyphAtlasFill measures rasterising and packing a line of glyphs
// into an empty atlas, the cost paid when text is first drawn.
func BenchmarkGlyphAtlasFill(b *testing.B) {
	glyphs := benchGlyphs(b, 1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		atlas := newGlyphGPUAtlas(glyphAtlasTexSize)
		for _, g := range glyphs {
			key := atlas.cacheKey(g.run, g.idx, 0, 14, 1)
			img, baseline := paint.RenderGlyphToImage(g.run, g.idx, 14, 1, 0)
			atlas.add(key, img, baseline)
		}
	}
}

// BenchmarkGlyphAtlasHit measures looking up glyphs already in the atlas, which
// is what every repeat drawing of the same text costs.
func BenchmarkGlyphAtlasHit(b *testing.B) {
	glyphs := benchGlyphs(b, 1)
	atlas := newGlyphGPUAtlas(glyphAtlasTexSize)
	for _, g := range glyphs {
		key := atlas.cacheKey(g.run, g.idx, 0, 14, 1)
		img, baseline := paint.RenderGlyphToImage(g.run, g.idx, 14, 1, 0)
		atlas.add(key, img, baseline)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, g := range glyphs {
			key := atlas.cacheKey(g.run, g.idx, 0, 14, 1)
			if _, ok := atlas.entries[key]; !ok {
				b.Fatal("expected a cached glyph")
			}
		}
	}
}

// BenchmarkAppendGlyphQuad measures building the vertices for a line of text.
func BenchmarkAppendGlyphQuad(b *testing.B) {
	p := &painter{pixScale: 1}
	p.glyphAtlas = newGlyphGPUAtlas(glyphAtlasTexSize)
	entry := glyphAtlasEntry{x: 4, y: 8, w: 9, h: 19}

	var points []float32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		points = points[:0]
		for x := 0; x < 50; x++ {
			points = p.appendGlyphQuad(points, entry, float32(x*9), 0, p.glyphAtlas)
		}
	}
}
