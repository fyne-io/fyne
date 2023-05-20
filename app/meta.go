package app

import (
	"fyne.io/fyne/v2"
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
	return meta
}
