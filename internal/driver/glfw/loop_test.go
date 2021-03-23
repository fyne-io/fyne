// +build !ci
// +build !mobile

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
	w := createWindow("Test").(*window)
	w.create()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runOnDraw(w, f)
	}
}
