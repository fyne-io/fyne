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
}

var customMetadata map[string]string
var customMetadataSet bool

// SetMetadata overrides the packaged application metadata.
// This data can be used in many places like notifications and about screens.
func SetMetadata(m fyne.AppMetadata) {
	meta = m
}

func (a *fyneApp) Metadata() fyne.AppMetadata {
	return meta
}

func (a *fyneApp) CustomMetadata(name string) (custom string, exist bool) {
	custom, exist = customMetadata[name]
	return
}

// SetCustomMetadata allow fyne to set the custom metadata map once at the start of your application
// You shouldn't use it if you do not know what you are doing
//
// Since: 2.3
func SetCustomMetadata(custom map[string]string) bool {
	if customMetadataSet {
		return false
	}
	customMetadata = custom
	customMetadataSet = true
	return true
}
