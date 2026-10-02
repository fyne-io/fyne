package cache

import (
	"strconv"
	"testing"

	"fyne.io/fyne/v2"
)

// BenchmarkGetFontMetrics is the per-frame lookup every canvas.Text does from
// MinSize while the canvas walks its tree.
func BenchmarkGetFontMetrics(b *testing.B) {
	texts := make([]string, 300)
	for i := range texts {
		texts[i] = strconv.FormatFloat(float64(i)*1.37, 'f', 2, 64)
		SetFontMetrics(texts[i], 11, fyne.TextStyle{}, nil, fyne.NewSize(20, 11), 9)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, t := range texts {
			GetFontMetrics(t, 11, fyne.TextStyle{}, nil)
		}
	}
}
