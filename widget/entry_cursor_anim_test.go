package widget

import (
	"image/color"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	col "fyne.io/fyne/v2/internal/color"
	_ "fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

func TestEntryCursorAnim(t *testing.T) {
	cursorOpaque := theme.Color(theme.ColorNamePrimary)
	r, g, b, _ := col.ToNRGBA(cursorOpaque)
	cursorDim := color.NRGBA{R: r, G: g, B: b, A: 0x16}

	alpha := func(c color.Color) uint8 {
		_, _, _, a := col.ToNRGBA(c)
		return a
	}

	cursor := canvas.NewRectangle(color.Black)
	a := newEntryCursorAnimation(cursor)
	a.start()

	assert.True(t, a.anim.AutoReverse)
	assert.Equal(t, fyne.AnimationRepeatForever, a.anim.RepeatCount)

	t.Run("animation changes from opaque to dimmed fading only a small time in between", func(t *testing.T) {
		a.anim.Tick(0.0)
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))
		a.anim.Tick(0.4)
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))
		a.anim.Tick(0.5)
		assert.InDelta(t, (alpha(cursorOpaque)-alpha(cursorDim))/2+alpha(cursorDim), alpha(a.cursor.FillColor), 1)
		a.anim.Tick(0.6)
		assert.Equal(t, alpha(cursorDim), alpha(a.cursor.FillColor))
		a.anim.Tick(1.0)
		assert.Equal(t, alpha(cursorDim), alpha(a.cursor.FillColor))
	})

	t.Run("interrupted animation is always opaque", func(t *testing.T) {
		a.interrupt()
		a.anim.Tick(0.0)
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))
		a.anim.Tick(0.5)
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))
		a.anim.Tick(1.0)
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))
	})

	t.Run("animation starts fading out again 300ms after interruption", func(t *testing.T) {
		timeNow = func() time.Time {
			return time.Now().Add(300 * time.Millisecond)
		}
		a.anim.Tick(0.0)
		a.anim.Tick(0.0) // first tick after interruption period creates a new animation which is ticked by the test driver to 1.0 directly
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))
		a.anim.Tick(0.4)
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))
		a.anim.Tick(0.5)
		assert.InDelta(t, (alpha(cursorOpaque)-alpha(cursorDim))/2+alpha(cursorDim), alpha(a.cursor.FillColor), 1)
		a.anim.Tick(0.6)
		assert.Equal(t, alpha(cursorDim), alpha(a.cursor.FillColor))
		a.anim.Tick(1.0)
		assert.Equal(t, alpha(cursorDim), alpha(a.cursor.FillColor))
	})

	t.Run("subsequent interruption works", func(t *testing.T) {
		timeNow = time.Now
		a.interrupt()
		a.anim.Tick(0.0)
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))

		timeNow = func() time.Time {
			return time.Now().Add(300 * time.Millisecond)
		}
		a.anim.Tick(0.0)
		a.anim.Tick(0.0) // first tick after interruption period creates a new animation which is ticked by the test driver to 1.0 directly
		assert.Equal(t, alpha(cursorOpaque), alpha(a.cursor.FillColor))
		a.anim.Tick(1.0)
		assert.Equal(t, alpha(cursorDim), alpha(a.cursor.FillColor))
	})

	a.stop()
	assert.Nil(t, a.anim)
}
