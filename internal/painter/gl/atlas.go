package gl

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	paint "fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/theme"
)

// This file is the GL half of the glyph atlas in internal/painter's atlas.go.
// Text and plain rectangles are queued as quads out of one atlas texture and
// drawn together, so a table of cells and their labels costs one draw call
// instead of two per cell.

// glyphVertexFloats is one batch vertex: position in clip space, atlas texture
// coordinate, and a straight-alpha colour.
const glyphVertexFloats = 2 + 2 + 4

// useProgram switches to the program of an object that draws on its own. Every
// such draw comes through here first, which is what keeps the batch in order:
// quads queued before the object reach the framebuffer before it does.
func (p *painter) useProgram(prog Program) {
	p.FlushGlyphs()
	p.ctx.UseProgram(prog)
}

// FlushGlyphs draws every queued quad in one draw call. It has to run before
// any other draw, before the scissor rectangle changes, before atlas space is
// reused and at the end of the frame; each of those would otherwise reorder,
// mis-clip or mis-sample the queue.
func (p *painter) FlushGlyphs() {
	vertices := len(p.glyphPending) / glyphVertexFloats
	if vertices == 0 {
		return
	}

	prog := p.programs.glyph
	p.ctx.UseProgram(prog.ref)
	p.ctx.BindBuffer(arrayBuffer, prog.buff)
	p.ctx.BufferData(arrayBuffer, p.glyphPending, staticDraw)
	p.UpdateVertexArray(prog, attrVertex, 2, glyphVertexFloats, 0)
	p.UpdateVertexArray(prog, attrVertexTextureCoordinates, 2, glyphVertexFloats, 2)
	p.UpdateVertexArray(prog, attrVertexColor, 4, glyphVertexFloats, 4)
	p.ctx.BlendFunc(srcAlpha, oneMinusSrcAlpha)
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.atlasTex)
	p.ctx.DrawArrays(triangles, 0, vertices)
	p.logError()

	if glDebug {
		p.countBatch(vertices / 6)
	}
	p.glyphPending = p.glyphPending[:0]
}

// UploadGlyph copies a bitmap from the shared atlas into the texture. The atlas
// allocates each one fresh, so its rows are tightly packed as GL expects.
func (p *painter) UploadGlyph(img *image.RGBA, x, y int) {
	p.ctx.ActiveTexture(texture0)
	p.ctx.BindTexture(texture2D, p.atlasTex)
	p.ctx.TexSubImage2D(texture2D, 0, x, y, img.Rect.Dx(), img.Rect.Dy(), colorFormatRGBA, unsignedByte, img.Pix)
	p.logError()
}

// ensureAtlas allocates the atlas texture the first time anything is batched.
//
// ponytail: RGBA, although only alpha is read. A single-channel texture would
// be a quarter of the 4MB, but needs format constants and an unpack alignment
// call on every GL binding; this needs neither.
func (p *painter) ensureAtlas() {
	if p.atlasTexValid {
		return
	}
	p.atlasTex = p.newTexture(canvas.ImageScaleSmooth)
	// Cleared rather than left undefined: nothing samples outside the glyphs,
	// but a driver that hands back stale memory would make any bug in that
	// promise show up as garbage instead of nothing.
	p.ctx.TexImage2D(texture2D, 0, paint.AtlasSize, paint.AtlasSize, colorFormatRGBA, unsignedByte,
		make([]uint8, paint.AtlasSize*paint.AtlasSize*4))
	p.logError()
	p.atlasTexValid = true
}

// drawTextFromAtlas queues one quad per glyph of the run, reporting false when
// the string has to go down the whole-run texture path instead.
func (p *painter) drawTextFromAtlas(text *canvas.Text, pos fyne.Position, frame fyne.Size) bool {
	if !glyphAtlasSupported {
		return false
	}
	p.ensureAtlas()
	var ok bool
	p.quadScratch, ok = p.atlas.TextQuads(p.quadScratch[:0], text, pos, p.pixScale, p)
	if !ok {
		return false
	}

	col := text.Color
	if col == nil {
		col = theme.Color(theme.ColorNameForeground)
	}
	r, g, b, a := getFragmentColor(col)
	fw, fh := p.scaleFrameSize(frame)
	for _, q := range p.quadScratch {
		p.pushQuad(q, [4]float32{r, g, b, a}, fw, fh)
	}
	return true
}

// pushSolidRect queues a plain filled rectangle on the batch instead of drawing
// it on its own. Only the square-cornered, unstroked, unshadowed case
// qualifies: for it the rectangle shader fills exactly the snapped pixels with
// no antialiasing, which a quad on the same pixels reproduces. A nonzero stroke
// always goes to the shader, which draws the rim in the stroke colour even when
// that colour is transparent. Reports false when the shader is needed.
func (p *painter) pushSolidRect(r *canvas.Rectangle, pos fyne.Position, frame fyne.Size) bool {
	if !glyphAtlasSupported || r.Aspect != 0 || r.StrokeWidth != 0 || paint.IsShadowVisible(r.Shadow) {
		return false
	}
	if r.FillColor == nil || r.FillColor == color.Transparent {
		return true // nothing to draw, and no stroke or shadow either
	}
	p.ensureAtlas()
	u, v, ok := p.atlas.White(p)
	if !ok {
		return false
	}

	// The snapping vecRectCoords applies, without the antialiasing skirt it adds
	// for the shader's quad (which the shader then clips back to bounds).
	size := r.Size()
	x := roundToPixel(pos.X, p.pixScale)
	y := roundToPixel(pos.Y, p.pixScale)
	w := roundToPixel(size.Width, p.pixScale)
	h := roundToPixel(size.Height, p.pixScale)
	x1, x2, y1, y2 := p.scaleRectCoords(x, x+w, y, y+h)
	fw, fh := p.scaleFrameSize(frame)
	cr, cg, cb, ca := getFragmentColor(r.FillColor)
	p.pushQuad(paint.GlyphQuad{X1: x1, Y1: y1, X2: x2, Y2: y2, U1: u, V1: v, U2: u, V2: v},
		[4]float32{cr, cg, cb, ca}, fw, fh)
	return true
}

// pushQuad queues q, given in destination pixels, in a straight-alpha colour.
// It is two triangles of six vertices that each carry the colour: GL 2.1 and
// ES 2 have no instancing, and the whole queue is still a single draw.
func (p *painter) pushQuad(q paint.GlyphQuad, col [4]float32, fw, fh float32) {
	// Pixels to clip space: x doubles and shifts, y additionally flips.
	x1, y1 := q.X1/fw*2-1, 1-q.Y1/fh*2
	x2, y2 := q.X2/fw*2-1, 1-q.Y2/fh*2
	r, g, b, a := col[0], col[1], col[2], col[3]
	p.glyphPending = append(p.glyphPending,
		x1, y1, q.U1, q.V1, r, g, b, a,
		x2, y1, q.U2, q.V1, r, g, b, a,
		x1, y2, q.U1, q.V2, r, g, b, a,
		x1, y2, q.U1, q.V2, r, g, b, a,
		x2, y1, q.U2, q.V1, r, g, b, a,
		x2, y2, q.U2, q.V2, r, g, b, a,
	)
}
