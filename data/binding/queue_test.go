package binding

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"

	"fyne.io/fyne/v2/internal/async"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	testQueueLazyInit()
	os.Exit(m.Run())
}

func TestQueueItem(t *testing.T) {
	called := 0
	queueItem(func() {
		called++
	})
	queueItem(func() {
		called++
	})

	waitForItems()
	assert.Equal(t, 2, called)
}

func TestMakeInfiniteQueue(t *testing.T) {
	var wg sync.WaitGroup
	queue := async.NewUnboundedFuncChan()

	wg.Add(1)
	c := 0
	go func() {
		for range queue.Out() {
			c++
		}
		wg.Done()
	}()

	for i := 0; i < 2048; i++ {
		queue.In() <- func() {}
	}
	close(queue.In())

	wg.Wait()
	assert.Equal(t, 2048, c)
}

type mainTest struct {
	name string
}

func newMainTest(name string) *mainTest {
	t := &mainTest{name: name}
	fmt.Printf("=== RUN   %s\n", name)
	return t
}

func (m *mainTest) exitError() {
	fmt.Printf("--- FAIL: %s (0.00s)\n", m.name)
	fmt.Println("FAIL")
	os.Exit(1)
}

func (m *mainTest) Errorf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		fmt.Fprintf(os.Stderr, "\t%s:%d\n", file, line)
	}
	fmt.Fprintf(os.Stderr, strings.TrimLeft(format, "\n"), args...)
	m.exitError()
}

func testQueueLazyInit() {
	t := newMainTest("TestQueueLazyInit")
	initialGoRoutines := runtime.NumGoroutine()

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go queueItem(func() { wg.Done() })
	}
	wg.Wait()
	assert.Equal(t, initialGoRoutines+2, runtime.NumGoroutine())
}
