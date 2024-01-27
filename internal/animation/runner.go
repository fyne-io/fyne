package animation

import (
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
)

// Runner is the main driver for animations package
type Runner struct {
	animationMutex    sync.Mutex
	animations        []*anim
	pendingAnimations []*anim

	runnerStarted atomic.Bool
}

// Start will register the passed application and initiate its ticking.
func (r *Runner) Start(a *fyne.Animation) {
	r.animationMutex.Lock()
	defer r.animationMutex.Unlock()

	if r.runnerStarted.CompareAndSwap(false, true) {
		r.animations = append(r.animations, newAnim(a))
		go r.runAnimations()
	} else {
		r.pendingAnimations = append(r.pendingAnimations, newAnim(a))
	}
}

// Stop causes an animation to stop ticking (if it was still running) and removes it from the runner.
func (r *Runner) Stop(a *fyne.Animation) {
	// Since the runner needs to lock for the whole duration of a tick, which invokes user code,
	// we must stop asynchronously to avoid possible deadlock if Stop is called within an animation tick callback.
	// Since stopping animations should occur much less frequently than ticking them, the performance
	// penalty of spawning a goroutine for stop should be acceptable to achieve a zero-allocation tick implementation.
	go func() {
		r.animationMutex.Lock()
		defer r.animationMutex.Unlock()

		// use technique from https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
		// to filter the animation slice without allocating a new slice
		newList := r.animations[:0]
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

		newList = r.pendingAnimations[:0]
		for _, item := range r.pendingAnimations {
			if item.a != a {
				newList = append(newList, item)
			} else {
				item.setStopped()
			}
		}
		r.pendingAnimations = newList
	}()
}

func (r *Runner) runAnimations() {
	draw := time.NewTicker(time.Second / 60)
	for done := false; !done; {
		<-draw.C
		// use technique from https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
		// to filter the still-running animations for the next iteration without allocating a new slice
		r.animationMutex.Lock()
		newList := r.animations[:0]
		for _, a := range r.animations {
			if !a.isStopped() && r.tickAnimation(a) {
				newList = append(newList, a)
			}
		}
		r.animations = append(newList, r.pendingAnimations...)
		r.pendingAnimations = r.pendingAnimations[:0]
		done = len(r.animations) == 0
		r.animationMutex.Unlock()
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
