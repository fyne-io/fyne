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

	rDelta := int(r2) - int(r1)
	gDelta := int(g2) - int(g1)
	bDelta := int(b2) - int(b1)
	aDelta := int(a2) - int(a1)

	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {

			fn(color.NRGBA{R: scaleChannel(r1, rDelta, done), G: scaleChannel(g1, gDelta, done),
				B: scaleChannel(b1, bDelta, done), A: scaleChannel(a1, aDelta, done)})
		}}
}

// NewPositionAnimation sets up a new animation that will transition from the start to stop Position over
// the specified Duration. The content of fn should apply the position value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
func NewPositionAnimation(start, stop fyne.Position, d time.Duration, fn func(fyne.Position)) *fyne.Animation {
	xDelta := stop.X - start.X
	yDelta := stop.Y - start.Y

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
	widthDelta := stop.Width - start.Width
	heightDelta := stop.Height - start.Height

	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(fyne.NewSize(scaleInt(start.Width, widthDelta, done), scaleInt(start.Height, heightDelta, done)))
		}}
}

func scaleChannel(start uint32, diff int, done float32) uint8 {
	shift := int(float32(diff) * done)

	return uint8((int(start) + shift) >> 8)
}

func scaleInt(start, delta int, done float32) int {
	shift := int(float32(delta) * done)
	return start + shift
}
