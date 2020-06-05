package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne"
)

// Declare conformity to command interface
var _ command = (*installer)(nil)

type installer struct {
	installDir, srcDir, icon string
	packager                 *packager
}

func (i *installer) addFlags() {
	flag.StringVar(&i.installDir, "installDir", "", "A specific location to install to, rather than the OS default")
	flag.StringVar(&i.icon, "icon", "Icon.png", "The name of the application icon file")
}

func (i *installer) printHelp(indent string) {
	fmt.Println(indent, "The install command packages an application for the current platform and copies it")
	fmt.Println(indent, "into the system location for applications. This can be overridden with installDir")
	fmt.Println(indent, "Command usage: fyne install [parameters]")
}

func (i *installer) validate() error {
	i.packager = &packager{os: runtime.GOOS, install: true, srcDir: i.srcDir}
	i.packager.icon = i.icon
	return i.packager.validate()
}

func (i *installer) install() error {
	p := i.packager

	if i.installDir == "" {
		switch runtime.GOOS {
		case "darwin":
			i.installDir = "/Applications"
		case "linux":
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
