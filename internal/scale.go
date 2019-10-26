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

// UnscaleInt converts a screen coordinate for a given canvas to a fyne coordinate
func UnscaleInt(c fyne.Canvas, v int) int {
	switch c.Scale() {
	case 1.0:
		return v
	default:
		return int(float32(v) / c.Scale())
	}
}
