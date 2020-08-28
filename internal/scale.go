package internal

import (
	"math"

	"fyne.io/fyne"
)

// ScaleInt converts a fyne coordinate in the given canvas to a screen coordinate
func ScaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(math.Round(float64(v) * float64(c.Scale())))
	}
}

// ScaleSize converts a fyne.Size (unit coordinate) to a fyne.PixelSize (screen coordinate)
func ScaleSize(c fyne.Canvas, s fyne.Size) PixelSize {
	return NewPixelSize(ScaleInt(c, s.Width), ScaleInt(c, s.Height))
}

// UnscaleInt converts a screen coordinate for a given canvas to a fyne coordinate
func UnscaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) / c.Scale())
	}
}

// UnscaleSize converts a fyne.PixelSize (screen coordinate) to fyne.Size (unit coordinate)
func UnscaleSize(c fyne.Canvas, s PixelSize) fyne.Size {
	return fyne.NewSize(UnscaleInt(c, s.WidthPx), UnscaleInt(c, s.HeightPx))
}
