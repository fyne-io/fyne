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
	// Coverage is one byte per pixel, so each atlas costs 4 MiB.
	glyphAtlasTexSize = 2048
	glyphAtlasPad     = 1

	// Colour glyphs such as emoji cannot reduce to coverage and need their own.
	glyphColourAtlasTexSize = 1024

	// Device sizes at which fewer sub-pixel positions are kept. A quarter of a
	// pixel is a visible slice of a small stroke and nothing on a large one.
	subpixelFullSizePx = 24
	subpixelHalfSizePx = 48

	// GL's default unpack alignment, which slot placement respects so that
	// single byte uploads never need it changed.
	glyphAtlasRowAlign = 4

	// Most sub-pixel positions a glyph is rasterised at, so that kerning is not
	// rounded to whole pixels.
	subpixelPhases = 4
)

// subpixelPhasesFor returns how many sub-pixel positions to keep for glyphs
// with an em box of emPx device pixels. Larger glyphs need fewer.
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

// glyphAtlasKey identifies a glyph bitmap. Colour is absent because bitmaps
// hold coverage and are tinted when drawn.
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

// glyphGPUAtlas packs rasterised glyphs into one GPU texture with a shelf
// packer. It is per-painter, so it lives and dies with one GL context.
type glyphGPUAtlas struct {
	// pixels is the CPU copy of the texture, row major, bpp bytes per pixel:
	// one for a coverage atlas, four for a colour one.
	pixels  []uint8
	bpp     int
	texture Texture
	texSize int
	entries map[glyphAtlasKey]glyphAtlasEntry

	// generation counts resets, so callers can tell that entries they hold are stale.
	generation int

	// resetPending records that a glyph would not fit. The caller resets between
	// strings, where nothing half built depends on the layout.
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

// isColour reports whether a rasterised glyph carries colour of its own.
// Glyphs render white, so an outline has every channel equal to its coverage.
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

// add packs an already rasterised glyph, returning its entry and the region
// of the atlas to upload.
func (a *glyphGPUAtlas) add(key glyphAtlasKey, glyphImg *image.RGBA, baseline int) (glyphAtlasEntry, image.Rectangle) {
	w, h := glyphImg.Bounds().Dx(), glyphImg.Bounds().Dy()

	// A glyph larger than the atlas cannot be packed at any offset.
	if w > a.texSize || h > a.texSize {
		return glyphAtlasEntry{}, image.Rectangle{}
	}

	// Slots start and end on the unpack alignment so that single byte rows
	// upload correctly; misaligned rows come out sheared.
	slotW := (w + glyphAtlasPad + glyphAtlasRowAlign - 1) / glyphAtlasRowAlign * glyphAtlasRowAlign

	// Advance to a new shelf if the glyph does not fit in the current row.
	if a.shelfX+slotW > a.texSize {
		a.shelfY += a.shelfH + glyphAtlasPad
		a.shelfX = 0
		a.shelfH = 0
	}
	if a.shelfY+h > a.texSize {
		// Full. Resetting now would move entries the caller already holds, so ask
		// for a reset between strings instead.
		a.resetPending = true
		return glyphAtlasEntry{}, image.Rectangle{}
	}

	// A coverage atlas keeps only alpha, a colour atlas every channel.
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

	// Upload the whole slot so that offset and width stay aligned.
	return entry, image.Rect(atlasX, atlasY, atlasX+slotW, atlasY+h)
}

// glyphEntry returns the atlas entry for a glyph, rasterising and filing it on
// a miss. Colour glyphs go to the colour atlas, the rest to the coverage one.
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

// uploadAtlas sends the whole atlas to the GPU, which a reset requires.
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

// Two triangles per quad: a strip cannot describe disjoint quads in one draw.
const (
	verticesPerGlyph = 6
	floatsPerGlyph   = verticesPerGlyph * coordinateSize2DWithTexture
)

// textVertices is the cached glyph geometry for one string, in its own GPU
// buffer. Vertices are device pixels relative to the string origin, so moving
// the text needs no new upload.
type textVertices struct {
	buffer Buffer
	// Coverage glyphs come first in the buffer and colour glyphs after, so each
	// group can be drawn from its own atlas without reordering at draw time.
	vertices       int // coverage vertices, from the start of the buffer
	colourVertices int // colour vertices, following them

	generation int     // atlas generation the texture coordinates came from
	pixScale   float32 // scale the glyphs were laid out and rasterised at
}

// usable reports whether cached geometry can still be drawn: not after an
// atlas reset, which moves every texture coordinate, nor after a scale change.
func (v *textVertices) usable(generation int, pixScale float32) bool {
	return v != nil && v.generation == generation && v.pixScale == pixScale
}

// appendGlyphQuad adds one glyph's two triangles to points. offX and offY are
// its whole-pixel offset from the string origin.
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

// glyphCacheKey identifies cached geometry by the text rather than the object
// drawing it, so a refresh that leaves the words alone keeps the geometry.
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

// glyphGeometry returns the glyph quads for text, building them on a miss.
func (p *painter) glyphGeometry(text *canvas.Text, face *paint.FontCacheItem) *textVertices {
	key := glyphCacheKey(text, p.canvas)
	if cached := p.textCache[key]; cached.usable(p.glyphAtlas.generation, p.pixScale) {
		cache.GetTextTexture(key) // keep the expiry marker alive while still drawn
		return cached
	}

	// A full atlas asks to be reset rather than resetting itself, so the quads
	// gathered so far stay valid. One retry is enough unless a single string
	// needs more than the whole atlas.
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
				// The bitmap's baseline is its own run's ascent, not the line's; the
				// difference keeps mixed faces on one baseline (#6448).
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

		// Registered with the text cache for lifetime only: it holds no texture,
		// but expiry calls back so the buffer can be freed.
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

	// The only upload until these words leave the screen. Skipped when empty,
	// which BufferData cannot take.
	if len(p.textBatch) > 0 {
		p.ctx.BindBuffer(arrayBuffer, cached.buffer)
		p.ctx.BufferData(arrayBuffer, p.textBatch, staticDraw)
		p.logError()
	}
	return cached
}

// drawGlyphBatch draws one string's glyph quads. Vertices are relative to the
// string origin, which pos and frame place in clip space, and col tints the
// coverage glyphs.
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

	// The blend is premultiplied, so the colour must be too.
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

	// Point the attribute arrays back at a buffer that outlives this one: without
	// VAOs they are global state, and this buffer is freed when the text changes.
	p.ctx.BindBuffer(arrayBuffer, p.programs.text.buff)
	p.UpdateVertexArray(p.programs.text, attrVertex, coordinateSize2D, coordinateSize2DWithTexture, 0)
	p.UpdateVertexArray(p.programs.text, attrVertexTextureCoordinates, coordinateSize2D, coordinateSize2DWithTexture, coordinateSize2D)
	p.logError()
}
