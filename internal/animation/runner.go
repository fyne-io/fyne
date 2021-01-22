package animation

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// Runner is the main driver for animations package
type Runner struct {
	animationMutex sync.RWMutex
	animations     []*anim
}

// Start will register the passed application and initiate its ticking.
func (r *Runner) Start(a *fyne.Animation) {
	r.animationMutex.Lock()
	defer r.animationMutex.Unlock()
	wasStopped := len(r.animations) == 0

	r.animations = append(r.animations, newAnim(a))
	if wasStopped {
		r.runAnimations()
	}
}

// Stop causes an animation to stop ticking (if it was still running) and removes it from the runner.
func (r *Runner) Stop(a *fyne.Animation) {
	r.animationMutex.Lock()
	defer r.animationMutex.Unlock()
	oldList := r.animations
	var newList []*anim
	for _, item := range oldList {
		if item.a != a {
			newList = append(newList, item)
		}
	}
	r.animations = newList
}

func (r *Runner) runAnimations() {
	draw := time.NewTicker(time.Second / 60)

	go func() {
		for done := false; !done; {
			<-draw.C
			r.animationMutex.Lock()
			oldList := r.animations
			r.animations = nil // clear the list so we can append any new ones after processing
			r.animationMutex.Unlock()
			var newList []*anim
			for _, a := range oldList {
				if r.tickAnimation(a) {
					newList = append(newList, a)
				}
			}
			r.animationMutex.Lock()
			r.animations = append(newList, r.animations...)
			done = len(r.animations) == 0
			r.animationMutex.Unlock()
		}
		draw.Stop()
	}()
}

// tickAnimation will process a frame of animation and return true if this should continue animating
func (r *Runner) tickAnimation(a *anim) bool {
	if time.Now().After(a.end) {
		if a.reverse {
			a.a.Tick(0.0)
			if a.repeatsLeft == 0 {
				return false
			}
			a.reverse = false
		} else {
			a.a.Tick(1.0)
			if a.a.AutoReverse {
				a.reverse = true
			}
		}
		if !a.reverse {
			if a.repeatsLeft == 0 {
				return false
			}
			if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
		}

		a.start = time.Now()
		a.end = a.start.Add(a.a.Duration)
	}

	delta := time.Since(a.start).Nanoseconds() / 1000000 // TODO change this to Milliseconds() when we drop Go 1.12

	val := float32(delta) / float32(a.total)
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
