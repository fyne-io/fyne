package software

import (
	"image"
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	"golang.org/x/image/draw"
)

// benchCanvas implements only Scale; drawOblong does not call other Canvas
// methods for the opaque, unrounded, unstroked rectangle path exercised here.
type benchCanvas struct {
	fyne.Canvas
}

func (benchCanvas) Scale() float32 { return 1 }

func setupRectangleBench(w, h int, fill color.Color) (fyne.Canvas, *canvas.Rectangle, *image.NRGBA, image.Rectangle) {
	rect := canvas.NewRectangle(fill)
	rect.Resize(fyne.NewSize(float32(w), float32(h)))
	base := image.NewNRGBA(image.Rect(0, 0, w, h))
	clip := base.Rect
	return benchCanvas{}, rect, base, clip
}

func benchmarkOpaque(b *testing.B, w, h int) {
	c, rect, base, clip := setupRectangleBench(w, h, color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		drawRectangle(c, rect, fyne.NewPos(0, 0), base, clip)
	}
}

func benchmarkGeneric(b *testing.B, w, h int) {
	fill := color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff}
	_, _, base, clip := setupRectangleBench(w, h, fill)
	uniform := image.NewUniform(fill)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		draw.Draw(base, clip, uniform, image.Point{}, draw.Over)
	}
}

func BenchmarkDrawRectangle_Opaque_1280x800(b *testing.B)  { benchmarkOpaque(b, 1280, 800) }
func BenchmarkDrawRectangle_Generic_1280x800(b *testing.B) { benchmarkGeneric(b, 1280, 800) }
func BenchmarkDrawRectangle_Opaque_200x100(b *testing.B)   { benchmarkOpaque(b, 200, 100) }
func BenchmarkDrawRectangle_Generic_200x100(b *testing.B)  { benchmarkGeneric(b, 200, 100) }
