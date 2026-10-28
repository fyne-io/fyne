package widget_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	internalWidget "fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// countingWidget is a minimal widget that records how often the scroll container
// re-measures, re-lays-out or re-paints it. Scrolling is a paint-time translation of
// the content, so none of these should happen while only the offset changes.
type countingWidget struct {
	internalWidget.Base
	size fyne.Size
}

var (
	countingMinSizeCalls uint64
	countingResizeCalls  uint64
	countingRefreshCalls uint64
)

func newCountingWidget(w, h float32) *countingWidget {
	c := &countingWidget{size: fyne.NewSize(w, h)}
	c.ExtendBaseWidget(c)
	return c
}

// countingRenderer is an object-free renderer; countingWidget measures itself so
// nothing needs to actually be painted.
type countingRenderer struct {
	internalWidget.BaseRenderer
}

func (*countingRenderer) Layout(fyne.Size) {}

func (*countingRenderer) MinSize() fyne.Size { return fyne.Size{} }

func (*countingRenderer) Refresh() {}

func (*countingWidget) CreateRenderer() fyne.WidgetRenderer {
	return &countingRenderer{BaseRenderer: internalWidget.NewBaseRenderer(nil)}
}

func (c *countingWidget) MinSize() fyne.Size {
	atomic.AddUint64(&countingMinSizeCalls, 1)
	return c.size
}

func (c *countingWidget) Resize(size fyne.Size) {
	if size != c.size {
		atomic.AddUint64(&countingResizeCalls, 1)
	}
	c.Base.Resize(size)
}

func (c *countingWidget) Refresh() {
	atomic.AddUint64(&countingRefreshCalls, 1)
	c.Base.Refresh()
}

func resetCounting() {
	atomic.StoreUint64(&countingMinSizeCalls, 0)
	atomic.StoreUint64(&countingResizeCalls, 0)
	atomic.StoreUint64(&countingRefreshCalls, 0)
}

func countingTotals() (minSize, resize, refresh uint64) {
	return atomic.LoadUint64(&countingMinSizeCalls),
		atomic.LoadUint64(&countingResizeCalls),
		atomic.LoadUint64(&countingRefreshCalls)
}

func newCountingContent(rows int) *fyne.Container {
	objects := make([]fyne.CanvasObject, rows)
	for i := range objects {
		objects[i] = newCountingWidget(80, 20)
	}
	return container.NewVBox(objects...)
}

func newLabelContent(rows int) *fyne.Container {
	objects := make([]fyne.CanvasObject, rows)
	for i := range objects {
		objects[i] = widget.NewLabel(fmt.Sprintf("Scrollable row number %d with some text", i))
	}
	return container.NewVBox(objects...)
}

func setupScroll(t testing.TB, rows int) (*internalWidget.Scroll, fyne.Canvas) {
	scroll := internalWidget.NewVScroll(newLabelContent(rows))
	w := test.NewTempWindow(t, scroll)
	w.Resize(fyne.NewSize(400, 400))
	return scroll, w.Canvas()
}

func scrollOnce(scroll *internalWidget.Scroll) {
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -10)})
}

// BenchmarkScrollScrolled measures the cost of handling one scroll input event on a
// scroll container holding a large content tree.
func BenchmarkScrollScrolled(b *testing.B) {
	for _, rows := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("rows=%d", rows), func(b *testing.B) {
			scroll, _ := setupScroll(b, rows)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if scroll.Offset.Y >= float32(rows*20)-400 {
					// reset directly: ScrollToTop goes through the refresh path and
					// would measure that instead of the scroll handling
					scroll.Offset = fyne.Position{}
				}
				scrollOnce(scroll)
			}
		})
	}
}

// BenchmarkScrollPaint measures a full paint cycle (walk, ensure minimum size, draw)
// after a scroll event, which is what a real driver does once per frame.
func BenchmarkScrollPaint(b *testing.B) {
	for _, rows := range []int{100, 1000} {
		b.Run(fmt.Sprintf("rows=%d", rows), func(b *testing.B) {
			scroll, canvas := setupScroll(b, rows)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if scroll.Offset.Y >= float32(rows*20)-400 {
					scroll.Offset = fyne.Position{}
				}
				scrollOnce(scroll)
				canvas.Capture()
			}
		})
	}
}

// BenchmarkContentMinSizeWalk measures one full minimum size sweep over the content
// tree. Canvas.EnsureMinSize does this for every object of every tree on every single
// paint, so this is the per frame cost that dominates long scroll views.
func BenchmarkContentMinSizeWalk(b *testing.B) {
	for _, rows := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("rows=%d", rows), func(b *testing.B) {
			content := newLabelContent(rows)
			w := test.NewTempWindow(b, content)
			w.Resize(fyne.NewSize(400, 400))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				content.MinSize()
			}
		})
	}
}

// TestScrollScrolled_DoesNotMeasureContent asserts that handling a scroll event does
// not walk the content tree. Scrolling only changes the paint position of the
// content and the length of the scroll bar thumbs, so re-measuring the content makes
// scrolling slower the more content there is.
func TestScrollScrolled_DoesNotMeasureContent(t *testing.T) {
	const rows = 10
	scroll := internalWidget.NewVScroll(newCountingContent(rows))
	w := test.NewTempWindow(t, scroll)
	w.Resize(fyne.NewSize(200, 100))

	// settle the layout before counting
	scroll.ScrollToTop()

	resetCounting()
	scrollOnce(scroll)

	minSize, resize, refresh := countingTotals()
	assert.Zero(t, minSize, "scrolling must not measure the content tree")
	assert.Zero(t, resize, "scrolling must not resize the content")
	assert.Zero(t, refresh, "scrolling must not refresh the content")
}

// TestScrollScrolled_CostIsIndependentOfContentSize is the regression guard for the
// above: the work done per event must not grow with the number of content children.
func TestScrollScrolled_CostIsIndependentOfContentSize(t *testing.T) {
	measure := func(rows int) uint64 {
		scroll := internalWidget.NewVScroll(newCountingContent(rows))
		w := test.NewTempWindow(t, scroll)
		w.Resize(fyne.NewSize(200, 100))
		scroll.ScrollToTop()

		resetCounting()
		scrollOnce(scroll)

		minSize, _, _ := countingTotals()
		return minSize
	}

	small := measure(4)
	large := measure(400)

	assert.Equal(t, small, large,
		"per scroll event measurement must not scale with the content size (4 rows: %d, 400 rows: %d)",
		small, large)
}

// TestScrollScrolled_StillRepositions guards that the cheap path did not drop the work
// that scrolling actually has to do.
func TestScrollScrolled_StillRepositions(t *testing.T) {
	content := newCountingContent(100)
	scroll := internalWidget.NewVScroll(content)
	w := test.NewTempWindow(t, scroll)
	w.Resize(fyne.NewSize(200, 100))

	scrollOnce(scroll)
	assert.Equal(t, fyne.NewPos(0, -10), content.Position(), "content must move with the offset")

	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -20)})
	assert.Equal(t, fyne.NewPos(0, -30), content.Position())

	// scrolling past the end stops at the content bound
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -100000)})
	assert.Equal(t, content.Size().Height-scroll.Size().Height, scroll.Offset.Y)
}
