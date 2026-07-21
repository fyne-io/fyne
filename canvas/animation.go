package canvas

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
)

const (
	// DurationStandard is the time a standard interface animation will run.
	//
	// Since: 2.0
	DurationStandard = time.Millisecond * 300
	// DurationShort is the time a subtle or small transition should use.
	//
	// Since: 2.0
	DurationShort = time.Millisecond * 150
)

// shaderMaxFrameDelta caps the time advanced per animation frame. A gap longer than this
// means animation was paused (or stalled) and should resume without time progression for smoothness.
const shaderMaxFrameDelta = 100 * time.Millisecond

// NewColorRGBAAnimation sets up a new animation that will transition from the start to stop Color over
// the specified Duration. The colour transition will move linearly through the RGB colour space.
// The content of fn should apply the color values to an object and refresh it.
// You should call Start() on the returned animation to start it.
//
// Since: 2.0
func NewColorRGBAAnimation(start, stop color.Color, d time.Duration, fn func(color.Color)) *fyne.Animation {
	c1, _ := color.RGBAModel.Convert(start).(color.RGBA)
	c2, _ := color.RGBAModel.Convert(stop).(color.RGBA)
	rDelta := int(c2.R) - int(c1.R)
	gDelta := int(c2.G) - int(c1.G)
	bDelta := int(c2.B) - int(c1.B)
	aDelta := int(c2.A) - int(c1.A)

	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(color.RGBA{
				R: scaleChannel(c1.R, rDelta, done),
				G: scaleChannel(c1.G, gDelta, done),
				B: scaleChannel(c1.B, bDelta, done),
				A: scaleChannel(c1.A, aDelta, done),
			})
		},
	}
}

// NewPositionAnimation sets up a new animation that will transition from the start to stop Position over
// the specified Duration. The content of fn should apply the position value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
//
// Since: 2.0
func NewPositionAnimation(start, stop fyne.Position, d time.Duration, fn func(fyne.Position)) *fyne.Animation {
	xDelta := stop.X - start.X
	yDelta := stop.Y - start.Y

	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(fyne.NewPos(scaleVal(start.X, xDelta, done), scaleVal(start.Y, yDelta, done)))
		},
	}
}

// NewSizeAnimation sets up a new animation that will transition from the start to stop Size over
// the specified Duration. The content of fn should apply the size value to an object for the change
// to be visible. You should call Start() on the returned animation to start it.
//
// Since: 2.0
func NewSizeAnimation(start, stop fyne.Size, d time.Duration, fn func(fyne.Size)) *fyne.Animation {
	widthDelta := stop.Width - start.Width
	heightDelta := stop.Height - start.Height

	return &fyne.Animation{
		Duration: d,
		Tick: func(done float32) {
			fn(fyne.NewSize(scaleVal(start.Width, widthDelta, done), scaleVal(start.Height, heightDelta, done)))
		},
	}
}

// NewShaderAnimation sets up a new animation that continuously redraws the given
// shader, advancing its "time" uniform each frame so the fragment shader can
// produce motion. You should call Start() on the returned animation to begin and
// Stop() to freeze the shader at its current frame. Stopping and starting the same
// animation again resumes from where it left off, without counting the paused time.
//
// Since: 2.8
func NewShaderAnimation(s *Shader) *fyne.Animation {
	var elapsed time.Duration
	var lastTick time.Time

	return &fyne.Animation{
		Duration:    time.Second,
		Curve:       fyne.AnimationLinear,
		RepeatCount: fyne.AnimationRepeatForever,
		Tick: func(float32) {
			elapsed, lastTick = advanceShaderTime(elapsed, lastTick, time.Now())
			if s.Uniforms == nil {
				s.Uniforms = make(map[string]float32, 1)
			}
			s.Uniforms["time"] = float32(elapsed.Seconds())
			s.Refresh()
		},
	}
}

// advanceShaderTime accumulates animation time for a single frame ending at now.
// A gap larger than shaderMaxFrameDelta is treated as a pause/resume and adds no
// time. It returns the updated elapsed time and the tick to measure from next.
func advanceShaderTime(elapsed time.Duration, lastTick, now time.Time) (time.Duration, time.Time) {
	if !lastTick.IsZero() {
		if delta := now.Sub(lastTick); delta <= shaderMaxFrameDelta {
			elapsed += delta
		}
	}
	return elapsed, now
}

func scaleChannel(start uint8, diff int, done float32) uint8 {
	return start + uint8(float32(diff)*done)
}

func scaleVal(start float32, delta, done float32) float32 {
	return start + delta*done
}
