package fyne

import "time"

type Animation struct {
	Duration time.Duration
	Repeat   bool
	Tick     func(float32)
}

func NewAnimation(d time.Duration, fn func(float32)) *Animation {
	return &Animation{Duration: d, Tick: fn}
}

func (a *Animation) Start() {
	d := CurrentApp().Driver()
	d.StartAnimation(a)
}
