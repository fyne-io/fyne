package animation

import (
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
)

// Runner is the main driver for animations package
type Runner struct {
	pendingAnimationsMutex sync.Mutex
	pendingAnimations      []*anim
	runnerStarted          atomic.Bool

	animations []*anim // accessed only by runAnimations
}

// Start will register the passed application and initiate its ticking.
func (r *Runner) Start(a *fyne.Animation) {
	r.pendingAnimationsMutex.Lock()
	r.pendingAnimations = append(r.pendingAnimations, newAnim(a))
	r.pendingAnimationsMutex.Unlock()

	if r.runnerStarted.CompareAndSwap(false, true) {
		go r.runAnimations()
	}
}

func (r *Runner) runAnimations() {
	draw := time.NewTicker(time.Second / 60)

	for done := false; !done; {
		<-draw.C

		// tick currently running animations
		// use technique from https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
		// to filter the still-running animations for the next iteration without allocating a new slice
		newList := r.animations[:0]
		for _, a := range r.animations {
			if stopped := a.a.State() == fyne.AnimationStateStopped; !stopped && r.tickAnimation(a) {
				newList = append(newList, a) // still running
			} else if !stopped {
				a.a.Stop() // mark as stopped (completed running)
			}
		}

		// bring in all pending animations
		r.pendingAnimationsMutex.Lock()
		for i, a := range r.pendingAnimations {
			newList = append(newList, a)
			r.pendingAnimations[i] = nil
		}
		r.pendingAnimations = r.pendingAnimations[:0]
		r.pendingAnimationsMutex.Unlock()

		done = len(newList) == 0
		for i := len(newList); i < len(r.animations); i++ {
			r.animations[i] = nil // nil out extra slice capacity
		}
		r.animations = newList
	}
	r.runnerStarted.Store(false)
	draw.Stop()
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
