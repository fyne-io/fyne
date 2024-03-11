package app

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/metadata"
)

var meta = fyne.AppMetadata{
	ID:      "",
	Name:    "",
	Version: "0.0.1",
	Build:   1,
	Release: false,
	Custom:  map[string]string{},
}

// SetMetadata overrides the packaged application metadata.
// This data can be used in many places like notifications and about screens.
func SetMetadata(m fyne.AppMetadata) {
	meta = m

	if meta.Custom == nil {
		meta.Custom = map[string]string{}
	}
}

func (a *fyneApp) Metadata() fyne.AppMetadata {
	if meta.ID == "" && meta.Name == "" {
		file, _ := filepath.Abs("FyneApp.toml")
		ref, err := os.Open(file)
		if err != nil { // no worries, this is just an optional fallback
			return meta
		}

		data, err := metadata.Load(ref)
		if err != nil || data == nil {
			fyne.LogError("failed to parse FyneApp.toml", err)
			return meta
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

	return meta
}
