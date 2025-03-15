package async

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mainRoutineID uint64

func init() {
	mainRoutineID = goroutineID()
}

func TestGoroutineID(t *testing.T) {
	assert.Equal(t, uint64(1), mainRoutineID)

	var childID1, childID2 uint64
	testID1 := goroutineID()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		childID1 = goroutineID()
		wg.Done()
	}()
	go func() {
		childID2 = goroutineID()
		wg.Done()
	}()
	wg.Wait()
	testID2 := goroutineID()

	assert.Equal(t, testID1, testID2)
	assert.Positive(t, childID1)
	assert.NotEqual(t, testID1, childID1)
	assert.Positive(t, childID2)
	assert.NotEqual(t, childID1, childID2)
}
