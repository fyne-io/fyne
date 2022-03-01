package commands

import (
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	version "github.com/mcuadros/go-version"
)

type builder struct {
	os, srcdir, target string
	release            bool
	tags               []string

	runner runner

	icon, id, name, version string
	buildNum                int
}

func checkVersion(output string, versionConstraint *version.ConstraintGroup) error {
	split := strings.Split(output, " ")
	// We are expecting something like: `go version goX.Y OS`
	if len(split) != 4 || split[0] != "go" || split[1] != "version" || len(split[2]) < 5 || split[2][:2] != "go" {
		return fmt.Errorf("invalid output for `go version`: `%s`", output)
	}

	normalized := version.Normalize(split[2][2 : len(split[2])-2])
	if !versionConstraint.Match(normalized) {
		return fmt.Errorf("expected go version %v got `%v`", versionConstraint.GetConstraints(), normalized)
	}

	return nil
}

func isWeb(goos string) bool {
	return goos == "gopherjs" || goos == "wasm"
}

func checkGoVersion(runner runner, versionConstraint *version.ConstraintGroup) error {
	if versionConstraint == nil {
		return nil
	}

	goVersion, err := runner.runOutput("version")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(goVersion))
		return err
	}

	return checkVersion(string(goVersion), versionConstraint)
}

func (b *builder) build() error {
	var versionConstraint *version.ConstraintGroup

	goos := b.os
	if goos == "" {
		goos = targetOS()
	}

	if b.runner == nil {
		if goos != "gopherjs" {
			b.runner = newCommand("go")
		} else {
			b.runner = newCommand("gopherjs")
		}
	}

	args := []string{"build"}
	env := os.Environ()

	if goos == "darwin" {
		env = append(env, "CGO_CFLAGS=-mmacosx-version-min=10.11", "CGO_LDFLAGS=-mmacosx-version-min=10.11")
	}

	meta := b.generateMetaLDFlags()
	if !isWeb(goos) {
		env = append(env, "CGO_ENABLED=1") // in case someone is trying to cross-compile...

		if goos == "windows" {
			if b.release {
				args = append(args, "-ldflags", "-s -w -H=windowsgui "+meta)
			} else {
				args = append(args, "-ldflags", "-H=windowsgui "+meta)
			}
		} else if b.release {
			args = append(args, "-ldflags", "-s -w "+meta)
		} else if meta != "" {
			args = append(args, "-ldflags", meta)
		}
	} else if meta != "" {
		args = append(args, "-ldflags", meta)
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
		if goos == "gopherjs" {
			args = append(args, "--tags")
		} else {
			args = append(args, "-tags")
		}
		args = append(args, strings.Join(tags, ","))
	}

	if goos != "ios" && goos != "android" && !isWeb(goos) {
		env = append(env, "GOOS="+goos)
	} else if goos == "wasm" {
		versionConstraint = version.NewConstrainGroupFromString(">=1.17")
		env = append(env, "GOARCH=wasm")
		env = append(env, "GOOS=js")
	}

	if err := checkGoVersion(b.runner, versionConstraint); err != nil {
		return err
	}

	b.runner.setDir(b.srcdir)
	b.runner.setEnv(env)
	out, err := b.runner.runOutput(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(out))
	}
	return err
}

func (b *builder) generateMetaLDFlags() string {
	buildIDString := ""
	if b.buildNum > 0 {
		buildIDString = strconv.Itoa(b.buildNum)
	}
	iconBytes := ""
	if b.icon != "" {
		res, err := fyne.LoadResourceFromPath(b.icon)
		if err != nil {
			fyne.LogError("Unable to load file "+b.icon, err)
		} else {
			staticRes, ok := res.(*fyne.StaticResource)
			if !ok {
				fyne.LogError("Unable to format resource", fmt.Errorf("unexpected resource type %T", res))
			} else {
				iconBytes = base64.StdEncoding.EncodeToString(staticRes.StaticContent)
			}
		}
	}

	inserts := [][2]string{
		{"MetaID", b.id},
		{"MetaName", b.name},
		{"MetaVersion", b.version},
		{"MetaBuild", buildIDString},
		{"MetaIcon", iconBytes},
	}

	var vars []string
	for _, item := range inserts {
		if item[1] != "" {
			vars = append(vars, "-X 'fyne.io/fyne/v2/internal/app."+item[0]+"="+item[1]+"'")
		}
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
