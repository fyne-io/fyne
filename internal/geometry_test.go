package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

// The helpers exist only to mirror fyne.Size.Max/Min without boxing the
// argument into their Vector2 parameter, so each test checks the element-wise
// result and that the helper stays in sync with the public method it shadows.

func TestMaxSizes(t *testing.T) {
	// width taken from b, height from a — the maximum is element-wise,
	// not whichever whole Size is "bigger"
	a := fyne.NewSize(1, 40)
	b := fyne.NewSize(30, 2)

	assert.Equal(t, fyne.NewSize(30, 40), MaxSizes(a, b))
	assert.Equal(t, a.Max(b), MaxSizes(a, b), "must match fyne.Size.Max")
}

func TestMinSizes(t *testing.T) {
	// width taken from a, height from b — element-wise, as above
	a := fyne.NewSize(1, 40)
	b := fyne.NewSize(30, 2)

	assert.Equal(t, fyne.NewSize(1, 2), MinSizes(a, b))
	assert.Equal(t, a.Min(b), MinSizes(a, b), "must match fyne.Size.Min")
}

// BenchmarkMaxSizes guards the reason the helpers exist: 0 allocs/op.
// fyne.Size.Max costs one heap allocation per call (interface boxing); if a
// change reintroduces an allocation here, the per-frame layout win is gone.
func BenchmarkMaxSizes(b *testing.B) {
	s := fyne.NewSize(10, 10)
	for i := 0; i < b.N; i++ {
		s = MaxSizes(s, fyne.NewSize(20, 20))
	}
	// use the result so the loop cannot be optimized away
	if s.Width != 20 {
		b.Fatal("unexpected result")
	}
}
