package widget

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestSelectRenderer_TapAnimation(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	test.ApplyTheme(t, test.NewTheme())
	sel := NewSelect([]string{"one"}, func(s string) {})
	w := test.NewWindow(sel)
	defer w.Close()
	w.Resize(sel.MinSize().Add(fyne.NewSize(10, 10)))
	sel.Resize(sel.MinSize())
	sel.Refresh()

	render1 := test.WidgetRenderer(sel).(*selectRenderer)
	test.Tap(sel)
	sel.popUp.Hide()
	sel.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "select/tap_animation.png", w.Canvas().Capture())

	cache.DestroyRenderer(sel)
	sel.Refresh()

	render2 := test.WidgetRenderer(sel).(*selectRenderer)

	assert.NotEqual(t, render1, render2)

	test.Tap(sel)
	sel.popUp.Hide()
	sel.tapAnim.Tick(0.5)
	test.AssertImageMatches(t, "select/tap_animation.png", w.Canvas().Capture())
}
