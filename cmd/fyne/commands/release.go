package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/cmd/fyne/internal/mobile"
	"fyne.io/fyne/cmd/fyne/internal/templates"
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
	flag.StringVar(&r.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&r.icon, "icon", "Icon.png", "The name of the application icon file")
	flag.StringVar(&r.appID, "appID", "", "For ios or darwin targets an appID is required, for ios this must \nmatch a valid provisioning profile")
	flag.StringVar(&r.appVersion, "appVersion", "", "Version number in the form x, x.y or x.y.z semantic version")
	flag.IntVar(&r.appBuild, "appBuild", 0, "Build number, should be greater than 0 and incremented for each build")
	flag.StringVar(&r.keyStore, "keyStore", "", "Android: location of .keystore file containing signing information")
	flag.StringVar(&r.certificate, "certificate", "Apple Distribution", "iOS/macOS/Windows: name of the certificate to sign the build")
	flag.StringVar(&r.profile, "profile", "", "iOS/macOS: name of the provisioning profile for this release build")
	// password (for (windows?) certificate)
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
	if r.os == "ios" {
		team, err := mobile.DetectIOSTeamID(r.certificate)
		if err != nil {
			return errors.New("failed to determine team ID")
		}

		payload := filepath.Join(r.dir, "Payload")
		os.Mkdir(payload, 0750)
		defer os.Remove(payload)
		appName := mobile.AppOutputName(r.os, r.name)
		payloadAppDir := filepath.Join(payload, appName)
		os.Rename(filepath.Join(r.dir, appName), payloadAppDir)

		copyResizeIcon("120", payloadAppDir)
		copyResizeIcon("76", payloadAppDir)
		copyResizeIcon("152", payloadAppDir)

		entitlementPath := filepath.Join(r.dir, "entitlements.plist")
		defer os.Remove(entitlementPath)
		entitlements, _ := os.Create(entitlementPath)
		entitlementData := struct{ TeamID, AppID string }{
			TeamID: team,
			AppID:  r.appID,
		}
		err = templates.EntitlementsDarwin.Execute(entitlements, entitlementData)
		if err != nil {
			return errors.New("failed to write entitlements plist template")
		}

		cmd := exec.Command("codesign", "-f", "-vv", "-s", r.certificate, "--entitlements",
			"entitlements.plist", "Payload/"+appName+"/")
		if err := cmd.Run(); err != nil {
			fyne.LogError("Codesign failed", err)
			return errors.New("unable to codesign application bundle")
		}

		return exec.Command("zip", "-r", appName[:len(appName)-4]+".ipa", "Payload/").Run()
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
	} else if r.os == "ios" {
		if r.profile == "" {
			return errors.New("missing required -profile parameter for iOS release")
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

func copyResizeIcon(size, dir string) error {
	cmd := exec.Command("sips", "-o", dir+"/Icon_"+size+".png", "-Z", "120", dir+"/Icon.png")
	return cmd.Run()
}
