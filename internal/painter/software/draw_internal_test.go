package software

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/draw"
)

// genericFill reproduces the pre-fast-path behaviour: it always goes through
// image/draw's generic Porter-Duff loop, which is the only path taken for an
// *image.NRGBA destination before fillRectFastPath was introduced.
func genericFill(base *image.NRGBA, bounds image.Rectangle, fill color.Color) {
	draw.Draw(base, bounds, image.NewUniform(fill), image.Point{}, draw.Over)
}

func TestFillRectFastPath_MatchesGenericPath(t *testing.T) {
	const w, h = 40, 30
	fillColor := color.NRGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xff}

	cases := map[string]image.Rectangle{
		"full":              image.Rect(0, 0, w, h),
		"interior":          image.Rect(5, 5, 20, 18),
		"clipped-right":     image.Rect(35, 5, 60, 20),
		"clipped-left":      image.Rect(-10, 5, 10, 20),
		"clipped-top":       image.Rect(5, -10, 20, 5),
		"negative-position": image.Rect(-5, -5, 5, 5),
		"zero-size":         image.Rect(10, 10, 10, 20),
		"zero-height":       image.Rect(10, 10, 20, 10),
		"out-of-bounds":     image.Rect(100, 100, 120, 120),
	}

	for name, r := range cases {
		t.Run(name, func(t *testing.T) {
			bounds := image.Rect(0, 0, w, h).Intersect(r)

			want := image.NewNRGBA(image.Rect(0, 0, w, h))
			genericFill(want, bounds, fillColor)

			got := image.NewNRGBA(image.Rect(0, 0, w, h))
			applied := fillRectFastPath(got, bounds, fillColor)

			if bounds.Empty() {
				if applied {
					t.Fatalf("fillRectFastPath reported success on an empty bounds %v", bounds)
				}
				return
			}

			if !applied {
				t.Fatalf("fillRectFastPath declined an opaque fill for bounds %v", bounds)
			}

			for i := range want.Pix {
				if want.Pix[i] != got.Pix[i] {
					t.Fatalf("pixel byte %d differs: generic=%d fastpath=%d (bounds=%v)", i, want.Pix[i], got.Pix[i], bounds)
				}
			}
		})
	}
}

func TestFillRectFastPath_DeclinesNonOpaque(t *testing.T) {
	base := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	translucent := color.NRGBA{R: 0x11, G: 0x22, B: 0x33, A: 0x80}

	if fillRectFastPath(base, base.Rect, translucent) {
		t.Fatalf("fillRectFastPath accepted a non-opaque fill")
	}
}

func TestFillRectFastPath_DeclinesNilFill(t *testing.T) {
	base := image.NewNRGBA(image.Rect(0, 0, 10, 10))

	if fillRectFastPath(base, base.Rect, nil) {
		t.Fatalf("fillRectFastPath accepted a nil fill")
	}
}

// TestFillRectFastPath_DetectsMutation is a canary proving the pixel-parity
// check above is sensitive to a wrong implementation: it swaps in a
// deliberately broken fill (off-by-one in the blue channel) and confirms the
// same comparison loop used above catches it. This documents that the tests
// in this file would have failed before fillRectFastPath's byte order and
// bounds handling were corrected during development.
func TestFillRectFastPath_DetectsMutation(t *testing.T) {
	bounds := image.Rect(2, 2, 8, 8)
	fillColor := color.NRGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xff}

	want := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	genericFill(want, bounds, fillColor)

	got := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	brokenFillRectFastPath(got, bounds, fillColor)

	mismatch := false
	for i := range want.Pix {
		if want.Pix[i] != got.Pix[i] {
			mismatch = true
			break
		}
	}
	if !mismatch {
		t.Fatalf("expected the deliberately broken fill to be detected as a mismatch")
	}
}

// brokenFillRectFastPath is a deliberately incorrect variant (blue channel
// off by one) used only to prove TestFillRectFastPath_DetectsMutation's
// comparison loop is sensitive to real regressions.
func brokenFillRectFastPath(base *image.NRGBA, bounds image.Rectangle, fill color.Color) {
	r, g, b, _ := fill.RGBA()
	nr, ng, nb := uint8(r>>8), uint8(g>>8), uint8(b>>8)+1
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		off := base.PixOffset(bounds.Min.X, y)
		end := base.PixOffset(bounds.Max.X, y)
		row := base.Pix[off:end]
		for p := 0; p < len(row); p += 4 {
			row[p] = nr
			row[p+1] = ng
			row[p+2] = nb
			row[p+3] = 0xff
		}
	}
}
