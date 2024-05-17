package app

import (
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/build"
	"fyne.io/fyne/v2/internal/metadata"
)

func checkLocalMetadata() {
	if build.NoMetadata || build.Mode == fyne.BuildRelease {
		return
	}

	dir := getProjectPath()
	file := filepath.Join(dir, "FyneApp.toml")
	ref, err := os.Open(file)
	if err != nil { // no worries, this is just an optional fallback
		return
	}
	defer ref.Close()

	data, err := metadata.Load(ref)
	if err != nil || data == nil {
		fyne.LogError("failed to parse FyneApp.toml", err)
		return
	}

	meta.ID = data.Details.ID
	meta.Name = data.Details.Name
	meta.Version = data.Details.Version
	meta.Build = data.Details.Build

	if data.Details.Icon != "" {
		res, err := fyne.LoadResourceFromPath(data.Details.Icon)
		if err == nil {
			meta.Icon = res
		}
	}

	meta.Release = false
	meta.Custom = data.Development
}

func getProjectPath() string {
	exe, err := os.Executable()
	work, _ := os.Getwd()

	if err != nil {
		fyne.LogError("failed to lookup build executable", err)
		return work
	}

	temp := os.TempDir()
	if strings.Contains(exe, temp) { // this happens with "go run"
		return work
	}

	// we were called with an executable from "go build"
	return filepath.Dir(exe)
}
