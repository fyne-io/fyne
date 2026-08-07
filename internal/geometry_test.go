package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestMaxSizes(t *testing.T) {
	a := fyne.NewSize(1, 40)
	b := fyne.NewSize(30, 2)

	assert.Equal(t, fyne.NewSize(30, 40), MaxSizes(a, b))
	assert.Equal(t, a.Max(b), MaxSizes(a, b), "must match fyne.Size.Max")
}

func TestMinSizes(t *testing.T) {
	a := fyne.NewSize(1, 40)
	b := fyne.NewSize(30, 2)

	assert.Equal(t, fyne.NewSize(1, 2), MinSizes(a, b))
	assert.Equal(t, a.Min(b), MinSizes(a, b), "must match fyne.Size.Min")
}

func BenchmarkMaxSizes(b *testing.B) {
	s := fyne.NewSize(10, 10)
	for i := 0; i < b.N; i++ {
		s = MaxSizes(s, fyne.NewSize(20, 20))
	}
	if s.Width != 20 {
		b.Fatal("unexpected result")
	}
}
