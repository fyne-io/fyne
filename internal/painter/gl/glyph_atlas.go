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
	x, y int // top-left position in the atlas texture
	w, h int // dimensions in pixels
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

	glyphImg := paint.RenderGlyphToImage(run, idx, fontSize, scale, col)
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

	entry := glyphAtlasEntry{x: atlasX, y: atlasY, w: w, h: h}
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
	points, insets := p.rectCoords(size, pos, frame, canvas.ImageFillStretch, 1, 0)
	inner, _ := rectInnerCoords(size, pos, canvas.ImageFillStretch, 1)

	// Override UV coordinates to sample only the glyph's region in the atlas.
	// rectCoords vertex layout (each 5 floats: x,y,z,u,v):
	//   vertex 0 [0..4]  – screen bottom-left: points[3]=u, points[4]=v
	//   vertex 1 [5..9]  – screen top-left:    points[8]=u, points[9]=v
	//   vertex 2 [10..14]– screen bottom-right: points[13]=u, points[14]=v
	//   vertex 3 [15..19]– screen top-right:    points[18]=u, points[19]=v
	af := float32(p.glyphAtlas.texSize)
	uMin := float32(entry.x) / af
	vMin := float32(entry.y) / af
	uMax := float32(entry.x+entry.w) / af
	vMax := float32(entry.y+entry.h) / af

	points[3], points[4] = uMin, vMax
	points[8], points[9] = uMin, vMin
	points[13], points[14] = uMax, vMax
	points[18], points[19] = uMax, vMin

	p.ctx.UseProgram(p.program.ref)
	p.updateBuffer(p.program.buff, points)
	p.UpdateVertexArray(p.program, "vert", 3, 5, 0)
	p.UpdateVertexArray(p.program, "vertTexCoord", 2, 5, 3)

	p.SetUniform1f(p.program, "cornerRadius", 0)
	p.SetUniform2f(p.program, "size", inner.Width*p.pixScale, inner.Height*p.pixScale)
	p.SetUniform4f(p.program, "inset", insets[0], insets[1], insets[2], insets[3])
	p.SetUniform1f(p.program, "alpha", 1.0)

	p.ctx.BlendFunc(one, oneMinusSrcAlpha)
	p.logError()

	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.glyphAtlas.texture)
	p.logError()

	p.ctx.DrawArrays(triangleStrip, 0, 4)
	p.logError()
}
