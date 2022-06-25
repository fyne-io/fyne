package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	version "github.com/mcuadros/go-version"
	"github.com/urfave/cli/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/metadata"
	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
)

// Builder generate the executables.
type Builder struct {
	*appData
	os, srcdir, target string
	goPackage          string
	release            bool
	tags               []string
	tagsToParse        string

	runner runner
}

// Build returns the cli command for building fyne applications
func Build() *cli.Command {
	b := &Builder{appData: &appData{}}

	return &cli.Command{
		Name:        "build",
		Usage:       "Build an application.",
		Description: "You can specify --target to define the OS to build for. The executable file will default to an appropriate name but can be overridden using -o.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"os"},
				Usage:       "The mobile platform to target (android, android/arm, android/arm64, android/amd64, android/386, ios, iossimulator).",
				Destination: &b.os,
			},
			&cli.StringFlag{
				Name:        "sourceDir",
				Aliases:     []string{"src"},
				Usage:       "The directory to package, if executable is not set.",
				Destination: &b.srcdir,
			},
			&cli.StringFlag{
				Name:        "tags",
				Usage:       "A comma-separated list of build tags.",
				Destination: &b.tagsToParse,
			},
			&cli.BoolFlag{
				Name:        "release",
				Usage:       "Enable installation in release mode (disable debug etc).",
				Destination: &b.release,
			},
			&cli.StringFlag{
				Name:        "o",
				Usage:       "Specify a name for the output file, default is based on the current directory.",
				Destination: &b.target,
			},
		},
		Action: func(ctx *cli.Context) error {
			argCount := ctx.Args().Len()
			if argCount > 0 {
				if argCount != 1 {
					return fmt.Errorf("incorrect amount of path provided")
				}
				b.goPackage = ctx.Args().First()
			}

			return b.Build()
		},
	}
}

// Build parse the tags and start building
func (b *Builder) Build() error {
	if b.srcdir != "" {
		dirStat, err := os.Stat(b.srcdir)
		if err != nil {
			return err
		}
		if !dirStat.IsDir() {
			return fmt.Errorf("specified source directory is not a valid directory")
		}
	}
	if b.tagsToParse != "" {
		b.tags = strings.Split(b.tagsToParse, ",")
	}

	return b.build()
}

func checkVersion(output string, versionConstraint *version.ConstraintGroup) error {
	split := strings.Split(output, " ")
	// We are expecting something like: `go version goX.Y OS`
	if len(split) != 4 || split[0] != "go" || split[1] != "version" || len(split[2]) < 5 || split[2][:2] != "go" {
		return fmt.Errorf("invalid output for `go version`: `%s`", output)
	}

	normalized := version.Normalize(split[2][2:len(split[2])])
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

type goModEdit struct {
	Module struct {
		Path string
	}
	Require []struct {
		Path    string
		Version string
	}
}

func getFyneGoModVersion(runner runner) (string, error) {
	dependenciesOutput, err := runner.runOutput("mod", "edit", "-json")
	if err != nil {
		return "", err
	}

	var parsed goModEdit
	err = json.Unmarshal(dependenciesOutput, &parsed)
	if err != nil {
		return "", err
	}

	if parsed.Module.Path == "fyne.io/fyne/v2" {
		return "master", nil
	}

	for _, dep := range parsed.Require {
		if dep.Path == "fyne.io/fyne/v2" {
			return dep.Version, nil
		}
	}

	return "", fmt.Errorf("fyne version not found")
}

func (b *Builder) createMetadataInitFile() (func(), error) {
	data, err := metadata.LoadStandard(b.srcdir)
	if err == nil {
		mergeMetadata(b.appData, data)
	}

	metadataInitFilePath := filepath.Join(b.srcdir, "fyne_metadata_init.go")
	metadataInitFile, err := os.Create(metadataInitFilePath)
	if err != nil {
		return func() {}, err
	}
	defer metadataInitFile.Close()

	err = templates.FyneMetadataInit.Execute(metadataInitFile, b.appData)
	if err == nil {
		if b.icon != "" {
			writeResource(b.icon, "fyneMetadataIcon", metadataInitFile)
		}
	}

	return func() { os.Remove(metadataInitFilePath) }, err
}

func (b *Builder) injectMetadataIfPossible(runner runner, createMetadataInitFile func() (func(), error)) (func(), error) {
	fyneGoModVersion, err := getFyneGoModVersion(runner)
	if err != nil {
		return nil, err
	}

	fyneGoModVersionNormalized := version.Normalize(fyneGoModVersion)
	fyneGoModVersionConstraint := version.NewConstrainGroupFromString(">=2.2")
	if fyneGoModVersion != "master" && !fyneGoModVersionConstraint.Match(fyneGoModVersionNormalized) {
		return nil, nil
	}

	return createMetadataInitFile()
}

func (b *Builder) build() error {
	var versionConstraint *version.ConstraintGroup

	goos := b.os
	if goos == "" {
		goos = targetOS()
	}

	fyneGoModRunner := b.runner
	if b.runner == nil {
		fyneGoModRunner = newCommand("go")
		if goos != "gopherjs" {
			b.runner = newCommand("go")
		} else {
			b.runner = newCommand("gopherjs")
		}
	}

	close, err := b.injectMetadataIfPossible(fyneGoModRunner, b.createMetadataInitFile)
	if err != nil {
		fyne.LogError("Failed to inject metadata init file, omitting metadata", err)
	} else if close != nil {
		defer close()
	}

	args := []string{"build"}
	env := os.Environ()

	if goos == "darwin" {
		env = append(env, "CGO_CFLAGS=-mmacosx-version-min=10.11", "CGO_LDFLAGS=-mmacosx-version-min=10.11")
	}

	if !isWeb(goos) {
		env = append(env, "CGO_ENABLED=1") // in case someone is trying to cross-compile...

		if goos == "windows" {
			if b.release {
				args = append(args, "-ldflags", "-s -w -H=windowsgui ")
			} else {
				args = append(args, "-ldflags", "-H=windowsgui ")
			}
		} else if b.release {
			args = append(args, "-ldflags", "-s -w ")
		}
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

	if b.goPackage != "" {
		args = append(args, b.goPackage)
	}

	if goos != "ios" && goos != "android" && !isWeb(goos) {
		env = append(env, "GOOS="+goos)
	} else if goos == "wasm" {
		versionConstraint = version.NewConstrainGroupFromString(">=1.17")
		env = append(env, "GOARCH=wasm")
		env = append(env, "GOOS=js")
	} else if goos == "gopherjs" {
		_, err := b.runner.runOutput("version")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not execute `gopherjs version`. Please do `go install github.com/gopherjs/gopherjs@latest`.\n")
			return err
		}
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

func targetOS() string {
	osEnv, ok := os.LookupEnv("GOOS")
	if ok {
		return osEnv
	}

	return runtime.GOOS
}
