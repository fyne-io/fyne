//go:build windows && directx

package dx

import (
	"image"
	"unsafe"

	"fyne.io/fyne/v2"
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
	tex *gpuTexture
	paint.GlyphAtlas
}

func (p *Painter) initAtlas() error {
	tex, err := p.newAtlasTexture(paint.AtlasSize)
	if err != nil {
		return err
	}
	p.atlas = glyphAtlas{tex: tex}
	return nil
}

func (p *Painter) newAtlasTexture(size int) (*gpuTexture, error) {
	tex, err := p.g.dev.CreateTexture2D(&texture2DDesc{
		Width: uint32(size), Height: uint32(size), MipLevels: 1, ArraySize: 1,
		Format: formatR8Unorm, SampleDesc: dxgiSampleDesc{Count: 1},
		Usage: usageDefault, BindFlags: bindShaderResource,
	}, nil)
	if err != nil {
		return nil, err
	}
	srv, err := p.g.dev.CreateShaderResourceView(tex)
	if err != nil {
		tex.Release()
		return nil, err
	}
	return &gpuTexture{tex: tex, srv: srv, width: uint32(size), height: uint32(size)}, nil
}

// ResizeAtlas swaps the atlas texture for an empty one of the new size, keeping
// the shared atlas's bookkeeping, which has already emptied itself.
func (p *Painter) ResizeAtlas(size int) {
	tex, err := p.newAtlasTexture(size)
	if err != nil {
		// A 2048 R8 texture is within every feature level, so this only fails
		// with the device gone, when nothing is drawn until it is rebuilt.
		fyne.LogError("directx: growing the glyph atlas", err)
		return
	}
	p.destroyTexture(p.atlas.tex)
	p.atlas.tex = tex
}

// UploadGlyph copies the alpha channel of a rasterised glyph into its atlas
// slot. Only coverage is kept: the text colour comes from the instance data, so
// the three colour channels carry nothing the shader would read.
func (p *Painter) UploadGlyph(img *image.RGBA, x, y int) {
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
}

func (a *glyphAtlas) release() {
	if a.tex == nil {
		return
	}
	releaseCOM(&a.tex.srv)
	releaseCOM(&a.tex.tex)
	*a = glyphAtlas{}
}
