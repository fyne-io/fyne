//go:build !windows || !ci

package gl

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	paint "fyne.io/fyne/v2/internal/painter"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

// Text measurement reads the theme through the current app, and the app that
// gl_test.go starts is not built under the ci tag, so start one here.
func TestMain(m *testing.M) {
	test.NewApp()
	m.Run()
}

// testPainter builds a painter driving a fake GL context, so the text drawing
// path can run without a window.
func testPainter() (*painter, *fakeContext) {
	ctx := &fakeContext{}
	p := &painter{pixScale: 1, ctx: ctx}
	p.programs = &programs{}
	// Every program needs its lookup maps before it can take a uniform.
	for _, prog := range []*programState{
		&p.programs.arbitraryPolygon, &p.programs.arc, &p.programs.bezierCurve,
		&p.programs.blur, &p.programs.ellipse, &p.programs.line, &p.programs.polygon,
		&p.programs.rectangle, &p.programs.roundRectangle, &p.programs.simple,
		&p.programs.text,
	} {
		prog.uniforms = make(map[string]*uniformState)
		prog.attributes = make(map[string]Attribute)
	}
	return p, ctx
}

func TestSubpixelPhasesFor(t *testing.T) {
	assert.Equal(t, subpixelPhases, subpixelPhasesFor(14), "small text keeps every phase")
	assert.Equal(t, 2, subpixelPhasesFor(42), "a phone-sized glyph halves them")
	assert.Equal(t, 1, subpixelPhasesFor(200), "a whole pixel is a small part of a large stroke")

	// Never zero, or a glyph would divide by it when picking its bitmap.
	for em := float32(1); em < 300; em += 1.5 {
		assert.Positive(t, subpixelPhasesFor(em), "phases must stay positive at em %v", em)
	}
}

func TestEnsureGlyphAtlas(t *testing.T) {
	p, ctx := testPainter()
	p.maxTextureSize = 4096
	p.ensureGlyphAtlas()

	require.NotNil(t, p.glyphAtlas)
	require.NotNil(t, p.glyphColourAtlas)
	assert.Equal(t, glyphAtlasTexSize, p.glyphAtlas.texSize)
	assert.Equal(t, 1, p.glyphAtlas.bpp, "coverage atlas holds one byte per pixel")
	assert.Equal(t, 4, p.glyphColourAtlas.bpp, "colour atlas holds four")
	assert.NotZero(t, ctx.imageBytes, "both atlas textures should have been allocated")

	// A second call must not rebuild them.
	before := ctx.imageBytes
	p.ensureGlyphAtlas()
	assert.Equal(t, before, ctx.imageBytes)
}

func TestEnsureGlyphAtlasRespectsMaxTexture(t *testing.T) {
	p, _ := testPainter()
	p.maxTextureSize = 512 // an older GPU
	p.ensureGlyphAtlas()

	assert.Equal(t, 512, p.glyphAtlas.texSize, "atlas must not exceed what the GPU accepts")
	assert.Equal(t, 512, p.glyphColourAtlas.texSize)
}

func newTestText(s string) *canvas.Text {
	t := canvas.NewText(s, color.White)
	t.TextSize = 14
	return t
}

func TestGlyphGeometryCachesByContent(t *testing.T) {
	p, ctx := testPainter()
	p.maxTextureSize = 4096
	p.ensureGlyphAtlas()
	face := paint.CachedFontFace(fyne.TextStyle{}, nil, nil)

	first := p.glyphGeometry(newTestText("hello"), face)
	require.NotNil(t, first)
	assert.Positive(t, first.vertices, "a drawn string should produce quads")
	uploads := ctx.bufferUploads

	// The same words on a different object reuse the geometry and upload nothing.
	again := p.glyphGeometry(newTestText("hello"), face)
	assert.Same(t, first, again, "same text should share cached geometry")
	assert.Equal(t, uploads, ctx.bufferUploads, "a cache hit must not upload")

	// Different words build their own.
	other := p.glyphGeometry(newTestText("world"), face)
	assert.NotSame(t, first, other)
	assert.Greater(t, ctx.bufferUploads, uploads, "new text has to be uploaded once")
}

func TestGlyphGeometryInvalidatedByScale(t *testing.T) {
	p, _ := testPainter()
	p.maxTextureSize = 4096
	p.ensureGlyphAtlas()
	face := paint.CachedFontFace(fyne.TextStyle{}, nil, nil)

	text := newTestText("scale me")
	first := p.glyphGeometry(text, face)
	require.NotNil(t, first)

	p.pixScale = 2
	assert.False(t, first.usable(p.glyphAtlas.generation, p.pixScale),
		"geometry rasterised for one scale cannot be drawn at another")
}

func TestDrawGlyphBatch(t *testing.T) {
	p, ctx := testPainter()
	p.maxTextureSize = 4096
	p.ensureGlyphAtlas()
	face := paint.CachedFontFace(fyne.TextStyle{}, nil, nil)

	cached := p.glyphGeometry(newTestText("draw"), face)
	require.NotNil(t, cached)

	before := ctx.draws
	p.drawGlyphBatch(cached, color.White, fyne.NewPos(10, 20), fyne.NewSize(300, 200))
	assert.Equal(t, before+1, ctx.draws, "a string of plain glyphs draws in one call")
	assert.Equal(t, cached.vertices, ctx.vertices, "every quad should be drawn")
}

func TestDrawGlyphBatchEmpty(t *testing.T) {
	p, ctx := testPainter()
	p.maxTextureSize = 4096
	p.ensureGlyphAtlas()

	before := ctx.draws
	p.drawGlyphBatch(&textVertices{}, color.White, fyne.NewPos(0, 0), fyne.NewSize(100, 100))
	assert.Equal(t, before, ctx.draws, "nothing to draw means no draw call")
}

// requireEmojiFont skips a test on builds made with -tags no_emoji, where there
// is no colour font to put in the colour atlas.
func requireEmojiFont(t *testing.T) {
	t.Helper()
	if theme.DefaultEmojiFont() == nil {
		t.Skip("built without an emoji font")
	}
}

func TestGlyphEntryUsesColourAtlasForEmoji(t *testing.T) {
	requireEmojiFont(t)
	p, _ := testPainter()
	p.maxTextureSize = 4096
	p.ensureGlyphAtlas()

	plain := p.glyphGeometry(newTestText("A"), paint.CachedFontFace(fyne.TextStyle{}, nil, nil))
	require.NotNil(t, plain)
	assert.Zero(t, plain.colourVertices, "an outline glyph is coverage only")
	assert.NotEmpty(t, p.glyphAtlas.entries, "it belongs in the coverage atlas")

	emoji := p.glyphGeometry(newTestText("\U0001F600"), paint.CachedFontFace(fyne.TextStyle{}, nil, nil))
	require.NotNil(t, emoji)
	assert.Positive(t, emoji.colourVertices, "an emoji keeps its own colours")
	assert.NotEmpty(t, p.glyphColourAtlas.entries, "so it belongs in the colour atlas")
}

func TestDrawGlyphBatchSplitsByAtlas(t *testing.T) {
	requireEmojiFont(t)
	p, ctx := testPainter()
	p.maxTextureSize = 4096
	p.ensureGlyphAtlas()

	mixed := p.glyphGeometry(newTestText("hi \U0001F600"), paint.CachedFontFace(fyne.TextStyle{}, nil, nil))
	require.NotNil(t, mixed)
	require.Positive(t, mixed.vertices)
	require.Positive(t, mixed.colourVertices)

	before := ctx.draws
	p.drawGlyphBatch(mixed, color.White, fyne.NewPos(0, 0), fyne.NewSize(300, 200))
	assert.Equal(t, before+2, ctx.draws, "coverage and colour glyphs come from different textures")
}

// TestGlyphGeometrySurvivesAtlasReset drives the path where the atlas runs out
// of room part way through a string and the caller has to start over.
func TestGlyphGeometrySurvivesAtlasReset(t *testing.T) {
	p, _ := testPainter()
	p.maxTextureSize = 64
	p.ensureGlyphAtlas() // a 64px atlas holds very little
	face := paint.CachedFontFace(fyne.TextStyle{}, nil, nil)

	start := p.glyphAtlas.generation
	cached := p.glyphGeometry(newTestText("the quick brown fox jumps"), face)

	require.NotNil(t, cached, "a full atlas must still return usable geometry")
	assert.GreaterOrEqual(t, p.glyphAtlas.generation, start)
	assert.Equal(t, p.glyphAtlas.generation, cached.generation,
		"cached geometry must record the layout it was built against")
}

func TestDrawTextSkipsNothingToDraw(t *testing.T) {
	p, ctx := testPainter()
	p.maxTextureSize = 4096
	frame := fyne.NewSize(300, 200)

	for name, text := range map[string]*canvas.Text{
		"empty string": newTestText(""),
		"lone space":   newTestText(" "),
	} {
		t.Run(name, func(t *testing.T) {
			before := ctx.draws
			p.drawText(text, fyne.NewPos(0, 0), frame, nil)
			assert.Equal(t, before, ctx.draws, "%s has no glyphs worth a draw call", name)
		})
	}
}

func TestDrawText(t *testing.T) {
	p, ctx := testPainter()
	p.maxTextureSize = 4096

	text := newTestText("hello")
	text.Resize(text.MinSize())

	p.drawText(text, fyne.NewPos(5, 5), fyne.NewSize(300, 200), nil)
	assert.Equal(t, 1, ctx.draws, "a plain string is one draw call")
	require.NotNil(t, p.glyphAtlas, "drawing text should have built the atlas")

	// Drawing it again reuses the cached geometry, so nothing new is uploaded.
	uploads := ctx.bufferUploads
	p.drawText(text, fyne.NewPos(5, 5), fyne.NewSize(300, 200), nil)
	assert.Equal(t, 2, ctx.draws)
	assert.Equal(t, uploads, ctx.bufferUploads, "a redraw must not re-upload geometry")
}

func TestDrawTextAlignment(t *testing.T) {
	frame := fyne.NewSize(300, 200)

	// Each alignment puts the same string at a different x, so the origin
	// uniform has to differ while the geometry stays shared.
	origins := map[fyne.TextAlign]float32{}
	for _, align := range []fyne.TextAlign{fyne.TextAlignLeading, fyne.TextAlignCenter, fyne.TextAlignTrailing} {
		p, _ := testPainter()
		p.maxTextureSize = 4096

		text := newTestText("hello")
		text.Alignment = align
		text.Resize(fyne.NewSize(200, text.MinSize().Height)) // room to move within

		p.drawText(text, fyne.NewPos(0, 0), frame, nil)

		u := p.programs.text.uniforms[attrOrigin]
		require.NotNil(t, u, "drawing should have set the origin uniform")
		origins[align] = u.prev[0]
	}

	assert.Less(t, origins[fyne.TextAlignLeading], origins[fyne.TextAlignCenter],
		"centred text should start further right than leading")
	assert.Less(t, origins[fyne.TextAlignCenter], origins[fyne.TextAlignTrailing],
		"trailing text should start further right again")
}

func TestDrawTextDecorations(t *testing.T) {
	frame := fyne.NewSize(300, 200)

	for name, style := range map[string]fyne.TextStyle{
		"underline":     {Underline: true},
		"strikethrough": {Strikethrough: true},
		"both":          {Underline: true, Strikethrough: true},
	} {
		t.Run(name, func(t *testing.T) {
			p, ctx := testPainter()
			p.maxTextureSize = 4096

			text := newTestText("hello")
			text.TextStyle = style
			text.Resize(text.MinSize())

			p.drawText(text, fyne.NewPos(0, 0), frame, nil)

			want := 2 // glyphs plus one rule
			if style.Underline && style.Strikethrough {
				want = 3
			}
			assert.Equal(t, want, ctx.draws, "%s should add a line per decoration", name)
		})
	}
}

// A decorated space has no glyphs but must still draw its rule.
func TestDrawTextDecoratedSpace(t *testing.T) {
	p, ctx := testPainter()
	p.maxTextureSize = 4096

	text := newTestText(" ")
	text.TextStyle = fyne.TextStyle{Underline: true}
	text.Resize(text.MinSize())

	p.drawText(text, fyne.NewPos(0, 0), fyne.NewSize(300, 200), nil)
	assert.Positive(t, ctx.draws, "an underlined space still has a line to draw")
}
