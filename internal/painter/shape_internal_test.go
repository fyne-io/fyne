package painter

import (
	"image"
	"image/color"
	"strconv"
	"testing"

	"github.com/go-text/typesetting/shaping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
)

// TestWalkString_MeasureAndDrawShareShaping checks that a string measured and
// then drawn, at another scale, is shaped once.
func TestWalkString_MeasureAndDrawShareShaping(t *testing.T) {
	ClearFontCache()
	face := CachedFontFace(fyne.TextStyle{}, nil, nil)

	MeasureString(face.Fonts, "Shaped once", 14, fyne.TextStyle{})
	require.Len(t, shapeCache, 1)
	key := shapeKey{text: "Shaped once", faces: face.Fonts.(*dynamicFontMap), size: float32ToFixed266(14)}
	shaped := shapeCache[key]
	require.NotNil(t, shaped)

	DrawString(image.NewRGBA(image.Rect(0, 0, 200, 40)), "Shaped once", color.White, face.Fonts, 14, 2, fyne.TextStyle{})
	assert.Len(t, shapeCache, 1, "drawing reuses the measured shaping")
	assert.Same(t, shaped, shapeCache[key])
}

// TestWalkString_CachedMatchesFresh compares measuring and drawing from the
// cache with shaping afresh, for the strings whose layout walkString does more
// than place one run: tab stops, which depend on the scale, and emoji, which
// split into runs of their own.
func TestWalkString_CachedMatchesFresh(t *testing.T) {
	for _, style := range []fyne.TextStyle{{}, {Monospace: true}} {
		for _, s := range []string{"Hello", "a\tb\tc", "\tlead and trail\t", "emoji 😀 and text", "0️⃣ keycap"} {
			name := strconv.Quote(s)
			face := CachedFontFace(style, nil, nil)
			clear(shapeCache)
			freshSize, freshAdvance := MeasureString(face.Fonts, s, 14, style)
			clear(shapeCache)
			freshImg := drawn(face.Fonts, s, style)
			require.Len(t, shapeCache, 1, name)

			size, advance := MeasureString(face.Fonts, s, 14, style)
			assert.Equal(t, freshSize, size, name)
			assert.Equal(t, freshAdvance, advance, name)
			assert.Equal(t, freshImg.Pix, drawn(face.Fonts, s, style).Pix, name)
		}
	}
}

// BenchmarkWalkString_MeasureThenDraw is a label getting a new value: the
// string is measured for layout and then drawn, at 2x.
func BenchmarkWalkString_MeasureThenDraw(b *testing.B) {
	face := CachedFontFace(fyne.TextStyle{}, nil, nil)
	img := image.NewRGBA(image.Rect(0, 0, 200, 40))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(i)
		MeasureString(face.Fonts, s, 14, fyne.TextStyle{})
		DrawString(img, s, color.White, face.Fonts, 14, 2, fyne.TextStyle{})
	}
}

// BenchmarkWalkString_ShapeTwice is the shaping part of the above: a new string
// walked for measuring and again at 2x for drawing, without rasterising, which
// is what drawing costs when the glyphs are already cached.
func BenchmarkWalkString_ShapeTwice(b *testing.B) {
	face := CachedFontFace(fyne.TextStyle{}, nil, nil)
	size := float32ToFixed266(14)
	noop := func(shaping.Output, float32, float32) {}
	var advance float32
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strconv.Itoa(i)
		walkString(face.Fonts, s, size, fyne.TextStyle{}, &advance, 1, noop)
		walkString(face.Fonts, s, size, fyne.TextStyle{}, &advance, 2, noop)
	}
}

func drawn(faces shaping.Fontmap, s string, style fyne.TextStyle) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 400, 60))
	DrawString(img, s, color.White, faces, 14, 1.5, style)
	return img
}
