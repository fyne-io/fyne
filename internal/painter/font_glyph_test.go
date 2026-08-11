package painter_test

import (
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/go-text/typesetting/shaping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/painter"
)

type walkedGlyph struct {
	run         shaping.Output
	idx         int
	penX, baseY float32
	xOff, yOff  float32
}

// walkGlyphs collects every glyph of s so tests can assert over the whole run.
func walkGlyphs(t *testing.T, s string, size float32, scale float32) []walkedGlyph {
	t.Helper()

	face := painter.CachedFontFace(fyne.TextStyle{}, nil, nil)
	require.NotNil(t, face, "bundled font should load")

	var got []walkedGlyph
	painter.WalkStringGlyphs(face.Fonts, s, size, fyne.TextStyle{}, scale,
		func(run shaping.Output, idx int, penX, baseY, xOff, yOff float32) {
			got = append(got, walkedGlyph{run: run, idx: idx, penX: penX, baseY: baseY, xOff: xOff, yOff: yOff})
		})
	return got
}

func TestWalkStringGlyphs(t *testing.T) {
	glyphs := walkGlyphs(t, "Hello", 14, 1)
	require.Len(t, glyphs, 5, "one callback per glyph of a simple ASCII string")

	// The pen only ever moves forwards, and every glyph of a single-face string
	// shares the one baseline.
	for i, g := range glyphs {
		if i > 0 {
			assert.Greater(t, g.penX, glyphs[i-1].penX, "pen should advance at glyph %d", i)
		}
		assert.Equal(t, glyphs[0].baseY, g.baseY, "baseline should be shared at glyph %d", i)
		assert.Positive(t, g.baseY, "baseline should sit below the top of the line")
	}
}

func TestWalkStringGlyphs_Empty(t *testing.T) {
	assert.Empty(t, walkGlyphs(t, "", 14, 1), "an empty string has no glyphs")
}

// TestWalkStringGlyphs_PositionsAreExact guards the regression that made letter
// spacing uneven: positions used to be rounded to whole pixels inside the walk,
// which discarded the fractional part of every advance. Landing on the pixel
// grid is the caller's decision, so the positions handed out here must be exact.
func TestWalkStringGlyphs_PositionsAreExact(t *testing.T) {
	// A long mixed string at a scale that will not divide evenly, so at least
	// one glyph is certain to land off the pixel grid.
	glyphs := walkGlyphs(t, "The quick brown fox jumps over the lazy dog", 13, 1.25)
	require.NotEmpty(t, glyphs)

	fractional := 0
	for _, g := range glyphs {
		if _, frac := math.Modf(float64(g.penX)); frac != 0 {
			fractional++
		}
	}
	assert.NotZero(t, fractional, "pen positions should keep their fractional part, got all whole pixels")
}

func TestRenderGlyphToImage(t *testing.T) {
	glyphs := walkGlyphs(t, "M", 24, 1)
	require.Len(t, glyphs, 1)
	g := glyphs[0]

	img, baseline := painter.RenderGlyphToImage(g.run, g.idx, 24, 1, 0, color.White)
	require.NotNil(t, img)

	assert.Positive(t, img.Bounds().Dx(), "glyph bitmap should have width")
	assert.Positive(t, img.Bounds().Dy(), "glyph bitmap should have height")
	assert.Positive(t, baseline, "baseline should sit below the top of the bitmap")
	assert.LessOrEqual(t, baseline, img.Bounds().Dy(), "baseline should fall inside the bitmap")
	assert.NotZero(t, inkPixels(img), "a rendered 'M' should put ink in the bitmap")
}

// TestRenderGlyphToImage_SubpixelShiftsInk checks that the sub-pixel offset is
// actually rasterised into the bitmap rather than ignored, which is what lets a
// glyph sit on a fractional position without being resampled at draw time.
func TestRenderGlyphToImage_SubpixelShiftsInk(t *testing.T) {
	glyphs := walkGlyphs(t, "l", 32, 1)
	require.Len(t, glyphs, 1)
	g := glyphs[0]

	atZero, baseZero := painter.RenderGlyphToImage(g.run, g.idx, 32, 1, 0, color.White)
	atHalf, baseHalf := painter.RenderGlyphToImage(g.run, g.idx, 32, 1, 0.5, color.White)
	require.NotNil(t, atZero)
	require.NotNil(t, atHalf)

	assert.Equal(t, baseZero, baseHalf, "sub-pixel offset is horizontal, it must not move the baseline")
	assert.Equal(t, atZero.Bounds(), atHalf.Bounds(), "sub-pixel variants should be the same size")
	assert.NotEqual(t, atZero.Pix, atHalf.Pix, "half a pixel of offset should change the rasterisation")

	// The shift is to the right, so ink should not start further left than it
	// did without any offset.
	assert.GreaterOrEqual(t, firstInkColumn(atHalf), firstInkColumn(atZero),
		"a positive sub-pixel offset should not move ink left")
}

func TestRenderGlyphToImage_Colour(t *testing.T) {
	glyphs := walkGlyphs(t, "X", 24, 1)
	require.Len(t, glyphs, 1)
	g := glyphs[0]

	red, _ := painter.RenderGlyphToImage(g.run, g.idx, 24, 1, 0, color.NRGBA{R: 0xff, A: 0xff})
	require.NotNil(t, red)

	// Find the most opaque pixel and check it carries the requested colour.
	var bestR, bestG, bestB, bestA uint32
	for y := red.Bounds().Min.Y; y < red.Bounds().Max.Y; y++ {
		for x := red.Bounds().Min.X; x < red.Bounds().Max.X; x++ {
			r, g, b, a := red.At(x, y).RGBA()
			if a > bestA {
				bestR, bestG, bestB, bestA = r, g, b, a
			}
		}
	}
	require.NotZero(t, bestA, "glyph should have rendered some ink")
	assert.Greater(t, bestR, bestG, "ink should carry the red requested")
	assert.Greater(t, bestR, bestB, "ink should carry the red requested")
}

func inkPixels(img *image.RGBA) int {
	count := 0
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if _, _, _, a := img.At(x, y).RGBA(); a > 0 {
				count++
			}
		}
	}
	return count
}

// firstInkColumn returns the leftmost column holding any ink, or the width when
// the bitmap is empty.
func firstInkColumn(img *image.RGBA) int {
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			if _, _, _, a := img.At(x, y).RGBA(); a > 0 {
				return x
			}
		}
	}
	return img.Bounds().Dx()
}
