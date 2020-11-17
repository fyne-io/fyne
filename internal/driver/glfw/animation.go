package glfw

import (
	"time"

	"fyne.io/fyne"
)

type anim struct {
	a       *fyne.Animation
	end     time.Time
	reverse bool
	start   time.Time
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
			done = len(d.animations) == 0
			d.animationMutex.Unlock()
		}
	}()
}

// tickAnimation will process a frame of animation and return true if this should continue animating
func (d *gLDriver) tickAnimation(a *anim) bool {
	if time.Now().After(a.end) {
		if a.reverse {
			a.a.Tick(0.0)
			if !a.a.Repeat {
				return false
			}
			a.reverse = false
		} else {
			a.a.Tick(1.0)
			if a.a.AutoReverse {
				a.reverse = true
			}
		}
		if !a.a.Repeat && !a.reverse {
			return false
		}

		a.start = time.Now()
		a.end = a.start.Add(a.a.Duration)
	}

	total := a.end.Sub(a.start).Nanoseconds() / 1000000 // TODO change this to Milliseconds() when we drop Go 1.12
	delta := time.Since(a.start).Nanoseconds() / 1000000

	val := float32(delta) / float32(total)
	curve := a.a.Curve
	if curve == nil {
		curve = fyne.AnimationEaseInOut
	}
	if a.reverse {
		a.a.Tick(curve(1 - val))
	} else {
		a.a.Tick(curve(val))
	}

	return true
}
