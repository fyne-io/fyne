package gl

import (
	"image"
	"image/color"
	"math"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/shaping"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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

	if !p.textBufferValid {
		// Sized for a typical line; drawGlyphBatch reallocates when a string
		// needs more, so this is a starting point rather than a limit.
		p.textBuffer = p.createBuffer(floatsPerGlyph * 128)
		p.textBufferValid = true
	}
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

// appendGlyphQuad adds one glyph's two triangles to points and returns the
// extended slice. The quad samples the sub-region [entry.x, entry.y,
// entry.x+entry.w, entry.y+entry.h] of the shared atlas texture. pos and size
// are in logical (fyne) coordinate space; frame is the canvas size.
//
// The geometry deliberately matches what rectCoords produces for a stretched
// fill with no padding, so batching changes only the number of draw calls and
// not a single pixel of output.
func (p *painter) appendGlyphQuad(points []float32, entry glyphAtlasEntry, pos fyne.Position, size, frame fyne.Size) []float32 {
	pixelSize, pixelPos := roundToPixelCoords(size, pos, p.pixScale)

	x1 := -1 + pixelPos.X/frame.Width*2
	x2 := -1 + (pixelPos.X+pixelSize.Width)/frame.Width*2
	y1 := 1 - pixelPos.Y/frame.Height*2
	y2 := 1 - (pixelPos.Y+pixelSize.Height)/frame.Height*2

	af := float32(p.glyphAtlas.texSize)
	uMin := float32(entry.x) / af
	vMin := float32(entry.y) / af
	uMax := float32(entry.x+entry.w) / af
	vMax := float32(entry.y+entry.h) / af

	// Corners are named as rectCoords names them; v runs opposite to screen y,
	// so the two "top" vertices take vMax. The winding matches the triangle
	// strip this replaced: (top left, bottom left, top right) then
	// (top right, bottom left, bottom right).
	return append(points,
		x1, y2, uMin, vMax, // top left
		x1, y1, uMin, vMin, // bottom left
		x2, y2, uMax, vMax, // top right

		x2, y2, uMax, vMax, // top right
		x1, y1, uMin, vMin, // bottom left
		x2, y1, uMax, vMin, // bottom right
	)
}

// drawGlyphBatch issues every glyph quad collected for one string as a single
// draw call. The simple shader ignores size and inset while cornerRadius is 0,
// so one uniform set covers the whole batch however many glyphs it holds.
func (p *painter) drawGlyphBatch(points []float32) {
	if len(points) == 0 {
		return
	}

	p.ctx.UseProgram(p.programs.simple.ref)

	// The batch size varies per string, so this goes through BufferData rather
	// than p.updateBuffer: BufferSubData cannot grow the allocation it writes
	// into, and the mobile path of updateBuffer uses exactly that.
	p.ctx.BindBuffer(arrayBuffer, p.textBuffer)
	p.ctx.BufferData(arrayBuffer, points, staticDraw)
	p.logError()

	p.UpdateVertexArray(p.programs.simple, attrVertex, coordinateSize2D, coordinateSize2DWithTexture, 0)
	p.UpdateVertexArray(p.programs.simple, attrVertexTextureCoordinates, coordinateSize2D, coordinateSize2DWithTexture, coordinateSize2D)

	p.SetUniform1f(p.programs.simple, attrRadiusCorner, 0)
	p.SetUniform1f(p.programs.simple, attrAlpha, 1.0)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.glyphAtlas.texture)
	p.logError()

	p.ctx.DrawArrays(triangles, 0, len(points)/coordinateSize2DWithTexture)
	p.logError()
}
