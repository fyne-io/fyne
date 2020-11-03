package commands

import (
	"os/exec"
	"path/filepath"

	"fyne.io/fyne/cmd/fyne/internal/mobile"
	"fyne.io/fyne/cmd/fyne/internal/util"
)

func (p *packager) packageAndroid(arch string) error {
	return mobile.RunNewBuild(arch, p.appID, p.icon, p.name, p.appVersion, p.appBuild, p.release, "", "")
}

func (p *packager) packageIOS() error {
	err := mobile.RunNewBuild("ios", p.appID, p.icon, p.name, p.appVersion, p.appBuild, p.release, p.certificate, p.profile)
	if err != nil {
		return err
	}

	appDir := filepath.Join(p.dir, mobile.AppOutputName(p.os, p.name))
	iconPath := filepath.Join(appDir, "Icon.png")
	if err = util.CopyFile(p.icon, iconPath); err != nil {
		return err
	}

	if err = copyResizeIcon("120", appDir); err != nil {
		return err
	}
	if err = copyResizeIcon("76", appDir); err != nil {
		return err
	}
	return copyResizeIcon("152", appDir)
}

func copyResizeIcon(size, dir string) error {
	cmd := exec.Command("sips", "-o", dir+"/Icon_"+size+".png", "-Z", size, dir+"/Icon.png")
	return cmd.Run()
}
