// +build !ci
// +build !mobile

package glfw

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne"
)

func TestGLDriver_StartAnimation(t *testing.T) {
	ticks := 0
	a := &fyne.Animation{
		Duration: time.Millisecond * 100,
		Tick: func(done float32) {
			ticks++
		}}

	a.Start()
	assert.Greater(t, ticks, 0)
}

func TestGLDriver_StopAnimation(t *testing.T) {
	ticks := 0
	a := &fyne.Animation{
		Duration: time.Second * 10,
		Tick: func(done float32) {
			ticks++
		}}

	a.Start()
	a.Stop()
	assert.Greater(t, ticks, 0)
	assert.Zero(t, len(d.(*gLDriver).animations))
}
