package commands

import "fyne.io/fyne/v2/internal/metadata"

type appData struct {
	icon, Name        string
	AppID, AppVersion string
	AppBuild          int
	ResGoString       string
	Release, rawIcon  bool
	CustomMetadata    map[string]string
	Migrations        map[string]bool
	VersionAtLeast2_3 bool
}

func (a *appData) appendCustomMetadata(fromFile map[string]string) {
	if a.CustomMetadata == nil {
		a.CustomMetadata = map[string]string{}
	}

	for key, value := range fromFile {
		_, ok := a.CustomMetadata[key]
		if ok {
			continue
		}
		a.CustomMetadata[key] = value
	}
}

func (a *appData) mergeMetadata(data *metadata.FyneApp) {
	if a.icon == "" {
		a.icon = data.Details.Icon
	}
	if a.Name == "" {
		a.Name = data.Details.Name
	}
	if a.AppID == "" {
		a.AppID = data.Details.ID
	}
	if a.AppVersion == "" {
		a.AppVersion = data.Details.Version
	}
	if a.AppBuild == 0 {
		a.AppBuild = data.Details.Build
	}
	if a.Release {
		a.appendCustomMetadata(data.Release)
	} else {
		a.appendCustomMetadata(data.Development)
	}
	a.Migrations = data.Migrations
}
