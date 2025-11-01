package widget

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

var globalProgressInfRenderer fyne.WidgetRenderer

func BenchmarkProgressbarInf(b *testing.B) {
	var renderer fyne.WidgetRenderer
	widget := &ProgressBarInfinite{}
	b.ReportAllocs()

	for b.Loop() {
		renderer = widget.CreateRenderer()
	}

	// Avoid having the value optimized out by the compiler.
	globalProgressInfRenderer = renderer
}

func TestProgressBarInfinite_Creation(t *testing.T) {
	bar := NewProgressBarInfinite()

	// ticker should be disabled until the widget is shown
	assert.False(t, bar.Running())

	win := test.NewWindow(bar)
	defer win.Close()
	win.Show()

	// ticker should start automatically once the renderer is created
	assert.True(t, bar.Running())
}

func TestProgressBarInfinite_Destroy(t *testing.T) {
	bar := NewProgressBarInfinite()
	win := test.NewWindow(bar)
	defer win.Close()
	win.Show()

	assert.True(t, cache.IsRendered(bar))
	assert.True(t, bar.Running())

	// check that it stopped
	cache.DestroyRenderer(bar)
	assert.False(t, bar.Running())

	// and that the cache was removed
	assert.False(t, cache.IsRendered(bar))
}

func TestProgressBarInfinite_Reshown(t *testing.T) {
	bar := NewProgressBarInfinite()
	win := test.NewWindow(bar)
	defer win.Close()
	win.Show()

	assert.True(t, bar.Running())
	bar.Hide()
	assert.False(t, bar.Running())

	// make sure it restarts when re-shown
	bar.Show()
	// Show() starts a goroutine, so pause for it to initialize
	time.Sleep(10 * time.Millisecond)
	assert.True(t, bar.Running())
	bar.Hide()
	assert.False(t, bar.Running())
}

func TestInfiniteProgressRenderer_Layout(t *testing.T) {
	bar := NewProgressBarInfinite()
	width := float32(100.0)
	bar.Resize(fyne.NewSize(width, 10))

	render := test.TempWidgetRenderer(t, bar).(*infProgressRenderer)

	render.updateBar(0.0)
	// start at the smallest size
	assert.InEpsilon(t, width*minProgressBarInfiniteWidthRatio, render.bar.Size().Width, 0.0001)

	// make sure the inner progress bar grows in size
	// call updateBar() enough times to grow the inner bar
	maxWidth := width * maxProgressBarInfiniteWidthRatio
	render.updateBar(0.5)
	assert.Equal(t, maxWidth, render.bar.Size().Width)

	render.updateBar(1.0)
	// ends at the smallest size again
	assert.Equal(t, width*minProgressBarInfiniteWidthRatio, render.bar.Size().Width)

	bar.Hide()
}
