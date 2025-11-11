//go:build !no_glfw && !mobile

package glfw

import "testing"

// BenchmarkRunOnMain measures the cost of calling a function
// on the main thread.
func BenchmarkRunOnMain(b *testing.B) {
	f := func() {}

	b.ReportAllocs()

	for b.Loop() {
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

	for b.Loop() {
		w.RunWithContext(f)
	}
}
