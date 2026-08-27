package software

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"golang.org/x/image/draw"
)

// benchCanvas implements only Scale; drawOblong does not call other Canvas methods.
type benchCanvas struct {
	fyne.Canvas
}

func (benchCanvas) Scale() float32 { return 1 }

func setupOblongBench(fill color.Color) (fyne.Canvas, *canvas.Rectangle, *image.NRGBA, image.Rectangle) {
	const w, h = 1280, 800
	rect := canvas.NewRectangle(fill)
	rect.Resize(fyne.NewSize(float32(w), float32(h)))
	base := image.NewNRGBA(image.Rect(0, 0, w, h))
	clip := base.Rect
	return benchCanvas{}, rect, base, clip
}

func BenchmarkDrawRectangle_opaqueScanline(b *testing.B) {
	c, rect, base, clip := setupOblongBench(color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		drawRectangle(c, rect, fyne.NewPos(0, 0), base, clip)
	}
}

func BenchmarkDrawRectangle_opaqueUniformOver(b *testing.B) {
	fill := color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff}
	_, _, base, clip := setupOblongBench(fill)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		draw.Draw(base, clip, image.NewUniform(fill), image.Point{}, draw.Over)
	}
}
