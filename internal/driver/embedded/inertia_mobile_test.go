//go:build mobile

package embedded

import (
	"image"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	embeddeddriver "fyne.io/fyne/v2/driver/embedded"
	"fyne.io/fyne/v2/widget"
)

func Test_noosDriver_ScrollInertiaStopsOnNewTouch(t *testing.T) {
	n := NewNoOSDriver(func(image.Image) {}, func(f func()) { f() },
		make(chan embeddeddriver.Event), func() fyne.Size { return fyne.NewSize(40, 60) }).(*noosDriver)
	w := n.CreateWindow("main").(*noosWindow)
	scroll := container.NewVScroll(widget.NewLabel(strings.Repeat("Hi\n", 30)))
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(40, 60))

	// Swipe upwards so the content scrolls down from offset zero.
	n.handleTouchDown(&embeddeddriver.TouchDownEvent{Position: fyne.NewPos(10, 50)}, w)
	n.handleTouchMove(&embeddeddriver.TouchMoveEvent{Position: fyne.NewPos(10, 30)}, w)
	n.handleTouchUp(&embeddeddriver.TouchUpEvent{Position: fyne.NewPos(10, 10)}, w)
	assert.NotNil(t, n.inertia)

	n.stepInertia()
	moved := scroll.Offset.Y
	assert.Greater(t, moved, float32(0))

	// The next finger-down must stop the scroll before another frame is drawn.
	n.handleTouchDown(&embeddeddriver.TouchDownEvent{Position: fyne.NewPos(10, 10)}, w)
	assert.Nil(t, n.inertia)
	n.stepInertia()
	assert.Equal(t, moved, scroll.Offset.Y)
}
