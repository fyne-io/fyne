package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mcuadros/go-version"
	"github.com/urfave/cli/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/internal/metadata"
)

// Builder generate the executables.
type Builder struct {
	*appData
	os, srcdir, target string
	goPackage          string
	release            bool
	pprof              bool
	pprofPort          int
	tags               []string
	tagsToParse        string

	customMetadata keyValueFlag

	runner runner
}

// NewBuilder returns a command that can handle the build of GUI apps built using Fyne.
func NewBuilder() *Builder {
	return &Builder{appData: &appData{}}
}

// Build returns the cli command for building fyne applications
func Build() *cli.Command {
	b := NewBuilder()

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
			&cli.BoolFlag{
				Name:        "pprof",
				Usage:       "Enable pprof profiling.",
				Destination: &b.pprof,
			},
			&cli.IntFlag{
				Name:        "pprof-port",
				Usage:       "Specify the port to use for pprof profiling.",
				Value:       6060,
				Destination: &b.pprofPort,
			},
			&cli.GenericFlag{
				Name:  "metadata",
				Usage: "Specify custom metadata key value pair that you do not want to store in your FyneApp.toml (key=value)",
				Value: &b.customMetadata,
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
		b.srcdir = util.EnsureAbsPath(b.srcdir)
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
	b.appData.Release = b.release
	b.appData.CustomMetadata = b.customMetadata.m

	return b.build()
}

func isWeb(goos string) bool {
	return goos == "js" || goos == "wasm" || goos == "web"
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

func (b *Builder) build() error {
	goos := b.os
	if goos == "" {
		goos = targetOS()
	}

	fyneGoModRunner := b.updateAndGetGoExecutable(goos)

	srcdir, err := b.computeSrcDir(fyneGoModRunner)
	if err != nil {
		return err
	}

	b.updateToDefaultIconIfNotSet(srcdir)

	if b.pprof {
		close, err := injectPprofFile(fyneGoModRunner, srcdir, b.pprofPort)
		if err != nil {
			fyne.LogError("Failed to inject pprof file, omitting pprof", err)
		} else if close != nil {
			defer close()
		}
	}

	close, err := injectMetadataIfPossible(fyneGoModRunner, srcdir, b.appData, createMetadataInitFile)
	if err != nil {
		fyne.LogError("Failed to inject metadata init file, omitting metadata", err)
	} else if close != nil {
		defer close()
	}

	args := []string{"build"}
	env := os.Environ()

	if goos == "darwin" {
		appendEnv(&env, "CGO_CFLAGS", "-mmacosx-version-min=10.11")
		appendEnv(&env, "CGO_LDFLAGS", "-mmacosx-version-min=10.11")
	}

	ldFlags := extractLdflagsFromGoFlags()
	if !isWeb(goos) {
		env = append(env, "CGO_ENABLED=1") // in case someone is trying to cross-compile...

		if b.release {
			ldFlags += " -s -w"
			args = append(args, "-trimpath")
		}

		if goos == "windows" {
			ldFlags += " -H=windowsgui"
		}
	}

	if len(ldFlags) > 0 {
		args = append(args, "-ldflags", strings.TrimSpace(ldFlags))
	}

	if b.target != "" {
		args = append(args, "-o", b.target)
	}

	// handle build tags
	tags := b.tags
	if b.release {
		tags = append(tags, "release")
	}
	if ok, set := b.appData.Migrations["fyneDo"]; ok && set {
		tags = append(tags, "migrated_fynedo")
	}
	if len(tags) > 0 {
		args = append(args, "-tags", strings.Join(tags, ","))
	}

	if b.goPackage != "" {
		args = append(args, b.goPackage)
	}

	if goos != "ios" && goos != "android" && !isWeb(goos) {
		env = append(env, "GOOS="+goos)
	} else if goos == "web" || goos == "wasm" {
		env = append(env, "GOARCH=wasm")
		env = append(env, "GOOS=js")
	}

	b.runner.setDir(b.srcdir)
	b.runner.setEnv(env)
	out, err := b.runner.runOutput(args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", string(out))
	}
	return err
}

func (b *Builder) computeSrcDir(fyneGoModRunner runner) (string, error) {
	if b.goPackage == "" || b.goPackage == "." {
		return b.srcdir, nil
	}

	srcdir := b.srcdir
	if strings.HasPrefix(b.goPackage, "."+string(os.PathSeparator)) ||
		strings.HasPrefix(b.goPackage, ".."+string(os.PathSeparator)) {
		srcdir = filepath.Join(srcdir, b.goPackage)
	} else if strings.HasPrefix(b.goPackage, string(os.PathSeparator)) {
		srcdir = b.goPackage
	} else {
		return "", fmt.Errorf("unrecognized go package: %s", b.goPackage)
	}
	return srcdir, nil
}

func (b *Builder) updateToDefaultIconIfNotSet(srcdir string) {
	if b.icon == "" {
		defaultIcon := filepath.Join(srcdir, "Icon.png")
		if util.Exists(defaultIcon) {
			b.icon = defaultIcon
		}
	}
}

func injectPprofFile(runner runner, srcdir string, port int) (func(), error) {
	pprofInitFilePath := filepath.Join(srcdir, "fyne_pprof.go")
	pprofInitFile, err := os.Create(pprofInitFilePath)
	if err != nil {
		return func() {}, err
	}
	defer pprofInitFile.Close()

	pprofInfo := struct {
		Port int
	}{
		Port: port,
	}

	err = templates.FynePprofInit.Execute(pprofInitFile, pprofInfo)
	if err != nil {
		os.Remove(pprofInitFilePath)
		return func() {}, err
	}

	return func() { os.Remove(pprofInitFilePath) }, nil
}

func (b *Builder) updateAndGetGoExecutable(goos string) runner {
	fyneGoModRunner := b.runner
	if b.runner == nil {
		fyneGoModRunner = newCommand("go")
		goBin := os.Getenv("GO")
		if goBin != "" {
			fyneGoModRunner = newCommand(goBin)
			b.runner = fyneGoModRunner
		} else {
			b.runner = newCommand("go")
		}
	}
	return fyneGoModRunner
}

func createMetadataInitFile(srcdir string, app *appData) (func(), error) {
	data, err := metadata.LoadStandard(srcdir)
	if err == nil {
		// When icon path specified in metadata file, we should make it relative to metadata file
		if data.Details.Icon != "" {
			data.Details.Icon = util.MakePathRelativeTo(srcdir, data.Details.Icon)
		}

		app.mergeMetadata(data)
	}

	metadataInitFilePath := filepath.Join(srcdir, "fyne_metadata_init.go")
	metadataInitFile, err := os.Create(metadataInitFilePath)
	if err != nil {
		return func() {}, err
	}
	defer metadataInitFile.Close()

	app.ResGoString = "nil"
	if app.icon != "" {
		res, err := fyne.LoadResourceFromPath(app.icon)
		if err != nil {
			fyne.LogError("Unable to load metadata icon file "+app.icon, err)
			return func() { os.Remove(metadataInitFilePath) }, err
		}

		res = metadata.ScaleIcon(res, 512)

		// The return type of fyne.LoadResourceFromPath is always a *fyne.StaticResource.
		app.ResGoString = res.(*fyne.StaticResource).GoString()
	}

	err = templates.FyneMetadataInit.Execute(metadataInitFile, app)
	if err != nil {
		fyne.LogError("Error executing metadata template", err)
	}

	return func() { os.Remove(metadataInitFilePath) }, err
}

func injectMetadataIfPossible(runner runner, srcdir string, app *appData,
	createMetadataInitFile func(srcdir string, app *appData) (func(), error),
) (func(), error) {
	fyneGoModVersion, err := getFyneGoModVersion(runner)
	if err != nil {
		return nil, err
	}

	fyneGoModVersion = normaliseVersion(fyneGoModVersion)
	fyneGoModVersionConstraint := version.NewConstrainGroupFromString(">=2.2")
	if fyneGoModVersion != "master" && !fyneGoModVersionConstraint.Match(fyneGoModVersion) {
		return nil, nil
	}

	fyneGoModVersionAtLeast2_3 := version.NewConstrainGroupFromString(">=2.3")
	if fyneGoModVersionAtLeast2_3.Match(fyneGoModVersion) {
		app.VersionAtLeast2_3 = true
	}

	return createMetadataInitFile(srcdir, app)
}

func targetOS() string {
	osEnv, ok := os.LookupEnv("GOOS")
	if ok {
		return osEnv
	}

	return runtime.GOOS
}

func appendEnv(env *[]string, varName, value string) {
	for i := range *env {
		keyValue := strings.SplitN((*env)[i], "=", 2)

		if keyValue[0] == varName {
			(*env)[i] += " " + value
			return
		}
	}

	*env = append(*env, varName+"="+value)
}

func extractLdflagsFromGoFlags() string {
	goFlags := os.Getenv("GOFLAGS")

	ldFlags, goFlags := extractLdFlags(goFlags)
	if goFlags != "" {
		os.Setenv("GOFLAGS", goFlags)
	} else {
		os.Unsetenv("GOFLAGS")
	}

	return ldFlags
}

func extractLdFlags(goFlags string) (string, string) {
	if goFlags == "" {
		return "", ""
	}

	flags := strings.Fields(goFlags)
	ldflags := ""
	newGoFlags := ""

	for _, flag := range flags {
		if after, ok := strings.CutPrefix(flag, "-ldflags="); ok {
			ldflags += after + " "
		} else {
			newGoFlags += flag + " "
		}
	}

	ldflags = strings.TrimSpace(ldflags)
	newGoFlags = strings.TrimSpace(newGoFlags)

	return ldflags, newGoFlags
}

func normaliseVersion(str string) string {
	if str == "master" {
		return str
	}

	if pos := strings.Index(str, "-0.20"); pos != -1 {
		str = str[:pos] + "-dev"
	} else if pos = strings.Index(str, "-rc"); pos != -1 {
		str = str[:pos] + "-dev"
	}
	return version.Normalize(str)
}
