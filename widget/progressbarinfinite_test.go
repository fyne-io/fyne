package widget

import (
	"testing"
	"time"

	"fyne.io/fyne"
	"github.com/stretchr/testify/assert"
)

func TestProgressBarInfinite_Creation(t *testing.T) {
	bar := NewProgressBarInfinite()
	// ticker should be nil when created
	assert.Nil(t, bar.ticker)
}

func TestProgressBarInfinite_Ticker(t *testing.T) {
	bar := NewProgressBarInfinite()

	// Show() starts a goroutine, so pause for it to initialize
	time.Sleep(10 * time.Millisecond)
	assert.NotNil(t, bar.ticker)
	bar.Hide()
	assert.Nil(t, bar.ticker)

	// make sure it restarts when re-shown
	bar.Show()
	// Show() starts a goroutine, so pause for it to initialize
	time.Sleep(10 * time.Millisecond)
	assert.NotNil(t, bar.ticker)
	bar.Hide()
	assert.Nil(t, bar.ticker)
}

func TestInfiniteProgressRenderer_Layout(t *testing.T) {
	bar := NewProgressBarInfinite()
	width := 100
	bar.Resize(fyne.NewSize(width, 10))

	render := Renderer(bar).(*infProgressRenderer)

	// width of bar is one step size because updateBar() will have run once
	assert.Equal(t, width*progressBarInfiniteStepSizePercent/100, render.bar.Size().Width)

	// make sure the inner progress bar grows in size
	// call updateBar() enough times to grow the inner bar
	maxWidth := width * maxProgressBarInfiniteWidthPercent / 100
	for i := 0; i < maxWidth; i++ {
		render.updateBar()
	}

	// width of bar is 1/5 of total width of progress bar
	assert.Equal(t, maxWidth, render.bar.Size().Width)
}
