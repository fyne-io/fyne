//go:build !js
// +build !js

package glfw

import (
	"os"
	"os/signal"
	"syscall"
)

func catchTerm(d *gLDriver) {
	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	for range terminateSignals {
		d.Quit()
		break
	}
}
