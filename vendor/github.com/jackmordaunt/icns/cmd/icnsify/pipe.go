package main

import (
	"fmt"
	"os"
)

func init() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(fmt.Sprintf("getting info on stdin file descriptor: %v", err))
	}
	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		return
	}
	if info.Size() > 0 {
		piping = true
		input = os.Stdin
		output = os.Stdout
	}
}
