package gl

import (
	"image"
	"image/color"
	"math"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/shaping"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/cache"
	paint "fyne.io/fyne/v2/internal/painter"
)

const (
	// The atlas holds one byte of coverage per pixel rather than four of
	// colour, so this is 4 MiB on the GPU and the same again on the CPU, the
	// size a 1024 square RGBA atlas used to cost while holding four times as
	// many glyphs. Headroom matters most at phone pixel densities, where a
	// glyph covers nine times the area it does at 1x.
	glyphAtlasTexSize = 2048
	glyphAtlasPad     = 1

	// Colour glyphs, meaning emoji and anything else from a bitmap or COLR
	// face, cannot be reduced to coverage and need their own smaller atlas of
	// full colour. There are far fewer of them than there are letters.
	glyphColourAtlasTexSize = 1024

	// Device sizes at which fewer sub-pixel positions are kept. A quarter of a
	// pixel is a visible slice of a small stroke and nothing on a large one.
	subpixelFullSizePx = 24
	subpixelHalfSizePx = 48

	// GL's default unpack alignment. Glyph slots are placed and sized to it so
	// that uploading a single byte texture never needs it changed, which fyne's
	// mobile GL binding has no call for.
	glyphAtlasRowAlign = 4

	// Most horizontal sub-pixel positions a glyph is rasterised at. Kerning puts
	// glyphs on fractional positions, and a bitmap can only be drawn on whole
	// pixels, so the fraction is baked into the bitmap rather than rounded away
	// (uneven spacing) or resampled at draw time (blurred edges).
	subpixelPhases = 4
)

// subpixelPhasesFor returns how many sub-pixel positions to keep for glyphs
// whose em box is emPx device pixels. Each one costs another copy of every
// glyph, and what it buys shrinks as the glyphs grow: a quarter of a pixel is a
// visible slice of a small stroke and nothing at all on a large one. Large text
// therefore settles for halves, which is what makes a screen of CJK at phone
// density fit in the atlas at all.
func subpixelPhasesFor(emPx float32) int {
	switch {
	case emPx <= subpixelFullSizePx: // small text, where uneven spacing is easiest to see
		return subpixelPhases
	case emPx <= subpixelHalfSizePx:
		return 2
	default: // a whole pixel here is already a small fraction of a stroke
		return 1
	}
}

// glyphAtlasKey identifies a glyph bitmap. Colour is deliberately absent:
// bitmaps hold coverage and are tinted when drawn, so one entry serves every
// colour. Keying on colour instead multiplied the entry count by the number of
// colours in the theme, which overflowed the atlas at phone pixel densities.
type glyphAtlasKey struct {
	face    *font.Face // pointer identity; stable while fontCache is alive
	gid     font.GID
	pixSize int32 // round(fontSize * pixScale * 64), avoids float key issues
	phase   uint8 // horizontal sub-pixel position, 0 to subpixelPhases-1
}

// subpixelPhaseAt returns which sub-pixel position x falls into, and the whole
// pixel the glyph should then be drawn at.
func subpixelPhaseAt(x float32, phases int) (phase int, whole float32) {
	whole = float32(math.Floor(float64(x)))
	phase = int((x - whole) * float32(phases))
	if phase >= phases { // guard against rounding at the top of the range
		phase = phases - 1
	}
	return phase, whole
}

type glyphAtlasEntry struct {
	x, y     int  // top-left position in the atlas texture
	w, h     int  // dimensions in pixels
	baseline int  // baseline position in pixels down from the top of the glyph bitmap
	colour   bool // lives in the colour atlas and is drawn untinted
}

// glyphGPUAtlas packs pre-rasterised glyphs into a single GPU texture holding
// one byte of coverage per pixel. New glyphs are appended with a simple shelf
// packer. It is per-painter, so it lives and dies with one GL context.
type glyphGPUAtlas struct {
	// pixels is the CPU copy of the texture, row major, bpp bytes per pixel:
	// one for a coverage atlas, four for a colour one.
	pixels  []uint8
	bpp     int
	texture Texture
	texSize int
	entries map[glyphAtlasKey]glyphAtlasEntry

	// generation counts resets. Callers that collect several entries before
	// drawing compare it either side to notice that the entries they gathered
	// were invalidated part way through.
	generation int

	// resetPending records that a glyph could not be packed. The reset itself
	// waits for the caller to reach a point where no half-built string depends
	// on the current layout.
	resetPending bool

	// shelf packer cursor
	shelfX, shelfY, shelfH int
}

// reset empties the atlas so packing starts over. Every texture coordinate
// handed out before now becomes wrong, which the generation bump signals.
func (a *glyphGPUAtlas) reset() {
	clear(a.pixels)
	a.entries = make(map[glyphAtlasKey]glyphAtlasEntry)
	a.shelfX, a.shelfY, a.shelfH = 0, 0, 0
	a.resetPending = false
	a.generation++
}

func newGlyphGPUAtlas(texSize int) *glyphGPUAtlas {
	return newGlyphAtlasOfDepth(texSize, 1)
}

func newGlyphAtlasOfDepth(texSize, bpp int) *glyphGPUAtlas {
	return &glyphGPUAtlas{
		pixels:  make([]uint8, texSize*texSize*bpp),
		bpp:     bpp,
		texSize: texSize,
		entries: make(map[glyphAtlasKey]glyphAtlasEntry),
	}
}

// isColour reports whether a rasterised glyph carries colour of its own. Glyphs
// are rendered white, so an ordinary outline comes back with every channel
// equal to its coverage; an emoji does not.
func isColour(img *image.RGBA) bool {
	for i := 0; i+3 < len(img.Pix); i += 4 {
		a := img.Pix[i+3]
		if a == 0 {
			continue
		}
		if img.Pix[i] != a || img.Pix[i+1] != a || img.Pix[i+2] != a {
			return true
		}
	}
	return false
}

func (*glyphGPUAtlas) cacheKey(run shaping.Output, idx, phase int, fontSize, scale float32) glyphAtlasKey {
	return glyphAtlasKey{
		face:    run.Face,
		gid:     run.Glyphs[idx].GlyphID,
		pixSize: int32(math.Round(float64(fontSize * scale * 64))),
		phase:   uint8(phase), //gosec:disable G115 -- phase is always 0 to subpixelPhases-1
	}
}

// getOrAdd returns the atlas entry for a glyph, adding it on a miss.
// The second return value is the dirty rectangle written into cpuImg; it is
// empty when the entry was already cached so no GPU upload is needed.
// add packs an already rasterised glyph. The caller decides which atlas a glyph
// belongs in, since that depends on whether it came back with colour.
func (a *glyphGPUAtlas) add(key glyphAtlasKey, glyphImg *image.RGBA, baseline int) (glyphAtlasEntry, image.Rectangle) {
	w, h := glyphImg.Bounds().Dx(), glyphImg.Bounds().Dy()

	// A glyph bigger than the atlas cannot be packed at any offset, and writing
	// it anyway would run past the end of the texture. Report it as empty so it
	// is skipped rather than drawn wrongly; only text far larger than the atlas
	// is affected.
	if w > a.texSize || h > a.texSize {
		return glyphAtlasEntry{}, image.Rectangle{}
	}

	// Slots are a multiple of glyphAtlasRowAlign wide and start on the same
	// boundary, so every row of a sub-image upload begins on a boundary GL is
	// content to read a single byte texture from. Without that its default
	// unpack alignment misreads each row of a glyph whose width is not a
	// multiple of four, and the text comes out sheared.
	slotW := (w + glyphAtlasPad + glyphAtlasRowAlign - 1) / glyphAtlasRowAlign * glyphAtlasRowAlign

	// Advance to a new shelf if the glyph does not fit in the current row.
	if a.shelfX+slotW > a.texSize {
		a.shelfY += a.shelfH + glyphAtlasPad
		a.shelfX = 0
		a.shelfH = 0
	}
	if a.shelfY+h > a.texSize {
		// Full. Resetting here would move entries that the caller has already
		// collected for the string it is part way through, so the glyph is
		// dropped for now and the reset left for the caller to trigger between
		// passes, where it invalidates nothing mid-flight.
		a.resetPending = true
		return glyphAtlasEntry{}, image.Rectangle{}
	}

	// A coverage atlas keeps only the alpha channel, since the glyph was
	// rasterised in white and the colour channels then carry nothing. A colour
	// atlas keeps the lot.
	atlasX, atlasY := a.shelfX, a.shelfY
	for y := 0; y < h; y++ {
		src := glyphImg.PixOffset(0, y)
		dst := ((atlasY+y)*a.texSize + atlasX) * a.bpp
		if a.bpp == 1 {
			for x := 0; x < w; x++ {
				a.pixels[dst+x] = glyphImg.Pix[src+x*4+3]
			}
			continue
		}
		copy(a.pixels[dst:dst+w*4], glyphImg.Pix[src:src+w*4])
	}

	// The entry keeps the glyph's true width so quads and texture coordinates
	// are unaffected by the slot padding.
	entry := glyphAtlasEntry{x: atlasX, y: atlasY, w: w, h: h, baseline: baseline, colour: a.bpp == 4}
	a.entries[key] = entry
	if h > a.shelfH {
		a.shelfH = h
	}
	a.shelfX += slotW

	// The upload covers the whole slot, not just the glyph, so both its start
	// and its width stay on the alignment boundary. The padding columns were
	// left at zero coverage, so uploading them changes nothing on screen.
	return entry, image.Rect(atlasX, atlasY, atlasX+slotW, atlasY+h)
}

// glyphEntry finds a glyph in whichever atlas holds it, rasterising and filing
// it on a miss. Which atlas that is depends on the glyph: letters reduce to
// coverage and are tinted when drawn, emoji keep their own colours.
func (p *painter) glyphEntry(run shaping.Output, idx, phase, phases int, fontSize, scale float32) glyphAtlasEntry {
	key := p.glyphAtlas.cacheKey(run, idx, phase, fontSize, scale)
	if entry, ok := p.glyphAtlas.entries[key]; ok {
		return entry
	}
	if entry, ok := p.glyphColourAtlas.entries[key]; ok {
		return entry
	}

	subpixel := float32(phase) / float32(phases)
	glyphImg, baseline := paint.RenderGlyphToImage(run, idx, fontSize, scale, subpixel)

	atlas := p.glyphAtlas
	if isColour(glyphImg) {
		atlas = p.glyphColourAtlas
	}
	entry, dirty := atlas.add(key, glyphImg, baseline)
	if !dirty.Empty() {
		p.uploadAtlasRegion(atlas, dirty)
	}
	return entry
}

// ensureGlyphAtlas lazily allocates the GPU texture for the glyph atlas.
func (p *painter) ensureGlyphAtlas() {
	if p.glyphAtlas != nil {
		return
	}
	size := glyphAtlasTexSize
	if p.maxTextureSize > 0 && size > p.maxTextureSize {
		size = p.maxTextureSize // older GPUs cap below what we would like
	}
	p.glyphAtlas = newGlyphAtlasOfDepth(size, 1)
	p.glyphAtlas.texture = p.newTexture(canvas.ImageScaleSmooth)
	p.uploadAtlas(p.glyphAtlas)

	colourSize := glyphColourAtlasTexSize
	if p.maxTextureSize > 0 && colourSize > p.maxTextureSize {
		colourSize = p.maxTextureSize
	}
	p.glyphColourAtlas = newGlyphAtlasOfDepth(colourSize, 4)
	p.glyphColourAtlas.texture = p.newTexture(canvas.ImageScaleSmooth)
	p.uploadAtlas(p.glyphColourAtlas)
}

// uploadAtlas sends the whole atlas to the GPU. Needed after a reset, where the
// texture still holds the previous layout's pixels and the incremental uploads
// only cover glyphs added since.
func (p *painter) uploadAtlas(a *glyphGPUAtlas) {
	format := uint32(colorFormatAlpha)
	if a.bpp == 4 {
		format = colorFormatRGBA
	}
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, a.texture)
	p.ctx.TexImage2D(texture2D, 0, a.texSize, a.texSize, format, unsignedByte, a.pixels)
	p.logError()
}

// uploadAtlasRegion copies dirty pixels from the CPU atlas to the GPU texture
// via TexSubImage2D, avoiding a full texture re-upload.
func (p *painter) uploadAtlasRegion(atlas *glyphGPUAtlas, dirty image.Rectangle) {
	bpp := atlas.bpp
	format := uint32(colorFormatAlpha)
	if bpp == 4 {
		format = colorFormatRGBA
	}
	w, h := dirty.Dx(), dirty.Dy()
	pixels := make([]uint8, w*h*bpp)
	for y := 0; y < h; y++ {
		src := ((dirty.Min.Y+y)*atlas.texSize + dirty.Min.X) * bpp
		copy(pixels[y*w*bpp:(y+1)*w*bpp], atlas.pixels[src:src+w*bpp])
	}
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, atlas.texture)
	p.ctx.TexSubImage2D(texture2D, 0, dirty.Min.X, dirty.Min.Y, w, h, format, unsignedByte, pixels)
	p.logError()
}

// A glyph quad is emitted as two triangles rather than a strip, because a strip
// cannot describe several disjoint quads in one draw without degenerate
// vertices joining them.
const (
	verticesPerGlyph = 6
	floatsPerGlyph   = verticesPerGlyph * coordinateSize2DWithTexture
)

// textVertices is the cached glyph geometry for one canvas.Text, held in its
// own GPU buffer. The vertices are in device pixels relative to the string
// origin rather than in clip space, so neither the data nor the buffer needs
// touching while the object merely moves, which is what happens to every
// visible string on every frame of a scroll. Drawing a cached string uploads
// nothing at all.
type textVertices struct {
	buffer Buffer
	// Coverage glyphs come first in the buffer and colour glyphs after, so each
	// group can be drawn from its own atlas without reordering at draw time.
	vertices       int // coverage vertices, from the start of the buffer
	colourVertices int // colour vertices, following them

	generation int     // atlas generation the texture coordinates came from
	pixScale   float32 // scale the glyphs were laid out and rasterised at
}

// usable reports whether cached geometry can be drawn as it stands. It cannot
// once the atlas has reset, since every texture coordinate in it then points at
// the wrong place, nor once the scale has changed, since the glyphs were both
// laid out and rasterised for the old one.
func (v *textVertices) usable(generation int, pixScale float32) bool {
	return v != nil && v.generation == generation && v.pixScale == pixScale
}

// appendGlyphQuad adds one glyph's two triangles to points and returns the
// extended slice. The quad samples the sub-region [entry.x, entry.y,
// entry.x+entry.w, entry.y+entry.h] of the shared atlas texture.
//
// offX and offY are the glyph's device-pixel offset from the string origin.
// offX is already a whole pixel, its fractional part having been rasterised
// into the bitmap as a sub-pixel offset, so the quad maps one texel to one
// pixel and stays crisp while still sitting where kerning asked for it. offY
// is rounded here, since vertical sub-pixel positioning buys nothing on a
// shared baseline.
func (*painter) appendGlyphQuad(points []float32, entry glyphAtlasEntry, offX, offY float32, atlas *glyphGPUAtlas) []float32 {
	x1 := offX
	y1 := float32(math.Round(float64(offY)))
	x2 := x1 + float32(entry.w)
	y2 := y1 + float32(entry.h)

	af := float32(atlas.texSize)
	uMin := float32(entry.x) / af
	vMin := float32(entry.y) / af
	uMax := float32(entry.x+entry.w) / af
	vMax := float32(entry.y+entry.h) / af

	// Vertex y and texture v both increase downwards here, so the top edge of
	// the quad takes vMin and the bottom edge vMax.
	return append(points,
		x1, y2, uMin, vMax, // bottom left
		x1, y1, uMin, vMin, // top left
		x2, y2, uMax, vMax, // bottom right

		x2, y2, uMax, vMax, // bottom right
		x1, y1, uMin, vMin, // top left
		x2, y1, uMax, vMin, // top right
	)
}

// glyphCacheKey identifies cached geometry by what the glyphs are, not by which
// object is showing them, mirroring how the per-string textures were keyed
// before this. Colour is left unset because the atlas holds coverage and the
// tint is a uniform, so one set of quads serves the same words in every colour.
//
// Keying on the object instead would look right and behave badly: a widget that
// refreshes rebuilds every string it draws, so typing into an Entry rebuilt the
// whole visible text on each keystroke rather than the line that changed.
func glyphCacheKey(text *canvas.Text, c fyne.Canvas) cache.FontCacheEntry {
	source := ""
	if text.FontSource != nil {
		source = text.FontSource.Name()
	}

	ent := cache.FontCacheEntry{Canvas: c}
	ent.Text = text.Text
	ent.Size = text.TextSize
	ent.Style = text.TextStyle
	ent.Source = source
	return ent
}

// glyphGeometry returns the glyph quads for text, building them only when there
// is no usable cached copy. A cached entry survives both the object moving and
// the object being refreshed, and is discarded when the scale changes, when the
// atlas has reset and every texture coordinate in it has gone stale, or when
// the text it belongs to has not been drawn for long enough to expire.
func (p *painter) glyphGeometry(text *canvas.Text, face *paint.FontCacheItem) *textVertices {
	key := glyphCacheKey(text, p.canvas)
	if cached := p.textCache[key]; cached.usable(p.glyphAtlas.generation, p.pixScale) {
		cache.GetTextTexture(key) // keep the expiry marker alive while still drawn
		return cached
	}

	// The atlas can run out of room part way through a string. It does not
	// reset under us when that happens, it asks for one, so the quads gathered
	// so far all still address the layout they were built against. The pass is
	// then repeated against the emptied atlas. Only a string whose own glyphs
	// exceed the whole atlas can fail twice, and it settles for what fits
	// rather than looping or drawing against coordinates that have moved.
	phases := subpixelPhasesFor(text.TextSize * p.pixScale)
	var colourBatch []float32
	for attempt := 0; attempt < 2; attempt++ {
		p.glyphAtlas.resetPending = false
		p.glyphColourAtlas.resetPending = false
		p.textBatch = p.textBatch[:0]
		colourBatch = colourBatch[:0]

		paint.WalkStringGlyphs(face.Fonts, text.Text, text.TextSize, text.TextStyle, p.pixScale,
			func(run shaping.Output, idx int, penX, baseY, xOff, yOff float32) {
				// Split the exact position into the pixel the glyph is drawn on
				// and the sub-pixel remainder that is rasterised into it.
				phase, wholeX := subpixelPhaseAt(penX+xOff, phases)

				entry := p.glyphEntry(run, idx, phase, phases, text.TextSize, p.pixScale)
				if entry.w == 0 { // did not fit the atlas
					return
				}
				// The bitmap carries its own baseline, which is the run's ascent
				// rather than the line's. Offsetting by the difference keeps runs
				// from differently sized faces on the one baseline (see #6448).
				offY := baseY - float32(entry.baseline) - yOff
				if entry.colour {
					colourBatch = p.appendGlyphQuad(colourBatch, entry, wholeX, offY, p.glyphColourAtlas)
					return
				}
				p.textBatch = p.appendGlyphQuad(p.textBatch, entry, wholeX, offY, p.glyphAtlas)
			},
		)

		if !p.glyphAtlas.resetPending && !p.glyphColourAtlas.resetPending {
			break
		}
		if attempt == 0 {
			if p.glyphAtlas.resetPending {
				p.glyphAtlas.reset()
				p.uploadAtlas(p.glyphAtlas)
			}
			if p.glyphColourAtlas.resetPending {
				p.glyphColourAtlas.reset()
				p.uploadAtlas(p.glyphColourAtlas)
			}
		}
	}
	p.textBatch = append(p.textBatch, colourBatch...)
	generation := p.glyphAtlas.generation

	cached, ok := p.textCache[key]
	if !ok {
		cached = &textVertices{buffer: p.ctx.CreateBuffer()}
		if p.textCache == nil {
			p.textCache = make(map[cache.FontCacheEntry]*textVertices)
		}
		p.textCache[key] = cached

		// Registered with the text cache purely for its lifetime: it holds no
		// texture, but this is what expires entries whose words have not been
		// on screen for a while, and calls back so the buffer can go with them.
		cache.SetTextTexture(key, cache.NoTexture, p.canvas, func() {
			p.ctx.DeleteBuffer(cached.buffer)
			p.logError()
			delete(p.textCache, key)
		})
	}
	cached.vertices = (len(p.textBatch) - len(colourBatch)) / coordinateSize2DWithTexture
	cached.colourVertices = len(colourBatch) / coordinateSize2DWithTexture
	cached.generation = generation
	cached.pixScale = p.pixScale

	// The only upload these words need until they leave the screen. Skipped when
	// every glyph was too large to pack, both because there is nothing to send
	// and because the desktop path takes the address of the first element.
	if len(p.textBatch) > 0 {
		p.ctx.BindBuffer(arrayBuffer, cached.buffer)
		p.ctx.BufferData(arrayBuffer, p.textBatch, staticDraw)
		p.logError()
	}
	return cached
}

// drawGlyphBatch issues every glyph quad of one string as a single draw call.
// The vertices are relative to the string origin, so pos places them and frame
// converts device pixels into clip space. The atlas holds coverage rather than
// colour, so col tints the whole batch, which is why one string's glyphs can
// share an entry with the same glyphs drawn elsewhere in another colour.
func (p *painter) drawGlyphBatch(cached *textVertices, col color.Color, pos fyne.Position, frame fyne.Size) {
	if cached.vertices == 0 && cached.colourVertices == 0 {
		return
	}

	p.ctx.UseProgram(p.programs.text.ref)

	// The string's own buffer already holds its geometry, so there is nothing
	// to upload here however many glyphs it has.
	p.ctx.BindBuffer(arrayBuffer, cached.buffer)
	p.logError()

	p.UpdateVertexArray(p.programs.text, attrVertex, coordinateSize2D, coordinateSize2DWithTexture, 0)
	p.UpdateVertexArray(p.programs.text, attrVertexTextureCoordinates, coordinateSize2D, coordinateSize2DWithTexture, coordinateSize2D)

	// Snap the origin to the pixel grid so the whole string lands on whole
	// pixels; the glyph offsets within it are already integers.
	originX := roundToPixel(pos.X, p.pixScale)
	originY := roundToPixel(pos.Y, p.pixScale)
	p.SetUniform2f(p.programs.text, attrOrigin,
		-1+originX/frame.Width*2,
		1-originY/frame.Height*2)
	p.SetUniform2f(p.programs.text, attrPixelScale,
		2/(frame.Width*p.pixScale),
		-2/(frame.Height*p.pixScale))

	// getFragmentColor hands back straight colour with alpha alongside, which
	// suits the shape shaders because they blend with SRC_ALPHA. Text goes
	// through the premultiplied blend that the texture paths use, so the colour
	// has to be premultiplied to match, or anything drawn with alpha below one
	// comes out too bright.
	r, g, b, a := getFragmentColor(col)
	p.SetUniform4f(p.programs.text, attrColor, r*a, g*a, b*a, a)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	// Coverage glyphs first, tinted, then any that carry their own colour. Most
	// strings have none of the latter and so still draw in a single call.
	p.ctx.ActiveTexture(texture0)
	if cached.vertices > 0 {
		p.SetUniform1f(p.programs.text, attrOwnColor, 0)
		p.ctx.BindTexture(texture2D, p.glyphAtlas.texture)
		p.ctx.DrawArrays(triangles, 0, cached.vertices)
		p.logError()
	}
	if cached.colourVertices > 0 {
		p.SetUniform1f(p.programs.text, attrOwnColor, 1)
		p.ctx.BindTexture(texture2D, p.glyphColourAtlas.texture)
		p.ctx.DrawArrays(triangles, cached.vertices, cached.colourVertices)
		p.logError()
	}

	// Leave no vertex attribute array pointing into this object's buffer. Every
	// other program draws from a buffer that lives as long as the context, but
	// these are freed when their text is refreshed, and without VAOs the
	// pointers are global state: a later draw would read from a deleted buffer.
	// Typing does exactly that, freeing a buffer on every keystroke.
	p.ctx.BindBuffer(arrayBuffer, p.programs.text.buff)
	p.UpdateVertexArray(p.programs.text, attrVertex, coordinateSize2D, coordinateSize2DWithTexture, 0)
	p.UpdateVertexArray(p.programs.text, attrVertexTextureCoordinates, coordinateSize2D, coordinateSize2DWithTexture, coordinateSize2D)
	p.logError()
}
