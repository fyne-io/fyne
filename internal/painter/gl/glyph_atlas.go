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

// drawGlyphQuad renders one glyph quad by sampling the sub-region [entry.x,
// entry.y, entry.x+entry.w, entry.y+entry.h] of the shared GPU atlas texture.
// pos and size are in logical (fyne) coordinate space; frame is the canvas size.
func (p *painter) drawGlyphQuad(entry glyphAtlasEntry, pos fyne.Position, size fyne.Size, frame fyne.Size) {
	points, insets, inner := p.rectCoords(size, pos, frame, canvas.ImageFillStretch, 1, 0)

	// Replace the full-texture UVs that rectCoords produced with the glyph's
	// sub-region of the atlas. rectCoords lays out vertices in the order
	// top-left, bottom-left, top-right, bottom-right, and v runs opposite to
	// screen y, so the two "top" vertices take vMax.
	af := float32(p.glyphAtlas.texSize)
	uMin := float32(entry.x) / af
	vMin := float32(entry.y) / af
	uMax := float32(entry.x+entry.w) / af
	vMax := float32(entry.y+entry.h) / af

	uv := [vertexCountRectangle][2]float32{
		{uMin, vMax}, // top left
		{uMin, vMin}, // bottom left
		{uMax, vMax}, // top right
		{uMax, vMin}, // bottom right
	}
	for i, c := range uv {
		texCoord := i*coordinateSize2DWithTexture + coordinateSize2D
		points[texCoord], points[texCoord+1] = c[0], c[1]
	}

	p.ctx.UseProgram(p.programs.simple.ref)
	p.updateBuffer(p.programs.simple.buff, points[:])
	p.UpdateVertexArray(p.programs.simple, attrVertex, coordinateSize2D, coordinateSize2DWithTexture, 0)
	p.UpdateVertexArray(p.programs.simple, attrVertexTextureCoordinates, coordinateSize2D, coordinateSize2DWithTexture, coordinateSize2D)

	p.SetUniform1f(p.programs.simple, attrRadiusCorner, 0)
	p.SetUniform2f(p.programs.simple, attrSize, inner.Width*p.pixScale, inner.Height*p.pixScale)
	p.SetUniform4f(p.programs.simple, attrInset, insets[0], insets[1], insets[2], insets[3])
	p.SetUniform1f(p.programs.simple, attrAlpha, 1.0)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.glyphAtlas.texture)
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, vertexCountRectangle)
	p.logError()
}
