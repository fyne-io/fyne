package animation

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// Runner is the main driver for animations package
type Runner struct {
	animationMutex    sync.RWMutex
	animations        []*anim
	pendingAnimations []*anim
}

// Start will register the passed application and initiate its ticking.
func (r *Runner) Start(a *fyne.Animation) {
	r.animationMutex.Lock()
	defer r.animationMutex.Unlock()
	r.pendingAnimations = append(r.pendingAnimations, newAnim(a))
}

// Stop causes an animation to stop ticking (if it was still running) and removes it from the runner.
func (r *Runner) Stop(a *fyne.Animation) {
	r.animationMutex.Lock()
	defer r.animationMutex.Unlock()

	newList := make([]*anim, 0, len(r.animations))
	stopped := false
	for _, item := range r.animations {
		if item.a != a {
			newList = append(newList, item)
		} else {
			item.setStopped()
			stopped = true
		}
	}
	r.animations = newList
	if stopped {
		return
	}

	newList = make([]*anim, 0, len(r.pendingAnimations))
	for _, item := range r.pendingAnimations {
		if item.a != a {
			newList = append(newList, item)
		} else {
			item.setStopped()
		}
	}
	r.pendingAnimations = newList
}

func (r *Runner) TickAnimations() {
	r.animationMutex.Lock()
	if len(r.animations) == 0 && len(r.pendingAnimations) == 0 {
		r.animationMutex.Unlock()
		return
	}
	currList := r.animations
	r.animationMutex.Unlock()

	evictAnimations := false
	now := time.Now()
	for _, a := range currList {
		if a.isStopped() || !r.tickAnimation(a, now) {
			evictAnimations = true
		}
	}

	if evictAnimations {
		newList := make([]*anim, 0, len(currList)+len(r.pendingAnimations))
		for _, a := range currList {
			if !a.isStopped() && a.repeatsLeft != 0 {
				newList = append(newList, a)
			}
		}
		currList = newList
	}

	r.animationMutex.Lock()
	r.animations = append(currList, r.pendingAnimations...)
	r.pendingAnimations = nil
	r.animationMutex.Unlock()
}

// tickAnimation will process a frame of animation and return true if this should continue animating
func (r *Runner) tickAnimation(a *anim, now time.Time) bool {
	if now.After(a.end) {
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

		a.start = now
		a.end = a.start.Add(a.a.Duration)
		return true
	}

	delta := time.Since(a.start).Milliseconds()

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
