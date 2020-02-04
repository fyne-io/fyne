package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type builder struct {
	os, srcdir string
}

func (b *builder) build() error {
	goos := b.os
	if goos == "" {
		goos = runtime.GOOS
	}

	var cmd *exec.Cmd
	if goos == "windows" {
		cmd = exec.Command("go", "build", "-ldflags", "-H=windowsgui", ".")
	} else {
		cmd = exec.Command("go", "build", ".")
	}
	cmd.Dir = b.srcdir
	env := os.Environ()
	if goos != "ios" && goos != "android" {
		env = append(env, "GOOS="+goos)
	}
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(out))
	}
	return err
}
