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
	glyphAtlasTexSize = 1024
	glyphAtlasPad     = 1
)

type glyphAtlasKey struct {
	face       *font.Face // pointer identity; stable while fontCache is alive
	gid        font.GID
	pixSize    int32 // round(fontSize * pixScale * 64), avoids float key issues
	r, g, b, a uint8
}

type glyphAtlasEntry struct {
	x, y     int // top-left position in the atlas texture
	w, h     int // dimensions in pixels
	baseline int // baseline position in pixels down from the top of the glyph bitmap
}

// glyphGPUAtlas packs pre-rasterised glyph bitmaps into a single GPU texture
// (glyphAtlasTexSize × glyphAtlasTexSize RGBA).  New glyphs are appended with
// a simple shelf packer; the atlas resets (with a brief visual flash) when it
// fills up.  It is per-painter so it lives and dies with a single GL context.
type glyphGPUAtlas struct {
	cpuImg  *image.RGBA
	texture Texture
	texSize int
	entries map[glyphAtlasKey]glyphAtlasEntry

	// generation counts resets. Callers that collect several entries before
	// drawing compare it either side to notice that the entries they gathered
	// were invalidated part way through.
	generation int

	// shelf packer cursor
	shelfX, shelfY, shelfH int
}

func newGlyphGPUAtlas(texSize int) *glyphGPUAtlas {
	return &glyphGPUAtlas{
		cpuImg:  image.NewRGBA(image.Rect(0, 0, texSize, texSize)),
		texSize: texSize,
		entries: make(map[glyphAtlasKey]glyphAtlasEntry),
	}
}

func (a *glyphGPUAtlas) cacheKey(run shaping.Output, idx int, fontSize, scale float32, col color.Color) glyphAtlasKey {
	g := run.Glyphs[idx]
	r32, g32, b32, a32 := col.RGBA()
	pixSize := int32(math.Round(float64(fontSize * scale * 64)))
	return glyphAtlasKey{
		face:    run.Face,
		gid:     g.GlyphID,
		pixSize: pixSize,
		r:       uint8((r32 >> 8) & 0xff), //gosec:disable G115 -- value is always 0-255 after the shift and mask
		g:       uint8((g32 >> 8) & 0xff), //gosec:disable G115 -- value is always 0-255 after the shift and mask
		b:       uint8((b32 >> 8) & 0xff), //gosec:disable G115 -- value is always 0-255 after the shift and mask
		a:       uint8((a32 >> 8) & 0xff), //gosec:disable G115 -- value is always 0-255 after the shift and mask
	}
}

// getOrAdd returns the atlas entry for a glyph, adding it on a miss.
// The second return value is the dirty rectangle written into cpuImg; it is
// empty when the entry was already cached so no GPU upload is needed.
func (a *glyphGPUAtlas) getOrAdd(run shaping.Output, idx int, fontSize, scale float32, col color.Color) (glyphAtlasEntry, image.Rectangle) {
	key := a.cacheKey(run, idx, fontSize, scale, col)
	if entry, ok := a.entries[key]; ok {
		return entry, image.Rectangle{}
	}

	glyphImg, baseline := paint.RenderGlyphToImage(run, idx, fontSize, scale, col)
	w, h := glyphImg.Bounds().Dx(), glyphImg.Bounds().Dy()

	// Advance to a new shelf if the glyph does not fit in the current row.
	if a.shelfX+w+glyphAtlasPad > a.texSize {
		a.shelfY += a.shelfH + glyphAtlasPad
		a.shelfX = 0
		a.shelfH = 0
	}
	if a.shelfY+h > a.texSize {
		// Atlas is full; wipe it and start over.
		a.cpuImg = image.NewRGBA(image.Rect(0, 0, a.texSize, a.texSize))
		a.entries = make(map[glyphAtlasKey]glyphAtlasEntry)
		a.shelfX, a.shelfY, a.shelfH = 0, 0, 0
		a.generation++
	}

	atlasX, atlasY := a.shelfX, a.shelfY
	for y := 0; y < h; y++ {
		src := glyphImg.PixOffset(0, y)
		dst := a.cpuImg.PixOffset(atlasX, atlasY+y)
		copy(a.cpuImg.Pix[dst:dst+w*4], glyphImg.Pix[src:src+w*4])
	}

	entry := glyphAtlasEntry{x: atlasX, y: atlasY, w: w, h: h, baseline: baseline}
	a.entries[key] = entry
	if h > a.shelfH {
		a.shelfH = h
	}
	a.shelfX += w + glyphAtlasPad

	return entry, image.Rect(atlasX, atlasY, atlasX+w, atlasY+h)
}

// ensureGlyphAtlas lazily allocates the GPU texture for the glyph atlas.
func (p *painter) ensureGlyphAtlas() {
	if p.glyphAtlas != nil {
		return
	}
	p.glyphAtlas = newGlyphGPUAtlas(glyphAtlasTexSize)
	p.glyphAtlas.texture = p.newTexture(canvas.ImageScaleSmooth)
	p.ctx.TexImage2D(texture2D, 0, glyphAtlasTexSize, glyphAtlasTexSize, colorFormatRGBA, unsignedByte, p.glyphAtlas.cpuImg.Pix)
	p.logError()
}

// uploadAtlasRegion copies dirty pixels from the CPU atlas image to the GPU
// texture via TexSubImage2D, avoiding a full texture re-upload.
func (p *painter) uploadAtlasRegion(dirty image.Rectangle) {
	w, h := dirty.Dx(), dirty.Dy()
	pixels := make([]uint8, w*h*4)
	for y := 0; y < h; y++ {
		src := p.glyphAtlas.cpuImg.PixOffset(dirty.Min.X, dirty.Min.Y+y)
		copy(pixels[y*w*4:], p.glyphAtlas.cpuImg.Pix[src:src+w*4])
	}
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.glyphAtlas.texture)
	p.ctx.TexSubImage2D(texture2D, 0, dirty.Min.X, dirty.Min.Y, w, h, colorFormatRGBA, unsignedByte, pixels)
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

// appendGlyphQuad adds one glyph's two triangles to points and returns the
// extended slice. The quad samples the sub-region [entry.x, entry.y,
// entry.x+entry.w, entry.y+entry.h] of the shared atlas texture.
//
// offX and offY are the glyph's device-pixel offset from the string origin.
// They are rounded to whole pixels so that the glyph bitmap, which was
// rasterised at integer alignment, maps one texel to one pixel. Rounding
// relative to the origin rather than to the screen also keeps letter spacing
// stable as the string moves, where rounding absolute positions made it shift
// from frame to frame.
func (p *painter) appendGlyphQuad(points []float32, entry glyphAtlasEntry, offX, offY float32) []float32 {
	x1 := float32(math.Round(float64(offX)))
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
func (p *painter) glyphGeometry(text *canvas.Text, face *paint.FontCacheItem, col color.Color) *textVertices {
	if cached, ok := p.textCache[text]; ok &&
		cached.generation == p.glyphAtlas.generation && cached.pixScale == p.pixScale {
		cache.GetTexture(text) // keep the expiry marker alive while still drawn
		return cached
	}

	// Adding a glyph can fill the atlas and reset it, which invalidates the
	// quads gathered so far, so the pass repeats if the generation moved. A
	// second pass cannot reset again unless one string alone exceeds the whole
	// atlas, in which case the last attempt is the best available.
	var generation int
	for attempt := 0; attempt < 2; attempt++ {
		generation = p.glyphAtlas.generation
		p.textBatch = p.textBatch[:0]

		paint.WalkStringGlyphs(face.Fonts, text.Text, text.TextSize, text.TextStyle, p.pixScale,
			func(run shaping.Output, idx int, penX, baseY, xOff, yOff float32) {
				entry, dirty := p.glyphAtlas.getOrAdd(run, idx, text.TextSize, p.pixScale, col)
				if !dirty.Empty() {
					p.uploadAtlasRegion(dirty)
				}
				// The bitmap carries its own baseline, which is the run's ascent
				// rather than the line's. Offsetting by the difference keeps runs
				// from differently sized faces on the one baseline (see #6448).
				p.textBatch = p.appendGlyphQuad(p.textBatch, entry,
					penX+xOff, baseY-float32(entry.baseline)-yOff)
			},
		)

		if p.glyphAtlas.generation == generation {
			break
		}
	}

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

	// The only upload this string needs until its text changes.
	p.ctx.BindBuffer(arrayBuffer, cached.buffer)
	p.ctx.BufferData(arrayBuffer, p.textBatch, staticDraw)
	p.logError()

	// Register with the texture cache purely for liveness: it holds no texture,
	// but this is what drives expiry, and so eventually painter.Free, for text
	// objects that are discarded without ever being refreshed.
	cache.SetTexture(text, cache.NoTexture, p.canvas)
	return cached
}

// drawGlyphBatch issues every glyph quad of one string as a single draw call.
// The vertices are relative to the string origin, so pos places them and frame
// converts device pixels into clip space.
func (p *painter) drawGlyphBatch(cached *textVertices, pos fyne.Position, frame fyne.Size) {
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

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.glyphAtlas.texture)
	p.logError()

	p.ctx.DrawArrays(triangles, 0, cached.vertices)
	p.logError()
}
