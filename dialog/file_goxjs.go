//go:build wasm || js
// +build wasm js

package dialog

import (
	"fyne.io/fyne/v2"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return nil
}

func isHidden(file fyne.URI) bool {
	return false
}

func fileOpenOSOverride(f *FileDialog) bool {
	// TODO #2737
	return true
}

func fileSaveOSOverride(f *FileDialog) bool {
	// TODO #2738
	return true
}

func (f *fileDialog) getPlaces() []favoriteItem {
	return []favoriteItem{}
}

func getFavoriteLocations() (map[string]fyne.ListableURI, error) {
	favoriteLocations := make(map[string]fyne.ListableURI)

	return favoriteLocations, nil
}
