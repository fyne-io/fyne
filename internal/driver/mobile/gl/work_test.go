//go:build cgo && (darwin || linux || openbsd || freebsd)

package gl

import (
	"sync"
	"testing"
	"time"
)

// newTestContext mirrors NewContext's signalling wiring without the cgo GL
// calls, so these tests need no GL context. retvalue stays nil because no test
// here enqueues a blocking call.
func newTestContext() *context {
	ctx := &context{
		workAvailable: make(chan struct{}, 1),
		work:          make(chan call, workQueueLen),
	}
	ctx.consumerIdle.Store(true)
	return ctx
}

// TestWakeupsCoalescePerBurst covers N calls producing one wakeup, not N.
func TestWakeupsCoalescePerBurst(t *testing.T) {
	ctx := newTestContext()

	// Below the queue capacity, so the test stays synchronous. Backpressure is
	// covered by TestConcurrentProducersNeverStrandWork.
	const frameCalls = 200
	for i := 0; i < frameCalls; i++ {
		ctx.enqueue(call{})
	}

	if !ctx.HasWork() {
		t.Fatal("HasWork() = false after enqueueing work")
	}
	if got := len(ctx.workAvailable); got != 1 {
		t.Errorf("pending wakeups after %d enqueues = %d, want 1", frameCalls, got)
	}
	if got := len(ctx.work); got != frameCalls {
		t.Errorf("queued calls = %d, want %d", got, frameCalls)
	}
}

// TestSignalIsIdempotentWhileWorkerAwake covers that calls made while the worker
// is awake queue no further wakeups.
func TestSignalIsIdempotentWhileWorkerAwake(t *testing.T) {
	ctx := newTestContext()

	ctx.enqueue(call{})
	ctx.enqueue(call{})
	<-ctx.workAvailable // the worker is now awake and busy

	for i := 0; i < 50; i++ {
		ctx.enqueue(call{})
	}
	if got := len(ctx.workAvailable); got != 0 {
		t.Errorf("pending wakeups while worker is awake = %d, want 0", got)
	}

	// Publishing idle then re-checking must find the queued calls.
	if ctx.drained() {
		t.Error("drained() = true although calls are queued, work would be stranded")
	}
}

// TestNoLostWakeupOnGoingIdle covers the ordering in drained: a call enqueued
// after the worker's last receive is not signalled, so the worker has to find it.
func TestNoLostWakeupOnGoingIdle(t *testing.T) {
	ctx := newTestContext()

	ctx.enqueue(call{})
	<-ctx.workAvailable
	<-ctx.work // the worker consumed its only call

	if ctx.consumerIdle.Load() {
		t.Fatal("test setup: worker should be marked awake after being woken")
	}

	// This producer sees consumerIdle == false, so it must not signal.
	ctx.enqueue(call{})
	if got := len(ctx.workAvailable); got != 0 {
		t.Fatalf("producer signalled while worker awake, wakeups = %d, want 0", got)
	}

	if ctx.drained() {
		t.Fatal("drained() = true with a call enqueued while going idle, work is lost")
	}
	if !ctx.HasWork() {
		t.Error("HasWork() = false, the call disappeared")
	}
}

// TestIdlePublishesStateForNextEnqueue covers that a parked worker is woken again.
func TestIdlePublishesStateForNextEnqueue(t *testing.T) {
	ctx := newTestContext()

	ctx.enqueue(call{})
	<-ctx.workAvailable
	<-ctx.work

	if !ctx.drained() {
		t.Fatal("drained() = false on an empty queue")
	}
	if !ctx.consumerIdle.Load() {
		t.Error("consumerIdle = false after going idle, producers would not wake the worker")
	}

	ctx.enqueue(call{})
	if got := len(ctx.workAvailable); got != 1 {
		t.Errorf("pending wakeups after enqueue to a parked worker = %d, want 1", got)
	}
}

// TestConcurrentProducersNeverStrandWork drives several producers while a worker
// drains. The total exceeds the queue capacity, so the producers block and are
// woken by the drain. The worker only ever waits on workAvailable, so the
// watchdog turns a lost wakeup into a failure instead of a hang.
func TestConcurrentProducersNeverStrandWork(t *testing.T) {
	ctx := newTestContext()

	const (
		producers      = 4
		callsPerWriter = 500
		total          = producers * callsPerWriter
	)

	var (
		mu   sync.Mutex
		seen int
	)
	workerDone := make(chan struct{})
	go func() {
		defer close(workerDone)
		// Mirrors DoWork: drain, then decide whether to go idle.
		for {
			for {
				select {
				case <-ctx.work:
					mu.Lock()
					seen++
					mu.Unlock()
					continue
				default:
				}
				break
			}

			if ctx.drained() {
				mu.Lock()
				finished := seen >= total
				mu.Unlock()
				if finished {
					return
				}
				select {
				case <-ctx.workAvailable:
				case <-time.After(5 * time.Second):
					if ctx.HasWork() {
						t.Errorf("worker was not woken although %d calls are queued", len(ctx.work))
						return
					}
				}
			}
		}
	}()

	var wg sync.WaitGroup
	for p := 0; p < producers; p++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < callsPerWriter; i++ {
				ctx.enqueue(call{})
			}
		}()
	}
	wg.Wait()

	select {
	case <-workerDone:
	case <-time.After(30 * time.Second):
		t.Fatal("worker did not finish, work was stranded")
	}

	mu.Lock()
	defer mu.Unlock()
	if seen != total {
		t.Errorf("worker processed %d calls, want %d", seen, total)
	}
}
