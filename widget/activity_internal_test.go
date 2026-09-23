package widget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
)

func TestActivity_Animation(t *testing.T) {
	test.NewTempApp(t)
	test.ApplyTheme(t, test.NewTheme())

	a := NewActivity()
	w := test.NewWindow(a)
	defer w.Close()
	w.SetPadded(false)
	w.Resize(a.MinSize())

	render := test.TempWidgetRenderer(t, a).(*activityRenderer)
	render.anim.Tick(0)
	test.AssertImageMatches(t, "activity/animate_0.0.png", w.Canvas().Capture())
	render.anim.Tick(0.1)
	test.AssertImageMatches(t, "activity/animate_0.1.png", w.Canvas().Capture())
	render.anim.Tick(0.2)
	test.AssertImageMatches(t, "activity/animate_0.2.png", w.Canvas().Capture())
	render.anim.Tick(0.3)
	test.AssertImageMatches(t, "activity/animate_0.3.png", w.Canvas().Capture())
	render.anim.Tick(0.4)
	test.AssertImageMatches(t, "activity/animate_0.4.png", w.Canvas().Capture())
	render.anim.Tick(0.5)
	test.AssertImageMatches(t, "activity/animate_0.5.png", w.Canvas().Capture())
	render.anim.Tick(0.6)
	test.AssertImageMatches(t, "activity/animate_0.6.png", w.Canvas().Capture())
	render.anim.Tick(0.7)
	test.AssertImageMatches(t, "activity/animate_0.7.png", w.Canvas().Capture())
	render.anim.Tick(0.8)
	test.AssertImageMatches(t, "activity/animate_0.8.png", w.Canvas().Capture())
	render.anim.Tick(0.9)
	test.AssertImageMatches(t, "activity/animate_0.9.png", w.Canvas().Capture())
	render.anim.Tick(1.0)
	// reset to loop
	test.AssertImageMatches(t, "activity/animate_0.0.png", w.Canvas().Capture())
}

func TestActivity_StaticEllipsisLayout(t *testing.T) {
	test.NewTempApp(t)
	test.ApplyTheme(t, test.NewTheme())

	a := NewActivity()
	w := test.NewWindow(a)
	defer w.Close()
	w.SetPadded(false)

	size := fyne.NewSize(40, 12)
	w.Resize(size)

	render := test.TempWidgetRenderer(t, a).(*activityRenderer)
	render.bound = size
	render.drawStaticEllipsis()

	dots := make([]*canvas.Circle, len(render.dots))
	for i, obj := range render.dots {
		dots[i] = obj.(*canvas.Circle)
	}

	// All three dots are the same non-zero size.
	d0 := dots[0].Size()
	assert.Greater(t, d0.Width, float32(0), "dot should have non-zero width")
	for _, d := range dots[1:] {
		assert.Equal(t, d0, d.Size(), "all dots should be the same size")
	}

	// All three sit on the same y-position — horizontal row.
	y := dots[0].Position().Y
	for _, d := range dots[1:] {
		assert.Equal(t, y, d.Position().Y, "dots must share a baseline")
	}

	// X-positions advance left to right with equal spacing.
	gap1 := dots[1].Position().X - dots[0].Position().X
	gap2 := dots[2].Position().X - dots[1].Position().X
	assert.Greater(t, gap1, float32(0), "second dot must be right of the first")
	assert.InDelta(t, gap1, gap2, 0.001, "dots should be evenly spaced")

	// The row should be horizontally centered in the bound.
	rowLeft := dots[0].Position().X
	rowRight := dots[2].Position().X + dots[2].Size().Width
	leftMargin := rowLeft
	rightMargin := size.Width - rowRight
	assert.InDelta(t, leftMargin, rightMargin, 0.001, "row should be centered horizontally")
}
