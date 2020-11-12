package fyne

import "time"

// AnimationCurve represents the different animation algorithms for animation timeline.
type AnimationCurve int

const (
	// AnimationEaseInOut is the default easing, it starts slowly, accelerates to the middle and slows to the end.
	AnimationEaseInOut AnimationCurve = iota
	// AnimationEaseIn starts slowly and accelerates to the end.
	AnimationEaseIn
	// AnimationEaseOut starts at speed and slows to the end.
	AnimationEaseOut
	// AnimationLinear is a linear mapping for animations that progress uniformly through their duration.
	AnimationLinear
)

// Animation represents an animated element within a Fyne canvas.
// These animations may control individual objects or entire scenes.
type Animation struct {
	Duration time.Duration
	Curve    AnimationCurve
	Repeat   bool
	Tick     func(float32)
}

// NewAnimation creates a very basic animation where the callback function will be called for every
// rendered frame between time.Now() and the specified duration. The callback values start at 0.0 and
// will be 1.0 when the animation completes.
func NewAnimation(d time.Duration, fn func(float32)) *Animation {
	return &Animation{Duration: d, Tick: fn}
}

// Start registers the animation with the application run-loop and starts its execution.
func (a *Animation) Start() {
	CurrentApp().Driver().StartAnimation(a)
}
