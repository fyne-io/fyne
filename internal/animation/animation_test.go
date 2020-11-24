package animation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
)

func TestGLDriver_StartAnimation(t *testing.T) {
	run := &Runner{}
	ticks := 0
	a := &fyne.Animation{
		Duration: time.Millisecond * 100,
		Tick: func(done float32) {
			ticks++
		}}

	run.Start(a)
	time.Sleep(time.Millisecond * 20)
	assert.Greater(t, ticks, 0)
}

func TestGLDriver_StopAnimation(t *testing.T) {
	run := &Runner{}
	ticks := 0
	a := &fyne.Animation{
		Duration: time.Second * 10,
		Tick: func(done float32) {
			ticks++
		}}

	run.Start(a)
	time.Sleep(time.Millisecond * 20)
	run.Stop(a)
	assert.Greater(t, ticks, 0)
	assert.Zero(t, len(run.animations))
}
