//go:build mobile && (!windows || !ci)

package mobile

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Test_driver_ScrollInertiaStopsOnNewTouch(t *testing.T) {
	d := &driver{}
	w := d.CreateWindow("main").(*window)
	scroll := container.NewVScroll(widget.NewLabel(strings.Repeat("Hi\n", 30)))
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(40, 60))

	// Swipe upwards so the content scrolls down from offset zero.
	d.tapDownCanvas(w, 10, 50, 0)
	d.tapMoveCanvas(w, 10, 30, 0)
	d.tapUpCanvas(w, 10, 10, 0)
	assert.NotNil(t, d.inertia)

	d.stepInertia()
	moved := scroll.Offset.Y
	assert.Greater(t, moved, float32(0))

	// The next finger-down must stop the scroll before another frame is drawn.
	d.tapDownCanvas(w, 10, 10, 0)
	assert.Nil(t, d.inertia)
	d.stepInertia()
	assert.Equal(t, moved, scroll.Offset.Y)
}
