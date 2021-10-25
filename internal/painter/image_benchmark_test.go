package painter

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

func benchmarkTextLoadImage(img fyne.Resource, imgW, imgH int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		PaintImage(&canvas.Image{Resource: img}, nil, imgW, imgH)
	}
}

func BenchmarkTextLoadPngImage(b *testing.B) {
	benchmarkTextLoadImage(theme.FyneLogo(), 256, 256, b)
}

func BenchmarkTextLoadSvgImage(b *testing.B) {
	benchmarkTextLoadImage(theme.ContentCutIcon(), 256, 256, b)
}
