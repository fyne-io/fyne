//go:build !js
// +build !js

package glfw

func (d *gLDriver) Run() {
	if goroutineID() != mainGoroutineID {
		panic("Run() or ShowAndRun() must be called from main goroutine")
	}
	d.runGL()
}
