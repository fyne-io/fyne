package painter_test

import (
	"image"
	"image/color"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter"
)

// fakeAtlasTexture records what a GlyphAtlas asks of its texture. live is what
// the texture holds that nothing has been drawn against yet: every upload since
// the last flush.
type fakeAtlasTexture struct {
	uploads, flushes int
	live             []image.Rectangle
}

func (f *fakeAtlasTexture) FlushGlyphs() {
	f.flushes++
	f.live = nil
}

func (f *fakeAtlasTexture) UploadGlyph(img *image.RGBA, x, y int) {
	f.uploads++
	f.live = append(f.live, img.Bounds().Add(image.Pt(x, y)))
}

// slot is the atlas rectangle a quad samples.
func slot(q painter.GlyphQuad) image.Rectangle {
	return image.Rect(int(q.U1*painter.AtlasSize+0.5), int(q.V1*painter.AtlasSize+0.5),
		int(q.U2*painter.AtlasSize+0.5), int(q.V2*painter.AtlasSize+0.5))
}

func TestGlyphAtlas_TextQuadsCachesGlyphs(t *testing.T) {
	text := canvas.NewText("Hello 12", color.White)
	var a painter.GlyphAtlas
	tex := &fakeAtlasTexture{}

	quads, ok := a.TextQuads(nil, text, fyne.NewPos(10, 20), 1, tex)
	require.True(t, ok)
	assert.Len(t, quads, 7, "one quad per inked glyph; the space has none")
	assert.Equal(t, 6, tex.uploads, "the second l reuses the first")
	for _, q := range quads {
		// Glyphs are drawn one texel to one pixel, never resampled.
		s := slot(q)
		assert.Equal(t, float32(s.Dx()), q.X2-q.X1)
		assert.Equal(t, float32(s.Dy()), q.Y2-q.Y1)
	}

	again, ok := a.TextQuads(nil, text, fyne.NewPos(10, 20), 1, tex)
	require.True(t, ok)
	assert.Equal(t, quads, again)
	assert.Equal(t, 6, tex.uploads, "a string seen before uploads nothing")
}

// TestGlyphAtlas_TextQuadsSurvivesReset fills the atlas with glyphs big enough
// that it resets part way through a string. Slots resolved before the reset
// are stale by the end of it, so every quad handed out has to point at a glyph
// uploaded since the last flush.
func TestGlyphAtlas_TextQuadsSurvivesReset(t *testing.T) {
	var a painter.GlyphAtlas
	tex := &fakeAtlasTexture{}
	for _, s := range []string{"abcdefghijklmnop", "ABCDEFGHIJKLMNOP", "qrstuvwxyz012345", "QRSTUVWXYZ6789&?"} {
		text := canvas.NewText(s, color.White)
		text.TextSize = 250

		quads, ok := a.TextQuads(nil, text, fyne.Position{}, 1, tex)
		require.True(t, ok, s)
		for _, q := range quads {
			assert.True(t, slices.Contains(tex.live, slot(q)), "%q has a quad sampling %v, which is not live", s, slot(q))
		}
	}
	assert.Positive(t, tex.flushes, "the glyphs should have overflowed the atlas")
}

func TestGlyphAtlas_White(t *testing.T) {
	var a painter.GlyphAtlas
	tex := &fakeAtlasTexture{}

	u, v, ok := a.White(tex)
	require.True(t, ok)
	require.Len(t, tex.live, 1)
	block := tex.live[0]
	assert.Equal(t, 3, block.Dx())
	assert.Equal(t, 3, block.Dy())
	assert.Equal(t, (float32(block.Min.X)+1.5)/painter.AtlasSize, u, "the centre texel of the block")
	assert.Equal(t, (float32(block.Min.Y)+1.5)/painter.AtlasSize, v)

	_, _, _ = a.White(tex)
	assert.Equal(t, 1, tex.uploads, "the block is packed once")
}
