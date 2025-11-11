package fyne_test

import (
	"testing"

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

var (
	benchmarkSize fyne.Size
	benchmarkPos  fyne.Position
)

func benchmarkPositionAdd(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; b.Loop(); n++ {
		pos = pos.Add(fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkPos = pos
}

func benchmarkPositionAddXY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; b.Loop(); n++ {
		pos = pos.AddXY(float32(n), float32(n))
	}
	benchmarkPos = pos
}

func benchmarkPositionSubtract(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; b.Loop(); n++ {
		pos = pos.Subtract(fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkPos = pos
}

func benchmarkPositionSubtractXY(b *testing.B) {
	pos := fyne.NewPos(10, 10)
	for n := 0; b.Loop(); n++ {
		pos = pos.SubtractXY(float32(n), float32(n))
	}
	benchmarkPos = pos
}

func benchmarkSizeAdd(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; b.Loop(); n++ {
		size = size.Add(fyne.NewPos(float32(n), float32(n)))
	}
	benchmarkSize = size
}

func benchmarkSizeAddWidthHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; b.Loop(); n++ {
		size = size.AddWidthHeight(float32(n), float32(n))
	}
	benchmarkSize = size
}

func benchmarkSizeSubtract(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; b.Loop(); n++ {
		size = size.Subtract(fyne.NewSize(float32(n), float32(n)))
	}
	benchmarkSize = size
}

func benchmarkSizeSubtractWidthHeight(b *testing.B) {
	size := fyne.NewSize(10, 10)
	for n := 0; b.Loop(); n++ {
		size = size.SubtractWidthHeight(float32(n), float32(n))
	}
	benchmarkSize = size
}
