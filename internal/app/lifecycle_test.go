package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLifecycle(t *testing.T) {
	life := &Lifecycle{}

	// Not setting anything should not panic.
	assert.Nil(t, life.OnEnteredForeground())
	assert.Nil(t, life.OnExitedForeground())
	assert.Nil(t, life.OnStarted())
	assert.Nil(t, life.OnStopped())

	var entered, exited, start, stop, hookedStop, called bool
	life.InitEventQueue()
	go life.RunEventQueue()
	life.QueueEvent(func() { called = true })
	life.SetOnEnteredForeground(func() { entered = true })
	life.OnEnteredForeground()()
	assert.True(t, entered)

	life.SetOnExitedForeground(func() { exited = true })
	life.OnExitedForeground()()
	assert.True(t, exited)

	life.SetOnStarted(func() { start = true })
	life.OnStarted()()
	assert.True(t, start)

	life.SetOnStopped(func() { stop = true })
	life.OnStopped()()
	assert.True(t, stop)

	stop = false
	life.SetOnStoppedHookExecuted(func() { hookedStop = true })
	life.OnStopped()()
	assert.True(t, stop && hookedStop)

	// Setting back to nil should not panic.
	life.SetOnEnteredForeground(nil)
	life.SetOnExitedForeground(nil)
	life.SetOnStarted(nil)
	life.SetOnStopped(nil)
	life.SetOnStoppedHookExecuted(nil)

	assert.Nil(t, life.OnEnteredForeground())
	assert.Nil(t, life.OnExitedForeground())
	assert.Nil(t, life.OnStarted())
	assert.Nil(t, life.OnStopped())

	life.WaitForEvents()
	life.DestroyEventQueue()
	assert.True(t, called)
}
