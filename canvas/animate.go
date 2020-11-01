package canvas

import (
	"image/color"
	"time"

	"fyne.io/fyne"
)

func NewColorAnimator(start, stop color.Color, d time.Duration, fn func(color.Color)) *fyne.Animation {
	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			r1, g1, b1, a1 := start.RGBA()
			r2, g2, b2, a2 := stop.RGBA()

			fn(color.NRGBA{R: scaleColor(r1, r2, done), G: scaleColor(g1, g2, done), B: scaleColor(b1, b2, done),
				A: scaleColor(a1, a2, done)})
		}}
}

func scaleColor(start uint32, end uint32, done float32) uint8 {
	diff := int(end) - int(start)
	shift := int(float32(diff) * done)

	return uint8((int(start) + shift) >> 8)
}
