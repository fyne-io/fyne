package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLifecycle(t *testing.T) {
	life := &Lifecycle{}

	// Not setting anything should not panic.
	assert.NotPanics(t, life.TriggerEnteredForeground)
	assert.NotPanics(t, life.TriggerExitedForeground)
	assert.NotPanics(t, life.TriggerStarted)
	assert.NotPanics(t, life.TriggerStopped)

	var entered, exited, start, stop, hookedStop bool
	life.SetOnEnteredForeground(func() { entered = true })
	life.TriggerEnteredForeground()
	assert.True(t, entered)

	life.SetOnExitedForeground(func() { exited = true })
	life.TriggerExitedForeground()
	assert.True(t, exited)

	life.SetOnStarted(func() { start = true })
	life.TriggerStarted()
	assert.True(t, start)

	life.SetOnStopped(func() { stop = true })
	life.TriggerStopped()
	assert.True(t, stop)

	stop = false
	life.SetOnStoppedHookExecuted(func() { hookedStop = true })
	life.TriggerStopped()
	assert.True(t, stop && hookedStop)

	// Setting back to nil should not panic.
	life.SetOnEnteredForeground(nil)
	life.SetOnExitedForeground(nil)
	life.SetOnStarted(nil)
	life.SetOnStopped(nil)
	life.SetOnStoppedHookExecuted(nil)

	assert.NotPanics(t, life.TriggerEnteredForeground)
	assert.NotPanics(t, life.TriggerExitedForeground)
	assert.NotPanics(t, life.TriggerStarted)
	assert.NotPanics(t, life.TriggerStopped)
}
