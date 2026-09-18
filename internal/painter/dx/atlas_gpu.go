//go:build windows && directx

package dx

import (
	"image"
	"image/color"
	"math"
	"unsafe"

	"github.com/go-text/typesetting/font"

	paint "fyne.io/fyne/v2/internal/painter"
)

// glyphAtlas keeps every glyph the canvas has drawn in one single-channel
// texture, so a string costs one quad per glyph out of a shared texture instead
// of a whole texture of its own.
//
// That replaces three costs at once: the per-string rasterisation (a label whose
// text changes every frame re-rasterises nothing but glyphs it has already seen),
// the video memory of thousands of run textures, and the draw call per text
// object - one texture means consecutive text batches into a single draw.
type glyphAtlas struct {
	tex     *gpuTexture
	packer  atlasPacker
	entries map[glyphKey]glyphEntry

	// whiteU/whiteV point at the centre of a solid 3x3 block, which lets a
	// plain filled rectangle ride the glyph batch as a quad with a constant
	// UV: coverage samples 1 and the instance colour is the fill. whiteOK is
	// false until the block is packed, and again after every atlas reset.
	whiteU, whiteV float32
	whiteOK        bool
}

func (p *Painter) initAtlas() error {
	tex, err := p.g.dev.CreateTexture2D(&texture2DDesc{
		Width: atlasSize, Height: atlasSize, MipLevels: 1, ArraySize: 1,
		Format: formatR8Unorm, SampleDesc: dxgiSampleDesc{Count: 1},
		Usage: usageDefault, BindFlags: bindShaderResource,
	}, nil)
	if err != nil {
		return err
	}
	srv, err := p.g.dev.CreateShaderResourceView(tex)
	if err != nil {
		tex.Release()
		return err
	}
	p.atlas = glyphAtlas{
		tex:     &gpuTexture{tex: tex, srv: srv, width: atlasSize, height: atlasSize},
		entries: make(map[glyphKey]glyphEntry),
	}
	return nil
}

// glyphEntry returns where a glyph lives in the atlas, rasterising and
// uploading it on first use. The second result is false when the glyph could
// not be placed at all, which sends the caller to the whole-run fallback.
func (p *Painter) glyphEntry(pg paint.PlacedGlyph, fontSize float32) (glyphEntry, bool) {
	key := newGlyphKey(pg.Face, pg.Glyph.GlyphID, fontSize, p.pixScale)
	if e, ok := p.atlas.entries[key]; ok {
		return e, !e.unsupported
	}

	// Colour glyphs - emoji, in practice - carry bitmap or SVG data a coverage
	// atlas cannot represent. The answer is cached along with everything else,
	// because GlyphData is far too expensive to ask once a frame.
	if _, ok := pg.Face.GlyphData(pg.Glyph.GlyphID).(font.GlyphOutline); !ok {
		e := glyphEntry{unsupported: true}
		p.atlas.entries[key] = e
		return e, false
	}

	// Ink bounds in destination pixels, relative to the pen position and
	// baseline. Height is measured downwards from YBearing, so it is negative.
	relX := fixedToPixels(int32(pg.Glyph.XBearing), p.pixScale)
	relY := fixedToPixels(int32(pg.Glyph.YBearing), p.pixScale)
	inkW := fixedToPixels(int32(pg.Glyph.Width), p.pixScale)
	inkH := -fixedToPixels(int32(pg.Glyph.Height), p.pixScale)
	if inkW <= 0 || inkH <= 0 { // a space, or anything else with no ink
		e := glyphEntry{}
		p.atlas.entries[key] = e
		return e, true
	}

	// Rasterise at a whole-pixel origin chosen so the ink, plus a pixel of room
	// for the antialiased edge, lands inside the bitmap.
	originX := 1 - int(math.Floor(float64(relX)))
	originY := 1 + int(math.Ceil(float64(relY)))
	w := int(math.Ceil(float64(inkW))) + 3
	h := int(math.Ceil(float64(inkH))) + 3

	x, y, ok := p.atlas.packer.add(w, h)
	if !ok {
		// Full. Pending quads still point at the old contents, so they have to
		// be drawn before anything is overwritten.
		p.flushGlyphs()
		p.atlas.packer.reset()
		clear(p.atlas.entries)
		p.atlas.whiteOK = false
		if x, y, ok = p.atlas.packer.add(w, h); !ok {
			return glyphEntry{}, false // larger than the atlas itself
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	paint.RasteriseGlyph(img, pg, color.White, fontSize, p.pixScale, originX, originY)
	if !p.uploadGlyph(img, x, y) {
		return glyphEntry{}, false
	}

	e := glyphEntry{x: x, y: y, width: w, height: h, bearX: -originX, bearY: -originY}
	p.atlas.entries[key] = e
	return e, true
}

// uploadGlyph copies the alpha channel of a rasterised glyph into its atlas
// slot. Only coverage is kept: the text colour comes from the instance data, so
// the three colour channels carry nothing the shader would read.
func (p *Painter) uploadGlyph(img *image.RGBA, x, y int) bool {
	w, h := img.Rect.Dx(), img.Rect.Dy()
	cov := make([]byte, w*h)
	for row := 0; row < h; row++ {
		src := img.Pix[row*img.Stride:]
		dst := cov[row*w:]
		for col := 0; col < w; col++ {
			dst[col] = src[col*4+3]
		}
	}

	region := box{
		Left: uint32(x), Top: uint32(y), Front: 0,
		Right: uint32(x + w), Bottom: uint32(y + h), Back: 1,
	}
	p.g.ctx.UpdateTextureRegion(unsafe.Pointer(p.atlas.tex.tex),
		unsafe.Pointer(&cov[0]), uint32(w), &region)
	return true
}

// ensureAtlasWhite packs the solid block rectangle quads sample. Sampling the
// centre texel of a 3x3 keeps every bilinear neighbour inside the block, so
// the empty pad the packer leaves around it can never bleed in.
func (p *Painter) ensureAtlasWhite() bool {
	if p.atlas.whiteOK {
		return true
	}
	x, y, ok := p.atlas.packer.add(3, 3)
	if !ok {
		// As in glyphEntry: draw what still points at the old contents first.
		p.flushGlyphs()
		p.atlas.packer.reset()
		clear(p.atlas.entries)
		if x, y, ok = p.atlas.packer.add(3, 3); !ok {
			return false
		}
	}
	cov := [9]byte{255, 255, 255, 255, 255, 255, 255, 255, 255}
	region := box{
		Left: uint32(x), Top: uint32(y), Front: 0,
		Right: uint32(x + 3), Bottom: uint32(y + 3), Back: 1,
	}
	p.g.ctx.UpdateTextureRegion(unsafe.Pointer(p.atlas.tex.tex),
		unsafe.Pointer(&cov[0]), 3, &region)

	p.atlas.whiteU = (float32(x) + 1.5) / atlasSize
	p.atlas.whiteV = (float32(y) + 1.5) / atlasSize
	p.atlas.whiteOK = true
	return true
}

func (a *glyphAtlas) release() {
	if a.tex == nil {
		return
	}
	releaseCOM(&a.tex.srv)
	releaseCOM(&a.tex.tex)
	a.tex = nil
	a.entries = nil
	a.packer.reset()
}

// fixedToPixels converts a 26.6 fixed point font metric to destination pixels.
func fixedToPixels(v int32, scale float32) float32 {
	return float32(v) / 64 * scale
}

// shapeable reports whether a shaped string is a candidate for the atlas at
// all. Glyph 0 means the shaper found nothing, which DrawString renders by
// substituting a replacement character from another face - a second shaping
// path this deliberately does not grow. Whether the glyphs themselves fit is
// glyphEntry's answer, and a cached one.
//
// ponytail: whole-string fallback rather than a second colour atlas beside this
// one. Mixed emoji and text then costs one texture per such label, which is
// what every label cost before this existed; a colour atlas is only worth
// building for an app that is mostly emoji.
func shapeable(glyphs []paint.PlacedGlyph) bool {
	for _, pg := range glyphs {
		if pg.Glyph.GlyphID == 0 || pg.Face == nil {
			return false
		}
	}
	return true
}
