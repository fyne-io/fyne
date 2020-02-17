package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type builder struct {
	os, srcdir, target string
}

func (b *builder) build() error {
	goos := b.os
	if goos == "" {
		goos = runtime.GOOS
	}

	var cmd *exec.Cmd
	if goos == "windows" {
		if b.target == "" {
			cmd = exec.Command("go", "build", "-ldflags", "-H=windowsgui", ".")
		} else {
			cmd = exec.Command("go", "build", "-ldflags", "-H=windowsgui", "-o", b.target, ".")
		}
	} else {
		if b.target == "" {
			cmd = exec.Command("go", "build", ".")
		} else {
			cmd = exec.Command("go","build", "-o", b.target, ".")
		}
	}
	cmd.Dir = b.srcdir
	env := os.Environ()
	env = append(env, "CGO_ENABLED=1") // in case someone is trying to cross-compile...

	if goos != "ios" && goos != "android" {
		env = append(env, "GOOS="+goos)
	}
	if goos == "wasm" {
		env = append(env, "GOARCH=wasm")
		env = append(env, "GOOS=js")
	}
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(out))
	}
	return err
}
