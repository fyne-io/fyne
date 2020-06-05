package widget

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/cache"
	"fyne.io/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestProgressBarInfinite_Creation(t *testing.T) {
	bar := NewProgressBarInfinite()
	// ticker should start automatically
	assert.True(t, bar.Running())
}

func TestProgressBarInfinite_Destroy(t *testing.T) {
	bar := NewProgressBarInfinite()
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
	width := 100.0
	bar.Resize(fyne.NewSize(int(width), 10))

	render := test.WidgetRenderer(bar).(*infProgressRenderer)

	// width of bar is one step size because updateBar() will have run once
	assert.Equal(t, int(width*progressBarInfiniteStepSizeRatio), render.bar.Size().Width)

	// make sure the inner progress bar grows in size
	// call updateBar() enough times to grow the inner bar
	maxWidth := int(width * maxProgressBarInfiniteWidthRatio)
	for i := 0; i < maxWidth; i++ {
		render.updateBar()
	}

	// width of bar is 1/5 of total width of progress bar
	assert.Equal(t, maxWidth, render.bar.Size().Width)
	bar.Hide()
}
