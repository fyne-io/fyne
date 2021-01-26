package widget

import (
	"image/color"
	"runtime"
	"sync"
	"testing"
	"time"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestEntryCursorAnim(t *testing.T) {
	cursor := canvas.NewRectangle(color.Black)
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor, counter: &safeCounter{}}
	a.inverted = false
	a.state = cursorStateStopped
	sleeper := make(chan int)
	cycle := make(chan int)
	a.sleepFn = func(d time.Duration) {
		cycle <- 1
		<-sleeper
	}

	// start animation
	a.start()
	<-cycle
	assert.False(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	assert.Zero(t, a.counter.Value())

	// pass some time
	for i := 0; i < 10; i++ {
		sleeper <- 1
		<-cycle
	}

	// entry cursor animation must have the same values as before
	assert.False(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	assert.Zero(t, a.counter.Value())

	// now call a TemporaryStop()
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateTemporarilyStopped, a.state)
	assert.NotNil(t, a.anim)
	assert.Zero(t, a.counter.Value())
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	// make some steps in time less than cursorStopTimex10ms
	// cursor.FillColor should be PrimaryColor always
	for i := 0; i < (cursorStopTimex10ms - 1); i++ {
		sleeper <- 1
		<-cycle
		assert.True(t, a.inverted)
		assert.Equal(t, i+1, a.counter.Value())
		assert.Equal(t, cursorStateTemporarilyStopped, a.state)
		assert.NotNil(t, a.anim)
		assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)
	}

	// advance one step more (equals to cursorPauseTimex10ms) and the animation
	// should start again
	sleeper <- 1
	<-cycle
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)

	// make some steps
	sleeper <- 1
	<-cycle
	sleeper <- 1
	<-cycle
	// animation should continue (not temp. stopped)
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	// counter value should be cursorStopTimex10ms (it shouldn't increase)
	assert.Equal(t, cursorStopTimex10ms, a.counter.Value())

	// temporary stop again
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateTemporarilyStopped, a.state)
	assert.NotNil(t, a.anim)
	assert.Zero(t, a.counter.Value())
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	for i := 0; i < 10; i++ {
		sleeper <- 1
		<-cycle
		assert.True(t, a.inverted)
		assert.Equal(t, i+1, a.counter.Value())
		assert.Equal(t, cursorStateTemporarilyStopped, a.state)
		assert.NotNil(t, a.anim)
		assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)
	}

	// temporary stop again (counter should be resetted)
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateTemporarilyStopped, a.state)
	assert.NotNil(t, a.anim)
	assert.Zero(t, a.counter.Value())
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	for i := 0; i < 10; i++ {
		sleeper <- 1
		<-cycle
		assert.True(t, a.inverted)
		assert.Equal(t, i+1, a.counter.Value())
		assert.Equal(t, cursorStateTemporarilyStopped, a.state)
		assert.NotNil(t, a.anim)
		assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)
	}

	// stop the animation
	a.stop()
	sleeper <- 1
	time.Sleep(1 * time.Millisecond)
	runtime.Gosched()
	time.Sleep(5 * time.Millisecond)
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateStopped, a.state)
	assert.Nil(t, a.anim)
	assert.Equal(t, 10, a.counter.Value())

	// calling a.TemporaryStop() on stopped animation, does not do anything (just reset the counter)
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateStopped, a.state)
	assert.Nil(t, a.anim)
	assert.Zero(t, a.counter.Value())

	assert.NotPanics(t, func() { a.temporaryStop() })
	assert.NotPanics(t, func() { a.start() })
	assert.NotPanics(t, func() { a.stop() })
}
