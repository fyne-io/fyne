//go:build ios || android

package dialog

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/mobile"
	"fyne.io/fyne/v2/storage"
)

func (f *fileDialog) getPlaces() []favoriteItem {
	return []favoriteItem{}
}

func isHidden(file fyne.URI) bool {
	if file.Scheme() != "file" {
		fyne.LogError("Cannot check if non file is hidden", nil)
		return false
	}
	return false
}

func hideFile(filename string) error {
	return nil
}

func fileOpenOSOverride(f *FileDialog) bool {
	if f.isDirectory() {
		mobile.ShowFolderOpenPicker(f.callback.(func(fyne.ListableURI, error)))
	} else {
		mobile.ShowFileOpenPicker(f.callback.(func(fyne.URIReadCloser, error)), f.filter)
	}
	return true
}

func fileSaveOSOverride(f *FileDialog) bool {
	mobile.ShowFileSavePicker(f.callback.(func(fyne.URIWriteCloser, error)), f.filter, f.initialFileName)
	return true
}

func getFavoriteLocation(homeURI fyne.URI, name string) (fyne.URI, error) {
	return storage.Child(homeURI, name)
}
