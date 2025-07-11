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

var sink int

func BenchmarkPool(b *testing.B) {
	p := &Pool[int]{
		New: func() int {
			return 0
		},
	}

	b.Run("GetOnly", func(b *testing.B) {
		b.ReportAllocs()
		local := 0
		for i := 0; i < b.N; i++ {
			local = p.Get()
		}
		sink = local
	})

	b.Run("PutOnly", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			p.Put(i)
		}
	})
}
