package software

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BenchmarkCanvas_Capture(b *testing.B) {
	c := NewCanvas()
	c.SetContent(container.NewVBox(widget.NewLabel("Hello"), widget.NewButton("OK", nil)))
	c.Resize(fyne.NewSize(400, 300))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Capture()
	}
}
