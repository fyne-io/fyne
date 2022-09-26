//go:build !ci
// +build !ci

package fyne_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func BenchmarkPosition_Add(b *testing.B) {
	b.Run("Add()", benchmarkPositionAdd)
	b.Run("AddXY()", benchmarkPositionAddXY)
}

func BenchmarkPosition_Subtract(b *testing.B) {
	b.Run("Subtract()", benchmarkPositionSubtract)
	b.Run("SubtractXY()", benchmarkPositionSubtractXY)
}

func BenchmarkSize_Add(b *testing.B) {
	b.Run("Add()", benchmarkSizeAdd)
	b.Run("AddWidthHeight()", benchmarkSizeAddWidthHeight)
}

func BenchmarkSize_Subtract(b *testing.B) {
	b.Run("Subtract()", benchmarkSizeSubtract)
	b.Run("SubtractWidthHeight()", benchmarkSizeSubtractWidthHeight)
}

// This test prevents Position.Add to be simplified to `return p.AddXY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkPositionAdd)
	slowAdd := testing.Benchmark(benchmarkSlowPositionAdd)
	assert.Less(t, nsPerOpPrecise(add)*2, nsPerOpPrecise(slowAdd))
}

// This test prevents Position.Subtract to be simplified to `return p.SubtractXY(v.Components())`
// because this slows down the speed by factor 10.
func TestPosition_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkPositionSubtract)
	slowSubtract := testing.Benchmark(benchmarkSlowPositionSubtract)
	assert.Less(t, nsPerOpPrecise(subtract)*2, nsPerOpPrecise(slowSubtract))
}

// This test prevents Size.Add to be simplified to `return s.AddWidthHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Add_Speed(t *testing.T) {
	add := testing.Benchmark(benchmarkSizeAdd)
	slowAdd := testing.Benchmark(benchmarkSlowSizeAdd)
	assert.Less(t, nsPerOpPrecise(add)*2, nsPerOpPrecise(slowAdd))
}

// This test prevents Size.Subtract to be simplified to `return s.SubtractWidthHeight(v.Components())`
// because this slows down the speed by factor 10.
func TestSize_Subtract_Speed(t *testing.T) {
	subtract := testing.Benchmark(benchmarkSizeSubtract)
	slowSubtract := testing.Benchmark(benchmarkSlowSizeSubtract)
	assert.Less(t, nsPerOpPrecise(subtract)*2, nsPerOpPrecise(slowSubtract))
}

var benchmarkResult interface{}

func benchmarkPositionAdd(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.Add(fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = pos
}

func benchmarkPositionAddXY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.AddXY(float32(n), float32(n))
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

func benchmarkPositionSubtractXY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = pos.SubtractXY(float32(n), float32(n))
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

func benchmarkSizeAddWidthHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.AddWidthHeight(float32(n), float32(n))
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

func benchmarkSizeSubtractWidthHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = size.SubtractWidthHeight(float32(n), float32(n))
	}
	benchmarkResult = size
}

func benchmarkSlowPositionAdd(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = slowPositionAdd(pos, fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = pos
}

func benchmarkSlowPositionSubtract(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; n < b.N; n++ {
		pos = slowPositionSubtract(pos, fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkResult = pos
}

func benchmarkSlowSizeAdd(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = slowSizeAdd(size, fyne.NewSize(float32(n), float32(n)))
	}
	benchmarkResult = size
}

func benchmarkSlowSizeSubtract(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; n < b.N; n++ {
		size = slowSizeSubtract(size, fyne.NewSize(float32(n), float32(n)))
	}
	benchmarkResult = size
}

func nsPerOpPrecise(b testing.BenchmarkResult) float64 {
	return float64(b.T.Nanoseconds()) / float64(b.N)
}

func slowPositionAdd(p fyne.Position, v fyne.Vector2) fyne.Position {
	x, y := v.Components()
	return p.AddXY(x, y)
}

func slowPositionSubtract(p fyne.Position, v fyne.Vector2) fyne.Position {
	x, y := v.Components()
	return p.SubtractXY(x, y)
}

func slowSizeAdd(s fyne.Size, v fyne.Vector2) fyne.Size {
	w, h := v.Components()
	return s.AddWidthHeight(w, h)
}

func slowSizeSubtract(s fyne.Size, v fyne.Vector2) fyne.Size {
	w, h := v.Components()
	return s.SubtractWidthHeight(w, h)
}
