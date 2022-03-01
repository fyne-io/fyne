package commands

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/sys/execabs"
)

type builder struct {
	os, srcdir string
	release    bool
	tags       []string

	id, name, version string
	buildNum          int
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

	meta := b.generateMetaLDFlags()
	if goos == "windows" {
		if b.release {
			args = append(args, "-ldflags", "-s -w -H=windowsgui "+meta)
		} else {
			args = append(args, "-ldflags", "-H=windowsgui "+meta)
		}
	} else if b.release {
		args = append(args, "-ldflags", "-s -w "+meta)
	} else {
		args = append(args, "-ldflags", meta)
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

func (b *builder) generateMetaLDFlags() string {
	var vars []string
	if b.id != "" {
		vars = append(vars, "-X 'fyne.io/fyne/v2/internal/app.MetaID="+b.id+"'")
	}
	if b.name != "" {
		vars = append(vars, "-X 'fyne.io/fyne/v2/internal/app.MetaName="+b.name+"'")
	}
	if b.version != "" {
		vars = append(vars, "-X 'fyne.io/fyne/v2/internal/app.MetaVersion="+b.version+"'")
	}
	if b.buildNum != 0 {
		vars = append(vars, fmt.Sprintf("-X 'fyne.io/fyne/v2/internal/app.MetaBuild=%d'", b.buildNum))
	}

	return strings.Join(vars, " ")
}

func targetOS() string {
	osEnv, ok := os.LookupEnv("GOOS")
	if ok {
		return osEnv
	}

	return runtime.GOOS
}
