package commands

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg" // import image encodings
	"image/png"    // import image encodings
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/fyne-io/image/ico" // import image encodings
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/internal/metadata"
)

const (
	defaultAppBuild   = 1
	defaultAppVersion = "0.0.1"
)

// Packager wraps executables into full GUI app packages.
type Packager struct {
	*appData
	srcDir, dir, exe, os           string
	install, release, distribution bool
	certificate, profile           string // optional flags for releasing
	tags, category                 string
	tempDir                        string
	langs                          []string

	linuxAndBSDMetadata *metadata.LinuxAndBSD
	sourceMetadata      *metadata.AppSource
}

// NewPackager returns a command that can handle the packaging a GUI apps built using Fyne from local source code.
func NewPackager() *Packager {
	return &Packager{appData: &appData{}}
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

func (p *Packager) buildPackage(runner runner, tags []string) ([]string, error) {
	target := p.exe

	b := &Builder{
		os:      p.os,
		srcdir:  p.srcDir,
		target:  target,
		release: p.release,
		tags:    tags,
		runner:  runner,

		appData: p.appData,
	}

	return []string{target}, b.build()
}

func (p *Packager) combinedVersion() string {
	versions := strings.Split(p.AppVersion, ".")
	for len(versions) < 3 {
		versions = append(versions, "0")
	}
	appVersion := strings.Join(versions, ".")

	return fmt.Sprintf("%s.%d", appVersion, p.AppBuild)
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

	var tags []string
	if p.tags != "" {
		tags = strings.Split(p.tags, ",")
	}

	if !util.Exists(p.exe) && !util.IsMobile(p.os) {
		files, err := p.buildPackage(runner, tags)
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
		return p.packageWindows(tags)
	case "android/arm", "android/arm64", "android/amd64", "android/386", "android":
		return p.packageAndroid(p.os, tags)
	case "ios", "iossimulator":
		return p.packageIOS(p.os, tags)
	case "web", "wasm":
		return p.packageWasm()
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
	p.tempDir, err = os.MkdirTemp("", "fyne-package-*")
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
	} else {
		if p.os == "ios" || p.os == "android" {
			return errors.New("parameter -sourceDir is currently not supported for mobile builds. " +
				"Change directory to the main package and try again")
		}
		p.srcDir = util.EnsureAbsPath(p.srcDir)
	}
	os.Chdir(p.srcDir)

	p.appData.Release = p.release

	data, err := metadata.LoadStandard(p.srcDir)
	if err == nil {
		// When icon path specified in metadata file, we should make it relative to metadata file
		if data.Details.Icon != "" {
			data.Details.Icon = util.MakePathRelativeTo(p.srcDir, data.Details.Icon)
		}

		p.appData.mergeMetadata(data)
		p.sourceMetadata = data.Source
		p.langs = data.Languages

		p.linuxAndBSDMetadata = data.LinuxAndBSD
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

func calculateExeName(sourceDir, osys string) string {
	exeName := filepath.Base(sourceDir)
	/* #nosec */
	if data, err := os.ReadFile(filepath.Join(sourceDir, "go.mod")); err == nil {
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

	if osys == "windows" {
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

	out, err := os.CreateTemp(p.tempDir, "fyne-ico-*.png")
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
