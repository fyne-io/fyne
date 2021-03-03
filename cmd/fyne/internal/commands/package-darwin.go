package commands

import (
	"image"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
	"github.com/jackmordaunt/icns"
	"github.com/pkg/errors"
)

type darwinData struct {
	Name, ExeName string
	AppID         string
	Version       string
	Build         int
	Category      string
}

func (p *Packager) packageDarwin() (err error) {
	appDir := util.EnsureSubDir(p.dir, p.name+".app")
	exeName := filepath.Base(p.exe)

	contentsDir := util.EnsureSubDir(appDir, "Contents")
	info := filepath.Join(contentsDir, "Info.plist")
	infoFile, err := os.Create(info)
	if err != nil {
		return errors.Wrap(err, "Failed to write plist template")
	}
	defer func() {
		if r := infoFile.Close(); r != nil && err == nil {
			err = r
		}
	}()

	tplData := darwinData{Name: p.name, ExeName: exeName, AppID: p.appID, Version: p.appVersion, Build: p.appBuild,
		Category: strings.ToLower(p.category)}
	if err := templates.InfoPlistDarwin.Execute(infoFile, tplData); err != nil {
		return errors.Wrap(err, "Failed to write plist template")
	}

	macOSDir := util.EnsureSubDir(contentsDir, "MacOS")
	binName := filepath.Join(macOSDir, exeName)
	if err := util.CopyExeFile(p.exe, binName); err != nil {
		return errors.Wrap(err, "Failed to copy exe file")
	}

	resDir := util.EnsureSubDir(contentsDir, "Resources")
	icnsPath := filepath.Join(resDir, "icon.icns")

	img, err := os.Open(p.icon)
	if err != nil {
		return errors.Wrapf(err, "Failed to open source image \"%s\"", p.icon)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return errors.Wrapf(err, "Failed to decode source image")
	}
	dest, err := os.Create(icnsPath)
	if err != nil {
		return errors.Wrap(err, "Failed to open destination file")
	}
	defer func() {
		if r := dest.Close(); r != nil && err == nil {
			err = r
		}
	}()
	if err := icns.Encode(dest, srcImg); err != nil {
		return errors.Wrap(err, "Failed to encode icns")
	}

	return nil
}
