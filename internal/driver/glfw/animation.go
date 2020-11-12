package glfw

import (
	"time"

	"fyne.io/fyne"
)

type anim struct {
	a     *fyne.Animation
	start time.Time
	end   time.Time
}

func (d *gLDriver) runAnimations() {
	draw := time.NewTicker(time.Second / 60)

	go func() {
		done := false
		for !done {

			<-draw.C
			d.animationMutex.Lock()
			oldList := d.animations
			d.animations = nil // clear the list so we can append any new ones after processing
			d.animationMutex.Unlock()
			var newList []*anim
			for _, a := range oldList {
				if d.tickAnimation(a) {
					newList = append(newList, a)
				}
			}
			d.animationMutex.Lock()
			d.animations = append(newList, d.animations...)
			done = len(newList) == 0
			d.animationMutex.Unlock()
		}
	}()
}

// tickAnimation will process a frame of animation and return true if this should continue animating
func (d *gLDriver) tickAnimation(a *anim) bool {
	if time.Now().After(a.end) {
		a.a.Tick(1.0)
		if !a.a.Repeat {
			return false
		}

		a.start = time.Now()
		a.end = a.start.Add(a.a.Duration)
	}

	total := a.end.Sub(a.start).Nanoseconds() / 1000000 // TODO change this to Milliseconds() when we drop Go 1.12
	delta := time.Since(a.start).Nanoseconds() / 1000000

	val := float32(delta) / float32(total)
	a.a.Tick(animationCurve(val, a.a.Curve))

	return true
}

func animationCurve(val float32, curve fyne.AnimationCurve) float32 {
	switch curve {
	case fyne.AnimationLinear:
		return val
	case fyne.AnimationEaseIn:
		return animationEaseIn(val)
	case fyne.AnimationEaseOut:
		return animationEaseOut(val)
	default: // EaseInOut is default
		return animationEaseInOut(val)
	}
}

func animationEaseIn(val float32) float32 {
	return val * val
}

func animationEaseInOut(val float32) float32 {
	if val <= 0.5 {
		return val * val * 2
	}

	return -1 + (4-val*2)*val
}

func animationEaseOut(val float32) float32 {
	return val * (2 - val)
}
