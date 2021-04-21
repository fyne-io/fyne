// +build ios android

package dialog

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/gomobile"
	"fyne.io/fyne/v2/storage"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return nil
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
		gomobile.ShowFolderOpenPicker(f.callback.(func(fyne.ListableURI, error)))
	} else {
		gomobile.ShowFileOpenPicker(f.callback.(func(fyne.URIReadCloser, error)), f.filter)
	}
	return true
}

func fileSaveOSOverride(f *FileDialog) bool {
	gomobile.ShowFileSavePicker(f.callback.(func(fyne.URIWriteCloser, error)), f.filter, f.initialFileName)

	return true
}

func getFavoriteLocations() (map[string]fyne.ListableURI, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homeURI := storage.NewFileURI(homeDir)

	favoriteNames := getFavoriteOrder()
	favoriteLocations := make(map[string]fyne.ListableURI)
	for _, favName := range favoriteNames {
		var uri fyne.URI
		var err1 error
		if favName == "Home" {
			uri = homeURI
		} else {
			uri, err1 = storage.Child(homeURI, favName)
		}
		if err1 != nil {
			err = err1
			continue
		}

		listURI, err1 := storage.ListerForURI(uri)
		if err1 != nil {
			err = err1
			continue
		}
		favoriteLocations[favName] = listURI
	}

	return favoriteLocations, err
}
