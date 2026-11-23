package gl

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"
)

var benchmarkBytes []byte

func TestToLEByteOrder(t *testing.T) {
	values := []float32{0, 1, -1, 0.5, math.Float32frombits(0x12345678)}
	got := toLEByteOrder(values...)
	if len(got) != 4*len(values) {
		t.Fatalf("len = %d, want %d", len(got), 4*len(values))
	}
	for i, v := range values {
		want := math.Float32bits(v)
		if have := binary.LittleEndian.Uint32(got[4*i:]); have != want {
			t.Errorf("value %d: got %#08x, want %#08x", i, have, want)
		}
	}

	// The gomobile GL worker reads queued payloads after the call returns, so
	// separate uploads must not share a backing array.
	first := toLEByteOrder(1, 2, 3)
	second := toLEByteOrder(4, 5, 6)
	first[0] = 0xff
	if have, want := binary.LittleEndian.Uint32(second), math.Float32bits(4); have != want {
		t.Errorf("uploads share a buffer: got %#08x, want %#08x", have, want)
	}

	if got := toLEByteOrder(); len(got) != 0 {
		t.Errorf("empty input: got %v, want an empty slice", got)
	}
}

// BenchmarkToLEByteOrder compares the memmove-based conversion with the former
// byte-by-byte loop. A glyph quad is 48 floats, so the largest case is roughly
// a screenful of text.
func BenchmarkToLEByteOrder(b *testing.B) {
	for _, size := range []int{48, 480, 4800, 48000} {
		b.Run(fmt.Sprintf("N=%d", size), func(b *testing.B) {
			values := make([]float32, size)
			for i := range values {
				values[i] = float32(i) * 0.5
			}

			b.Run("memmove", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkBytes = toLEByteOrder(values...)
				}
			})
			b.Run("scalar", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkBytes = scalarLEByteOrder(values...)
				}
			})
		})
	}
}

// scalarLEByteOrder is the former implementation, kept as a benchmark baseline.
func scalarLEByteOrder(values ...float32) []byte {
	b := make([]byte, 4*len(values))
	for i, v := range values {
		u := math.Float32bits(v)
		b[4*i+0] = byte(u >> 0)  //gosec:disable G115 -- extracting one byte of a float bit pattern
		b[4*i+1] = byte(u >> 8)  //gosec:disable G115 -- extracting one byte of a float bit pattern
		b[4*i+2] = byte(u >> 16) //gosec:disable G115 -- extracting one byte of a float bit pattern
		b[4*i+3] = byte(u >> 24)
	}
	return b
}
