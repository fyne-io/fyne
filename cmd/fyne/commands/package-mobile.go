package commands

import (
	"path/filepath"

	"fyne.io/fyne/cmd/fyne/internal/mobile"
	"fyne.io/fyne/cmd/fyne/internal/util"
)

func (p *packager) packageAndroid(arch string) error {
	return mobile.RunNewBuild(arch, p.appID, p.icon, p.name, p.appVersion, p.appBuild, p.release)
}

func (p *packager) packageIOS() error {
	err := mobile.RunNewBuild("ios", p.appID, p.icon, p.name, p.appVersion, p.appBuild, p.release)
	if err != nil {
		return err
	}

	appDir := filepath.Join(p.dir, mobile.AppOutputName("os", p.name))
	iconPath := filepath.Join(appDir, "Icon.png")
	return util.CopyFile(p.icon, iconPath)
}
