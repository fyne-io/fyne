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
	installDir string
}

func (i *installer) addFlags() {
	flag.StringVar(&i.installDir, "installDir", "", "A specific location to install to, rather than the OS default")
}

func (i *installer) printHelp(indent string) {
	fmt.Println(indent, "The install command packages an application for the current platform and copies it")
	fmt.Println(indent, "into the system location for applications. This can be overridden with installDir")
	fmt.Println(indent, "Command usage: fyne install [parameters]")
}

func (i *installer) install() error {
	p := &packager{os: runtime.GOOS, install: true}
	err := p.validate()
	if err != nil {
		return err
	}

	if i.installDir == "" {
		switch runtime.GOOS {
		case "darwin":
			i.installDir = "/Applications"
		case "linux":
			i.installDir = "/" // the tarball contains the structure starting at usr/local
		case "windows":
			i.installDir = filepath.Join(os.Getenv("ProgramFiles"), p.name)
			runAsAdminWindows("mkdir", "\"\""+filepath.Join(os.Getenv("ProgramFiles"), p.name)+"\"\"")
		default:
			return errors.New("Unsupported terget operating system \"" + p.os + "\"")
		}
	}

	p.dir = i.installDir
	err = p.doPackage()
	return err
}

func (i *installer) run(args []string) {
	if len(args) != 0 {
		fyne.LogError("Unexpected parameter after flags", nil)
		return
	}

	err := i.install()
	if err != nil {
		fyne.LogError("Unable to install application", err)
		os.Exit(1)
	}
}
