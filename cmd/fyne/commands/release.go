package commands

import (
	"errors"
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

	keyStore string
}

// NewReleaser returns a command that can adapt app packages for distribution
func NewReleaser() Command {
	return &releaser{}
}

func (r *releaser) AddFlags() {
	flag.StringVar(&r.os, "os", "", "The operating system to target (android, android/arm, android/arm64, android/amd64, android/386, darwin, freebsd, ios, linux, netbsd, openbsd, windows)")
	flag.StringVar(&r.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&r.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.StringVar(&r.appVersion, "appVersion", "", "Version number in the form x, x.y or x.y.z semantic version")
	flag.IntVar(&r.appBuild, "appBuild", 0, "Build number, should be greater than 0 and incremented for each build")
	flag.StringVar(&r.keyStore, "keyStore", "", "Android: location of .keystore file containing signing information")
}

func (r *releaser) PrintHelp(indent string) {
	fmt.Println(indent, "The release command prepares an application for public distribution.")
}

func (r *releaser) Run(params []string) {
	if err := r.validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	r.packager.release = true

	if err := r.beforePackage(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	r.packager.Run(params)
	if err := r.afterPackage(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func (r *releaser) afterPackage() error {
	if util.IsAndroid(r.os) {
		target := mobile.AppOutputName(r.os, r.packager.name)
		apk := filepath.Join(r.dir, target)
		if err := r.zipAlign(apk); err != nil {
			return err
		}
		return r.sign(apk)
	}

	return nil
}

func (r *releaser) beforePackage() error {
	if util.IsAndroid(r.os) {
		if err := util.RequireAndroidSDK(); err != nil {
			return err
		}
	}

	return nil
}

func (r *releaser) sign(path string) error {
	signer := filepath.Join(util.AndroidBuildToolsPath(), "/apksigner")
	cmd := exec.Command(signer, "sign", "--ks", r.keyStore, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (r *releaser) validate() error {
	if util.IsMobile(r.os) || r.os == "windows" {
		if r.appVersion == "" { // Here it is required, if provided then package validate will check format
			return errors.New("missing required -appVersion parameter")
		}
		if r.appBuild <= 0 {
			return errors.New("missing required -appBuild parameter")
		}
	}
	if util.IsAndroid(r.os) {
		if r.keyStore == "" {
			return errors.New("missing required -keyStore parameter for android release")
		}
	}
	return nil
}

func (r *releaser) zipAlign(path string) error {
	unaligned := filepath.Join(filepath.Dir(path), "unaligned.apk")
	err := os.Rename(path, unaligned)
	if err != nil {
		return nil
	}

	cmd := filepath.Join(util.AndroidBuildToolsPath(), "zipalign")
	err = exec.Command(cmd, "4", unaligned, path).Run()
	if err != nil {
		_ = os.Rename(path, unaligned) // ignore error, return previous
		return err
	}
	return os.Remove(unaligned)
}
