// +build !ci

package fyne_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func BenchmarkPosition_Add(b *testing.B) {
	b.Run("Add()", benchmarkPositionAdd)
	b.Run("AddXAndY()", benchmarkPositionAddXAndY)
}

func BenchmarkPosition_Subtract(b *testing.B) {
	b.Run("Subtract()", benchmarkPositionSubtract)
	b.Run("SubtractXAndY()", benchmarkPositionSubtractXAndY)
}

func BenchmarkSize_Add(b *testing.B) {
	b.Run("Add()", benchmarkSizeAdd)
	b.Run("AddWidthAndHeight()", benchmarkSizeAddWidthAndHeight)
}

func BenchmarkSize_Subtract(b *testing.B) {
	b.Run("Subtract()", benchmarkSizeSubtract)
	b.Run("SubtractWidthAndHeight()", benchmarkSizeSubtractWidthAndHeight)
}

// This test prevents Position.Add to be simplified to `return p.AddXAndY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkPositionAdd)
	addXAndY := testing.Benchmark(benchmarkPositionAddXAndY)
	assert.Less(t, add.NsPerOp(), int64(5))
	assert.Less(t, addXAndY.NsPerOp(), int64(1))
}

// This test prevents Position.Subtract to be simplified to `return p.SubtractXAndY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkPositionSubtract)
	subtractXAndY := testing.Benchmark(benchmarkPositionSubtractXAndY)
	assert.Less(t, subtract.NsPerOp(), int64(5))
	assert.Less(t, subtractXAndY.NsPerOp(), int64(1))
}

// This test prevents Size.Add to be simplified to `return s.AddWidthAndHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkSizeAdd)
	addWidthAndHeight := testing.Benchmark(benchmarkSizeAddWidthAndHeight)
	assert.Less(t, add.NsPerOp(), int64(5))
	assert.Less(t, addWidthAndHeight.NsPerOp(), int64(1))
}

// This test prevents Size.Subtract to be simplified to `return s.SubtractWidthAndHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkSizeSubtract)
	subtractWidthAndHeight := testing.Benchmark(benchmarkSizeSubtractWidthAndHeight)
	assert.Less(t, subtract.NsPerOp(), int64(5))
	assert.Less(t, subtractWidthAndHeight.NsPerOp(), int64(1))
}

func benchmarkPositionAdd(b *testing.B) {
	pos1 := fyne.NewPos(10, 10)
	pos2 := fyne.NewPos(25, 25)
	for n := 0; n < b.N; n++ {
		pos1.Add(pos2)
	}
}

func benchmarkPositionAddXAndY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos.AddXAndY(25, 25)
	}
}

func benchmarkPositionSubtract(b *testing.B) {
	pos1 := fyne.NewPos(10, 10)
	pos2 := fyne.NewPos(25, 25)
	for n := 0; n < b.N; n++ {
		pos1.Subtract(pos2)
	}
}

func benchmarkPositionSubtractXAndY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos.SubtractXAndY(25, 25)
	}
}

func benchmarkSizeAdd(b *testing.B) {
	size1 := fyne.NewSize(10, 10)
	size2 := fyne.NewSize(25, 25)
	for n := 0; n < b.N; n++ {
		size1.Add(size2)
	}
}

func benchmarkSizeAddWidthAndHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size.AddWidthAndHeight(25, 25)
	}
}

func benchmarkSizeSubtract(b *testing.B) {
	size1 := fyne.NewSize(10, 10)
	size2 := fyne.NewSize(25, 25)
	for n := 0; n < b.N; n++ {
		size1.Subtract(size2)
	}
}

func benchmarkSizeSubtractWidthAndHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size.SubtractWidthAndHeight(25, 25)
	}
}
