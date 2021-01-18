package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/mobile"
	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
	"github.com/pkg/errors"
)

func (p *packager) packageAndroid(arch string) error {
	return mobile.RunNewBuild(arch, p.appID, p.icon, p.name, p.appVersion, p.appBuild, p.release, "", "")
}

func (p *packager) packageIOS() error {
	err := mobile.RunNewBuild("ios", p.appID, p.icon, p.name, p.appVersion, p.appBuild, p.release, p.certificate, p.profile)
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
		return errors.Wrap(err, "Failed to write xcassets content template")
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
	cmd := exec.Command("xcrun", "actool", "Images.xcassets", "--compile", appDir, "--platform",
		"iphoneos", "--target-device", "iphone", "--minimum-deployment-target", "9.0", "--app-icon", "AppIcon",
		"--output-partial-info-plist", "/dev/null")
	return cmd.Run()
}

func copyResizeIcon(size int, dir, source string) error {
	path := fmt.Sprintf("%s/Icon_%d.png", dir, size)
	strSize := fmt.Sprintf("%d", size)
	return exec.Command("sips", "-o", path, "-Z", strSize, source).Run()
}
