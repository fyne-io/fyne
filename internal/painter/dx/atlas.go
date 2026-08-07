//go:build windows && directx

package dx

import "github.com/go-text/typesetting/font"

// This file holds the parts of the glyph atlas that need no Direct3D; the GPU
// half lives in atlas_gpu.go.
//
// Caching whole rasterised strings costs one texture per distinct string, which
// for a UI whose labels change - a dashboard, a log player, anything counting -
// is unbounded. Caching glyphs instead bounds the whole thing at one texture,
// because an app draws from a fixed set of characters however much its text
// changes.

const (
	// atlasSize is the side of the square coverage texture every glyph is packed
	// into. 1024x1024 at one byte per pixel is 1MB and holds several thousand
	// glyphs at UI sizes, which is every character of every font a normal app
	// uses several times over.
	atlasSize = 1024
	// atlasPad separates neighbouring glyphs so linear sampling at the edge of
	// one cannot pick up the next.
	atlasPad = 1
	// glyphBatchMax is the number of glyph quads one DrawInstanced can carry,
	// and must match the gGlyphs array length in common.hlsl. A longer run
	// simply flushes and starts a new batch, so this is a memory choice rather
	// than a limit worth tuning: 256 instances is 12KB of constant buffer.
	glyphBatchMax = 256
)

// glyphInst mirrors the `GlyphInst` struct in common.hlsl, one per glyph quad.
// The two declarations must stay in the same order; TestGlyphInstMatchesShader
// pins them, because a mirror that has drifted is exactly the mistake no
// compiler will catch.
type glyphInst struct {
	NDC   [4]float32
	UV    [4]float32
	Color [4]float32
}

// glyphKey identifies one rasterised glyph. The size folds in the canvas scale
// because that is what the rasteriser is given, and it is fixed point so two
// sizes a hair apart cannot collide through float rounding.
type glyphKey struct {
	face *font.Face
	gid  font.GID
	size int32 // fontSize * pixScale in 26.6 fixed point
}

func newGlyphKey(face *font.Face, gid font.GID, fontSize, pixScale float32) glyphKey {
	return glyphKey{face: face, gid: gid, size: int32(fontSize * pixScale * 64)}
}

// glyphEntry is where a glyph landed in the atlas and how its bitmap sits
// relative to the pen position and baseline it was rasterised from.
type glyphEntry struct {
	x, y          int // top-left in atlas pixels
	width, height int
	// bearX and bearY offset the quad from the pen position and baseline. They
	// are whole pixels because the bitmap was rasterised at a whole pixel origin.
	bearX, bearY int
	// unsupported marks a glyph a coverage atlas cannot hold - a colour emoji,
	// in practice. It is cached like any other entry so the decision is made
	// once per glyph rather than once per glyph per frame: finding out costs a
	// GlyphData lookup, far too expensive to repeat every frame.
	unsupported bool
}

// empty reports a glyph with no ink - a space, most often - which needs no quad.
func (e glyphEntry) empty() bool { return e.width == 0 || e.height == 0 }

// atlasPacker fills the atlas with the shelf algorithm: rectangles are laid
// left to right in a row whose height is set by the tallest one placed in it,
// and a full row starts a new one above.
//
// ponytail: shelves, not a skyline or a full bin packer. Glyphs of one font
// size are all much the same height, so a shelf wastes very little, and the
// whole structure is three ints. If mixed body and heading sizes ever fragment
// this badly enough to matter, the atlas simply resets and repacks - the
// upgrade path is a skyline packer, not a rewrite of the callers.
type atlasPacker struct {
	penX        int // left edge of the next free slot on the current shelf
	shelfY      int // top of the current shelf
	shelfHeight int
}

// add reserves a w by h slot, reporting where it went. It fails only when the
// atlas is full, which the caller answers by resetting and starting again.
//
// The reservation is a pixel wider and taller than asked for, so neighbours are
// always separated whether they sit side by side or on adjacent shelves.
func (p *atlasPacker) add(w, h int) (x, y int, ok bool) {
	slotW, slotH := w+atlasPad, h+atlasPad
	if w <= 0 || h <= 0 || slotW > atlasSize || slotH > atlasSize {
		return 0, 0, false
	}
	if p.penX+slotW > atlasSize { // shelf full, open the next one
		p.shelfY += p.shelfHeight
		p.shelfHeight = 0
		p.penX = 0
	}
	if slotH > p.shelfHeight {
		// A taller glyph grows the shelf rather than opening a new one, which
		// would waste everything already placed on this one. Glyphs already on
		// the shelf start at its top, so growing it downwards cannot disturb them.
		if p.shelfY+slotH > atlasSize {
			return 0, 0, false
		}
		p.shelfHeight = slotH
	}
	x, y = p.penX, p.shelfY
	p.penX += slotW
	return x, y, true
}

// reset empties the atlas. Callers must draw anything already queued against
// the old contents first, and forget every glyphEntry they hold.
func (p *atlasPacker) reset() {
	p.penX, p.shelfY, p.shelfHeight = 0, 0, 0
}
