package commands

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"

	"github.com/jackmordaunt/icns/v2"
)

type darwinData struct {
	Name, ExeName string
	AppID         string
	Version       string
	Build         int
	Category      string
	Languages     []string
}

func darwinLangs(langs []string) []string {
	r := make([]string, len(langs))
	for n, lang := range langs {
		r[n] = strings.Replace(lang, "-", "_", 1)
	}
	return r
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
		Category: strings.ToLower(p.category), Languages: darwinLangs(p.langs)}
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
	if !p.rawIcon {
		srcImg = processMacOSIcon(srcImg)
	}
	if err := icns.Encode(dest, srcImg); err != nil {
		return fmt.Errorf("failed to encode icns: %w", err)
	}

	return nil
}

func processMacOSIcon(in image.Image) image.Image {
	size := 1024
	border := 100
	radius := 185.4

	innerSize := int(float64(size) - float64(border*2)) // how many pixels inside border
	sized := resize.Resize(uint(innerSize), uint(innerSize), in, resize.Lanczos3)

	dc := gg.NewContext(size, size)
	dc.DrawRoundedRectangle(float64(border), float64(border), float64(innerSize), float64(innerSize), radius)
	dc.SetColor(color.Black)
	dc.Fill()
	mask := dc.AsMask()

	dc = gg.NewContext(size, size)
	_ = dc.SetMask(mask) // ignore error if size was not equal, as it is
	dc.DrawImage(sized, border, border)

	return dc.Image()
}
