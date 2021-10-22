package commands

import (
	"image"
	"os"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	ico "github.com/Kodeworks/golang-image-ico"
	"github.com/josephspurrier/goversioninfo"
	"github.com/pkg/errors"
	"golang.org/x/sys/execabs"
)

type windowsData struct {
	Name            string
	CombinedVersion string
}

func (p *Packager) packageWindows() error {
	exePath := filepath.Dir(p.exe)

	// convert icon
	img, err := os.Open(p.icon)
	if err != nil {
		return errors.Wrap(err, "Failed to open source image")
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return errors.Wrap(err, "Failed to decode source image")
	}

	icoPath := filepath.Join(exePath, p.name+".ico")
	file, err := os.Create(icoPath)
	if err != nil {
		return errors.Wrap(err, "Failed to open image file")
	}

	err = ico.Encode(file, srcImg)
	if err != nil {
		return errors.Wrap(err, "Failed to encode icon")
	}

	err = file.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close image file")
	}

	// write manifest
	manifest := p.exe + ".manifest"
	manifestGenerated := false
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		manifestGenerated = true
		manifestFile, _ := os.Create(manifest)

		tplData := windowsData{
			Name:            p.name,
			CombinedVersion: p.combinedVersion(),
		}
		err := templates.ManifestWindows.Execute(manifestFile, tplData)
		if err != nil {
			return errors.Wrap(err, "Failed to write manifest template")
		}
	}

	// launch rsrc to generate the object file
	outPath := filepath.Join(exePath, "fyne.syso")

	vi := &goversioninfo.VersionInfo{}
	vi.ProductName = p.name
	vi.IconPath = icoPath
	vi.ManifestPath = manifest

	vi.Build()
	vi.Walk()

	arch, ok := os.LookupEnv("GOARCH")
	if !ok {
		arch = runtime.GOARCH
	}

	err = vi.WriteSyso(outPath, arch)
	if err != nil {
		return errors.Wrap(err, "Failed to write .syso file")
	}
	defer os.Remove(outPath)

	err = os.Remove(icoPath)
	if err != nil {
		return errors.Wrap(err, "Failed to remove icon")
	} else if manifestGenerated {
		err := os.Remove(manifest)
		if err != nil {
			return errors.Wrap(err, "Failed to remove manifest")
		}
	}

	err = p.buildPackage()
	if err != nil {
		return errors.Wrap(err, "Failed to rebuild after adding metadata")
	}

	appPath := p.exe
	appName := filepath.Base(p.exe)
	if filepath.Base(p.exe) != p.name {
		appName = p.name
		if filepath.Ext(p.name) != ".exe" {
			appName = appName + ".exe"
		}
		appPath = filepath.Join(filepath.Dir(p.exe), appName)
		os.Rename(filepath.Base(p.exe), appName)
	}

	if p.install {
		err := runAsAdminWindows("copy", "\"\""+appPath+"\"\"", "\"\""+filepath.Join(p.dir, appName)+"\"\"")
		if err != nil {
			return errors.Wrap(err, "Failed to run as administrator")
		}
	}
	return nil
}

func runAsAdminWindows(args ...string) error {
	cmd := "\"/c\""

	for _, arg := range args {
		cmd += ",\"" + arg + "\""
	}

	return execabs.Command("powershell.exe", "Start-Process", "cmd.exe", "-Verb", "runAs", "-ArgumentList", cmd).Run()
}
