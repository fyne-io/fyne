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
	// AtlasSize is the side of the square texture every glyph is packed into
	// at first. It holds several thousand glyphs at UI sizes, which is every
	// character of every font a normal app uses several times over.
	AtlasSize = 1024
	// AtlasMaxSize is the side the atlas grows to the first time it fills. Past
	// 3x a page mixing body text, bold, italic, code and headings no longer fits
	// in AtlasSize, and a working set larger than the atlas would reset it every
	// frame, rasterising and uploading all of it again.
	AtlasMaxSize = 2 * AtlasSize
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
	// subpixelUnits is the precision a glyph's horizontal sub-pixel offset is
	// keyed at: 64ths of a pixel, the 26.6 fixed point the shaper works in.
	subpixelUnits = 64
	// subpixelPerPoint is how many horizontal sub-pixel positions a glyph is
	// cached at per logical pixel; see subpixelSteps.
	subpixelPerPoint = 4
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
	// ResizeAtlas replaces the texture with an empty one size texels square.
	// Everything queued against the old one has been flushed first.
	ResizeAtlas(size int)
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
// dst. Glyphs wholly outside the device pixel columns minX to maxX are left
// out, so a string far wider than the screen costs what is visible of it. It
// reports false with dst unchanged when the string has to take the whole-run
// texture path instead: a character the shaper found no glyph for, or a glyph
// a coverage atlas cannot hold, such as a colour emoji.
func (a *GlyphAtlas) TextQuads(dst []GlyphQuad, text *canvas.Text, pos fyne.Position, pixScale, minX, maxX float32, tex AtlasTexture) ([]GlyphQuad, bool) {
	glyphs, ok := a.shape(text, pixScale)
	if !ok {
		return dst, false
	}

	// The run's pixel origin, rounded the way the whole-run path rounds its quad.
	originX := float32(math.Round(float64(pos.X * pixScale)))
	originY := float32(math.Round(float64(pos.Y * pixScale)))
	steps := subpixelSteps(pixScale)

	// Every glyph is resolved before any quad is handed out. Resolving can reset
	// a full atlas, which strands the slots resolved before it, so a reset means
	// resolving the string again; only a string whose glyphs alone overflow the
	// atlas resets twice.
	for attempt := 0; ; attempt++ {
		resets := a.resets
		a.slotScratch = a.slotScratch[:0]
		for _, pg := range glyphs {
			// A compare per glyph per frame, even for a 100k glyph line. Pen
			// positions could be bisected instead if that ever shows up.
			if !pg.inColumns(originX+pg.X, pixScale, minX, maxX) {
				a.slotScratch = append(a.slotScratch, glyphEntry{}) // no quad
				continue
			}
			_, sub := snapX(originX+pg.X, steps)
			e, ok := a.glyph(pg, sub, text.TextSize, pixScale, tex)
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

	size := float32(a.Size())
	for i, pg := range glyphs {
		e := a.slotScratch[i]
		if e.empty() {
			continue
		}
		// Quads land on whole device pixels, one texel to one pixel: a quad at a
		// fractional position would resample its bitmap and the glyph would come
		// out soft. The pen position's fraction is in the bitmap instead, which
		// was rasterised that far right of its origin, so spacing follows the
		// shaped advances as the whole-run path draws them. Vertically the
		// baseline is already whole, as it is for the whole-run path.
		ix, _ := snapX(originX+pg.X, steps)
		x1 := float32(ix + e.bearX)
		y1 := float32(math.Round(float64(originY+pg.Y))) + float32(e.bearY)
		dst = append(dst, GlyphQuad{
			X1: x1, Y1: y1, X2: x1 + float32(e.width), Y2: y1 + float32(e.height),
			U1: float32(e.x) / size, V1: float32(e.y) / size,
			U2: float32(e.x+e.width) / size, V2: float32(e.y+e.height) / size,
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
	size := float32(a.Size())
	return (float32(a.whiteX) + whiteCentre) / size, (float32(a.whiteY) + whiteCentre) / size, true
}

// Size is the side of the square texture the atlas packs into, which a
// painter allocates before anything is uploaded.
func (a *GlyphAtlas) Size() int {
	return a.packer.side()
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

// glyph returns where a glyph lives in the atlas when drawn sub/64 of a pixel
// right of a whole pixel, rasterising and uploading it on first use. The second
// result is false when the glyph could not be placed at all, which sends the
// caller to the whole-run fallback.
func (a *GlyphAtlas) glyph(pg PlacedGlyph, sub int, fontSize, pixScale float32, tex AtlasTexture) (glyphEntry, bool) {
	key := newGlyphKey(pg.Face, pg.Glyph.GlyphID, fontSize, pixScale, sub)
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

	// Rasterise at an origin chosen so the ink, shifted right by the sub-pixel
	// offset and with a pixel of room for the antialiased edge, lands inside the
	// bitmap. The offset is under a pixel and folded into the floor, so the
	// width below still has room for it.
	subX := float32(sub) / subpixelUnits
	originX := 1 - int(math.Floor(float64(relX+subX)))
	originY := 1 + int(math.Ceil(float64(relY)))
	w := int(math.Ceil(float64(inkW))) + 3
	h := int(math.Ceil(float64(inkH))) + 3

	x, y, ok := a.place(w, h, tex)
	if !ok {
		return glyphEntry{}, false // larger than the atlas itself
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	RasteriseGlyph(img, pg, color.White, fontSize, pixScale, originX, originY, subX)
	tex.UploadGlyph(img, x, y)

	e := glyphEntry{x: x, y: y, width: w, height: h, bearX: -originX, bearY: -originY}
	a.entries[key] = e
	return e, true
}

// place reserves w by h texels, emptying a full atlas first, and growing it
// if it has not grown yet. Quads already queued point at the old contents, so
// they are drawn before anything is overwritten.
func (a *GlyphAtlas) place(w, h int, tex AtlasTexture) (x, y int, ok bool) {
	if x, y, ok = a.packer.add(w, h); ok {
		return x, y, true
	}
	tex.FlushGlyphs()
	// It grows once, then resets as before. A working set larger than
	// AtlasMaxSize still resets every frame; evicting the least recently used
	// shelf would be the next step if that turns up.
	if side := a.packer.side(); side < AtlasMaxSize {
		a.packer.size = side * 2
		tex.ResizeAtlas(a.packer.size)
	}
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

// inColumns reports whether any of the glyph's ink, with its pen at x, can land
// between the device pixel columns minX and maxX. The pad covers the bitmap's
// antialiased border and the sub-pixel shift it is rasterised at.
func (pg PlacedGlyph) inColumns(x, pixScale, minX, maxX float32) bool {
	const pad = 2
	left := x + fixed266ToFloat32(pg.Glyph.XBearing)*pixScale
	right := left + fixed266ToFloat32(pg.Glyph.Width)*pixScale
	if right < left {
		left, right = right, left
	}
	return right+pad >= minX && left-pad <= maxX
}

// shapeable reports whether a shaped string is a candidate for the atlas at
// all. Glyph 0 means the shaper found nothing, which DrawString renders by
// substituting a replacement character from another face - a second shaping
// path this deliberately does not grow. Whether the glyphs themselves fit is
// the atlas's answer, and a cached one.
//
// Such strings fall back whole rather than getting a second colour atlas beside
// this one. Mixed emoji and text then costs one texture per such label, which is
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
// sizes a hair apart cannot collide through float rounding. The sub-pixel
// offset is kept in 64ths rather than as a phase index, because the number of
// phases depends on the scale and the same index would mean different offsets.
type glyphKey struct {
	face *font.Face
	gid  font.GID
	size int32 // fontSize * pixScale in 26.6 fixed point
	sub  int   // horizontal offset the bitmap was rasterised at, in 64ths of a pixel
}

func newGlyphKey(face *font.Face, gid font.GID, fontSize, pixScale float32, sub int) glyphKey {
	return glyphKey{face: face, gid: gid, size: int32(fontSize * pixScale * 64), sub: sub}
}

// subpixelSteps is how many horizontal positions within a device pixel a glyph
// is cached at. Snapping each glyph to the nearest one keeps spacing within an
// eighth of a logical pixel of where the shaper put it. Denser screens need
// fewer steps for that, and their glyphs cover more texels, so the steps
// shrink as the scale grows and the atlas holds about as much text at any
// density: four at 1x, two at 2x and 3x, one from 4x.
func subpixelSteps(pixScale float32) int {
	return max(1, min(subpixelPerPoint, int(math.Ceil(float64(subpixelPerPoint/pixScale)))))
}

// snapX splits a pen position in device pixels into the whole pixel its quad
// is placed at and the offset, in 64ths of a pixel, its bitmap is rasterised
// at: the nearest of steps evenly spaced positions. Rounding up past the last
// position moves on to the next whole pixel, and negative positions (text
// scrolled off the left edge) floor like any other.
func snapX(x float32, steps int) (pixel, sub int) {
	n := int(math.Round(float64(x) * float64(steps)))
	pixel, phase := n/steps, n%steps
	if phase < 0 {
		pixel, phase = pixel-1, phase+steps
	}
	return pixel, phase * subpixelUnits / steps
}

// glyphEntry is where a glyph landed in the atlas and how its bitmap sits
// relative to the pen position and baseline it was rasterised from.
type glyphEntry struct {
	x, y          int // top-left in atlas pixels
	width, height int
	// bearX and bearY offset the quad from the whole pixel the pen position was
	// snapped to and the baseline; any fraction is inside the bitmap.
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
// Shelves rather than a skyline or a full bin packer, because glyphs of one font
// size are all much the same height, so a shelf wastes very little, and the
// whole structure is three ints. If mixed body and heading sizes ever fragment
// this badly enough to matter, the atlas simply resets and repacks - the
// upgrade path is a skyline packer, not a rewrite of the callers.
type atlasPacker struct {
	size        int // side of the square being filled; zero means AtlasSize
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
	side := p.side()
	slotW, slotH := w+atlasPad, h+atlasPad
	if w <= 0 || h <= 0 || slotW > side || slotH > side {
		return 0, 0, false
	}
	if p.penX+slotW > side { // shelf full, open the next one
		p.shelfY += p.shelfHeight
		p.shelfHeight = 0
		p.penX = 0
	}
	if slotH > p.shelfHeight {
		// A taller glyph grows the shelf rather than opening a new one, which
		// would waste everything already placed on this one. Glyphs already on
		// the shelf start at its top, so growing it downwards cannot disturb them.
		if p.shelfY+slotH > side {
			return 0, 0, false
		}
		p.shelfHeight = slotH
	}
	x, y = p.penX, p.shelfY
	p.penX += slotW
	return x, y, true
}

// reset empties the atlas, keeping its size. Callers must draw anything
// already queued against the old contents first, and forget every glyphEntry
// they hold.
func (p *atlasPacker) reset() {
	p.penX, p.shelfY, p.shelfHeight = 0, 0, 0
}

func (p *atlasPacker) side() int {
	if p.size == 0 {
		return AtlasSize
	}
	return p.size
}
