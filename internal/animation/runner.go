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
	nextAnimations    []*anim // used by runner goroutine only

	runnerStarted bool
}

// Start will register the passed application and initiate its ticking.
func (r *Runner) Start(a *fyne.Animation) {
	r.animationMutex.Lock()
	defer r.animationMutex.Unlock()

	if !r.runnerStarted {
		r.runnerStarted = true
		if r.animations == nil {
			// initialize with excess capacity to avoid re-allocations
			// on subsequent Starts
			r.animations = make([]*anim, 0, 16)
		}
		r.animations = append(r.animations, newAnim(a))
		r.runAnimations()
	} else {
		if r.pendingAnimations == nil {
			// initialize with excess capacity to avoid re-allocations
			// on subsequent Starts
			r.pendingAnimations = make([]*anim, 0, 16)
		}
		r.pendingAnimations = append(r.pendingAnimations, newAnim(a))
	}
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

func (r *Runner) runAnimations() {
	draw := time.NewTicker(time.Second / 60)
	go func() {
		for done := false; !done; {
			<-draw.C
			done = r.runOneFrame()
		}
		r.animationMutex.Lock()
		r.runnerStarted = false
		r.animationMutex.Unlock()
		draw.Stop()
	}()
}

func (r *Runner) runOneFrame() (done bool) {
	r.animationMutex.Lock()
	oldList := r.animations
	r.animationMutex.Unlock()
	for _, a := range oldList {
		if !a.isStopped() && r.tickAnimation(a) {
			r.nextAnimations = append(r.nextAnimations, a)
		}
	}

	r.animationMutex.Lock()
	// nil out old r.animations for re-use as next r.nextAnimations
	tmp := r.animations
	for i := range tmp {
		tmp[i] = nil
	}
	r.animations = append(r.nextAnimations, r.pendingAnimations...)
	r.nextAnimations = tmp[:0]
	// nil out r.pendingAnimations
	for i := range r.pendingAnimations {
		r.pendingAnimations[i] = nil
	}
	r.pendingAnimations = r.pendingAnimations[:0]
	done = len(r.animations) == 0
	r.animationMutex.Unlock()
	return done
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
