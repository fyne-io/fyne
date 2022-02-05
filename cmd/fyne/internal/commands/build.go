package commands

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/sys/execabs"
)

type builder struct {
	os, srcdir, target string
	release            bool
	tags               []string
}

func (b *builder) build() error {
	goos := b.os
	if goos == "" {
		goos = targetOS()
	}

	args := []string{"build"}
	env := append(os.Environ(), "CGO_ENABLED=1") // in case someone is trying to cross-compile...
	if goos == "darwin" {
		env = append(env, "CGO_CFLAGS=-mmacosx-version-min=10.11", "CGO_LDFLAGS=-mmacosx-version-min=10.11")
	}

	if goos == "windows" {
		if b.release {
			args = append(args, "-ldflags", "-s -w -H=windowsgui")
		} else {
			args = append(args, "-ldflags", "-H=windowsgui")
		}
	} else if b.release {
		args = append(args, "-ldflags", "-s -w")
	}

	if b.target != "" {
		args = append(args, "-o", b.target)
	}

	// handle build tags
	tags := b.tags
	if b.release {
		tags = append(tags, "release")
	}
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}

	cmd := execabs.Command("go", args...)
	cmd.Dir = b.srcdir
	if goos != "ios" && goos != "android" && goos != "wasm" {
		env = append(env, "GOOS="+goos)
	} else if goos == "wasm" {
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

func targetOS() string {
	osEnv, ok := os.LookupEnv("GOOS")
	if ok {
		return osEnv
	}

	return runtime.GOOS
}
