package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne"
)

// Declare conformity to command interface
var _ command = (*installer)(nil)

type installer struct {
	installDir, srcDir, icon, os string
	packager                     *packager
}

func (i *installer) addFlags() {
	flag.StringVar(&i.installDir, "installDir", "", "A specific location to install to, rather than the OS default")
	flag.StringVar(&i.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&i.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, ios, linux, windows)")
}

func (i *installer) printHelp(indent string) {
	fmt.Println(indent, "The install command packages an application for the current platform and copies it")
	fmt.Println(indent, "into the system location for applications. This can be overridden with installDir")
	fmt.Println(indent, "Command usage: fyne install [parameters]")
}

func (i *installer) validate() error {
	if i.os == "" {
		i.os = runtime.GOOS
	}

	i.packager = &packager{os: runtime.GOOS, install: true, srcDir: i.srcDir}
	i.packager.icon = i.icon
	i.packager.os = i.os
	return i.packager.validate()
}

func (i *installer) install() error {
	p := i.packager

	switch i.os {
	}

	switch i.os {
	case "darwin":
		if i.installDir == "" {
			i.installDir = "/Applications"
		}
	case "linux":
		if i.installDir == "" {
			i.installDir = "/" // the tarball contains the structure starting at usr/local
		}
	case "windows":
		if i.installDir == "" {
			i.installDir = filepath.Join(os.Getenv("ProgramFiles"), p.name)
		}
		err := runAsAdminWindows("mkdir", "\"\""+i.installDir+"\"\"")
		if err != nil {
			fyne.LogError("Failed to run as windows administrator", err)
			return err
		}
	case "android/arm", "android/arm64", "android/amd64", "android/386", "android", "ios":
	default:
		return errors.New("Unsupported target operating system \"" + p.os + "\"")
	}

	p.dir = i.installDir
	err := p.doPackage()
	if err == nil {
		if strings.Index(i.os, "android") == 0 {
			err = i.installAndroid()
		} else if i.os == "ios" {
			err = i.installIOS()
		}
	}
	return err
}

func (i *installer) installAndroid() error {
	return nil
}

func (i *installer) installIOS() error {
	return nil
}

func (i *installer) run(args []string) {
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
