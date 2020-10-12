// +build ios android

package dialog

import (
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver/gomobile"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return nil
}

func isHidden(file, _ string) bool {
	return false
}

func fileOpenOSOverride(f *FileDialog) bool {
	gomobile.ShowFileOpenPicker(f.callback.(func(fyne.URIReadCloser, error)), f.filter)
	return true
}

func fileSaveOSOverride(f *FileDialog) bool {
	ShowInformation("File Save", "File save not available on mobile", f.parent)

	callback := f.callback.(func(fyne.URIWriteCloser, error))
	if callback != nil {
		callback(nil, nil)
	}

	return true
}

func getFavoriteLocations() (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Home":      homeDir,
		"Documents": filepath.Join(homeDir, "Documents"),
		"Downloads": filepath.Join(homeDir, "Downloads"),
	}, nil
}
