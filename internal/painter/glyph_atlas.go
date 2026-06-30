package painter

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/shaping"
)

// glyphAtlas caches pre-rendered glyph bitmaps so that the same glyph shape is
// rasterised at most once per (face, GlyphID, pixel-size, colour) combination.
//
// When compositing a string, callers retrieve the cached bitmap for each glyph
// and blit it into the destination image at the correct pen position.  HarfBuzz
// XOffset/YOffset (kerning adjustments) are excluded from the cache key and are
// applied by the caller during compositing, so the same cached bitmap is reused
// even when two occurrences of the same glyph have different kerning contexts.
//
// The atlas is a global singleton: it is shared across all painters and canvases.
// Call ResetGlyphAtlas after a font or theme change to drop stale entries.
type glyphAtlas struct {
	mu      sync.RWMutex
	entries map[glyphCacheKey]*image.RGBA
}

// glyphCacheKey uniquely identifies one pre-rendered glyph.
type glyphCacheKey struct {
	face       *font.Face // pointer identity; stable as long as fontCache is alive
	gid        font.GID
	pixSize    int32 // int32(round(fontSize * pixScale * 64)), avoids float key issues
	r, g, b, a uint8
}

// sharedGlyphAtlas is the process-wide glyph bitmap cache.
var sharedGlyphAtlas = &glyphAtlas{
	entries: make(map[glyphCacheKey]*image.RGBA),
}

// ResetGlyphAtlas discards all cached glyph bitmaps.
// Must be called whenever the font faces or display scale change.
func ResetGlyphAtlas() {
	sharedGlyphAtlas.mu.Lock()
	sharedGlyphAtlas.entries = make(map[glyphCacheKey]*image.RGBA)
	sharedGlyphAtlas.mu.Unlock()
}

func (a *glyphAtlas) cacheKey(face *font.Face, gid font.GID, fontSize, scale float32, col color.Color) glyphCacheKey {
	r32, g32, b32, a32 := col.RGBA()
	pixSize := int32(math.Round(float64(fontSize * scale * 64)))
	return glyphCacheKey{
		face:    face,
		gid:     gid,
		pixSize: pixSize,
		r:       uint8(r32 >> 8),
		g:       uint8(g32 >> 8),
		b:       uint8(b32 >> 8),
		a:       uint8(a32 >> 8),
	}
}

// glyphImage returns the cached RGBA bitmap for run.Glyphs[idx].
// On a cache miss it rasterises the glyph and stores the result.
//
// The glyph is rendered with XOffset and YOffset zeroed so the bitmap is
// independent of per-string kerning context.  Callers must apply the real
// offsets (g.XOffset, g.YOffset) when blitting into the destination image.
//
// The image height equals the full line height (ascent + |descent|); its
// baseline is at y = ascent pixels from the top, matching the coordinate
// convention used by DrawShapedRunAt.
//
// ren must not be used concurrently by another goroutine; each call to
// DrawStringOffsetAtlas creates a fresh Renderer so this invariant holds.
func (a *glyphAtlas) glyphImage(ren *render.Renderer, run shaping.Output, idx int, col color.Color, scale float32) *image.RGBA {
	g := run.Glyphs[idx]
	key := a.cacheKey(run.Face, g.GlyphID, ren.FontSize, scale, col)

	a.mu.RLock()
	if img, ok := a.entries[key]; ok {
		a.mu.RUnlock()
		return img
	}
	a.mu.RUnlock()

	a.mu.Lock()
	defer a.mu.Unlock()

	// Double-check after acquiring the write lock.
	if img, ok := a.entries[key]; ok {
		return img
	}

	ascent := int(math.Ceil(float64(fixed266ToFloat32(run.LineBounds.Ascent) * scale)))
	descent := int(math.Ceil(float64(-fixed266ToFloat32(run.LineBounds.Descent) * scale)))
	h := ascent + descent
	if h <= 0 {
		h = 1
	}

	// Width = advance + italic overspill guard (fontSize/5 at pixScale) + edge pad.
	// The italic overspill pad ensures that italic glyphs whose ink extends beyond
	// their advance metric are not clipped on the right side.
	italicPad := int(math.Ceil(float64(ren.FontSize * scale / 5)))
	w := int(math.Ceil(float64(fixed266ToFloat32(g.Advance)*scale))) + italicPad + 2
	if w <= 0 {
		w = 1
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Render with zeroed HarfBuzz offsets so the bitmap is kerning-context-free.
	noOffset := g
	noOffset.XOffset = 0
	noOffset.YOffset = 0
	singleRun := run
	singleRun.Glyphs = []shaping.Glyph{noOffset}
	ren.DrawShapedRunAt(singleRun, img, 0, ascent)

	a.entries[key] = img
	return img
}
