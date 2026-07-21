//go:build !ci || !darwin

package animation

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func tick(run *Runner) {
	time.Sleep(time.Millisecond * 5) // wait long enough that we are not at 0 time
	run.TickAnimations()
}

func TestGLDriver_StartAnimation(t *testing.T) {
	done := make(chan float32)
	run := &Runner{}
	a := &fyne.Animation{
		Duration: time.Millisecond * 100,
		Tick: func(d float32) {
			done <- d
		},
	}

	run.Start(a)
	go tick(run) // simulate a graphics draw loop
	select {
	case d := <-done:
		assert.Greater(t, d, float32(0))
	case <-time.After(100 * time.Millisecond):
		t.Error("animation was not ticked")
	}
}

func TestGLDriver_StopAnimation(t *testing.T) {
	done := make(chan float32)
	run := &Runner{}
	a := &fyne.Animation{
		Duration: time.Second * 10,
		Tick: func(d float32) {
			done <- d
		},
	}

	run.Start(a)
	go tick(run) // simulate a graphics draw loop
	select {
	case d := <-done:
		assert.Greater(t, d, float32(0))
	case <-time.After(time.Second):
		t.Error("animation was not ticked")
	}
	run.Stop(a)
	run.animationMutex.RLock()
	assert.Zero(t, len(run.animations))
	run.animationMutex.RUnlock()
}

func TestRunner_DurationIncreasedMidAnimation(t *testing.T) {
	progress := make(chan float32, 8)
	run := &Runner{}
	a := &fyne.Animation{
		Duration: time.Second,
		Curve:    fyne.AnimationLinear,
		Tick: func(d float32) {
			progress <- d
		},
	}
	run.Start(a)

	time.Sleep(200 * time.Millisecond)
	run.TickAnimations()
	var before float32
	select {
	case before = <-progress:
	case <-time.After(time.Second):
		t.Fatal("animation was not ticked")
	}

	// Extend the duration; progress should not snap backwards.
	a.Duration = 4 * time.Second
	run.TickAnimations()
	var after float32
	select {
	case after = <-progress:
	case <-time.After(time.Second):
		t.Fatal("animation was not ticked after duration change")
	}
	assert.InDelta(t, before, after, 0.1,
		"progress should be preserved when Duration grows: before=%v after=%v", before, after)

	time.Sleep(100 * time.Millisecond)
	run.TickAnimations()
	select {
	case d := <-progress:
		assert.Greater(t, d, after,
			"progress should continue forward after pin: after=%v d=%v", after, d)
		assert.Less(t, d, float32(0.6),
			"progress should pace against the new longer duration: d=%v", d)
	case <-time.After(time.Second):
		t.Fatal("animation did not continue ticking")
	}
}

func TestRunner_DurationDecreasedMidAnimation(t *testing.T) {
	progress := make(chan float32, 8)
	run := &Runner{}
	a := &fyne.Animation{
		Duration: time.Second,
		Curve:    fyne.AnimationLinear,
		Tick: func(d float32) {
			progress <- d
		},
	}
	run.Start(a)

	time.Sleep(200 * time.Millisecond)
	run.TickAnimations()
	var before float32
	select {
	case before = <-progress:
	case <-time.After(time.Second):
		t.Fatal("animation was not ticked")
	}

	// Shorten Duration but keep it longer than elapsed time.
	a.Duration = 400 * time.Millisecond
	run.TickAnimations()
	var after float32
	select {
	case after = <-progress:
	case <-time.After(time.Second):
		t.Fatal("animation was not ticked after duration change")
	}
	assert.InDelta(t, before, after, 0.1,
		"progress should be preserved when Duration shrinks: before=%v after=%v", before, after)

	time.Sleep(300 * time.Millisecond)
	run.TickAnimations()
	select {
	case d := <-progress:
		assert.Equal(t, float32(1.0), d, "animation should complete on the new shorter duration")
	case <-time.After(time.Second):
		t.Fatal("animation did not complete after duration change")
	}
}

func TestRunner_DurationShortenedBelowElapsed(t *testing.T) {
	progress := make(chan float32, 4)
	run := &Runner{}
	a := &fyne.Animation{
		Duration: time.Second,
		Curve:    fyne.AnimationLinear,
		Tick: func(d float32) {
			progress <- d
		},
	}
	run.Start(a)

	time.Sleep(200 * time.Millisecond)
	run.TickAnimations()
	select {
	case <-progress:
	case <-time.After(time.Second):
		t.Fatal("animation was not ticked")
	}

	// Shorten Duration below the elapsed time.
	a.Duration = 50 * time.Millisecond
	run.TickAnimations()
	select {
	case d := <-progress:
		assert.Equal(t, float32(1.0), d,
			"animation should complete when new duration is shorter than elapsed time")
	case <-time.After(time.Second):
		t.Fatal("animation did not complete")
	}
}

func TestGLDriver_StopAnimationImmediatelyAndInsideTick(t *testing.T) {
	var wg sync.WaitGroup
	run := &Runner{}

	// stopping an animation immediately after start, should be effectively removed
	// from the internal animation list (first one is added directly to animation list)
	a := &fyne.Animation{
		Duration: time.Second,
		Tick:     func(f float32) {},
	}
	run.Start(a)
	go tick(run) // simulate a graphics draw loop
	run.Stop(a)

	run = &Runner{}
	wg = sync.WaitGroup{}

	// stopping animation inside tick function
	for i := 0; i < 10; i++ {
		wg.Add(1)
		var b *fyne.Animation
		b = &fyne.Animation{
			Duration: time.Second,
			Tick: func(d float32) {
				run.Stop(b)
				wg.Done()
			},
		}
		run.Start(b)
	}

	run = &Runner{}
	wg = sync.WaitGroup{}

	// Similar to first part, but in this time this animation should be added and then removed
	// from pendingAnimation slice.
	c := &fyne.Animation{
		Duration: time.Second,
		Tick:     func(f float32) {},
	}
	run.Start(c)
	tick(run) // simulate a graphics draw loop

	run.Stop(c)
	tick(run) // simulate a graphics draw loop

	wg.Wait()
	// animations stopped inside tick are really stopped in the next runner cycle
	time.Sleep(time.Second/60 + 100*time.Millisecond)
	run.animationMutex.RLock()
	assert.Zero(t, len(run.animations))
	run.animationMutex.RUnlock()
}
