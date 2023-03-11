package commands

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"

	"github.com/jackmordaunt/icns/v2"
)

type darwinData struct {
	Name, ExeName string
	AppID         string
	Version       string
	Build         int
	Category      string
}

func (p *Packager) packageDarwin() (err error) {
	appDir := util.EnsureSubDir(p.dir, p.Name+".app")
	exeName := filepath.Base(p.exe)

	contentsDir := util.EnsureSubDir(appDir, "Contents")
	info := filepath.Join(contentsDir, "Info.plist")
	infoFile, err := os.Create(info)
	if err != nil {
		return fmt.Errorf("failed to create plist template: %w", err)
	}
	defer func() {
		if r := infoFile.Close(); r != nil && err == nil {
			err = r
		}
	}()

	tplData := darwinData{Name: p.Name, ExeName: exeName, AppID: p.AppID, Version: p.AppVersion, Build: p.AppBuild,
		Category: strings.ToLower(p.category)}
	if err := templates.InfoPlistDarwin.Execute(infoFile, tplData); err != nil {
		return fmt.Errorf("failed to write plist template: %w", err)
	}

	macOSDir := util.EnsureSubDir(contentsDir, "MacOS")
	binName := filepath.Join(macOSDir, exeName)
	if err := util.CopyExeFile(p.exe, binName); err != nil {
		return fmt.Errorf("failed to copy executable: %w", err)
	}

	resDir := util.EnsureSubDir(contentsDir, "Resources")
	icnsPath := filepath.Join(resDir, "icon.icns")

	img, err := os.Open(p.icon)
	if err != nil {
		return fmt.Errorf("failed to open source image \"%s\": %w", p.icon, err)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return fmt.Errorf("failed to decode source image: %w", err)
	}
	dest, err := os.Create(icnsPath)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %w", err)
	}
	defer func() {
		if r := dest.Close(); r != nil && err == nil {
			err = r
		}
	}()
	if err := icns.Encode(dest, srcImg); err != nil {
		return fmt.Errorf("failed to encode icns: %w", err)
	}

	return nil
}
