package glfw

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNextMainLoopEventReady(t *testing.T) {
	for _, native := range []bool{false, true} {
		t.Run(map[bool]string{false: "blocking", true: "native"}[native], func(t *testing.T) {
			var wake *mainLoopWake
			if native {
				wake = &mainLoopWake{wait: func(time.Duration) { t.Fatal("waited with an event ready") }}
			}
			done := make(chan struct{})
			close(done)
			require.Equal(t, mainLoopQuit, nextMainLoopEvent(done, nil, nil, time.Now(), wake).kind)

			functions := make(chan funcData, 1)
			functions <- funcData{}
			require.Equal(t, mainLoopFunction, nextMainLoopEvent(nil, functions, nil, time.Now(), wake).kind)

			frames := make(chan time.Time, 1)
			tick := time.Now()
			frames <- tick
			event := nextMainLoopEvent(nil, nil, frames, tick, wake)
			require.Equal(t, mainLoopFrame, event.kind)
			require.Equal(t, tick, event.frameTime)
		})
	}
}

func TestNextMainLoopEventWaitsUntilNextFrame(t *testing.T) {
	functions := make(chan funcData, 1)
	var timeout time.Duration
	wake := &mainLoopWake{wait: func(d time.Duration) {
		timeout = d
		functions <- funcData{}
	}}
	event := nextMainLoopEvent(nil, functions, nil, time.Now().Add(time.Second), wake)
	require.Equal(t, mainLoopFunction, event.kind)
	require.InDelta(t, time.Second, timeout, float64(100*time.Millisecond))
}

func TestNextMainLoopEventOverdueFrameDoesNotWait(t *testing.T) {
	frames := make(chan time.Time)
	wake := &mainLoopWake{wait: func(time.Duration) { t.Error("waited for native events with a frame overdue") }}
	tick := time.Now()
	go func() {
		time.Sleep(10 * time.Millisecond)
		frames <- tick
	}()
	event := nextMainLoopEvent(nil, nil, frames, tick.Add(-time.Millisecond), wake)
	require.Equal(t, mainLoopFrame, event.kind)
}

func TestMainLoopWake(t *testing.T) {
	var wake mainLoopWake
	calls := 0
	wake.notify()
	wake.start(func() { calls++ })
	wake.notify()
	wake.notify()
	require.Equal(t, 1, calls, "wakeups should coalesce until rearmed")
	wake.rearm()
	wake.notify()
	require.Equal(t, 2, calls)
	wake.stop()
	wake.rearm()
	wake.notify()
	require.Equal(t, 2, calls)
}
