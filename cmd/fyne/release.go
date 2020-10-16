package main

import (
	"flag"
	"fmt"
	"os"
)

// Declare conformity to command interface
var _ command = (*releaser)(nil)

type releaser struct {
	packager
}

func (r *releaser) addFlags() {
	flag.StringVar(&r.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, freebsd, ios, linux, netbsd, openbsd, windows)")
	flag.StringVar(&r.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&r.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
}

func (r *releaser) afterPackage() error {
	return nil
}

func (r *releaser) beforePackage() error {
	return nil
}

func (r *releaser) printHelp(indent string) {
	fmt.Println(indent, "The release command prepares an application for public distribution.")
}

func (r *releaser) run(params []string) {
	r.packager.release = true

	err := r.beforePackage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	r.packager.run(params)
	err = r.afterPackage()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
