package canvas

import (
	"image/color"
	"time"

	"fyne.io/fyne"
)

const (
	// DurationStandard is the time a standard interface animation will run.
	DurationStandard = time.Millisecond * 300
	// DurationShort is the time a subtle or small transition should use.
	DurationShort = time.Millisecond * 150
)

// NewColorRGBAAnimation sets up a new animation that will transition from the start to stop Color over
// the specified Duration. The colour transition will move linearly through the RGB colour space.
// The content of fn should apply the color values to an object and refresh it.
// You should call Start() on the returned animation to start it.
func NewColorRGBAAnimation(start, stop color.Color, d time.Duration, fn func(color.Color)) *fyne.Animation {
	r1, g1, b1, a1 := start.RGBA()
	r2, g2, b2, a2 := stop.RGBA()

	rStart := int(r1 >> 8)
	gStart := int(g1 >> 8)
	bStart := int(b1 >> 8)
	aStart := int(a1 >> 8)
	rDelta := float32(int(r2>>8) - rStart)
	gDelta := float32(int(g2>>8) - gStart)
	bDelta := float32(int(b2>>8) - bStart)
	aDelta := float32(int(a2>>8) - aStart)

	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(color.RGBA{R: scaleChannel(rStart, rDelta, done), G: scaleChannel(gStart, gDelta, done),
				B: scaleChannel(bStart, bDelta, done), A: scaleChannel(aStart, aDelta, done)})
		}}
}

// NewPositionAnimation sets up a new animation that will transition from the start to stop Position over
// the specified Duration. The content of fn should apply the position value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
func NewPositionAnimation(start, stop fyne.Position, d time.Duration, fn func(fyne.Position)) *fyne.Animation {
	xDelta := float32(stop.X - start.X)
	yDelta := float32(stop.Y - start.Y)

	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(fyne.NewPos(scaleInt(start.X, xDelta, done), scaleInt(start.Y, yDelta, done)))
		}}
}

// NewSizeAnimation sets up a new animation that will transition from the start to stop Size over
// the specified Duration. The content of fn should apply the size value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
func NewSizeAnimation(start, stop fyne.Size, d time.Duration, fn func(fyne.Size)) *fyne.Animation {
	widthDelta := float32(stop.Width - start.Width)
	heightDelta := float32(stop.Height - start.Height)

	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(fyne.NewSize(scaleInt(start.Width, widthDelta, done), scaleInt(start.Height, heightDelta, done)))
		}}
}

func scaleChannel(start int, diff, done float32) uint8 {
	return uint8(start + int(diff*done))
}

func scaleInt(start int, delta, done float32) int {
	return start + int(delta*done)
}
