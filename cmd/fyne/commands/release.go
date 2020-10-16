package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/cmd/fyne/internal/mobile"
	"fyne.io/fyne/cmd/fyne/internal/util"
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
	if util.IsAndroid(r.os) {
		target := mobile.AppOutputName(r.os, r.packager.name)
		return r.zipAlign(filepath.Join(r.dir, target))
	}

	return nil
}

func (r *releaser) beforePackage() error {
	if util.IsAndroid(r.os) {
		return util.RequireAndroidSDK()
	}

	return nil
}

func (r *releaser) zipAlign(path string) error {
	buildTools := util.AndroidBuildToolsPath()
	unaligned := filepath.Join(filepath.Dir(path), "unaligned.apk")
	err := os.Rename(path, unaligned)
	if err != nil {
		return nil
	}

	cmd := filepath.Join(buildTools, "zipalign")
	err = exec.Command(cmd, "4", unaligned, path).Run()
	if err != nil {
		_ = os.Rename(path, unaligned) // ignore error, return previous
		return err
	}
	return os.Remove(unaligned)
}
