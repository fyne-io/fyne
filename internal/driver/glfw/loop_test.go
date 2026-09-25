//go:build !no_glfw && !mobile

package glfw

import (
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	intapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/async"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const runLifecycleShutdownHelperEnv = "FYNE_TEST_RUN_LIFECYCLE_SHUTDOWN_HELPER"

var (
	runLifecycleShutdownReady    = make(chan struct{})
	runLifecycleShutdownReturned = make(chan struct{})
)

type driverOverrideApp struct {
	fyne.App
	driver fyne.Driver
}

func (a *driverOverrideApp) Driver() fyne.Driver {
	return a.driver
}

func TestRunLifecycleShutdown(t *testing.T) {
	cmd := exec.Command(os.Args[0], //nolint:gosec // os.Args[0] is the current test binary.
		"-test.run=^TestRunLifecycleShutdownHelper$",
		"-test.timeout=10s")
	cmd.Env = append(os.Environ(), runLifecycleShutdownHelperEnv+"=1")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", output)
}

func TestRunLifecycleShutdownHelper(t *testing.T) {
	if os.Getenv(runLifecycleShutdownHelperEnv) != "1" {
		t.Skip("only run as a subprocess with fresh package state")
	}

	baseApp := test.NewApp()
	fyne.SetCurrentApp(&driverOverrideApp{App: baseApp, driver: d})

	lifecycle := baseApp.Lifecycle().(*intapp.Lifecycle)
	lifecycle.InitEventQueue()
	go lifecycle.RunEventQueue(d.DoFromGoroutine)

	var (
		orderMu sync.Mutex
		order   []string

		userCalls     atomic.Int32
		internalCalls atomic.Int32
		cancelledRan  atomic.Bool
	)
	doAndWaitReturned := make(chan struct{})

	record := func(name string, calls *atomic.Int32) {
		assert.False(t, drained.Load(), "%s callback ran after the main queue was drained", name)
		assert.True(t, async.IsMainGoroutine(), "%s callback did not run on the main goroutine", name)
		calls.Add(1)

		orderMu.Lock()
		order = append(order, name)
		orderMu.Unlock()
	}

	lifecycle.SetOnStopped(func() {
		record("user", &userCalls)

		go func() {
			fyne.DoAndWait(func() {
				cancelledRan.Store(true)
			})
			close(doAndWaitReturned)
		}()

		deadline := time.Now().Add(time.Second)
		for len(funcQueue.Out()) == 0 && time.Now().Before(deadline) {
			time.Sleep(time.Millisecond)
		}
		assert.NotZero(t, len(funcQueue.Out()), "DoAndWait was not queued during shutdown")
	})
	lifecycle.SetOnStoppedHookExecuted(func() {
		record("internal", &internalCalls)
	})

	go func() {
		for !running.Load() {
			time.Sleep(time.Millisecond)
		}
		d.Quit()
	}()

	close(runLifecycleShutdownReady)
	<-runLifecycleShutdownReturned

	require.Eventually(t, func() bool {
		select {
		case <-doAndWaitReturned:
			return true
		default:
			return false
		}
	}, time.Second, time.Millisecond, "shutdown-pending DoAndWait did not return")
	assert.False(t, cancelledRan.Load(), "shutdown-pending DoAndWait callback was executed")
	assert.EqualValues(t, 1, userCalls.Load())
	assert.EqualValues(t, 1, internalCalls.Load())

	orderMu.Lock()
	assert.Equal(t, []string{"user", "internal"}, order)
	orderMu.Unlock()
}

// BenchmarkRunOnMain measures the cost of calling a function
// on the main thread.
func BenchmarkRunOnMain(b *testing.B) {
	f := func() {}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runOnMain(f)
	}
}

// BenchmarkRunOnDraw measures the cost of calling a function
// on the draw thread.
func BenchmarkRunOnDraw(b *testing.B) {
	f := func() {}
	w := createWindow("Test")
	w.create()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.RunWithContext(f)
	}
}
