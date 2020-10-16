package commands

import (
	"flag"
	"fmt"
	"os"
)

// Declare conformity to Command interface
var _ Command = (*releaser)(nil)

type releaser struct {
	packager
}

// NewReleaser returns a command that can adapt app packages for distribution
func NewReleaser() Command {
	return &releaser{}
}

func (r *releaser) AddFlags() {
	flag.StringVar(&r.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, freebsd, ios, linux, netbsd, openbsd, windows)")
	flag.StringVar(&r.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&r.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
}

func (r *releaser) PrintHelp(indent string) {
	fmt.Println(indent, "The release command prepares an application for public distribution.")
}

func (r *releaser) Run(params []string) {
	r.packager.release = true

	err := r.beforePackage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	r.packager.Run(params)
	err = r.afterPackage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func (r *releaser) afterPackage() error {
	return nil
}

func (r *releaser) beforePackage() error {
	return nil
}
