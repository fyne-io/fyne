//go:build !no_glfw && !mobile

package glfw

import "testing"

// BenchmarkRunOnMain measures the cost of calling a function
// on the main thread.
func BenchmarkRunOnMain(b *testing.B) {
	f := func() {}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runOnMain(f)
	}
}

// BenchmarkRunOnDraw measures the cost of calling a function
// on the draw thread.
func BenchmarkRunOnDraw(b *testing.B) {
	f := func() {}
	w := createWindow("Test")
	w.create()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.RunWithContext(f)
	}
}
