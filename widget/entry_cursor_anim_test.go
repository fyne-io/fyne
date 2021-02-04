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
	cycle := make(chan int)
	sleep := make(chan int)
	a.ticker.mockWait = func() (reset bool) {
		cycle <- 1
		<-sleep
		return
	}
	flushReset := func() { <-a.ticker.rstCh }

	// start animation
	a.start()
	<-cycle
	assert.False(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, a.ticker.Started())

	// pass 1 cursorTicker cycle
	sleep <- 1
	<-cycle

	// entry cursor animation must have the same values as before
	assert.False(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, a.ticker.Started())

	// now call a TemporaryStop()
	go flushReset()
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateInterrupted, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, a.ticker.Started())
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	// after 1 cursorTicker cycle, the animation should start again
	sleep <- 1
	<-cycle
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.True(t, a.ticker.Started())
	assert.NotNil(t, a.anim)

	// after 1 cursorTicker cycle, the animation should continue (not interrupted)
	sleep <- 1
	<-cycle
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.True(t, a.ticker.Started())
	assert.NotNil(t, a.anim)

	// temporary stop again
	go flushReset()
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateInterrupted, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, a.ticker.Started())
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	// stop the animation
	a.stop()
	sleep <- 1
	time.Sleep(2 * time.Millisecond)
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateStopped, a.state)
	assert.Nil(t, a.anim)

	// calling a.TemporaryStop() on stopped animation, does not do anything
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateStopped, a.state)
	assert.Nil(t, a.anim)
	assert.False(t, a.ticker.Started())

	assert.NotPanics(t, func() { a.temporaryStop() })
	assert.NotPanics(t, func() { a.start() })
	assert.NotPanics(t, func() { a.stop() })
}
