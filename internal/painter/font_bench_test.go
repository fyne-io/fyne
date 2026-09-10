package painter_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/go-text/typesetting/shaping"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/painter"
)

const benchLine = "The quick brown fox jumps over the lazy dog 0123456789"

// BenchmarkWalkStringGlyphs measures shaping a line and reporting its glyphs.
func BenchmarkWalkStringGlyphs(b *testing.B) {
	face := painter.CachedFontFace(fyne.TextStyle{}, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		painter.WalkStringGlyphs(face.Fonts, benchLine, 14, fyne.TextStyle{}, 1,
			func(shaping.Output, int, float32, float32, float32, float32) {})
	}
}

// BenchmarkRenderGlyphToImage measures rasterising a single glyph, which is the
// cost of a cache miss.
func BenchmarkRenderGlyphToImage(b *testing.B) {
	face := painter.CachedFontFace(fyne.TextStyle{}, nil, nil)
	var run shaping.Output
	var idx int
	painter.WalkStringGlyphs(face.Fonts, "M", 14, fyne.TextStyle{}, 1,
		func(r shaping.Output, i int, _, _, _, _ float32) { run, idx = r, i })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		painter.RenderGlyphToImage(run, idx, 14, 1, 0)
	}
}

// BenchmarkDrawString measures the software path drawing the same line, for
// comparison with the glyph based benchmarks above.
func BenchmarkDrawString(b *testing.B) {
	face := painter.CachedFontFace(fyne.TextStyle{}, nil, nil)
	img := newBenchImage(400, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		painter.DrawString(img, benchLine, benchWhite, face.Fonts, 14, 1, fyne.TextStyle{})
	}
}

func newBenchImage(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

var benchWhite = color.White
