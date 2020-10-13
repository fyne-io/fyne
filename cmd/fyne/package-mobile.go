package main

import (
	"path/filepath"

	"fyne.io/fyne/cmd/fyne/internal/mobile"
)

func (p *packager) packageAndroid(arch string) error {
	return mobile.RunNewBuild(arch, p.appID, p.icon, p.name, p.release)
}

func (p *packager) packageIOS() error {
	err := mobile.RunNewBuild("ios", p.appID, p.icon, p.name, p.release)
	if err != nil {
		return err
	}

	appDir := filepath.Join(p.dir, p.name+".app")
	iconPath := filepath.Join(appDir, "Icon.png")
	return copyFile(p.icon, iconPath)
}
