package async

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	pool := Pool[int]{}

	item := pool.Get()
	assert.Equal(t, 0, item)

	pool.New = func() int {
		return -1
	}

	assert.Equal(t, -1, pool.Get())
}
