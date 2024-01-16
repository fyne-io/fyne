package async_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/async"
	"github.com/stretchr/testify/assert"
)

var globalSize fyne.Size

func BenchmarkAtomicLoadAndStore(b *testing.B) {
	local := fyne.Size{}
	size := async.Size{}

	// Each loop iteration replicates how .Resize() in BaseWidget checks size changes.
	for i := 0; i < b.N; i++ {
		new := fyne.Size{Width: float32(i), Height: float32(i)}
		if new == size.Load() {
			continue
		}

		size.Store(new)
	}

	globalSize = local
}

func TestSize(t *testing.T) {
	size := async.Size{}
	assert.Equal(t, fyne.Size{}, size.Load())

	square := fyne.NewSquareSize(100)
	size.Store(square)
	assert.Equal(t, square, size.Load())

	uneven := fyne.NewSize(125, 600)
	size.Store(uneven)
	assert.Equal(t, uneven, size.Load())

	floats := fyne.NewSize(-22.565, 133.333)
	size.Store(floats)
	assert.Equal(t, floats, size.Load())
}

func TestPosition(t *testing.T) {
	pos := async.Position{}
	assert.Equal(t, fyne.Position{}, pos.Load())

	even := fyne.NewSquareOffsetPos(100)
	pos.Store(even)
	assert.Equal(t, even, pos.Load())

	uneven := fyne.NewPos(125, 600)
	pos.Store(uneven)
	assert.Equal(t, uneven, pos.Load())

	floats := fyne.NewPos(-22.565, 133.333)
	pos.Store(floats)
	assert.Equal(t, floats, pos.Load())
}
