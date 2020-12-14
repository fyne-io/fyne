package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
