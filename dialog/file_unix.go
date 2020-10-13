// +build !windows,!android,!ios

package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return []fyne.CanvasObject{makeFavoriteButton("Computer", theme.ComputerIcon(), func() {
		lister, err := storage.ListerForURI(storage.NewURI("file:///"))
		if err != nil {
			fyne.LogError("could not create lister for /", err)
			return
		}
		f.setLocation(lister)
	})}
}

func isHidden(file fyne.URI) bool {
	if file.Scheme() != "file" {
		fyne.LogError("Cannot check if non file is hidden", nil)
		return false
	}
	filePath := file.String()[len(file.Scheme()+3):]
	return len(filePath) == 0 || filePath[0] == '.'
}

func fileOpenOSOverride(*FileDialog) bool {
	return false
}

func fileSaveOSOverride(*FileDialog) bool {
	return false
}
