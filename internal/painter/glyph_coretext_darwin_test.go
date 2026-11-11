//go:build darwin

package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CoreText draws a glyph with font smoothing, which thickens its strokes:
// the same glyph, at the same place, carries more ink than go-text draws.
func TestRasteriseAtlasGlyph_SmoothsThemeAndSystemFonts(t *testing.T) {
	fonts := CachedFontFace(fyne.TextStyle{}, nil, nil).Fonts
	// CJK comes from a system font, which CoreText opens from its file.
	for _, s := range []string{"m", "g", "水"} {
		var placed []PlacedGlyph
		WalkGlyphs(fonts, s, 14, 2, fyne.TextStyle{}, func(pg PlacedGlyph) { placed = append(placed, pg) })
		require.Len(t, placed, 1, s)
		pg := placed[0]
		require.NotZero(t, pg.Glyph.GlyphID, "%q has no glyph", s)

		goText := image.NewRGBA(image.Rect(0, 0, 48, 48))
		RasteriseGlyph(goText, pg, color.White, 14, 2, 8, 36, 0)
		coreText := image.NewRGBA(image.Rect(0, 0, 48, 48))
		require.True(t, rasteriseGlyphNative(coreText, pg, 14, 2, 8, 36, 0), "%q was not drawn by CoreText", s)

		thin, thinBox := ink(goText)
		thick, box := ink(coreText)
		assert.Greater(t, float64(thick), 1.15*float64(thin), "%q is not smoothed", s)
		// Smoothing grows a stroke by about a pixel; more is a misplaced glyph.
		assert.InDelta(t, thinBox.Min.X, box.Min.X, 1, "%q left", s)
		assert.InDelta(t, thinBox.Min.Y, box.Min.Y, 1, "%q top", s)
		assert.InDelta(t, thinBox.Max.X, box.Max.X, 1, "%q right", s)
		assert.InDelta(t, thinBox.Max.Y, box.Max.Y, 1, "%q bottom", s)
	}
}

// A GPU painter draws a string its atlas cannot hold with DrawStringNative,
// which must have the atlas's weight; the software painter, which Fyne's
// tests render with, keeps DrawString, the same on every platform.
func TestDrawStringNative_IsSmoothedAndDrawStringIsNot(t *testing.T) {
	face := CachedFontFace(fyne.TextStyle{}, nil, nil).Fonts
	draw1 := func(native bool) int {
		img := image.NewRGBA(image.Rect(0, 0, 200, 40))
		draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)
		if native {
			DrawStringNative(img, "Tolerably wavy", color.White, face, theme.TextSize(), 2, fyne.TextStyle{}, 0)
		} else {
			DrawString(img, "Tolerably wavy", color.White, face, theme.TextSize(), 2, fyne.TextStyle{})
		}
		n, _ := ink(img)
		return n
	}
	assert.Greater(t, float64(draw1(true)), 1.15*float64(draw1(false)))
}

func ink(img *image.RGBA) (int, image.Rectangle) {
	var box image.Rectangle
	sum := 0
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			a := img.RGBAAt(x, y).A
			if a == 0 {
				continue
			}
			sum += int(a)
			box = box.Union(image.Rect(x, y, x+1, y+1))
		}
	}
	return sum, box
}

// A theme change clears the font cache; the fonts CoreText opened from it go
// too, and drawing afterwards opens them again.
func TestClearFontCache_ReleasesCoreTextFonts(t *testing.T) {
	face := CachedFontFace(fyne.TextStyle{}, nil, nil).Fonts
	var pg PlacedGlyph
	WalkGlyphs(face, "m", 14, 2, fyne.TextStyle{}, func(p PlacedGlyph) { pg = p })
	img := image.NewRGBA(image.Rect(0, 0, 48, 48))
	require.True(t, rasteriseGlyphNative(img, pg, 14, 2, 8, 36, 0))

	ClearFontCache()
	coreText.Lock()
	assert.Empty(t, coreText.fonts)
	assert.Empty(t, coreText.descriptors)
	assert.Empty(t, coreText.sources)
	coreText.Unlock()

	face = CachedFontFace(fyne.TextStyle{}, nil, nil).Fonts
	WalkGlyphs(face, "m", 14, 2, fyne.TextStyle{}, func(p PlacedGlyph) { pg = p })
	assert.True(t, rasteriseGlyphNative(image.NewRGBA(image.Rect(0, 0, 48, 48)), pg, 14, 2, 8, 36, 0))
}
