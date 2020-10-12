package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne"
	"github.com/pkg/errors"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func ensureSubDir(parent, name string) string {
	path := filepath.Join(parent, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fyne.LogError("Failed to create dirrectory", err)
		}
	}
	return path
}

func copyFile(src, tgt string) error {
	return copyFileMode(src, tgt, 0644)
}

func copyExeFile(src, tgt string) error {
	return copyFileMode(src, tgt, 0755)
}

func copyFileMode(src, tgt string, perm os.FileMode) error {
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	source, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.OpenFile(tgt, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return err
}

// Declare conformity to command interface
var _ command = (*packager)(nil)

type packager struct {
	name, srcDir, dir, exe, icon string
	os, appID                    string
	install, release             bool
}

func (p *packager) addFlags() {
	flag.StringVar(&p.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, freebsd, ios, linux, netbsd, openbsd, windows)")
	flag.StringVar(&p.exe, "executable", "", "The path to the executable, default is the current dir main binary")
	flag.StringVar(&p.srcDir, "sourceDir", "", "The directory to package, if executable is not set")
	flag.StringVar(&p.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&p.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&p.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.BoolVar(&p.release, "release", false, "Should this package be prepared for release? (disable debug etc)")
}

func (*packager) printHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for distribution.")
	fmt.Println(indent, "You may specify the -executable to package, otherwise -sourceDir will be built.")
	fmt.Println(indent, "Command usage: fyne package [parameters]")
}

func (p *packager) buildPackage() error {
	b := &builder{
		os:      p.os,
		srcdir:  p.srcDir,
		release: p.release,
	}

	return b.build()
}

func (p *packager) run(_ []string) {
	err := p.validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = p.doPackage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func (p *packager) validate() error {
	if p.os == "" {
		p.os = runtime.GOOS
	}
	baseDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "Unable to get the current directory, needed to find main executable")
	}
	if p.dir == "" {
		p.dir = baseDir
	}
	if p.srcDir == "" {
		p.srcDir = baseDir
	} else if p.os == "ios" || p.os == "android" {
		return errors.New("Parameter -sourceDir is currently not supported for mobile builds.\n" +
			"Change directory to the main package and try again.")
	}

	exeName := calculateExeName(p.srcDir, p.os)

	if p.exe == "" {
		p.exe = filepath.Join(p.srcDir, exeName)
	} else if p.os == "ios" || p.os == "android" {
		_, _ = fmt.Fprint(os.Stderr, "Parameter -executable is ignored for mobile builds.\n")
	}

	if p.name == "" {
		p.name = exeName
	}
	if p.icon == "" || p.icon == "Icon.png" {
		p.icon = filepath.Join(p.srcDir, "Icon.png")
	}
	if !exists(p.icon) {
		return errors.New("Missing application icon at \"" + p.icon + "\"")
	}

	// only used for iOS and macOS
	if p.appID == "" {
		if p.os == "darwin" {
			p.appID = "com.example." + p.name
		} else if p.os == "ios" || p.os == "android" {
			return errors.New("Missing appID parameter for mobile package")
		}
	}

	return nil
}

func (p *packager) doPackage() error {
	if !exists(p.exe) && !p.isMobile() {
		err := p.buildPackage()
		if err != nil {
			return errors.Wrap(err, "Error building application")
		}
		if !exists(p.exe) {
			return fmt.Errorf("unable to build directory to expected executable, %s", p.exe)
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
	case "ios":
		return p.packageIOS()
	default:
		return fmt.Errorf("unsupported target operating system \"%s\"", p.os)
	}
}

func (p *packager) isMobile() bool {
	return p.os == "ios" || strings.HasPrefix(p.os, "android")
}
