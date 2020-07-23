// +build !windows,!android,!ios

package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
	return []fyne.CanvasObject{makeFavoriteButton("Computer", theme.ComputerIcon(), func() {
		f.setDirectory("/")
	})}
}

func isHidden(file, _ string) bool {
	return len(file) == 0 || file[0] == '.'
}

func fileOpenOSOverride(*FileDialog) bool {
	return false
}

func fileSaveOSOverride(*FileDialog) bool {
	return false
}
