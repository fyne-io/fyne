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
	assert.Less(t, add.NsPerOp()/addXAndY.NsPerOp(), int64(5))
}

// This test prevents Position.Subtract to be simplified to `return p.SubtractXAndY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkPositionSubtract)
	subtractXAndY := testing.Benchmark(benchmarkPositionSubtractXAndY)
	assert.Less(t, subtract.NsPerOp()/subtractXAndY.NsPerOp(), int64(5))
}

// This test prevents Size.Add to be simplified to `return s.AddWidthAndHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkSizeAdd)
	addWidthAndHeight := testing.Benchmark(benchmarkSizeAddWidthAndHeight)
	assert.Less(t, add.NsPerOp()/addWidthAndHeight.NsPerOp(), int64(5))
}

// This test prevents Size.Subtract to be simplified to `return s.SubtractWidthAndHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkSizeSubtract)
	subtractWidthAndHeight := testing.Benchmark(benchmarkSizeSubtractWidthAndHeight)
	assert.Less(t, subtract.NsPerOp()/subtractWidthAndHeight.NsPerOp(), int64(5))
}

var benchmarkResult interface{}

func benchmarkPositionAdd(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.Add(fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = pos
}

func benchmarkPositionAddXAndY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.AddXAndY(float32(n), float32(n))
	}
	benchmarkResult = pos
}

func benchmarkPositionSubtract(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.Subtract(fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = pos
}

func benchmarkPositionSubtractXAndY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.SubtractXAndY(float32(n), float32(n))
	}
	benchmarkResult = pos
}

func benchmarkSizeAdd(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.Add(fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = size
}

func benchmarkSizeAddWidthAndHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.AddWidthAndHeight(float32(n), float32(n))
	}
	benchmarkResult = size
}

func benchmarkSizeSubtract(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.Subtract(fyne.NewSize(float32(n), float32(n)))
	}
	benchmarkResult = size
}

func benchmarkSizeSubtractWidthAndHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.SubtractWidthAndHeight(float32(n), float32(n))
	}
	benchmarkResult = size
}
