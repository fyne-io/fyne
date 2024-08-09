package async

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	pool := Pool[int]{}

	item := pool.Get()
	assert.Equal(t, 0, item)

	item = 5
	pool.Put(item)
	assert.Equal(t, item, pool.Get())

	pool.New = func() int {
		return -1
	}

	assert.Equal(t, -1, pool.Get())
}
