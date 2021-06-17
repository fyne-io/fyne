package binding

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

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
	in, out := makeInfiniteQueue()

	wg.Add(1)
	c := 0
	go func() {
		for range out {
			c++
		}
		wg.Done()
	}()

	for i := 0; i < 2048; i++ {
		in <- &itemData{}
	}
	close(in)

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
	assert.Nil(t, itemQueueIn)
	assert.Nil(t, itemQueueOut)

	initialGoRoutines := runtime.NumGoroutine()

	for i := 0; i < 1000; i++ {
		go queueItem(func() {})
	}

	waitForItems()
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, initialGoRoutines+2, runtime.NumGoroutine())
	assert.NotNil(t, itemQueueIn)
	assert.NotNil(t, itemQueueOut)
}
