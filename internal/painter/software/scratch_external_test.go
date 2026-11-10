package software_test

import (
	"bytes"
	"image"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal/painter/software"
	"fyne.io/fyne/v2/test"
)

func TestPaintScratch_MatchesPaintAcrossSizes(t *testing.T) {
	p := software.NewPainter()
	for _, size := range []fyne.Size{{Width: 64, Height: 48}, {Width: 32, Height: 96}, {Width: 128, Height: 16}} {
		c := test.NewCanvas()
		rect := canvas.NewRectangle(image.White)
		rect.Resize(fyne.NewSize(size.Width/2, size.Height/2))
		c.SetContent(rect)
		c.Resize(size)

		want := p.Paint(c).(*image.NRGBA)
		got := p.PaintScratch(c)
		if got.Rect != want.Rect || got.Stride != want.Stride {
			t.Fatalf("size %v: geometry differs: got rect=%v stride=%d, want rect=%v stride=%d", size, got.Rect, got.Stride, want.Rect, want.Stride)
		}
		if !bytes.Equal(got.Pix, want.Pix) {
			t.Fatalf("size %v: pixels differ between Paint and PaintScratch", size)
		}
		p.ReleaseScratch(got)
	}
}
