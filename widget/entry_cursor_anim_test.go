package widget

import (
	"image/color"
	"sync"
	"testing"
	"time"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestEntryCursorAnim(t *testing.T) {
	cursor := canvas.NewRectangle(color.Black)
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor}
	a.inverted = false
	a.state = cursorStateStopped
	a.ticker = &cursorTicker{}
	beforeNextTick := make(chan int)
	afterTick := make(chan int)
	a.ticker.mockWait = func() (reset bool) {
		beforeNextTick <- 1
		<-afterTick
		return
	}
	flushReset := func() { <-a.ticker.resetChan }

	// start animation
	a.start()
	// unblock waitTick() to start testing
	<-beforeNextTick
	assert.False(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, a.ticker.started())

	// pass 1 cursorTicker tick
	afterTick <- 1   // unblock to execute code below waitTick()
	<-beforeNextTick // wait until the code below waitTick was effectively executed

	// entry cursor animation must have the same values as before
	assert.False(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, a.ticker.started())

	// now call a TemporaryStop()
	go flushReset()
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateInterrupted, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, a.ticker.started())
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	// after 1 cursorTicker tick, the animation should start again
	afterTick <- 1   // unblock to execute code below waitTick()
	<-beforeNextTick // wait until the code below waitTick was effectively executed
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.True(t, a.ticker.started())
	assert.NotNil(t, a.anim)

	// after 1 cursorTicker tick, the animation should continue (not interrupted)
	afterTick <- 1   // unblock to execute code below waitTick()
	<-beforeNextTick // wait until the code below waitTick was effectively executed
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.True(t, a.ticker.started())
	assert.NotNil(t, a.anim)

	// temporary stop again
	go flushReset()
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateInterrupted, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, a.ticker.started())
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	// stop the animation
	a.stop()
	// unblock to execute code below waitTick() and so it should break the for-loop, and stop
	// ticker and animation
	afterTick <- 1
	time.Sleep(2 * time.Millisecond)
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateStopped, a.state)
	assert.Nil(t, a.anim)
	assert.False(t, a.ticker.started())

	// calling a.TemporaryStop() on stopped animation, does not do anything
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateStopped, a.state)
	assert.Nil(t, a.anim)
	assert.False(t, a.ticker.started())

	assert.NotPanics(t, func() { a.temporaryStop() })
	assert.NotPanics(t, func() { a.start() })
	assert.NotPanics(t, func() { a.stop() })
}
