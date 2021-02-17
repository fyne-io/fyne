package widget

import (
	"image/color"
	"runtime"
	"testing"
	"time"

	"fyne.io/fyne/v2/canvas"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestEntryCursorAnim(t *testing.T) {
	cursorOpaque := theme.PrimaryColor()
	r, g, b, _ := theme.PrimaryColor().RGBA()
	cursorDim := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0x16}

	alphaEquals := func(color1, color2 color.Color) bool {
		_, _, _, a1 := color1.RGBA()
		_, _, _, a2 := color2.RGBA()
		return a1 == a2
	}

	cursor := canvas.NewRectangle(color.Black)
	a := newEntryCursorAnimation(cursor)

	a.start()
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorDim, a.cursor.FillColor))
	a.anim.Tick(1.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))

	a.interrupt()
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))
	a.anim.Tick(0.5)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))
	a.anim.Tick(1.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))

	a.timeNow = func() time.Time {
		return time.Now().Add(cursorInterruptTime)
	}
	// animation should be restarted inverting the colors
	a.anim.Tick(0.0)
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond) // ensure go routine for restart animation is executed
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))
	a.anim.Tick(1.0)
	assert.True(t, alphaEquals(cursorDim, a.cursor.FillColor))

	a.timeNow = time.Now
	a.interrupt()
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))

	a.timeNow = func() time.Time {
		return time.Now().Add(cursorInterruptTime)
	}
	a.anim.Tick(0.0)
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond) // ensure go routine for restart animation is executed
	a.anim.Tick(0.0)
	assert.True(t, alphaEquals(cursorOpaque, a.cursor.FillColor))
	a.anim.Tick(1.0)
	assert.True(t, alphaEquals(cursorDim, a.cursor.FillColor))

	a.stop()
	assert.Nil(t, a.anim)
}
