package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/cmd/fyne/internal/mobile"
)

// Declare conformity to Command interface
var _ Command = (*installer)(nil)

type installer struct {
	installDir, srcDir, icon, os, appID string
	packager                            *packager
	release                             bool
}

// NewInstaller returns an install command that can install locally built Fyne apps.
func NewInstaller() Command {
	return &installer{}
}

func (i *installer) AddFlags() {
	flag.StringVar(&i.os, "os", "", "The mobile platform to target (android, android/arm, android/arm64, android/amd64, android/386, ios)")
	flag.StringVar(&i.installDir, "installDir", "", "A specific location to install to, rather than the OS default")
	flag.StringVar(&i.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&i.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.BoolVar(&i.release, "release", false, "Should this package be installed in release mode? (disable debug etc)")
}

func (i *installer) PrintHelp(indent string) {
	fmt.Println(indent, "The install command packages an application for the current platform and copies it")
	fmt.Println(indent, "into the system location for applications. This can be overridden with installDir")
	fmt.Println(indent, "Command usage: fyne install [parameters]")
}

func (i *installer) Run(args []string) {
	if len(args) != 0 {
		fyne.LogError("Unexpected parameter after flags", nil)
		return
	}

	err := i.validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	err = i.install()
	if err != nil {
		fyne.LogError("Unable to install application", err)
		os.Exit(1)
	}
}

func (i *installer) install() error {
	p := i.packager

	if i.os != "" {
		if i.os == "ios" {
			return i.installIOS()
		} else if strings.Index(i.os, "android") == 0 {
			return i.installAndroid()
		}

		return errors.New("Unsupported target operating system \"" + i.os + "\"")
	}

	if i.installDir == "" {
		switch p.os {
		case "darwin":
			i.installDir = "/Applications"
		case "linux", "openbsd", "freebsd", "netbsd":
			i.installDir = "/" // the tarball contains the structure starting at usr/local
		case "windows":
			i.installDir = filepath.Join(os.Getenv("ProgramFiles"), p.name)
			err := runAsAdminWindows("mkdir", "\"\""+filepath.Join(os.Getenv("ProgramFiles"), p.name)+"\"\"")
			if err != nil {
				fyne.LogError("Failed to run as windows administrator", err)
				return err
			}
		default:
			return errors.New("Unsupported target operating system \"" + p.os + "\"")
		}
	}

	p.dir = i.installDir
	return p.doPackage()
}

func (i *installer) installAndroid() error {
	target := mobile.AppOutputName(i.os, i.packager.name)

	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		err := i.packager.doPackage()
		if err != nil {
			return nil
		}
	}

	return i.runMobileInstall("adb", target, "install")
}

func (i *installer) installIOS() error {
	target := mobile.AppOutputName(i.os, i.packager.name)
	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		err := i.packager.doPackage()
		if err != nil {
			return nil
		}
	}

	return i.runMobileInstall("ios-deploy", target, "--bundle")
}

func (i *installer) runMobileInstall(tool, target string, args ...string) error {
	_, err := exec.LookPath(tool)
	if err != nil {
		return err
	}

	cmd := exec.Command(tool, append(args, target)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (i *installer) validate() error {
	os := i.os
	if os == "" {
		os = targetOS()
	}
	i.packager = &packager{appID: i.appID, os: os, install: true, srcDir: i.srcDir}
	i.packager.icon = i.icon
	i.packager.release = i.release
	return i.packager.validate()
}
