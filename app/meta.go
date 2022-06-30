package app

import (
	"fyne.io/fyne/v2"
)

var meta = fyne.AppMetadata{
	ID:      "",
	Name:    "",
	Version: "0.0.1",
	Build:   1,
}

// SetMetadata overrides the packaged application metadata.
// This data can be used in many places like notifications and about screens.
func SetMetadata(m fyne.AppMetadata) {
	meta = m
}

func (a *fyneApp) Metadata() fyne.AppMetadata {
	return meta
}
