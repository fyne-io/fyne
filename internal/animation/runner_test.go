package animation

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
)

func BenchmarkRunnerAllocs(b *testing.B) {
	r := Runner{}
	var fl float32
	// setup some animations
	for range 10 {
		r.pendingAnimations = append(r.pendingAnimations, newAnim(
			fyne.NewAnimation(1000*time.Second, func(f float32) {
				fl = f
			})))
	}
	for b.Loop() {
		r.runOneFrame()
	}

	b.ReportAllocs()
	fl = fl + 1 // dummy use of variable
}
