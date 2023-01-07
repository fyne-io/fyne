package commands

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg" // import image encodings
	"image/png"    // import image encodings
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/fyne-io/image/ico" // import image encodings
	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/cmd/fyne/internal/metadata"
)

const (
	defaultAppBuild   = 1
	defaultAppVersion = "0.0.1"
)

type appData struct {
	icon, Name        string
	AppID, AppVersion string
	AppBuild          int
	ResGoString       string
	Release           bool
	CustomMetadata    map[string]string
	VersionAtLeast2_3 bool
}

// Package returns the cli command for packaging fyne applications
func Package() *cli.Command {
	p := &Packager{appData: &appData{}}

	return &cli.Command{
		Name:        "package",
		Usage:       "Packages an application for distribution.",
		Description: "You may specify the -executable to package, otherwise -sourceDir will be built.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "target",
				Aliases:     []string{"os"},
				Usage:       "The mobile platform to target (android, android/arm, android/arm64, android/amd64, android/386, ios, iossimulator, wasm, gopherjs, web).",
				Destination: &p.os,
			},
			&cli.StringFlag{
				Name:        "executable",
				Aliases:     []string{"exe"},
				Usage:       "The path to the executable, default is the current dir main binary",
				Destination: &p.exe,
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "The name of the application, default is the executable file name",
				Destination: &p.Name,
			},
			&cli.StringFlag{
				Name:        "tags",
				Usage:       "A comma-separated list of build tags.",
				Destination: &p.tags,
			},
			&cli.StringFlag{
				Name:        "appVersion",
				Usage:       "Version number in the form x, x.y or x.y.z semantic version",
				Destination: &p.AppVersion,
			},
			&cli.IntFlag{
				Name:        "appBuild",
				Usage:       "Build number, should be greater than 0 and incremented for each build",
				Destination: &p.AppBuild,
			},
			&cli.StringFlag{
				Name:        "sourceDir",
				Aliases:     []string{"src"},
				Usage:       "The directory to package, if executable is not set.",
				Destination: &p.srcDir,
			},
			&cli.StringFlag{
				Name:        "icon",
				Usage:       "The name of the application icon file.",
				Value:       "",
				Destination: &p.icon,
			},
			&cli.StringFlag{
				Name:        "appID",
				Aliases:     []string{"id"},
				Usage:       "For Android, darwin, iOS and Windows targets an appID in the form of a reversed domain name is required, for ios this must match a valid provisioning profile",
				Destination: &p.AppID,
			},
			&cli.StringFlag{
				Name:        "certificate",
				Aliases:     []string{"cert"},
				Usage:       "iOS/macOS/Windows: name of the certificate to sign the build",
				Destination: &p.certificate,
			},
			&cli.StringFlag{
				Name:        "profile",
				Usage:       "iOS/macOS: name of the provisioning profile for this build",
				Destination: &p.profile,
			},
			&cli.BoolFlag{
				Name:        "release",
				Usage:       "Enable installation in release mode (disable debug etc).",
				Destination: &p.release,
			},
			&cli.GenericFlag{
				Name:  "metadata",
				Usage: "Specify custom metadata key value pair that you do not want to store in your FyneApp.toml (key=value)",
				Value: &p.customMetadata,
			},
		},
		Action: func(_ *cli.Context) error {
			if p.customMetadata.m == nil {
				p.customMetadata.m = map[string]string{}
			}

			return p.Package()
		},
	}
}

// Packager wraps executables into full GUI app packages.
type Packager struct {
	*appData
	srcDir, dir, exe, os           string
	install, release, distribution bool
	certificate, profile           string // optional flags for releasing
	tags, category                 string
	tempDir                        string

	customMetadata keyValueFlag
}

// AddFlags adds the flags for interacting with the package command.
//
// Deprecated: Access to the individual cli commands are being removed.
func (p *Packager) AddFlags() {
	flag.StringVar(&p.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, freebsd, ios, linux, netbsd, openbsd, windows, wasm)")
	flag.StringVar(&p.exe, "executable", "", "Specify an existing binary instead of building before package")
	flag.StringVar(&p.srcDir, "sourceDir", "", "The directory to package, if executable is not set")
	flag.StringVar(&p.Name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&p.icon, "icon", "", "The name of the application icon file")
	flag.StringVar(&p.AppID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.StringVar(&p.AppVersion, "appVersion", "", "Version number in the form x, x.y or x.y.z semantic version")
	flag.IntVar(&p.AppBuild, "appBuild", 0, "Build number, should be greater than 0 and incremented for each build")
	flag.BoolVar(&p.release, "release", false, "Should this package be prepared for release? (disable debug etc)")
	flag.StringVar(&p.tags, "tags", "", "A comma-separated list of build tags")
}

// PrintHelp prints the help for the package command.
//
// Deprecated: Access to the individual cli commands are being removed.
func (*Packager) PrintHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for installation and testing.")
	fmt.Println(indent, "You may specify the -executable to package, otherwise -sourceDir will be built.")
	fmt.Println(indent, "Command usage: fyne package [parameters]")
}

// Run runs the package command.
//
// Deprecated: A better version will be exposed in the future.
func (p *Packager) Run(_ []string) {
	err := p.validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = p.doPackage(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

// Package starts the packaging process
func (p *Packager) Package() error {
	err := p.validate()
	if err != nil {
		return err
	}

	return p.packageWithoutValidate()
}

func (p *Packager) packageWithoutValidate() error {
	err := p.doPackage(nil)
	if err != nil {
		return err
	}

	data, err := metadata.LoadStandard(p.srcDir)
	if err != nil {
		return nil // no metadata to update
	}
	data.Details.Build++
	return metadata.SaveStandard(data, p.srcDir)
}

func (p *Packager) buildPackage(runner runner) ([]string, error) {
	var tags []string
	if p.tags != "" {
		tags = strings.Split(p.tags, ",")
	}
	if p.os != "web" {
		b := &Builder{
			os:      p.os,
			srcdir:  p.srcDir,
			target:  p.exe,
			release: p.release,
			tags:    tags,
			runner:  runner,

			appData: p.appData,
		}

		return []string{p.exe}, b.build()
	}

	bWasm := &Builder{
		os:      "wasm",
		srcdir:  p.srcDir,
		target:  p.exe + ".wasm",
		release: p.release,
		tags:    tags,
		runner:  runner,

		appData: p.appData,
	}

	err := bWasm.build()
	if err != nil {
		return nil, err
	}

	if runtime.GOOS == "windows" {
		return []string{bWasm.target}, nil
	}

	bGopherJS := &Builder{
		os:      "gopherjs",
		srcdir:  p.srcDir,
		target:  p.exe + ".js",
		release: p.release,
		tags:    tags,
		runner:  runner,

		appData: p.appData,
	}

	err = bGopherJS.build()
	if err != nil {
		return nil, err
	}

	return []string{bWasm.target, bGopherJS.target}, nil
}

func (p *Packager) combinedVersion() string {
	return fmt.Sprintf("%s.%d", p.AppVersion, p.AppBuild)
}

func (p *Packager) doPackage(runner runner) error {
	// sensible defaults - validation deemed them optional
	if p.AppVersion == "" {
		p.AppVersion = defaultAppVersion
	}
	if p.AppBuild <= 0 {
		p.AppBuild = defaultAppBuild
	}
	defer os.RemoveAll(p.tempDir)

	if !util.Exists(p.exe) && !util.IsMobile(p.os) {
		files, err := p.buildPackage(runner)
		if err != nil {
			return fmt.Errorf("error building application: %w", err)
		}
		for _, file := range files {
			if p.os != "web" && !util.Exists(file) {
				return fmt.Errorf("unable to build directory to expected executable, %s", file)
			}
		}
		if p.os != "windows" {
			defer p.removeBuild(files)
		}
	}
	if util.IsMobile(p.os) { // we don't use the normal build command for mobile so inject before gomobile...
		close, err := injectMetadataIfPossible(newCommand("go"), p.dir, p.appData, createMetadataInitFile)
		if err != nil {
			fyne.LogError("Failed to inject metadata init file, omitting metadata", err)
		} else if close != nil {
			defer close()
		}
	}

	switch p.os {
	case "darwin":
		return p.packageDarwin()
	case "linux", "openbsd", "freebsd", "netbsd":
		return p.packageUNIX()
	case "windows":
		return p.packageWindows()
	case "android/arm", "android/arm64", "android/amd64", "android/386", "android":
		return p.packageAndroid(p.os)
	case "ios", "iossimulator":
		return p.packageIOS(p.os)
	case "wasm":
		return p.packageWasm()
	case "gopherjs":
		return p.packageGopherJS()
	case "web":
		return p.packageWeb()
	default:
		return fmt.Errorf("unsupported target operating system \"%s\"", p.os)
	}
}

func (p *Packager) removeBuild(files []string) {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			log.Println("Unable to remove temporary build file", p.exe)
		}
	}
}

func (p *Packager) validate() (err error) {
	p.tempDir, err = ioutil.TempDir("", "fyne-package-*")
	defer func() {
		if err != nil {
			_ = os.RemoveAll(p.tempDir)
		}
	}()

	if p.os == "" {
		p.os = targetOS()
	}
	baseDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get the current directory, needed to find main executable: %w", err)
	}
	if p.dir == "" {
		p.dir = baseDir
	}
	if p.srcDir == "" {
		p.srcDir = baseDir
	} else if p.os == "ios" || p.os == "android" {
		return errors.New("parameter -sourceDir is currently not supported for mobile builds. " +
			"Change directory to the main package and try again")
	}
	os.Chdir(p.srcDir)

	p.appData.CustomMetadata = p.customMetadata.m

	data, err := metadata.LoadStandard(p.srcDir)
	if err == nil {
		p.appData.Release = p.release
		mergeMetadata(p.appData, data)
	}

	exeName := calculateExeName(p.srcDir, p.os)

	if p.exe == "" {
		p.exe = filepath.Join(p.srcDir, exeName)

		if util.Exists(p.exe) { // the exe was not specified, assume stale
			p.removeBuild([]string{p.exe})
		}
	} else if p.os == "ios" || p.os == "android" {
		_, _ = fmt.Fprint(os.Stderr, "Parameter -executable is ignored for mobile builds.\n")
	}

	if p.Name == "" {
		p.Name = exeName
	}
	if p.icon == "" || p.icon == "Icon.png" {
		p.icon = filepath.Join(p.srcDir, "Icon.png")
	}
	if !util.Exists(p.icon) {
		return errors.New("Missing application icon at \"" + p.icon + "\"")
	}
	if strings.ToLower(filepath.Ext(p.icon)) != ".png" {
		tmp, err := p.normaliseIcon(p.icon)
		if err != nil {
			return err
		}
		p.icon = tmp
	}

	p.AppID, err = validateAppID(p.AppID, p.os, p.Name, p.release)
	if err != nil {
		return err
	}
	if p.AppVersion != "" && !isValidVersion(p.AppVersion) {
		return errors.New("invalid -appVersion parameter, integer and '.' characters only up to x.y.z")
	}

	return nil
}

func calculateExeName(sourceDir, os string) string {
	exeName := filepath.Base(sourceDir)
	/* #nosec */
	if data, err := ioutil.ReadFile(filepath.Join(sourceDir, "go.mod")); err == nil {
		modulePath := modfile.ModulePath(data)
		moduleName, _, ok := module.SplitPathVersion(modulePath)
		if ok {
			paths := strings.Split(moduleName, "/")
			name := paths[len(paths)-1]
			if name != "" {
				exeName = name
			}
		}
	}

	if os == "windows" {
		exeName = exeName + ".exe"
	}

	return exeName
}

func isValidVersion(ver string) bool {
	nums := strings.Split(ver, ".")
	if len(nums) == 0 || len(nums) > 3 {
		return false
	}
	for _, num := range nums {
		if _, err := strconv.Atoi(num); err != nil {
			return false
		}
	}
	return true
}

func appendCustomMetadata(customMetadata *map[string]string, fromFile map[string]string) {
	for key, value := range fromFile {
		_, ok := (*customMetadata)[key]
		if ok {
			continue
		}
		(*customMetadata)[key] = value
	}
}

func mergeMetadata(p *appData, data *metadata.FyneApp) {
	if p.icon == "" {
		p.icon = data.Details.Icon
	}
	if p.Name == "" {
		p.Name = data.Details.Name
	}
	if p.AppID == "" {
		p.AppID = data.Details.ID
	}
	if p.AppVersion == "" {
		p.AppVersion = data.Details.Version
	}
	if p.AppBuild == 0 {
		p.AppBuild = data.Details.Build
	}
	if p.Release {
		appendCustomMetadata(&p.CustomMetadata, data.Release)
	} else {
		appendCustomMetadata(&p.CustomMetadata, data.Development)
	}
}

// normaliseIcon takes a non-png image file and converts it to PNG for use in packaging.
// Successful conversion will return a path to the new file.
// Any errors that occur will be returned with an empty string for new path.
func (p *Packager) normaliseIcon(path string) (string, error) {
	// convert icon
	img, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open source image: %w", err)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return "", fmt.Errorf("failed to decode source image: %w", err)
	}

	out, err := ioutil.TempFile(p.tempDir, "fyne-ico-*.png")
	if err != nil {
		return "", fmt.Errorf("failed to open image output file: %w", err)
	}
	tmpPath := out.Name()
	defer out.Close()

	err = png.Encode(out, srcImg)
	return tmpPath, err
}

func validateAppID(appID, os, name string, release bool) (string, error) {
	// old darwin compatibility
	if os == "darwin" {
		if appID == "" {
			return "com.example." + name, nil
		}
	} else if os == "ios" || util.IsAndroid(os) || (os == "windows" && release) {
		// all mobile, and for windows when releasing, needs a unique id - usually reverse DNS style
		if appID == "" {
			return "", errors.New("missing appID parameter for package")
		} else if !strings.Contains(appID, ".") {
			return "", errors.New("appID must be globally unique and contain at least 1 '.'")
		} else if util.IsAndroid(os) {
			if strings.Contains(appID, "-") {
				return "", errors.New("appID can not contain '-'")
			}

			// appID package names can not start with '_' or a number
			packageNames := strings.Split(appID, ".")
			for _, name := range packageNames {
				if len(name) == 0 {
					continue
				}

				if name[0] == '_' {
					return "", fmt.Errorf("appID package names can not start with '_' (%s)", name)
				} else if name[0] >= '0' && name[0] <= '9' {
					return "", fmt.Errorf("appID package names can not start with a number (%s)", name)
				}
			}
		}
	}

	return appID, nil
}
