//go:build mobile && (!windows || !ci)

package mobile

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/driver/mobile/event/mouse"
	"fyne.io/fyne/v2/internal/scale"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

// scrollAt scrolls vertically at the given position, which is given in canvas
// coordinates while the mouse reports screen pixels.
func scrollAt(drv *driver, w *window, x, y, dy float32) {
	drv.scrollCanvas(w, mouse.ScrollEvent{
		X:       float32(scale.ToScreenCoordinate(w.canvas, x)),
		Y:       float32(scale.ToScreenCoordinate(w.canvas, y)),
		ScrollY: dy,
	})
}

func Test_driver_ScrollCanvas(t *testing.T) {
	scroll := container.NewScroll(widget.NewLabel(strings.Repeat("Hi\n", 20)))
	drv := &driver{}
	w := drv.CreateWindow("main").(*window)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(40, 100))
	assert.Equal(t, float32(0), scroll.Offset.Y)

	// a wheel notch down moves the view on by the scroll speed
	scrollAt(drv, w, 10, 10, -1)
	assert.Equal(t, common.ScrollSpeed, scroll.Offset.Y)

	// and a wheel notch up moves it back
	scrollAt(drv, w, 10, 10, 1)
	assert.Equal(t, float32(0), scroll.Offset.Y)
}

func Test_driver_ScrollCanvas_notScrollable(t *testing.T) {
	drv := &driver{}
	w := drv.CreateWindow("main").(*window)
	w.SetContent(widget.NewLabel(strings.Repeat("Hi\n", 20)))
	w.Resize(fyne.NewSize(40, 100))

	// scrolling content that cannot be scrolled should do nothing at all
	assert.NotPanics(t, func() {
		scrollAt(drv, w, 10, 10, -1)
	})
}
