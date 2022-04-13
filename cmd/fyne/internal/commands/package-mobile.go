package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/mobile"
	"fyne.io/fyne/v2/cmd/fyne/internal/templates"

	"golang.org/x/sys/execabs"
)

func (p *Packager) packageAndroid(arch string) error {
	return mobile.RunNewBuild(
		mobile.WithTarget(arch),
		mobile.WithBundleID(p.appID),
		mobile.WithIcon(p.icon),
		mobile.WithName(p.name),
		mobile.WithVersion(p.appVersion),
		mobile.WithBuild(p.appBuild),
		mobile.WithRelease(p.release),
		mobile.WithAndroidAPI(p.androidAPI),
	)
}

func (p *Packager) packageIOS(target string) error {
	err := mobile.RunNewBuild(
		mobile.WithTarget(target),
		mobile.WithBundleID(p.appID),
		mobile.WithIcon(p.icon),
		mobile.WithName(p.name),
		mobile.WithVersion(p.appVersion),
		mobile.WithBuild(p.appBuild),
		mobile.WithRelease(p.release),
		mobile.WithCert(p.certificate),
		mobile.WithProfile(p.profile),
		mobile.WithIOSVersion(p.iosVersion),
	)
	if err != nil {
		return err
	}

	assetDir := util.EnsureSubDir(p.dir, "Images.xcassets")
	defer os.RemoveAll(assetDir)
	err = ioutil.WriteFile(filepath.Join(assetDir, "Contents.json"), []byte(`{
  "info" : {
    "author" : "xcode",
    "version" : 1
  }
}`), 0644)
	if err != nil {
		fyne.LogError("Content err", err)
	}

	iconDir := util.EnsureSubDir(assetDir, "AppIcon.appiconset")
	contentFile, _ := os.Create(filepath.Join(iconDir, "Contents.json"))

	err = templates.XCAssetsDarwin.Execute(contentFile, nil)
	if err != nil {
		return fmt.Errorf("failed to write xcassets content template: %w", err)
	}

	if err = copyResizeIcon(1024, iconDir, p.icon); err != nil {
		return err
	}
	if err = copyResizeIcon(180, iconDir, p.icon); err != nil {
		return err
	}
	if err = copyResizeIcon(120, iconDir, p.icon); err != nil {
		return err
	}
	if err = copyResizeIcon(76, iconDir, p.icon); err != nil {
		return err
	}
	if err = copyResizeIcon(152, iconDir, p.icon); err != nil {
		return err
	}

	appDir := filepath.Join(p.dir, mobile.AppOutputName(p.os, p.name))
	return runCmdCaptureOutput("xcrun", "actool", "Images.xcassets", "--compile", appDir, "--platform",
		"iphoneos", "--target-device", "iphone", "--minimum-deployment-target", "9.0", "--app-icon", "AppIcon",
		"--output-format", "human-readable-text", "--output-partial-info-plist", "/dev/null")
}

func copyResizeIcon(size int, dir, source string) error {
	strSize := strconv.Itoa(size)
	path := dir + "/Icon_" + strSize + ".png"
	return runCmdCaptureOutput("sips", "-o", path, "-Z", strSize, source)
}

// runCmdCaptureOutput is a exec.Command wrapper that offers better error messages from stdout and stderr.
func runCmdCaptureOutput(name string, args ...string) error {
	var (
		outbuf = &bytes.Buffer{}
		errbuf = &bytes.Buffer{}
	)
	cmd := execabs.Command(name, args...)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf
	err := cmd.Run()
	if err != nil {
		outstr := outbuf.String()
		errstr := errbuf.String()
		if outstr != "" {
			err = fmt.Errorf(outbuf.String()+": %w", err)
		}
		if errstr != "" {
			err = fmt.Errorf(outbuf.String()+": %w", err)
		}
		return err
	}
	return nil
}
