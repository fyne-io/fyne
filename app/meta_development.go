//go:build !no_metadata

package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/metadata"
)

func checkLocalMetadata() {
	file, _ := filepath.Abs("FyneApp.toml")
	ref, err := os.Open(file)
	if err != nil { // no worries, this is just an optional fallback
		return
	}

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
