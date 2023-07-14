package commands

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"github.com/fyne-io/image/ico"
	"github.com/josephspurrier/goversioninfo"
	"golang.org/x/sys/execabs"
)

type windowsData struct {
	Name            string
	CombinedVersion string
}

func (p *Packager) packageWindows(tags []string) error {
	exePath := filepath.Dir(p.exe)

	// convert icon
	img, err := os.Open(p.icon)
	if err != nil {
		return fmt.Errorf("failed to open source image: %w", err)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		return fmt.Errorf("failed to decode source image: %w", err)
	}

	icoPath := filepath.Join(exePath, "FyneApp.ico")
	file, err := os.Create(icoPath)
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}

	err = ico.Encode(file, srcImg)
	if err != nil {
		return fmt.Errorf("failed to encode icon: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close image file: %w", err)
	}

	// write manifest
	manifest := p.exe + ".manifest"
	manifestGenerated := false
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		manifestGenerated = true
		manifestFile, _ := os.Create(manifest)

		tplData := windowsData{
			Name:            encodeXMLString(p.Name),
			CombinedVersion: p.combinedVersion(),
		}
		err := templates.ManifestWindows.Execute(manifestFile, tplData)
		if err != nil {
			return fmt.Errorf("failed to write manifest template: %w", err)
		}
		manifestFile.Close()
	}

	// launch rsrc to generate the object file
	outPath := filepath.Join(exePath, "fyne.syso")

	vi := &goversioninfo.VersionInfo{}
	vi.ProductName = p.Name
	vi.IconPath = icoPath
	vi.ManifestPath = manifest
	vi.StringFileInfo.ProductVersion = p.combinedVersion()
	vi.StringFileInfo.FileDescription = p.Name
	vi.FixedFileInfo.FileVersion = fixedVersionInfo(p.combinedVersion())

	vi.Build()
	vi.Walk()

	arch, ok := os.LookupEnv("GOARCH")
	if !ok {
		arch = runtime.GOARCH
	}

	err = vi.WriteSyso(outPath, arch)
	if err != nil {
		return fmt.Errorf("failed to write .syso file: %w", err)
	}
	defer os.Remove(outPath)

	err = os.Remove(icoPath)
	if err != nil {
		return fmt.Errorf("failed to remove icon: %w", err)
	} else if manifestGenerated {
		err := os.Remove(manifest)
		if err != nil {
			return fmt.Errorf("failed to remove manifest: %w", err)
		}
	}

	_, err = p.buildPackage(nil, tags)
	if err != nil {
		return fmt.Errorf("failed to rebuild after adding metadata: %w", err)
	}

	appName := filepath.Base(p.exe)
	if filepath.Base(p.exe) != p.Name {
		appName = p.Name
		if filepath.Ext(p.Name) != ".exe" {
			appName = appName + ".exe"
		}
		os.Rename(filepath.Base(p.exe), appName)
	}

	if p.install {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to locate current working directory")
		}
		appPath := filepath.Join(wd, appName)

		err = runAsAdminWindows("copy", appPath, filepath.Join(p.dir, appName))
		if err != nil {
			return fmt.Errorf("failed to run as administrator: %w", err)
		}
	}
	return nil
}

func escapePowerShellArguments(args ...string) string {
	cmd := "\"/c\""

	for idx, arg := range args {
		if idx == 0 {
			cmd += ",'" + arg
		} else {
			cmd += " \"" + arg + "\""
		}
	}
	if len(args) != 0 {
		cmd += "'"
	}

	return cmd
}

func runAsAdminWindows(args ...string) error {
	cmd := escapePowerShellArguments(args...)

	return execabs.Command("powershell.exe", "Start-Process", "cmd.exe", "-Verb", "runAs", "-ArgumentList", cmd).Run()
}

func fixedVersionInfo(ver string) (ret goversioninfo.FileVersion) {
	ret.Build = 1 // as 0,0,0,0 is not valid
	if len(ver) == 0 {
		return ret
	}
	split := strings.Split(ver, ".")
	setVersionField(&ret.Major, split[0])
	if len(split) > 1 {
		setVersionField(&ret.Minor, split[1])
	}
	if len(split) > 2 {
		setVersionField(&ret.Patch, split[2])
	}
	if len(split) > 3 {
		setVersionField(&ret.Build, split[3])
	}
	return ret
}

func setVersionField(to *int, ver string) {
	num, err := strconv.Atoi(ver)
	if err != nil {
		fyne.LogError("Failed to parse app version field", err)
		return
	}

	*to = num
}
