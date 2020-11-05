package fyne

import "time"

// Animation represents an animated element within a Fyne canvas.
// These animations may control individual objects or entire scenes.
type Animation struct {
	Duration time.Duration
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
