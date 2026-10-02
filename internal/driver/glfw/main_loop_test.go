package glfw

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNextMainLoopEvent(t *testing.T) {
	for _, native := range []bool{false, true} {
		t.Run(map[bool]string{false: "channel", true: "native"}[native], func(t *testing.T) {
			t.Run("function", func(t *testing.T) {
				functions := make(chan funcData, 1)
				called := false
				functions <- funcData{f: func() { called = true }}
				var wait func(time.Duration)
				if native {
					wait = func(time.Duration) { t.Fatal("waited with a function ready") }
				}
				event := nextMainLoopEvent(nil, functions, nil, time.Now(), wait)
				require.Equal(t, mainLoopFunction, event.kind)
				event.function.f()
				require.True(t, called)
			})
			t.Run("frame", func(t *testing.T) {
				frames := make(chan time.Time, 1)
				tick := time.Now()
				frames <- tick
				var wait func(time.Duration)
				if native {
					wait = func(time.Duration) { t.Fatal("waited with a frame ready") }
				}
				event := nextMainLoopEvent(nil, nil, frames, tick, wait)
				require.Equal(t, mainLoopFrame, event.kind)
				require.Equal(t, tick, event.frameTime)
			})
			t.Run("quit", func(t *testing.T) {
				done := make(chan struct{})
				close(done)
				var wait func(time.Duration)
				if native {
					wait = func(time.Duration) { t.Fatal("waited after shutdown") }
				}
				event := nextMainLoopEvent(done, nil, nil, time.Now(), wait)
				require.Equal(t, mainLoopQuit, event.kind)
			})
		})
	}
}

func TestNextMainLoopEventNativeWake(t *testing.T) {
	functions := make(chan funcData, 1)
	nextFrame := time.Now().Add(time.Second)
	waits := 0
	event := nextMainLoopEvent(nil, functions, nil, nextFrame, func(timeout time.Duration) {
		require.Positive(t, timeout)
		require.LessOrEqual(t, timeout, time.Second)
		waits++
		functions <- funcData{f: func() {}}
	})
	require.Equal(t, mainLoopFunction, event.kind)
	require.Equal(t, 1, waits)
}

func TestNextMainLoopEventNativeFrameDeadline(t *testing.T) {
	frames := make(chan time.Time, 1)
	nextFrame := time.Now().Add(time.Second)
	event := nextMainLoopEvent(nil, nil, frames, nextFrame, func(timeout time.Duration) {
		require.Positive(t, timeout)
		require.LessOrEqual(t, timeout, time.Second)
		frames <- nextFrame
	})
	require.Equal(t, mainLoopFrame, event.kind)
	require.Equal(t, nextFrame, event.frameTime)
}

func TestMainLoopWakeLifecycle(t *testing.T) {
	var wake mainLoopWake
	var calls atomic.Int32
	wake.notify()
	wake.start(func() { calls.Add(1) })
	wake.notify()
	require.EqualValues(t, 1, calls.Load())
	wake.stop()
	wake.notify()
	require.EqualValues(t, 1, calls.Load())
}

func TestMainLoopWakeStopWaitsForNotification(t *testing.T) {
	var wake mainLoopWake
	entered := make(chan struct{})
	release := make(chan struct{})
	notified := make(chan struct{})
	wake.start(func() {
		close(entered)
		<-release
	})
	go func() {
		wake.notify()
		close(notified)
	}()
	<-entered
	stopStarted := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		close(stopStarted)
		wake.stop()
		close(stopped)
	}()
	<-stopStarted
	select {
	case <-stopped:
		t.Fatal("stopped with a native wakeup still in flight")
	case <-time.After(10 * time.Millisecond):
	}
	close(release)
	<-stopped
	<-notified
	wake.notify()
}
