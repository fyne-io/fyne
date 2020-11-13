package canvas

import (
	"image/color"
	"time"

	"fyne.io/fyne"
)

// DurationStandard is the time a standard interface animation will run.
const DurationStandard = time.Millisecond * 300

// NewColorRGBAAnimation sets up a new animation that will transition from the start to stop Color over
// the specified Duration. The colour transition will move linearly through the RGB colour space.
// The content of fn should apply the color values to an object and refresh it.
// You should call Start() on the returned animation to start it.
func NewColorRGBAAnimation(start, stop color.Color, d time.Duration, fn func(color.Color)) *fyne.Animation {
	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			r1, g1, b1, a1 := start.RGBA()
			r2, g2, b2, a2 := stop.RGBA()

			fn(color.NRGBA{R: scaleColor(r1, r2, done), G: scaleColor(g1, g2, done), B: scaleColor(b1, b2, done),
				A: scaleColor(a1, a2, done)})
		}}
}

// NewPositionAnimation sets up a new animation that will transition from the start to stop Position over
// the specified Duration. The content of fn should apply the position value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
func NewPositionAnimation(start, stop fyne.Position, d time.Duration, fn func(fyne.Position)) *fyne.Animation {
	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(fyne.NewPos(scaleInt(start.X, stop.X, done), scaleInt(start.Y, stop.Y, done)))
		}}
}

// NewSizeAnimation sets up a new animation that will transition from the start to stop Size over
// the specified Duration. The content of fn should apply the size value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
func NewSizeAnimation(start, stop fyne.Size, d time.Duration, fn func(fyne.Size)) *fyne.Animation {
	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(fyne.NewSize(scaleInt(start.Width, stop.Width, done), scaleInt(start.Height, stop.Height, done)))
		}}
}

func scaleColor(start uint32, end uint32, done float32) uint8 {
	diff := int(end) - int(start)
	shift := int(float32(diff) * done)

	return uint8((int(start) + shift) >> 8)
}

func scaleInt(start, end int, done float32) int {
	shift := int(float32(end-start) * done)
	return start + shift
}
