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

type mockTicker struct {
	started bool
	ticks   int
	cycle   chan int
	tickCh  chan int
}

func (m *mockTicker) sendTick() {
	m.tickCh <- 1
	<-m.cycle
	m.ticks++
}
func (m *mockTicker) WaitTick() {
	m.cycle <- 1
	<-m.tickCh
}
func (m *mockTicker) Stop() {
	m.started = false
}
func (m *mockTicker) Reset() { m.ticks = 0 }
func (m *mockTicker) Start(d time.Duration) {
	m.started = true
	m.tickCh = make(chan int)
	m.cycle = make(chan int)
	go func() { <-m.cycle }()
	runtime.Gosched()
}
func (m *mockTicker) Started() bool { return m.started }

func TestEntryCursorAnim(t *testing.T) {
	cursor := canvas.NewRectangle(color.Black)
	a := &entryCursorAnimation{mu: &sync.RWMutex{}, cursor: cursor}
	a.inverted = false
	a.state = cursorStateStopped
	mticker := &mockTicker{}
	a.ticker = mticker

	// start animation
	a.start()
	assert.False(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, mticker.Started())

	// pass some time
	for i := 0; i < 10; i++ {
		assert.Equal(t, i, mticker.ticks)
		mticker.sendTick()
	}

	// entry cursor animation must have the same values as before
	assert.False(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, mticker.Started())

	// now call a TemporaryStop()
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateInterrupted, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, mticker.Started())
	assert.Zero(t, mticker.ticks)
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	// when ticks, the animation should start again
	mticker.sendTick()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.True(t, mticker.Started())
	assert.Equal(t, 1, mticker.ticks)
	assert.NotNil(t, a.anim)

	// make some steps
	mticker.sendTick()
	mticker.sendTick()
	// animation should continue (not interrupted)
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateRunning, a.state)
	assert.True(t, mticker.Started())
	assert.Equal(t, 3, mticker.ticks)
	assert.NotNil(t, a.anim)

	// temporary stop again
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateInterrupted, a.state)
	assert.NotNil(t, a.anim)
	assert.True(t, mticker.Started())
	assert.Zero(t, mticker.ticks)
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	// temporary stop again (counter should be resetted)
	mticker.ticks = 100 // just to ensure it is resetted below
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateInterrupted, a.state)
	assert.NotNil(t, a.anim)
	assert.Zero(t, mticker.ticks)
	assert.Equal(t, theme.PrimaryColor(), a.cursor.FillColor)

	// stop the animation
	a.stop()
	mticker.tickCh <- 1
	time.Sleep(5 * time.Millisecond)
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateStopped, a.state)
	assert.Nil(t, a.anim)

	// calling a.TemporaryStop() on stopped animation, does not do anything
	a.temporaryStop()
	assert.True(t, a.inverted)
	assert.Equal(t, cursorStateStopped, a.state)
	assert.Nil(t, a.anim)
	assert.False(t, mticker.Started())

	assert.NotPanics(t, func() { a.temporaryStop() })
	assert.NotPanics(t, func() { a.start() })
	assert.NotPanics(t, func() { a.stop() })
}
