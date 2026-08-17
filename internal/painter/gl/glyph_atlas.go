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

	// GL's default unpack alignment. Glyph slots are placed and sized to it so
	// that uploading a single byte texture never needs it changed, which fyne's
	// mobile GL binding has no call for.
	glyphAtlasRowAlign = 4

	// Horizontal sub-pixel positions each glyph is rasterised at. Kerning puts
	// glyphs on fractional positions, and a bitmap can only be drawn on whole
	// pixels, so the fraction is baked into the bitmap rather than rounded away
	// (uneven spacing) or resampled at draw time (blurred edges). Four costs
	// four entries per glyph and leaves the error under an eighth of a pixel.
	subpixelPhases = 4
)

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
func subpixelPhaseAt(x float32) (phase int, whole float32) {
	whole = float32(math.Floor(float64(x)))
	phase = int((x - whole) * subpixelPhases)
	if phase >= subpixelPhases { // guard against rounding at the top of the range
		phase = subpixelPhases - 1
	}
	return phase, whole
}

type glyphAtlasEntry struct {
	x, y     int // top-left position in the atlas texture
	w, h     int // dimensions in pixels
	baseline int // baseline position in pixels down from the top of the glyph bitmap
}

// glyphGPUAtlas packs pre-rasterised glyphs into a single GPU texture holding
// one byte of coverage per pixel. New glyphs are appended with a simple shelf
// packer. It is per-painter, so it lives and dies with one GL context.
type glyphGPUAtlas struct {
	// coverage is the CPU copy of the texture, one byte per pixel, row major.
	coverage []uint8
	texture  Texture
	texSize  int
	entries  map[glyphAtlasKey]glyphAtlasEntry

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
	clear(a.coverage)
	a.entries = make(map[glyphAtlasKey]glyphAtlasEntry)
	a.shelfX, a.shelfY, a.shelfH = 0, 0, 0
	a.resetPending = false
	a.generation++
}

func newGlyphGPUAtlas(texSize int) *glyphGPUAtlas {
	return &glyphGPUAtlas{
		coverage: make([]uint8, texSize*texSize),
		texSize:  texSize,
		entries:  make(map[glyphAtlasKey]glyphAtlasEntry),
	}
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
func (a *glyphGPUAtlas) getOrAdd(run shaping.Output, idx, phase int, fontSize, scale float32) (glyphAtlasEntry, image.Rectangle) {
	key := a.cacheKey(run, idx, phase, fontSize, scale)
	if entry, ok := a.entries[key]; ok {
		return entry, image.Rectangle{}
	}

	subpixel := float32(phase) / subpixelPhases
	glyphImg, baseline := paint.RenderGlyphToImage(run, idx, fontSize, scale, subpixel)
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

	// Keep only the alpha channel: the glyph was rasterised in white, so alpha
	// is its coverage and the colour channels carry nothing the atlas needs.
	atlasX, atlasY := a.shelfX, a.shelfY
	for y := 0; y < h; y++ {
		src := glyphImg.PixOffset(0, y)
		dst := (atlasY+y)*a.texSize + atlasX
		for x := 0; x < w; x++ {
			a.coverage[dst+x] = glyphImg.Pix[src+x*4+3]
		}
	}

	// The entry keeps the glyph's true width so quads and texture coordinates
	// are unaffected by the slot padding.
	entry := glyphAtlasEntry{x: atlasX, y: atlasY, w: w, h: h, baseline: baseline}
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

// ensureGlyphAtlas lazily allocates the GPU texture for the glyph atlas.
func (p *painter) ensureGlyphAtlas() {
	if p.glyphAtlas != nil {
		return
	}
	size := glyphAtlasTexSize
	if p.maxTextureSize > 0 && size > p.maxTextureSize {
		size = p.maxTextureSize // older GPUs cap below what we would like
	}
	p.glyphAtlas = newGlyphGPUAtlas(size)
	p.glyphAtlas.texture = p.newTexture(canvas.ImageScaleSmooth)
	p.uploadAtlas()
}

// uploadAtlas sends the whole atlas to the GPU. Needed after a reset, where the
// texture still holds the previous layout's pixels and the incremental uploads
// only cover glyphs added since.
func (p *painter) uploadAtlas() {
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.glyphAtlas.texture)
	p.ctx.TexImage2D(texture2D, 0, p.glyphAtlas.texSize, p.glyphAtlas.texSize,
		colorFormatAlpha, unsignedByte, p.glyphAtlas.coverage)
	p.logError()
}

// uploadAtlasRegion copies dirty pixels from the CPU atlas to the GPU texture
// via TexSubImage2D, avoiding a full texture re-upload.
func (p *painter) uploadAtlasRegion(dirty image.Rectangle) {
	atlas := p.glyphAtlas
	w, h := dirty.Dx(), dirty.Dy()
	pixels := make([]uint8, w*h)
	for y := 0; y < h; y++ {
		src := (dirty.Min.Y+y)*atlas.texSize + dirty.Min.X
		copy(pixels[y*w:(y+1)*w], atlas.coverage[src:src+w])
	}
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, atlas.texture)
	p.ctx.TexSubImage2D(texture2D, 0, dirty.Min.X, dirty.Min.Y, w, h, colorFormatAlpha, unsignedByte, pixels)
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
	buffer   Buffer
	vertices int // number of vertices in buffer, for the draw call

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
func (p *painter) appendGlyphQuad(points []float32, entry glyphAtlasEntry, offX, offY float32) []float32 {
	x1 := offX
	y1 := float32(math.Round(float64(offY)))
	x2 := x1 + float32(entry.w)
	y2 := y1 + float32(entry.h)

	af := float32(p.glyphAtlas.texSize)
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

// glyphGeometry returns the glyph quads for text, building them only when there
// is no usable cached copy. A cached entry survives the object moving, which is
// the common case while scrolling, and is discarded when the text is refreshed
// (see painter.Free), when the scale changes, or when the atlas has reset and
// every texture coordinate in it has become stale.
func (p *painter) glyphGeometry(text *canvas.Text, face *paint.FontCacheItem) *textVertices {
	if cached := p.textCache[text]; cached.usable(p.glyphAtlas.generation, p.pixScale) {
		cache.GetTexture(text) // keep the expiry marker alive while still drawn
		return cached
	}

	// The atlas can run out of room part way through a string. It does not
	// reset under us when that happens, it asks for one, so the quads gathered
	// so far all still address the layout they were built against. The pass is
	// then repeated against the emptied atlas. Only a string whose own glyphs
	// exceed the whole atlas can fail twice, and it settles for what fits
	// rather than looping or drawing against coordinates that have moved.
	for attempt := 0; attempt < 2; attempt++ {
		p.glyphAtlas.resetPending = false
		p.textBatch = p.textBatch[:0]

		paint.WalkStringGlyphs(face.Fonts, text.Text, text.TextSize, text.TextStyle, p.pixScale,
			func(run shaping.Output, idx int, penX, baseY, xOff, yOff float32) {
				// Split the exact position into the pixel the glyph is drawn on
				// and the sub-pixel remainder that is rasterised into it.
				phase, wholeX := subpixelPhaseAt(penX + xOff)

				entry, dirty := p.glyphAtlas.getOrAdd(run, idx, phase, text.TextSize, p.pixScale)
				if !dirty.Empty() {
					p.uploadAtlasRegion(dirty)
				}
				if entry.w == 0 { // did not fit the atlas
					return
				}
				// The bitmap carries its own baseline, which is the run's ascent
				// rather than the line's. Offsetting by the difference keeps runs
				// from differently sized faces on the one baseline (see #6448).
				p.textBatch = p.appendGlyphQuad(p.textBatch, entry,
					wholeX, baseY-float32(entry.baseline)-yOff)
			},
		)

		if !p.glyphAtlas.resetPending {
			break
		}
		if attempt == 0 {
			p.glyphAtlas.reset()
			p.uploadAtlas()
		}
	}
	generation := p.glyphAtlas.generation

	cached, ok := p.textCache[text]
	if !ok {
		cached = &textVertices{buffer: p.ctx.CreateBuffer()}
		if p.textCache == nil {
			p.textCache = make(map[*canvas.Text]*textVertices)
		}
		p.textCache[text] = cached
	}
	cached.vertices = len(p.textBatch) / coordinateSize2DWithTexture
	cached.generation = generation
	cached.pixScale = p.pixScale

	// The only upload this string needs until its text changes. Skipped when
	// every glyph was too large to pack, both because there is nothing to send
	// and because the desktop path takes the address of the first element.
	if len(p.textBatch) > 0 {
		p.ctx.BindBuffer(arrayBuffer, cached.buffer)
		p.ctx.BufferData(arrayBuffer, p.textBatch, staticDraw)
		p.logError()
	}

	// Register with the texture cache purely for liveness: it holds no texture,
	// but this is what drives expiry, and so eventually painter.Free, for text
	// objects that are discarded without ever being refreshed.
	cache.SetTexture(text, cache.NoTexture, p.canvas)
	return cached
}

// drawGlyphBatch issues every glyph quad of one string as a single draw call.
// The vertices are relative to the string origin, so pos places them and frame
// converts device pixels into clip space. The atlas holds coverage rather than
// colour, so col tints the whole batch, which is why one string's glyphs can
// share an entry with the same glyphs drawn elsewhere in another colour.
func (p *painter) drawGlyphBatch(cached *textVertices, col color.Color, pos fyne.Position, frame fyne.Size) {
	if cached.vertices == 0 {
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

	r, g, b, a := getFragmentColor(col)
	p.SetUniform4f(p.programs.text, attrColor, r, g, b, a)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.glyphAtlas.texture)
	p.logError()

	p.ctx.DrawArrays(triangles, 0, cached.vertices)
	p.logError()
}
