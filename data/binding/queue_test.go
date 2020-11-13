package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueItem(t *testing.T) {
	called := 0
	queueItem(func(DataItem) {
		called++
	}, nil)
	queueItem(func(DataItem) {
		called++
	}, nil)

	waitForItems()
	assert.Equal(t, 2, called)
}
