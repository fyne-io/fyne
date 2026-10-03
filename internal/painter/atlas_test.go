package painter_test

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"slices"
	"strings"
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
	uploads, flushes, resizes int
	live                      []image.Rectangle
}

func (f *fakeAtlasTexture) FlushGlyphs() {
	f.flushes++
	f.live = nil
}

func (f *fakeAtlasTexture) UploadGlyph(img *image.RGBA, x, y int) {
	f.uploads++
	f.live = append(f.live, img.Bounds().Add(image.Pt(x, y)))
}

func (f *fakeAtlasTexture) ResizeAtlas(int) {
	f.resizes++
	f.live = nil
}

// wide is a column range no test string reaches past, for calls that want
// every glyph.
const wide = 1e9

// slot is the rectangle of a's texture a quad samples.
func slot(a *painter.GlyphAtlas, q painter.GlyphQuad) image.Rectangle {
	size := float32(a.Size())
	return image.Rect(int(q.U1*size+0.5), int(q.V1*size+0.5), int(q.U2*size+0.5), int(q.V2*size+0.5))
}

func TestGlyphAtlas_TextQuadsCachesGlyphs(t *testing.T) {
	text := canvas.NewText("Hello 12", color.White)
	var a painter.GlyphAtlas
	tex := &fakeAtlasTexture{}

	quads, ok := a.TextQuads(nil, text, fyne.NewPos(10, 20), 1, -wide, wide, tex)
	require.True(t, ok)
	assert.Len(t, quads, 7, "one quad per inked glyph; the space has none")
	assert.LessOrEqual(t, tex.uploads, 7, "at most one bitmap per inked glyph")
	for _, q := range quads {
		// Glyphs are drawn one texel to one pixel, never resampled.
		s := slot(&a, q)
		assert.Equal(t, float32(s.Dx()), q.X2-q.X1)
		assert.Equal(t, float32(s.Dy()), q.Y2-q.Y1)
	}

	uploads := tex.uploads
	again, ok := a.TextQuads(nil, text, fyne.NewPos(10, 20), 1, -wide, wide, tex)
	require.True(t, ok)
	assert.Equal(t, quads, again)
	assert.Equal(t, uploads, tex.uploads, "a string seen before uploads nothing")
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
		text.TextSize = 400

		quads, ok := a.TextQuads(nil, text, fyne.Position{}, 1, -wide, wide, tex)
		require.True(t, ok, s)
		for _, q := range quads {
			assert.True(t, slices.Contains(tex.live, slot(&a, q)), "%q has a quad sampling %v, which is not live", s, slot(&a, q))
		}
	}
	assert.Equal(t, 1, tex.resizes, "the first overflow grows the atlas, once")
	assert.Equal(t, painter.AtlasMaxSize, a.Size())
	assert.Greater(t, tex.flushes, 1, "the glyphs should have overflowed the grown atlas too")
}

// TestGlyphAtlas_GrowsBeforeThrashing draws a page's worth of styles at 4x, a
// working set that overflows AtlasSize. Once grown, drawing it again has to
// find every glyph resident: an atlas that resets every frame re-rasterises
// and re-uploads all of its text every frame.
func TestGlyphAtlas_GrowsBeforeThrashing(t *testing.T) {
	var chars []rune
	for r := rune('!'); r <= '~'; r++ {
		chars = append(chars, r)
	}
	styles := []struct {
		size  float32
		style fyne.TextStyle
	}{
		{14, fyne.TextStyle{}}, {14, fyne.TextStyle{Bold: true}}, {14, fyne.TextStyle{Italic: true}},
		{14, fyne.TextStyle{Monospace: true}}, {24, fyne.TextStyle{Bold: true}}, {18, fyne.TextStyle{Bold: true}},
		{11, fyne.TextStyle{}},
	}
	const scale = 4
	var a painter.GlyphAtlas
	tex := &fakeAtlasTexture{}
	page := func() {
		for _, s := range styles {
			text := canvas.NewText(string(chars), color.White)
			text.TextSize, text.TextStyle = s.size, s.style
			_, ok := a.TextQuads(nil, text, fyne.Position{}, scale, -wide, wide, tex)
			require.True(t, ok)
		}
	}

	// Growing part way through the page empties the atlas, so the strings drawn
	// before that are uploaded again on the next page, and only then is it warm.
	page()
	page()
	require.Equal(t, 1, tex.resizes, "the page should overflow the initial atlas")
	flushes, uploads := tex.flushes, tex.uploads
	page()
	assert.Equal(t, flushes, tex.flushes, "the grown atlas holds the page")
	assert.Equal(t, uploads, tex.uploads, "a page drawn before uploads nothing")
}

// TestGlyphAtlas_TextQuadsSkipsHiddenGlyphs checks that a string much wider
// than the visible columns yields the visible quads, placed exactly as without
// the limit, and rasterises nothing for the rest.
func TestGlyphAtlas_TextQuadsSkipsHiddenGlyphs(t *testing.T) {
	text := canvas.NewText(strings.Repeat("0123456789 ", 200), color.White)
	const minX, maxX = 500, 700

	var all painter.GlyphAtlas
	allTex := &fakeAtlasTexture{}
	full, ok := all.TextQuads(nil, text, fyne.NewPos(-20, 0), 1, -wide, wide, allTex)
	require.True(t, ok)
	var visible painter.GlyphAtlas
	visibleTex := &fakeAtlasTexture{}
	culled, ok := visible.TextQuads(nil, text, fyne.NewPos(-20, 0), 1, minX, maxX, visibleTex)
	require.True(t, ok)

	var want []image.Rectangle
	for _, q := range full {
		if q.X2 >= minX && q.X1 <= maxX {
			want = append(want, image.Rect(int(q.X1), int(q.Y1), int(q.X2), int(q.Y2)))
		}
	}
	var got []image.Rectangle
	for _, q := range culled {
		got = append(got, image.Rect(int(q.X1), int(q.Y1), int(q.X2), int(q.Y2)))
	}
	require.NotEmpty(t, want)
	for _, r := range want {
		assert.Contains(t, got, r, "a glyph in view went missing")
	}
	assert.Less(t, len(culled), len(want)+4, "only glyphs at the edges may be kept beyond the columns")
	assert.Less(t, visibleTex.uploads, allTex.uploads, "hidden glyphs are not rasterised")
}

// BenchmarkGlyphAtlas_TextQuadsLongLine is a 10,000 character line, the kind a
// log or a minified document has, seen through a 1920 pixel window: what each
// frame costs once its visible glyphs are in the atlas.
func BenchmarkGlyphAtlas_TextQuadsLongLine(b *testing.B) {
	text := canvas.NewText(strings.Repeat("The quick brown fox jumps over the lazy dog. ", 222), color.White)
	var a painter.GlyphAtlas
	tex := &fakeAtlasTexture{}
	var quads []painter.GlyphQuad
	quads, _ = a.TextQuads(quads[:0], text, fyne.Position{}, 1, 0, 1920, tex)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		quads, _ = a.TextQuads(quads[:0], text, fyne.Position{}, 1, 0, 1920, tex)
	}
	b.ReportMetric(float64(len(quads)), "quads/op")
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

// TestGlyphAtlas_TextQuadsSubpixel checks that glyphs are cached per sub-pixel
// position. A run of one letter puts it at several fractional pen positions,
// which share at most one bitmap per position; moving the string by whole
// pixels reuses them all and moves every quad by exactly that much.
func TestGlyphAtlas_TextQuadsSubpixel(t *testing.T) {
	text := canvas.NewText("llllllllllll", color.White)
	var a painter.GlyphAtlas
	tex := &fakeAtlasTexture{}

	at10, ok := a.TextQuads(nil, text, fyne.NewPos(10, 0), 1, -wide, wide, tex)
	require.True(t, ok)
	require.Len(t, at10, 12)
	assert.Greater(t, tex.uploads, 1, "the l lands at more than one sub-pixel position")
	assert.LessOrEqual(t, tex.uploads, 4, "one bitmap per position at most, and 1x has four")

	uploads := tex.uploads
	at13, ok := a.TextQuads(nil, text, fyne.NewPos(13, 0), 1, -wide, wide, tex)
	require.True(t, ok)
	assert.Equal(t, uploads, tex.uploads, "a whole pixel move reuses every bitmap")
	for i := range at10 {
		assert.Equal(t, at10[i].X1+3, at13[i].X1)
		assert.Equal(t, slot(&a, at10[i]), slot(&a, at13[i]))
	}
}

// cpuAtlasTexture keeps what is uploaded in an image, so quads can be
// composited on the CPU the way the GPU draws them.
type cpuAtlasTexture struct{ img *image.RGBA }

func (*cpuAtlasTexture) FlushGlyphs() {}

func (c *cpuAtlasTexture) UploadGlyph(img *image.RGBA, x, y int) {
	draw.Draw(c.img, img.Bounds().Add(image.Pt(x, y)), img, image.Point{}, draw.Src)
}

func (c *cpuAtlasTexture) ResizeAtlas(size int) {
	c.img = image.NewRGBA(image.Rect(0, 0, size, size))
}

// TestGlyphAtlas_MatchesDrawString composites a string's quads and compares the
// coverage with DrawString, which is what the whole-run texture path draws.
// Spacing is where the two could part: a glyph snapped to a whole pixel sits up
// to half a pixel from where the shaper put it, and over a line of text that
// shows as uneven kerning.
func TestGlyphAtlas_MatchesDrawString(t *testing.T) {
	const s = "Kerning: AVAWAY Tolerably wavy text, 0123456789"
	for _, scale := range []float32{1, 1.5, 2, 3} {
		text := canvas.NewText(s, color.White)
		tex := &cpuAtlasTexture{img: image.NewRGBA(image.Rect(0, 0, painter.AtlasSize, painter.AtlasSize))}
		var a painter.GlyphAtlas
		quads, ok := a.TextQuads(nil, text, fyne.Position{}, scale, -wide, wide, tex)
		require.True(t, ok)

		size, _ := painter.RenderedTextSize(s, text.TextSize, text.TextStyle, nil)
		bounds := image.Rect(0, 0, int(math.Ceil(float64(size.Width*scale))), int(math.Ceil(float64(size.Height*scale))))
		got := image.NewRGBA(bounds)
		for _, q := range quads {
			draw.Draw(got, image.Rect(int(q.X1), int(q.Y1), int(q.X2), int(q.Y2)), tex.img, slot(&a, q).Min, draw.Over)
		}
		want := image.NewRGBA(bounds)
		face := painter.CachedFontFace(text.TextStyle, nil, nil)
		painter.DrawString(want, s, color.White, face.Fonts, text.TextSize, scale, text.TextStyle)

		var diff, ink float64
		for i := 3; i < len(want.Pix); i += 4 {
			diff += math.Abs(float64(got.Pix[i]) - float64(want.Pix[i]))
			ink += float64(want.Pix[i])
		}
		// Measured at 6-7% with sub-pixel positions, the residue of snapping to
		// the nearest of them; whole-pixel placement is 14-32% below 3x.
		assert.Less(t, diff/ink, 0.10, "at %vx the atlas differs from DrawString by %.1f%% of the ink", scale, diff/ink*100)
	}
}
