package async

import (
	"math"
	"sync/atomic"

	"fyne.io/fyne/v2"
)

// Position is an atomic version of fyne.Position.
// Loads and stores are guaranteed to happen using a single atomic operation.
type Position struct {
	pos atomic.Uint64
}

// Load performs an atomic load on the fyne.Position value.
func (p *Position) Load() fyne.Position {
	return fyne.NewPos(twoFloat32FromUint64(p.pos.Load()))
}

// Store performs an atomic store on the fyne.Position value.
func (p *Position) Store(pos fyne.Position) {
	p.pos.Store(uint64fromTwoFloat32(pos.X, pos.Y))
}

// Size is an atomic version of fyne.Size.
// Loads and stores are guaranteed to happen using a single atomic operation.
type Size struct {
	size atomic.Uint64
}

// Load performs an atomic load on the fyne.Size value.
func (s *Size) Load() fyne.Size {
	return fyne.NewSize(twoFloat32FromUint64(s.size.Load()))
}

// Store performs an atomic store on the fyne.Size value.
func (s *Size) Store(size fyne.Size) {
	s.size.Store(uint64fromTwoFloat32(size.Width, size.Height))
}

func uint64fromTwoFloat32(a, b float32) uint64 {
	x := uint64(math.Float32bits(a))
	y := uint64(math.Float32bits(b))
	return (y << 32) | x
}

func twoFloat32FromUint64(combined uint64) (float32, float32) {
	x := uint32(combined)
	y := uint32(combined >> 32)
	return math.Float32frombits(x), math.Float32frombits(y)
}
