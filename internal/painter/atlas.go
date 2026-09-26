package painter

import (
	"image"
	"image/color"
	"math"

	"github.com/go-text/typesetting/font"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// This file is the renderer-neutral half of the glyph atlas the GPU painters
// draw text through. It decides where glyphs live, rasterises them and lays
// strings out as one quad per glyph; each painter owns the texture itself and
// turns the quads into its own vertex format.
//
// Caching whole rasterised strings costs one texture per distinct string, which
// for a UI whose labels change - a dashboard, a log player, anything counting -
// is unbounded. Caching glyphs instead bounds the whole thing at one texture,
// because an app draws from a fixed set of characters however much its text
// changes. One texture is also what lets consecutive text draw as one batch.

const (
	// AtlasSize is the side of the square texture every glyph is packed into.
	// It holds several thousand glyphs at UI sizes, which is every character of
	// every font a normal app uses several times over.
	AtlasSize = 1024
	// atlasPad separates neighbouring glyphs so linear sampling at the edge of
	// one cannot pick up the next.
	atlasPad = 1
	// shapedCacheMax caps the shaped-run cache. A UI whose labels churn
	// (counters, live values) grows entries without bound; past the cap the
	// whole cache is dropped and rebuilt from live strings, the same
	// reset-and-repack answer the atlas uses. 4096 entries of a dozen glyphs is
	// roughly 2MB.
	shapedCacheMax = 4096
	// whiteBlock is the side of the solid block plain rectangles sample, and
	// whiteCentre the offset of its centre texel from the block's corner.
	whiteBlock  = 3
	whiteCentre = whiteBlock / 2.0
)

// AtlasTexture is a painter's GPU half of a GlyphAtlas.
type AtlasTexture interface {
	// FlushGlyphs draws every quad queued so far. The atlas calls it before it
	// reuses space, while the queued quads still point at the old contents.
	FlushGlyphs()
	// UploadGlyph copies a rasterised bitmap into the texture with its top-left
	// at x, y. Only the alpha channel carries anything: glyphs are coverage, and
	// their colour comes with each quad.
	UploadGlyph(img *image.RGBA, x, y int)
}

// GlyphQuad is one glyph to draw: its rectangle in destination pixels and the
// matching rectangle of the atlas in texture coordinates.
type GlyphQuad struct {
	X1, Y1, X2, Y2 float32
	U1, V1, U2, V2 float32
}

// GlyphAtlas tracks what is where in a painter's atlas texture. The zero value
// is an empty atlas.
type GlyphAtlas struct {
	packer  atlasPacker
	entries map[glyphKey]glyphEntry
	// shaped caches the result of shaping one string, keyed by everything the
	// shaper reads. Shaping (segmenting plus HarfBuzz) is the dominant cost of
	// text once rasterisation and draws are cached and batched, and a UI's
	// strings mostly repeat frame over frame.
	shaped map[shapedKey]shapedEntry
	// resets counts atlas resets, which is how a string being resolved notices
	// that the slots it already has went stale.
	resets int

	// whiteX and whiteY locate a solid 3x3 block, which lets a plain filled
	// rectangle ride the glyph batch as a quad with a constant texture
	// coordinate: coverage samples 1 and the quad colour is the fill. whiteOK
	// is false until the block is packed, and again after every reset.
	whiteX, whiteY int
	whiteOK        bool

	// glyphScratch and slotScratch collect one string's shaped glyphs and
	// their atlas slots, reused from string to string.
	glyphScratch []PlacedGlyph
	slotScratch  []glyphEntry
}

// TextQuads appends a quad per inked glyph of text, whose top-left is at pos, to
// dst. It reports false with dst unchanged when the string has to take the
// whole-run texture path instead: a character the shaper found no glyph for,
// or a glyph a coverage atlas cannot hold, such as a colour emoji.
func (a *GlyphAtlas) TextQuads(dst []GlyphQuad, text *canvas.Text, pos fyne.Position, pixScale float32, tex AtlasTexture) ([]GlyphQuad, bool) {
	glyphs, ok := a.shape(text, pixScale)
	if !ok {
		return dst, false
	}

	// Every glyph is resolved before any quad is handed out. Resolving can reset
	// a full atlas, which strands the slots resolved before it, so a reset means
	// resolving the string again; only a string whose glyphs alone overflow the
	// atlas resets twice.
	for attempt := 0; ; attempt++ {
		resets := a.resets
		a.slotScratch = a.slotScratch[:0]
		for _, pg := range glyphs {
			e, ok := a.glyph(pg, text.TextSize, pixScale, tex)
			if !ok {
				return dst, false
			}
			a.slotScratch = append(a.slotScratch, e)
		}
		if a.resets == resets {
			break
		}
		if attempt > 0 {
			return dst, false
		}
	}

	// The run's pixel origin, rounded the way the whole-run path rounds its quad.
	originX := float32(math.Round(float64(pos.X * pixScale)))
	originY := float32(math.Round(float64(pos.Y * pixScale)))
	for i, pg := range glyphs {
		e := a.slotScratch[i]
		if e.empty() {
			continue
		}
		// Snapped to whole device pixels. The bitmap was rasterised at a whole
		// pixel origin, so a quad at a fractional position would resample it and
		// the glyph would come out soft and fringed - the pen position drifts
		// fractional as the shaped advances accumulate, so this is every glyph
		// but the first, not an edge case.
		//
		// ponytail: whole-pixel placement, so within-run subpixel positioning is
		// lost and spacing can differ from the whole-run path by under a pixel.
		// The upgrade, if that is ever visible, is to cache each glyph at a few
		// horizontal subpixel phases and pick by the fractional pen position.
		x1 := float32(math.Round(float64(originX+pg.X))) + float32(e.bearX)
		y1 := float32(math.Round(float64(originY+pg.Y))) + float32(e.bearY)
		dst = append(dst, GlyphQuad{
			X1: x1, Y1: y1, X2: x1 + float32(e.width), Y2: y1 + float32(e.height),
			U1: float32(e.x) / AtlasSize, V1: float32(e.y) / AtlasSize,
			U2: float32(e.x+e.width) / AtlasSize, V2: float32(e.y+e.height) / AtlasSize,
		})
	}
	return dst, true
}

// White returns the texture coordinate at the centre of the solid block plain
// rectangles sample, packing it first if need be. Sampling the centre of a 3x3
// keeps every bilinear neighbour inside the block, so the empty pad the packer
// leaves around it can never bleed in.
func (a *GlyphAtlas) White(tex AtlasTexture) (u, v float32, ok bool) {
	if !a.whiteOK {
		x, y, ok := a.place(whiteBlock, whiteBlock, tex)
		if !ok {
			return 0, 0, false
		}
		img := image.NewRGBA(image.Rect(0, 0, whiteBlock, whiteBlock))
		for i := range img.Pix {
			img.Pix[i] = 0xff
		}
		tex.UploadGlyph(img, x, y)
		a.whiteX, a.whiteY, a.whiteOK = x, y, true
	}
	return (float32(a.whiteX) + whiteCentre) / AtlasSize, (float32(a.whiteY) + whiteCentre) / AtlasSize, true
}

// shape returns text's glyphs as WalkGlyphs places them, from the cache when the
// string has been shaped before.
func (a *GlyphAtlas) shape(text *canvas.Text, pixScale float32) ([]PlacedGlyph, bool) {
	face := CachedFontFace(text.TextStyle, text.FontSource, text)
	key := shapedKey{text: text.Text, face: face, size: text.TextSize, scale: pixScale, style: text.TextStyle}
	if e, ok := a.shaped[key]; ok {
		return e.glyphs, e.ok
	}

	a.glyphScratch = a.glyphScratch[:0]
	WalkGlyphs(face.Fonts, text.Text, text.TextSize, pixScale, text.TextStyle,
		func(g PlacedGlyph) { a.glyphScratch = append(a.glyphScratch, g) })
	var e shapedEntry
	e.ok = len(a.glyphScratch) > 0 && shapeable(a.glyphScratch)
	if e.ok {
		e.glyphs = append([]PlacedGlyph(nil), a.glyphScratch...)
	}
	if a.shaped == nil {
		a.shaped = make(map[shapedKey]shapedEntry)
	} else if len(a.shaped) >= shapedCacheMax {
		clear(a.shaped)
	}
	a.shaped[key] = e
	return e.glyphs, e.ok
}

// glyph returns where a glyph lives in the atlas, rasterising and uploading it
// on first use. The second result is false when the glyph could not be placed
// at all, which sends the caller to the whole-run fallback.
func (a *GlyphAtlas) glyph(pg PlacedGlyph, fontSize, pixScale float32, tex AtlasTexture) (glyphEntry, bool) {
	key := newGlyphKey(pg.Face, pg.Glyph.GlyphID, fontSize, pixScale)
	if e, ok := a.entries[key]; ok {
		return e, !e.unsupported
	}
	if a.entries == nil {
		a.entries = make(map[glyphKey]glyphEntry)
	}

	// Colour glyphs - emoji, in practice - carry bitmap or SVG data a coverage
	// atlas cannot represent. The answer is cached along with everything else,
	// because GlyphData is far too expensive to ask once a frame.
	if _, ok := pg.Face.GlyphData(pg.Glyph.GlyphID).(font.GlyphOutline); !ok {
		a.entries[key] = glyphEntry{unsupported: true}
		return glyphEntry{}, false
	}

	// Ink bounds in destination pixels, relative to the pen position and
	// baseline. Height is measured downwards from YBearing, so it is negative.
	relX := fixed266ToFloat32(pg.Glyph.XBearing) * pixScale
	relY := fixed266ToFloat32(pg.Glyph.YBearing) * pixScale
	inkW := fixed266ToFloat32(pg.Glyph.Width) * pixScale
	inkH := -fixed266ToFloat32(pg.Glyph.Height) * pixScale
	if inkW <= 0 || inkH <= 0 { // a space, or anything else with no ink
		a.entries[key] = glyphEntry{}
		return glyphEntry{}, true
	}

	// Rasterise at a whole-pixel origin chosen so the ink, plus a pixel of room
	// for the antialiased edge, lands inside the bitmap.
	originX := 1 - int(math.Floor(float64(relX)))
	originY := 1 + int(math.Ceil(float64(relY)))
	w := int(math.Ceil(float64(inkW))) + 3
	h := int(math.Ceil(float64(inkH))) + 3

	x, y, ok := a.place(w, h, tex)
	if !ok {
		return glyphEntry{}, false // larger than the atlas itself
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	RasteriseGlyph(img, pg, color.White, fontSize, pixScale, originX, originY)
	tex.UploadGlyph(img, x, y)

	e := glyphEntry{x: x, y: y, width: w, height: h, bearX: -originX, bearY: -originY}
	a.entries[key] = e
	return e, true
}

// place reserves w by h texels, emptying a full atlas first. Quads already
// queued point at the old contents, so they are drawn before anything is
// overwritten.
func (a *GlyphAtlas) place(w, h int, tex AtlasTexture) (x, y int, ok bool) {
	if x, y, ok = a.packer.add(w, h); ok {
		return x, y, true
	}
	tex.FlushGlyphs()
	a.packer.reset()
	clear(a.entries)
	a.whiteOK = false
	a.resets++
	return a.packer.add(w, h)
}

// shapedKey identifies one shaped string: the text and everything else the
// shaper reads. The face is the cached *FontCacheItem pointer, whose identity
// is stable until the font caches are cleared - a theme or font change hands
// out new pointers, orphaning (not corrupting) old entries.
type shapedKey struct {
	text  string
	face  *FontCacheItem
	size  float32
	scale float32
	style fyne.TextStyle
}

// shapedEntry is one cached shaping result. ok is false when the string has to
// take the whole-run texture path (an emoji, in practice), cached so the
// losing shape is not re-run every frame just to fail again.
type shapedEntry struct {
	glyphs []PlacedGlyph
	ok     bool
}

// shapeable reports whether a shaped string is a candidate for the atlas at
// all. Glyph 0 means the shaper found nothing, which DrawString renders by
// substituting a replacement character from another face - a second shaping
// path this deliberately does not grow. Whether the glyphs themselves fit is
// the atlas's answer, and a cached one.
//
// ponytail: whole-string fallback rather than a second colour atlas beside this
// one. Mixed emoji and text then costs one texture per such label, which is
// what every label cost before this existed; a colour atlas is only worth
// building for an app that is mostly emoji.
func shapeable(glyphs []PlacedGlyph) bool {
	for _, pg := range glyphs {
		if pg.Glyph.GlyphID == 0 || pg.Face == nil {
			return false
		}
	}
	return true
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
	// once per glyph rather than once per glyph per frame.
	unsupported bool
}

// empty reports a glyph with no ink - a space, most often - which needs no quad.
func (e glyphEntry) empty() bool { return e.width == 0 || e.height == 0 }

// atlasPacker fills the atlas with the shelf algorithm: rectangles are laid
// left to right in a row whose height is set by the tallest one placed in it,
// and a full row starts a new one below.
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
	if w <= 0 || h <= 0 || slotW > AtlasSize || slotH > AtlasSize {
		return 0, 0, false
	}
	if p.penX+slotW > AtlasSize { // shelf full, open the next one
		p.shelfY += p.shelfHeight
		p.shelfHeight = 0
		p.penX = 0
	}
	if slotH > p.shelfHeight {
		// A taller glyph grows the shelf rather than opening a new one, which
		// would waste everything already placed on this one. Glyphs already on
		// the shelf start at its top, so growing it downwards cannot disturb them.
		if p.shelfY+slotH > AtlasSize {
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
