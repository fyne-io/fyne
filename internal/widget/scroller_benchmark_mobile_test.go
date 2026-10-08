//go:build ci || no_glfw || android || ios || mobile

package widget_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
)

// dragOnce is the touch counterpart of scrollOnce. A drag moves the content with
// the finger while the offset grows the other way, so a negative delta scrolls
// further down, the same value scrollOnce uses.
func dragOnce(scroll *internalWidget.Scroll) {
	scroll.Dragged(&fyne.DragEvent{Dragged: fyne.NewDelta(0, -10)})
}

// BenchmarkScrollDragged measures the cost of handling one touch drag event. Dragged
// is the path a real mobile or tablet build takes, and it is built separately from
// Scrolled, so it needs its own numbers.
func BenchmarkScrollDragged(b *testing.B) {
	for _, rows := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("rows=%d", rows), func(b *testing.B) {
			scroll, _ := setupScroll(b, rows)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if scroll.Offset.Y >= float32(rows*20)-400 {
					// reset directly: ScrollToTop goes through the refresh path and
					// would measure that instead of the drag handling
					scroll.Offset = fyne.Position{}
				}
				dragOnce(scroll)
			}
		})
	}
}

// TestScrollDragged_DoesNotMeasureContent asserts that handling a drag event does not
// walk the content tree, the same guarantee TestScrollScrolled gives for wheel input.
func TestScrollDragged_DoesNotMeasureContent(t *testing.T) {
	const rows = 10
	scroll := internalWidget.NewVScroll(newCountingContent(rows))
	w := test.NewTempWindow(t, scroll)
	w.Resize(fyne.NewSize(200, 100))

	// settle the layout before counting
	scroll.ScrollToTop()

	resetCounting()
	dragOnce(scroll)

	minSize, resize, refresh := countingTotals()
	assert.Zero(t, minSize, "dragging must not measure the content tree")
	assert.Zero(t, resize, "dragging must not resize the content")
	assert.Zero(t, refresh, "dragging must not refresh the content")
}

// TestScrollDragged_StillRepositions guards that the cheap path did not drop the work
// that dragging actually has to do.
func TestScrollDragged_StillRepositions(t *testing.T) {
	content := newCountingContent(100)
	scroll := internalWidget.NewVScroll(content)
	w := test.NewTempWindow(t, scroll)
	w.Resize(fyne.NewSize(200, 100))

	dragOnce(scroll)
	assert.Equal(t, fyne.NewPos(0, -10), content.Position(), "content must move with the drag")

	dragOnce(scroll)
	assert.Equal(t, fyne.NewPos(0, -20), content.Position())

	// dragging past the end stops at the content bound
	scroll.Dragged(&fyne.DragEvent{Dragged: fyne.NewDelta(0, -100000)})
	assert.Equal(t, content.Size().Height-scroll.Size().Height, scroll.Offset.Y)
}
