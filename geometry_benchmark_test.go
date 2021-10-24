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

func BenchmarkPosition_Subtract(b *testing.B) {
	b.Run("Subtract()", benchmarkPositionSubtract)
	b.Run("SubtractX()", benchmarkPositionSubtractX)
	b.Run("SubtractXAndY()", benchmarkPositionSubtractXAndY)
	b.Run("SubtractY()", benchmarkPositionSubtractY)
}

func BenchmarkSize_Add(b *testing.B) {
	b.Run("Add()", benchmarkSizeAdd)
	b.Run("AddHeight()", benchmarkSizeAddHeight)
	b.Run("AddWidth()", benchmarkSizeAddWidth)
	b.Run("AddWidthAndHeight()", benchmarkSizeAddWidthAndHeight)
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

// This test prevents Position.Subtract to be simplified to `return p.SubtractXAndY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkPositionSubtract)
	subtractX := testing.Benchmark(benchmarkPositionSubtractX)
	subtractXAndY := testing.Benchmark(benchmarkPositionSubtractXAndY)
	subtractY := testing.Benchmark(benchmarkPositionSubtractY)
	assert.Less(t, subtract.NsPerOp(), int64(5))
	assert.Less(t, subtractX.NsPerOp(), int64(1))
	assert.Less(t, subtractXAndY.NsPerOp(), int64(1))
	assert.Less(t, subtractY.NsPerOp(), int64(1))
}

// This test prevents Size.Add to be simplified to `return s.AddWidthAndHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkSizeAdd)
	addHeight := testing.Benchmark(benchmarkSizeAddHeight)
	addWidthAndHeight := testing.Benchmark(benchmarkSizeAddWidth)
	addWidth := testing.Benchmark(benchmarkSizeAddWidthAndHeight)
	assert.Less(t, add.NsPerOp(), int64(5))
	assert.Less(t, addHeight.NsPerOp(), int64(1))
	assert.Less(t, addWidthAndHeight.NsPerOp(), int64(1))
	assert.Less(t, addWidth.NsPerOp(), int64(1))
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

func benchmarkPositionSubtract(b *testing.B) {
	pos1 := fyne.NewPos(10, 10)
	pos2 := fyne.NewPos(25, 25)
	for n := 0; n < b.N; n++ {
		pos1.Subtract(pos2)
	}
}

func benchmarkPositionSubtractX(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos.SubtractX(25)
	}
}

func benchmarkPositionSubtractXAndY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos.SubtractXAndY(25, 25)
	}
}

func benchmarkPositionSubtractY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos.SubtractY(25)
	}
}

func benchmarkSizeAdd(b *testing.B) {
	size1 := fyne.NewSize(10, 10)
	size2 := fyne.NewSize(25, 25)
	for n := 0; n < b.N; n++ {
		size1.Add(size2)
	}
}

func benchmarkSizeAddWidth(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size.AddWidth(25)
	}
}

func benchmarkSizeAddWidthAndHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size.AddWidthAndHeight(25, 25)
	}
}

func benchmarkSizeAddHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size.AddHeight(25)
	}
}
