// +build !ci !darwin

package animation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func TestGLDriver_StartAnimation(t *testing.T) {
	done := make(chan float32)
	run := &Runner{}
	a := &fyne.Animation{
		Duration: time.Millisecond * 100,
		Tick: func(d float32) {
			done <- d
		}}

	run.Start(a)
	select {
	case d := <-done:
		assert.Greater(t, d, float32(0))
	case <-time.After(100 * time.Millisecond):
		t.Error("animation was not ticked")
	}
}

func TestGLDriver_StopAnimation(t *testing.T) {
	done := make(chan float32)
	run := &Runner{}
	a := &fyne.Animation{
		Duration: time.Second * 10,
		Tick: func(d float32) {
			done <- d
		}}

	run.Start(a)
	select {
	case d := <-done:
		assert.Greater(t, d, float32(0))
	case <-time.After(time.Second):
		t.Error("animation was not ticked")
	}
	run.Stop(a)
	assert.Zero(t, len(run.animations))
}
