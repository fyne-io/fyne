//go:build !ci || !darwin
// +build !ci !darwin

package animation

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestGLDriver_StartAnimation(t *testing.T) {
	done := make(chan float32)
	run := &Runner{}
	a := &fyne.Animation{
		Duration: time.Millisecond * 100,
		Tick: func(d float32) {
			done <- d
		}}

	run.Start(a)
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
		}}

	run.Start(a)
	select {
	case d := <-done:
		assert.Greater(t, d, float32(0))
	case <-time.After(time.Second):
		t.Error("animation was not ticked")
	}
	a.Stop()
	run.animationMutex.Lock()
	assert.True(t, a.Stopped(), "animation was not stopped")
	run.animationMutex.Unlock()
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
	a.Stop()

	// stopping animation inside tick function
	for i := 0; i < 10; i++ {
		wg.Add(1)
		var b *fyne.Animation
		b = &fyne.Animation{
			Duration: time.Second,
			Tick: func(d float32) {
				b.Stop()
				wg.Done()
			}}
		run.Start(b)
	}

	// Similar to first part, but in this time this animation should be added and then removed
	// from pendingAnimation slice.
	c := &fyne.Animation{
		Duration: time.Second,
		Tick:     func(f float32) {},
	}
	run.Start(c)
	c.Stop()

	wg.Wait()
	// animations stopped inside tick are really stopped in the next runner cycle
	time.Sleep(time.Second/60 + 100*time.Millisecond)
	run.animationMutex.Lock()
	assert.Zero(t, len(run.animations))
	run.animationMutex.Unlock()
}
