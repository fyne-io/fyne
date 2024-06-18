package widget

import (
	"testing"

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

	render.anim.Tick(0.25)
	test.AssertImageMatches(t, "activity/animate_0.25.png", w.Canvas().Capture())

	render.anim.Tick(0.5)
	test.AssertImageMatches(t, "activity/animate_0.5.png", w.Canvas().Capture())

	// check reset to loop
	render.anim.Tick(1.0)
	test.AssertImageMatches(t, "activity/animate_0.0.png", w.Canvas().Capture())
}
