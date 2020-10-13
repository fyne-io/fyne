package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type builder struct {
	os, srcdir string
	release    bool
}

func (b *builder) build() error {
	goos := b.os
	if goos == "" {
		goos = targetOS()
	}

	var cmd *exec.Cmd
	if goos == "windows" {
		if b.release {
			cmd = exec.Command("go", "build", "-ldflags", "-s -w -H=windowsgui", ".")
		} else {
			cmd = exec.Command("go", "build", "-ldflags", "-H=windowsgui", ".")
		}
	} else {
		if b.release {
			cmd = exec.Command("go", "build", "-ldflags", "-s -w", ".")
		} else {
			cmd = exec.Command("go", "build", ".")
		}
	}
	cmd.Dir = b.srcdir
	env := os.Environ()
	env = append(env, "CGO_ENABLED=1") // in case someone is trying to cross-compile...

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

func targetOS() string {
	osEnv, ok := os.LookupEnv("GOOS")
	if ok {
		return osEnv
	}

	return runtime.GOOS
}
