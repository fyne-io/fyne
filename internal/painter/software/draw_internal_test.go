package software

import (
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/draw"
)

func genericFill(base *image.NRGBA, bounds image.Rectangle, fill color.Color) {
	draw.Draw(base, bounds, image.NewUniform(fill), image.Point{}, draw.Over)
}

func TestFillRectFastPath_MatchesGenericPath(t *testing.T) {
	const w, h = 40, 30
	fillColor := color.NRGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xff}

	cases := map[string]image.Rectangle{
		"full":          image.Rect(0, 0, w, h),
		"interior":      image.Rect(5, 5, 20, 18),
		"clipped-right": image.Rect(35, 5, 60, 20),
		"clipped-left":  image.Rect(-10, 5, 10, 20),
		"empty":         image.Rect(10, 10, 10, 20),
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
		t.Fatal("fillRectFastPath accepted a non-opaque fill")
	}

	if fillRectFastPath(base, base.Rect, nil) {
		t.Fatal("fillRectFastPath accepted a nil fill")
	}
}
