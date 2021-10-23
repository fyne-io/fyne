// +build !ci

package fyne_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func BenchmarkPosition_Add(b *testing.B) {
	b.Run("Add()", benchmarkPositionAdd)
	b.Run("AddX()", benchmarkPositionAddX)
	b.Run("AddXAndY()", benchmarkPositionAddXAndY)
	b.Run("AddY()", benchmarkPositionAddY)
}

// This test prevents Position.Add to be simplified to `return p.AddXAndY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkPositionAdd)
	addX := testing.Benchmark(benchmarkPositionAddX)
	addXAndY := testing.Benchmark(benchmarkPositionAddXAndY)
	addY := testing.Benchmark(benchmarkPositionAddY)
	assert.Less(t, add.NsPerOp(), int64(5))
	assert.Less(t, addX.NsPerOp(), int64(1))
	assert.Less(t, addXAndY.NsPerOp(), int64(1))
	assert.Less(t, addY.NsPerOp(), int64(1))
}

func benchmarkPositionAdd(b *testing.B) {
	pos1 := fyne.NewPos(10, 10)
	pos2 := fyne.NewPos(25, 25)
	for n := 0; n < b.N; n++ {
		pos1.Add(pos2)
	}
}

func benchmarkPositionAddX(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos.AddX(25)
	}
}

func benchmarkPositionAddXAndY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos.AddXAndY(25, 25)
	}
}

func benchmarkPositionAddY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos.AddY(25)
	}
}
